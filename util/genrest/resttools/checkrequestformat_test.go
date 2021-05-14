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
		wantError string
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
			label: "normal case, correct headers",
			json:  `{"fString": "hi"}`,
			header: http.Header{
				headerNameContentType: []string{headerValueContentTypeJSON},
				headerNameAPIClient:   []string{"foo/1 rest/foo gapic/2 blah"},
			},
		},
		{
			label:     "no api client header",
			json:      `{"fString": "hi"}`,
			header:    http.Header{headerNameContentType: []string{headerValueContentTypeJSON}},
			wantError: "(HeaderAPIClientError)",
		},
		{
			label: "no header rest token",
			json:  `{"fString": "hi"}`,
			header: http.Header{
				headerNameContentType: []string{headerValueContentTypeJSON},
				headerNameAPIClient:   []string{"foo/1 gapic/2 blah"},
			},
			wantError: "(HeaderTransportRESTError)",
		},
		{
			label: "no header gapic token",
			json:  `{"fString": "hi"}`,
			header: http.Header{
				headerNameContentType: []string{headerValueContentTypeJSON},
				headerNameAPIClient:   []string{"foo/1 rest/foo blah"},
			},
			wantError: "(HeaderClientGAPICError)",
		},
		{
			label: "no content-type header",
			json:  `{"fString": "hi"}`,
			header: http.Header{
				headerNameAPIClient: []string{"foo/1 rest/foo gapic/2 blah"},
			},
			wantError: "(HeaderContentTypeError)",
		},
		{
			label: "bad content-type header",
			json:  `{"fString": "hi"}`,
			header: http.Header{
				headerNameContentType: []string{"something " + headerValueContentTypeJSON},
				headerNameAPIClient:   []string{"foo/1 rest/foo gapic/2 blah"},
			},
			wantError: "(HeaderContentTypeError)",
		},

		{
			label: "random string does not fail",
			json:  `{"fString": "hi", "fKingdom": "MONACO"}`,
		},
		{
			label:     "numeric enum",
			json:      `{"fString": "hi", "fKingdom": 56}`,
			wantError: "(EnumEncodingError)",
		},
		{
			label:     "numeric optional enum",
			json:      `{"fString": "hi", "pKingdom": 57}`,
			wantError: "(EnumEncodingError)",
		},
		{
			label:     "stringy numeric enum",
			json:      `{"fString": "hi", "fKingdom": "56"}`,
			wantError: "(EnumEncodingError)",
		},
		{
			label:     "stringy optional numeric enum",
			json:      `{"fString": "hi", "fKingdom": "57"}`,
			wantError: "(EnumEncodingError)",
		},
		{
			label: "stringy number+letter enum",
			json:  `{"fString": "hi", "fKingdom": "56Abacus"}`,
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
		err = CheckRequestFormat(request.Body, request.Header, complianceData.ProtoReflect())
		if (err != nil) != (len(testCase.wantError) > 0) {
			t.Errorf("test case %d[%q] CheckRequestFormat(): expected error==%v, got: %v", idx, testCase.label, testCase.wantError, err)
			continue
		}
		if err == nil {
			continue
		}
		if got, want := err.Error(), testCase.wantError; !strings.Contains(got, want) {
			t.Errorf("test case %d[%q] CheckRequestFormat(): incorrect error: want: %q, got %q", idx, testCase.label, want, got)
		}
	}
}
