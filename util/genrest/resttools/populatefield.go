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
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func PopulateField(message protoreflect.Message, fieldPath string, value string) (out string, err error) {

	out += fmt.Sprintf("BEFORE %q  %q %#v\n", fieldPath, string(message.Descriptor().FullName()), message)

	// TODO: check for proto3

	levels := strings.Split(fieldPath, ".")
	lastLevel := len(levels) - 1
	for idx, fieldName := range levels {
		if len(fieldName) == 0 {
			return out, fmt.Errorf("segment %d of path field %q is empty", idx, fieldPath)
		}

		subFields := message.Descriptor().Fields()
		fieldDescriptor := subFields.ByName(protoreflect.Name(fieldName))
		if fieldDescriptor == nil {
			return out, fmt.Errorf("could not find %dth field (%q) in field path %q in message %q", idx, fieldName, fieldPath, message.Descriptor().FullName())
		}

		if idx != lastLevel {
			// non-terminal field: go to the next level

			if kind := fieldDescriptor.Kind(); kind != protoreflect.MessageKind {
				return out, fmt.Errorf("%dth field (%q) in field path %q in message %q is not itself a message but a %q",
					idx, fieldName, fieldPath, message.Descriptor().FullName(), kind)
			}
			message = message.Mutable(fieldDescriptor).Message()
			continue
		}

		// terminal field
		switch kind := fieldDescriptor.Kind(); kind {
		case protoreflect.MessageKind:
			return out, fmt.Errorf("terminal field %q of field path %q in message %q is a message type", fieldName, fieldPath, message.Descriptor().FullName())
		case protoreflect.StringKind:
			message.Set(fieldDescriptor, protoreflect.ValueOfString(value))
		default:
			// TODO: Handle additional kinds
			// TODO: Handle lists
			return out, fmt.Errorf("terminal field %q of field path %q is of type %q, which is not handled yet", fieldName, fieldPath, kind)
		}
	}

	out += fmt.Sprintf("\nAFTER %q  %#v\n", fieldPath, message)
	return out, nil
}
