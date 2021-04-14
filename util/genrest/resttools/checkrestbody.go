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

	"google.golang.org/protobuf/reflect/protoreflect"
)

// CheckRESTBody verifies that any enum fields in message are properly represented in the JSON
// payload carried by jsonReader: the fields must be either absent or have string values.
func CheckRESTBody(jsonReader io.Reader, message protoreflect.Message) error {
	jsonBytes, err := ioutil.ReadAll(jsonReader)
	if err != nil {
		return err
	}
	enumFields := GetEnumFields(message)
	return CheckJSONEnumFields(jsonBytes, enumFields)
}

// CheckJSONEnumFields verifies that each of the fields listed in fieldsToCheck, presumably all
// referring to enum fields, are encoded correctly in jsonBytes, meaning that the field is absent or
// its value is a string. Each element of fieldsToCheck is a qualified proto field name represented
// as a sequence of simple protoreflect.Name.
func CheckJSONEnumFields(jsonBytes []byte, fieldsToCheck [][]protoreflect.Name) error {
	// See See eg https://michaelheap.com/golang-encodedecode-arbitrary-json/

	var payload map[string]interface{}
	json.Unmarshal(jsonBytes, &payload)
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
// means that if fieldPath is a path to en enum field, the boolean will be false if the enum is
// encoded in payloadusing a non-string representation.
func CheckEnum(payload map[string]interface{}, fieldPath []protoreflect.Name) (fieldName string, ok bool) {
	nameParts := []string{}
	last := len(fieldPath) - 1
	var value string
	var found, isString bool
	for idx, pathSegment := range fieldPath {
		segment := string(pathSegment)
		nameParts = append(nameParts, segment)

		// TODO: For repeated fields, will need to recurse. Consider denoting repeated fields by appending a "*" to the pathSegment

		if idx < last {
			payload, found = payload[segment].(map[string]interface{})
			if !found {
				// Some elements of the field path are not populated
				break
			}
			continue
		}

		if _, found = payload[segment]; found {
			value, isString = payload[segment].(string)
		}
	}

	fieldName = strings.Join(nameParts, ".")
	if !found {
		return fieldName, true
	}

	if !isString {
		return fieldName, false
	}

	if _, err := strconv.Atoi(value); err == nil {
		// A string representation of an enum value should not be parseable as an int
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

// protoEnumFields is a list of fields paths (themselves represented as a list of nested field
// names) which represent enum fields. This is used to memoize calls to GetEnumFields.
var protoEnumFields map[protoreflect.FullName][][]protoreflect.Name

func init() {
	protoEnumFields = map[protoreflect.FullName][][]protoreflect.Name{}
}
