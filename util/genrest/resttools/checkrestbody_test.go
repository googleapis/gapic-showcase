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
	"reflect"
	"strings"
	"testing"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestComputeEnumFields(t *testing.T) {
	complianceData := &genprotopb.ComplianceData{}
	want := [][]protoreflect.Name{
		{"f_kingdom"},
		{"f_child", "f_continent"},
		{"f_child", "p_continent"},
		{"p_kingdom"},
		{"p_child", "f_continent"},
		{"p_child", "p_continent"},
	}

	got := ComputeEnumFields(complianceData.ProtoReflect())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("enum fields: got %v, want %v", got, want)
	}
}

func TestCheckRESTBody(t *testing.T) {
	for idx, testCase := range []struct {
		label     string
		json      string
		wantError bool
	}{
		{
			label: "normal case",
			json:  `{"f_string": "hi", "f_kingdom": "FUNGI"}`,
		},
		{
			label: "normal case, optional field",
			json:  `{"f_string": "hi", "p_kingdom": "FUNGI"}`,
		},
		{
			label: "normal case, nested message, optional field",
			json:  `{"f_string": "hi", "f_child": {"p_continent": "AFRICA"}}`,
		},
		{
			label: "normal case, nested optional message",
			json:  `{"f_string": "hi", "p_child": {"f_continent": "AFRICA"}}`,
		},
		{
			label: "normal case, no enum",
			json:  `{"f_string": "hi"}`,
		},
		{
			label: "random string does not fail",
			json:  `{"f_string": "hi", "f_kingdom": "MONACO"}`,
		},
		{
			label:     "numeric enum",
			json:      `{"f_string": "hi", "f_kingdom": "56"}`,
			wantError: true,
		},
		{
			label:     "number +letter enum",
			json:      `{"f_string": "hi", "f_kingdom": "56Abacus"}`,
			wantError: false,
		},
	} {
		complianceData := &genprotopb.ComplianceData{}
		if err := CheckRESTBody(strings.NewReader(testCase.json), complianceData.ProtoReflect()); (err != nil) != testCase.wantError {
			t.Errorf("test case %d[%q] text enum encoding: expected error==%v, got %v", idx, testCase.label, testCase.wantError, err)
		}
	}
}
