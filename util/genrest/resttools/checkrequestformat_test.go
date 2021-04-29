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
	"net/http"
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

func TestCheckRequestFormat(t *testing.T) {
	for idx, testCase := range []struct {
		label     string
		json      string
		header    http.Header // nil means standard headers
		wantError bool
	}{
		{
			label: "normal case",
			json:  `{"fString": "hi", "fKingdom": "FUNGI"}`,
		},
		{
			label: "normal case, optional field",
			json:  `{"fString": "hi", "pKingdom": "FUNGI"}`,
		},
		{
			label: "normal case, nested message, optional field",
			json:  `{"fString": "hi", "fChild": {"pContinent": "AFRICA"}}`,
		},
		{
			label: "normal case, nested optional message",
			json:  `{"fString": "hi", "pChild": {"fContinent": "AFRICA"}}`,
		},
		{
			label: "normal case, no enum",
			json:  `{"fString": "hi"}`,
		},
		{
			label: "random string does not fail",
			json:  `{"fString": "hi", "fKingdom": "MONACO"}`,
		},
		{
			label:     "numeric enum",
			json:      `{"fString": "hi", "fKingdom": 56}`,
			wantError: true,
		},
		{
			label:     "numeric optional enum",
			json:      `{"fString": "hi", "pKingdom": 57}`,
			wantError: true,
		},
		{
			label:     "stringy numeric enum",
			json:      `{"fString": "hi", "fKingdom": "56"}`,
			wantError: true,
		},
		{
			label:     "stringy optional numeric enum",
			json:      `{"fString": "hi", "fKingdom": "57"}`,
			wantError: true,
		},
		{
			label:     "stringy number+letter enum",
			json:      `{"fString": "hi", "fKingdom": "56Abacus"}`,
			wantError: false,
		},
	} {
		complianceData := &genprotopb.ComplianceData{}
		request, err := http.NewRequest("POST", "showcase.foo.com", strings.NewReader(testCase.json))
		if err != nil {
			t.Fatalf("test case %d[%q] could not create request: %s", idx, testCase.label, err)

		}
		request.Header = testCase.header
		if request.Header == nil {
			PopulateRequestHeaders(request)
		}
		if err := CheckRequestFormat(request.Body, request.Header, complianceData.ProtoReflect()); (err != nil) != testCase.wantError {
			t.Errorf("test case %d[%q] text enum encoding: expected error==%v, got: %v", idx, testCase.label, testCase.wantError, err)
		}
	}
}
