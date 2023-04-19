// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resttools

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// wellKnownTypes has a key for every well-known type message which JSON-encodes to an atomic symbol
// (like a number or a string) as opposed to a structured object. The value for each key is true iff
// the JSON encoding for the type is a string. We need both these data for well-known types so that
// we can properly decode them in URL paths and query strings.
var wellKnownTypes = map[string]bool{
	// == The following are the only three common types used in this API ==
	"google.protobuf.FieldMask": true,
	"google.protobuf.Timestamp": true,
	"google.protobuf.Duration":  true,
	// == End utilized types ==
	// TODO: When the following start being used in the Showcase API, add tests for each of
	// these. These types are defined in
	// https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/wrappers.proto
	"google.protobuf.DoubleValue": false,
	"google.protobuf.FloatValue":  false,
	"google.protobuf.Int64Value":  true,
	"google.protobuf.UInt64Value": true,
	"google.protobuf.Int32Value":  false,
	"google.protobuf.UInt32Value": false,
	"google.protobuf.BoolValue":   false,
	"google.protobuf.StringValue": true,
	"google.protobuf.BytesValue":  true,
	// TODO: Determine if the following are even viable as query params.
	"google.protobuf.Value":     false,
	"google.protobuf.ListValue": false,
}

// PopulateSingularFields sets the fields within protoMessage to the values provided in
// fieldValues. The fields and values are provided as a map of field paths to the string
// representation of their values. The field paths can refer to fields nested arbitrarily deep
// within protoMessage; each path element must be lower-camel-cased. This returns an error if any
// field path is not valid or if any value can't be parsed into the correct data type for the field.
func PopulateSingularFields(protoMessage proto.Message, fieldValues map[string]string) error {
	for name, value := range fieldValues {
		if err := PopulateOneField(protoMessage, name, []string{value}); err != nil {
			// TODO: accumulate errors so we report them all at once
			return err
		}
	}
	return nil
}

// PopulateFields sets the fields within protoMessage to the values provided in fieldValues. The
// fields and values are provided as a map of field paths to lists of string representations of
// their values. The field paths can refer to fields nested arbitrarily deep within protoMessage;
// each path element must be lower-camel-cased. `fieldValues` contains the list of values applicable
// to the field: a single value of singular fields, or ordered values for a repeated field. This
// returns an error if any field path is not valid or if any value can't be parsed into the correct
// data type for the field.
func PopulateFields(protoMessage proto.Message, fieldValues map[string][]string) error {
	for name, value := range fieldValues {
		if err := PopulateOneField(protoMessage, name, value); err != nil {
			// TODO: accumulate errors so we report them all at once
			return err
		}
	}
	return nil
}

