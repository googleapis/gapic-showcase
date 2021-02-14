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

func CheckRestBody(jsonReader io.Reader, message protoreflect.Message) error {
	jsonBytes, err := ioutil.ReadAll(jsonReader)
	if err != nil {
		return err
	}
	enumFields := GetEnumFields(message)
	return CheckJSONEnumFields(jsonBytes, enumFields)
}

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

func CheckEnum(payload map[string]interface{}, fieldPath []protoreflect.Name) (fieldName string, ok bool) {
	nameParts := []string{}
	last := len(fieldPath) - 1
	var value string
	found := false
	for idx, pathSegment := range fieldPath {
		segment := string(pathSegment)
		nameParts = append(nameParts, segment)

		// TODO: For repeated fields, will need to recurse. Consider denoting repeated fields by appending a "*" to the pathSegment

		if idx < last {
			payload, ok = payload[segment].(map[string]interface{})
			if !ok {
				// Some elements of the field path are not populated
				break
			}
			continue
		}
		value, ok = payload[segment].(string)
		found = true
	}
	fieldName = strings.Join(nameParts, ".")
	if !found {
		return fieldName, true
	}

	if ok {
		if _, err := strconv.Atoi(value); err == nil {
			// A string representation of an enum value should not be parseable as an int
			ok = false
		}
	}
	return fieldName, ok
}

func GetEnumFields(message protoreflect.Message) [][]protoreflect.Name {
	messageName := message.Descriptor().FullName()
	if fields, ok := protoEnumFields[messageName]; ok {
		return fields
	}

	fields := ComputeEnumFields(message)
	protoEnumFields[messageName] = fields
	return fields
}

func ComputeEnumFields(message protoreflect.Message) [][]protoreflect.Name {
	return computeEnumFields(message.Descriptor(), []protoreflect.Name{})
}

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

// protoEnumFields is a list of fields paths (themselves represented as a list of nested field names) which represent enum fields
var protoEnumFields map[protoreflect.FullName][][]protoreflect.Name

func init() {
	protoEnumFields = map[protoreflect.FullName][][]protoreflect.Name{}
}
