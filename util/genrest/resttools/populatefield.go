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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// PopulateFields sets the fields within protoMessage to the values provided in fieldValues. The
// fields and values are provided as a map of field paths to the string representation of their
// values. The fields paths can refer to fields nested arbitrarily deep within protoMessage. This
// returns an error if any field path is not valid or if any value can't be parsed into the correct
// data type for the field.
func PopulateFields(protoMessage proto.Message, fieldValues map[string]string) error {
	for name, value := range fieldValues {
		if err := PopulateOneField(protoMessage, name, value); err != nil {
			// TODO: accumulate errors so we report them all at once
			return err
		}
	}
	return nil
}

// PopulateOneField finds in protoMessage the field identified by fieldPath (which could refer to an
// arbitrarily nested field using dotted notation) and sets it to `value`. It returns an error if
// the fieldPath does not properly reference a field, or if `value` could not be parsed into the
// data type expected for the field.
func PopulateOneField(protoMessage proto.Message, fieldPath string, value string) error {
	message := protoMessage.ProtoReflect()

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
		fieldDescriptor := subFields.ByName(protoreflect.Name(fieldName))
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
			parseError = fmt.Errorf("terminal field %q of field path %q in message %q is a message type",
				fieldName, fieldPath, messageFullName)

		// reference for proto scalar types:
		// https://developers.google.com/protocol-buffers/docs/proto3#scalar

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
			// TODO: should we be stricter in what we accept? ParseBool accepts various
			// representations of "true" and "false" (https://golang.org/pkg/strconv/#ParseBool)
			parsedValue, err := strconv.ParseBool(value)
			parseError, protoValue = err, protoreflect.ValueOfBool(parsedValue)

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