// PopulateOneField finds in protoMessage the field identified by fieldPath (which could refer to an
// arbitrarily nested field using dotted notation; each path element must be lower-camel-cased) and
// sets it to `value`. It returns an error if the fieldPath does not properly reference a field, or
// if `fieldValues` could not be parsed into a list of the data type expected for the
// field. `fieldValues` is a slice to allow passing a series of repeated field entries.
func PopulateOneField(protoMessage proto.Message, fieldPath string, fieldValues []string) error {
	message := protoMessage.ProtoReflect()

	// TODO: Support repeated fields.
	if len(fieldValues) > 1 {
		return fmt.Errorf("repeated field values are not supported yet (culprit: %q: %v)", fieldPath, fieldValues)
	}
	value := fieldValues[0]

	levels := strings.Split(fieldPath, ".")
	lastLevel := len(levels) - 1
	for idx, fieldName := range levels {
		messageDescriptor := message.Descriptor()
		messageFullName := messageDescriptor.FullName()

		if messageDescriptor.Syntax() != protoreflect.Proto3 {
			return fmt.Errorf("cannot process %q as it does not use proto3 syntax", messageFullName)
		}

		if len(fieldName) == 0 {
			return fmt.Errorf("segment %d of path field %q is empty", idx, fieldPath)
		}

		// find field
		subFields := messageDescriptor.Fields()
		fieldDescriptor, err := findProtoFieldByJSONName(fieldName, subFields)
		if err != nil {
			return err
		}
		if fieldDescriptor == nil {
			return fmt.Errorf("could not find %dth field (%q) in field path %q in message %q",
				idx, fieldName, fieldPath, messageFullName)
		}

		if idx != lastLevel {
			// non-terminal field: go to the next level

			if kind := fieldDescriptor.Kind(); kind != protoreflect.MessageKind {
				return fmt.Errorf("%dth field (%q) in field path %q in message %q is not itself a message but a %q",
					idx, fieldName, fieldPath, messageFullName, kind)
			}
			message = message.Mutable(fieldDescriptor).Message()
			continue
		}

		// terminal field
		var (
			protoValue protoreflect.Value
			parseError error
		)
		kind := fieldDescriptor.Kind()
		switch kind {
		case protoreflect.MessageKind:
			// The only `MessageKind`s we accept in URL paths and query params are
			// well-known types where the message as a whole encodes into a format
			// similar to a regular terminal field. Normal messages conveyed via URL
			// paths or query params must be non-repeated and are represented by listing
			// their terminal primitive fields, which will trigger the other cases
			// below instead of this one.
			if pval, err := parseWellKnownType(message, fieldDescriptor, value); err != nil {
				parseError = err
			} else if pval != nil {
				protoValue = *pval
			} else {
				parseError = fmt.Errorf("terminal field %q of field path %q in message %q is a message type",
					fieldName, fieldPath, messageFullName)
			}

		// reference for proto scalar types:
		// https://developers.google.com/protocol-buffers/docs/proto3#scalar

		// TODO: see https://pkg.go.dev/google.golang.org/protobuf@v1.25.0/reflect/protoreflect#Kind

		case protoreflect.StringKind:
			protoValue = protoreflect.ValueOfString(value)
		case protoreflect.BytesKind:
			protoValue = protoreflect.ValueOfBytes([]byte(value))

		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			parsedValue, err := strconv.ParseInt(value, 10, 32)
			parseError, protoValue = err, protoreflect.ValueOfInt32(int32(parsedValue))
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			parseError, protoValue = err, protoreflect.ValueOfUint32(uint32(parsedValue))

		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			parsedValue, err := strconv.ParseInt(value, 10, 64)
			parseError, protoValue = err, protoreflect.ValueOfInt64(parsedValue)
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			parsedValue, err := strconv.ParseUint(value, 10, 64)
			parseError, protoValue = err, protoreflect.ValueOfUint64(parsedValue)

		case protoreflect.FloatKind:
			parsedValue, err := strconv.ParseFloat(value, 32)
			parseError, protoValue = err, protoreflect.ValueOfFloat32(float32(parsedValue))
		case protoreflect.DoubleKind:
			parsedValue, err := strconv.ParseFloat(value, 64)
			parseError, protoValue = err, protoreflect.ValueOfFloat64(parsedValue)

		case protoreflect.BoolKind:
			parsedValue, err := parseBool(value)
			parseError, protoValue = err, protoreflect.ValueOfBool(parsedValue)

		case protoreflect.EnumKind:
			var parsedValue protoreflect.EnumNumber
			if numericValue, err := strconv.ParseFloat(value, 32); err == nil {
				parsedValue = protoreflect.EnumNumber(numericValue)
			} else {
				enum := fieldDescriptor.Enum().Values().ByName(protoreflect.Name(value))
				if enum == nil {
					return fmt.Errorf("(UnknownEnumStringError) unknown enum symbol %q for field path %q", value, fieldPath)
				}
				parsedValue = enum.Number()
			}
			parseError, protoValue = nil, protoreflect.ValueOfEnum(parsedValue)

		default:
			// TODO: Handle lists
			// TODO: Handle oneofs
			return fmt.Errorf("terminal field %q of field path %q is of type %q, which is not handled yet", fieldName, fieldPath, kind)
		}
		if parseError != nil {
			return fmt.Errorf("terminal field %q of field path %q is of type %q with value string %q, which could not be parsed: %s",
				fieldName, fieldPath, kind, value, parseError)
		}
		message.Set(fieldDescriptor, protoValue)
	}

	return nil
}

// parseBool parses a proper JSON representation of a bool value (either of the strings "true" or
// "false") into a bool data type. Other values cause an error. These are stricter parsing semantics
// than those of strconv.ParseBool and adhere to the JSON standard.
func parseBool(asString string) (bool, error) {
	switch asString {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("could not parse %q as a bool", asString)
	}
}

// findFieldByProtoName queries `fields` for the unique field whose lower-camel-cased name is given
// by `fieldName` and returns the corresponding FieldDescriptor or nil if none can be found. This
// returns an error if more than one entry in `fields` has the same lower-camel-cased field
// representations (as would occur, for example, for fields with proto names `foo5_bar` and
// `foo_5_bar`).
func findProtoFieldByJSONName(fieldName string, fields protoreflect.FieldDescriptors) (result protoreflect.FieldDescriptor, err error) {
	if ToDottedLowerCamel(fieldName) != fieldName {
		return nil, fmt.Errorf("(QueryParameterNameIncorrectlyCasedError)field name %q is not lower-camel-cased", fieldName)
	}
	numFields := fields.Len()
	for idx := 0; idx < numFields; idx++ {
		candidate := fields.Get(idx)
		if candidate.JSONName() == fieldName {
			if result != nil {
				return nil, fmt.Errorf("internal error: more than one field name has JSON lower-camel-case representation %q", fieldName)
			}
			result = candidate
		}
	}
	return result, nil
}

func parseWellKnownType(message protoreflect.Message, fieldDescriptor protoreflect.FieldDescriptor, value string) (*protoreflect.Value, error) {
	messageFieldTypes, ok := message.Type().(protoreflect.MessageFieldTypes)
	if !ok {
		return nil, fmt.Errorf("expected message to implement protoreflect.MessageFieldTypes but it did not")
	}
	fieldMsg := messageFieldTypes.Message(fieldDescriptor.Index())
	fullName := string(fieldMsg.Descriptor().FullName())
	stringEncoded, isWellKnown := wellKnownTypes[fullName]
	if !isWellKnown {
		return nil, nil
	}

	quoted := value[0] == '"' && value[len(value)-1] == '"'
	if stringEncoded && !quoted {
		value = fmt.Sprintf("%q", value)
	}

	msgValue := fieldMsg.New()
	err := protojson.Unmarshal([]byte(value), msgValue.Interface())
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal field %q in message %q: %w",
			string(fieldDescriptor.Name()), string(message.Descriptor().FullName()), err)
	}

	v := protoreflect.ValueOfMessage(msgValue)
	return &v, nil
}
