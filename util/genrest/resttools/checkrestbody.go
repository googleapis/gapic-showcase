// Copyright 2021 Google LLC
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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// CheckRESTBody verifies that any enum fields in message are properly represented in the JSON
// payload carried by jsonReader: the fields must be either absent or have lower-camel-cased names.
func CheckRESTBody(jsonReader io.Reader, message protoreflect.Message) error {
	jsonBytes, err := ioutil.ReadAll(jsonReader)
	if err != nil {
		return err
	}

	var payload jsonPayload
	json.Unmarshal(jsonBytes, &payload)

	if err := CheckFieldNames(payload); err != nil {
		return err
	}

	enumFields := GetEnumFields(message)
	return CheckJSONEnumFields(payload, enumFields)
}

// CheckFieldNames checks that the field names in the JSON request body are properly formatted
// (lower-camel-cased).
func CheckFieldNames(payload jsonPayload) error {
	for fieldName, value := range payload {
		rune, _ := utf8.DecodeRuneInString(fieldName)
		if strings.ContainsAny(fieldName, "_- ") || !unicode.IsLower(rune) {
			return fmt.Errorf("%s field name is not lower-camel-cased; probably want be %q", fieldName, strcase.ToLowerCamel(fieldName))
		}
		if nested, ok := value.(map[string]interface{}); ok {
			if err := CheckFieldNames(nested); err != nil {
				return fmt.Errorf("%s.%s", fieldName, err)
			}
		}
	}
	return nil
}

// CheckJSONEnumFields verifies that each of the fields listed in fieldsToCheck, presumably all
// referring to enum fields, are encoded correctly in the parsed JSON payload, meaning that the
// field is absent or its value is a string. Each element of fieldsToCheck is a qualified proto
// field name represented as a sequence of simple protoreflect.Name.
func CheckJSONEnumFields(payload jsonPayload, fieldsToCheck [][]protoreflect.Name) error {
	badFields := []string{}
	for _, fieldPath := range fieldsToCheck {
		if field, ok := CheckEnum(payload, fieldPath); !ok {
			badFields = append(badFields, field)
		}
	}
	if len(badFields) > 0 {
		return fmt.Errorf("badly transcoded enum values in fields: %v", badFields)
	}
	return nil
}

// CheckEnum verifies whether the field whose qualified name is captured in the elements of
// fieldPath has a string value, if it exists, in the JSON representation captured by payload. This
// returns the qualified field name (as present as it is in payload) as a single string, and a
// boolean that is true only if either fieldPath is not present or if its value is a string. This
// means that if fieldPath is a path to an enum field, the boolean will be false if the enum is
// encoded in the payload using a non-string representation.
func CheckEnum(payload jsonPayload, fieldPath []protoreflect.Name) (fieldName string, ok bool) {
	nameParts := []string{}
	last := len(fieldPath) - 1
	var value string
	var found, isString bool
	for idx, pathSegment := range fieldPath {
		segment := strcase.ToLowerCamel(string(pathSegment))
		nameParts = append(nameParts, segment)

		// TODO: For repeated fields, will need to recurse. Consider denoting repeated fields by appending a "*" to the pathSegment

		if idx < last {
			payload, found = payload[segment].(jsonPayload)
			if !found {
				// Some elements of the field path are not populated
				break
			}
			continue
		}

		if _, found = payload[segment]; found {
			// We found the field specified by fieldPath.
			value, isString = payload[segment].(string)
		}
	}
	fieldName = strings.Join(nameParts, ".")

	if !found {
		// We did not find the field denoted by fieldPath, so there is no error.
		return fieldName, true
	}

	if !isString {
		// We found the enum field denoted by fieldPath, but its value is not a string. This
		// is an error: we require all enum values to be REST-encoded via their string
		// representations for REST transport.
		return fieldName, false
	}

	if _, err := strconv.Atoi(value); err == nil {
		// We found the enum field denoted by fieldPath, and its JSON value is of string
		// type, but the value of the string merely represents a number. This is an error: a
		// string representation of an enum value should not be parseable as an int, as it
		// must contain letters, typically forming words.
		return fieldName, false
	}

	return fieldName, true
}

// GetEnumFields returns a list of any arbitrarily nested fields in message that are enums. Each
// member of the returned list is a qualified field name, itself represented as a list of
// simple protoreflect.Name.
func GetEnumFields(message protoreflect.Message) [][]protoreflect.Name {
	messageName := message.Descriptor().FullName()
	if fields, ok := protoEnumFields[messageName]; ok {
		return fields
	}

	fields := ComputeEnumFields(message)
	protoEnumFields[messageName] = fields
	return fields
}

// ComputeEnumFields determines which fields in message or its submessages are enums, and returns a
// list of those qualified field names (each one of those being a list of simple
// protoreflect.Name).
func ComputeEnumFields(message protoreflect.Message) [][]protoreflect.Name {
	return computeEnumFields(message.Descriptor(), []protoreflect.Name{})
}

// computeEnumFields determines which fields in message or its submessages are enums, and returns a
// list of those qualified field names (each one of those being a list of simple
// protoreflect.Name). currentPath must be the fully qualified name for message.
func computeEnumFields(message protoreflect.MessageDescriptor, currentPath []protoreflect.Name) [][]protoreflect.Name {
	results := [][]protoreflect.Name{}
	allFields := message.Fields()
	for idx := 0; idx < allFields.Len(); idx++ {
		field := allFields.Get(idx)
		switch field.Kind() {
		case protoreflect.EnumKind:
			results = append(results, append(currentPath, field.Name()))
		case protoreflect.MessageKind:
			results = append(results, computeEnumFields(field.Message(), append(currentPath, field.Name()))...)
		default:
			continue

		}
	}
	return results
}

// jsonPayload is used for unmarshaling arbitrary JSON.
type jsonPayload map[string]interface{}

// protoEnumFields is a list of fields paths (themselves represented as a list of nested field
// names) which represent enum fields. This is used to memoize calls to GetEnumFields.
var protoEnumFields map[protoreflect.FullName][][]protoreflect.Name

func init() {
	protoEnumFields = map[protoreflect.FullName][][]protoreflect.Name{}
}
