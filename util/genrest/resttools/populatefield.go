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

func PopulateOneField(protoMessage proto.Message, fieldPath string, value string) error {
	message := protoMessage.ProtoReflect()

	// TODO: check for proto3

	levels := strings.Split(fieldPath, ".")
	lastLevel := len(levels) - 1
	for idx, fieldName := range levels {
		if len(fieldName) == 0 {
			return fmt.Errorf("segment %d of path field %q is empty", idx, fieldPath)
		}

		subFields := message.Descriptor().Fields()
		fieldDescriptor := subFields.ByName(protoreflect.Name(fieldName))
		if fieldDescriptor == nil {
			return fmt.Errorf("could not find %dth field (%q) in field path %q in message %q", idx, fieldName, fieldPath, message.Descriptor().FullName())
		}

		if idx != lastLevel {
			// non-terminal field: go to the next level

			if kind := fieldDescriptor.Kind(); kind != protoreflect.MessageKind {
				return fmt.Errorf("%dth field (%q) in field path %q in message %q is not itself a message but a %q",
					idx, fieldName, fieldPath, message.Descriptor().FullName(), kind)
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
			parseError = fmt.Errorf("terminal field %q of field path %q in message %q is a message type", fieldName, fieldPath, message.Descriptor().FullName())

			// reference for proto scalar types: https://developers.google.com/protocol-buffers/docs/proto3#scalar

		case protoreflect.StringKind:
			protoValue = protoreflect.ValueOfString(value)

		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			v64, err := strconv.ParseInt(value, 10, 32)
			parseError, protoValue = err, protoreflect.ValueOfInt32(int32(v64))
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			v64, err := strconv.ParseUint(value, 10, 32)
			parseError, protoValue = err, protoreflect.ValueOfUint32(uint32(v64))

		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			v64, err := strconv.ParseInt(value, 10, 64)
			parseError, protoValue = err, protoreflect.ValueOfInt64(v64)
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			v64, err := strconv.ParseUint(value, 10, 64)
			parseError, protoValue = err, protoreflect.ValueOfUint64(v64)

		case protoreflect.FloatKind:
			v64, err := strconv.ParseFloat(value, 32)
			parseError, protoValue = err, protoreflect.ValueOfFloat32(float32(v64))
		case protoreflect.DoubleKind:
			v64, err := strconv.ParseFloat(value, 64)
			parseError, protoValue = err, protoreflect.ValueOfFloat64(v64)

		case protoreflect.BoolKind:
			// TODO: should we be stricter in what we accept? ParseBool accepts various
			// representations of "true" and "false" (https://golang.org/pkg/strconv/#ParseBool)
			vBool, err := strconv.ParseBool(value)
			parseError, protoValue = err, protoreflect.ValueOfBool(vBool)

		case protoreflect.BytesKind:
			protoValue = protoreflect.ValueOfBytes([]byte(value))

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

func PopulateFields(protoMessage proto.Message, fieldValues map[string]string) error {
	for name, value := range fieldValues {
		if err := PopulateOneField(protoMessage, name, value); err != nil {
			// TODO: accumulate error
			return err
		}
	}
	return nil
}
