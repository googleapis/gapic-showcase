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

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestRESTCalls tests that arbitrary rest calls received by the Showcase REST server are handled
// correctly.
func TestRESTCalls(t *testing.T) {
	server := httptest.NewUnstartedServer(nil)
	backend := createBackends()
	restServer := newEndpointREST(nil, backend)

	server.Config = restServer.server
	server.Start()
	defer server.Close()

	for idx, testCase := range []struct {
		verb       string
		path       string
		body       string
		fullJSON   bool
		want       string
		statusCode int // 0 taken to mean 200 for simplicity
	}{
		{
			verb: "GET",
			path: "/hello",
			want: "GAPIC Showcase: HTTP/REST endpoint using gorilla/mux\n",
		},
		{
			verb: "POST",
			path: "/v1beta1/repeat:body",
			body: `{"info":{"fString":"jonas^ mila"}}`,
			want: `{"info":{"fString":"jonas^ mila"}}`,
		},
		{
			verb: "GET",
			path: "/v1beta1/repeat:query?info.f_string=jonas+mila",
			want: `{"info":{"fString":"jonas mila"}}`,
		},
		{
			verb: "GET",
			path: "/v1beta1/repeat:query?info.f_string=jonas^mila",

			// TODO: Fix so that this returns an error, because `^` is not URL-escaped
			statusCode: 200,
			want:       `{"info":{"fString":"jonas^mila"}}`,
		},
		{
			verb:       "GET",
			path:       "/v1beta1/repeat:query?info.f_string=jonas mila",
			statusCode: 400, // unescaped space in query param
		},

		{
			// Test responses:
			//   1. unset optional field absent
			//   2. zero-set optional field present
			//   3. unset non-optional field present
			//   4. enum field is symbolic rather than numeric
			verb:     "POST",
			path:     "/v1beta1/repeat:body",
			body:     `{"info":{"fString":"jonas^ mila", "p_double": 0}}`,
			fullJSON: true,
			want: `{
                          "info": {
                            "fString": "jonas^ mila",
                            "fInt32": 0,
                            "fSint32": 0,
                            "fSfixed32": 0,
                            "fUint32": 0,
                            "fFixed32": 0,
                            "fInt64": "0",
                            "fSint64": "0",
                            "fSfixed64": "0",
                            "fUint64": "0",
                            "fFixed64": "0",
                            "fDouble": 0,
                            "fFloat": 0,
                            "fBool": false,
                            "fBytes": "",
                            "fKingdom": "LIFE_KINGDOM_UNSPECIFIED",
                            "fChild": null,
                            "pDouble": 0
                          }
                        }`,
		},
	} {

		var jsonOptions *resttools.JSONMarshalOptions
		if testCase.fullJSON {
			jsonOptions = allowFullJSON()
		} else {
			jsonOptions = allowCompactJSON()
		}

		request, err := http.NewRequest(testCase.verb, server.URL+testCase.path, strings.NewReader(testCase.body))
		if err != nil {
			jsonOptions.Restore()
			t.Fatal(err)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			jsonOptions.Restore()
			t.Fatal(err)
		}

		want := testCase.statusCode
		if want == 0 {
			want = 200
		}
		if got := response.StatusCode; got != want {
			t.Errorf("testcase %2d: status code: got %d, want %d", idx, got, want)
			t.Errorf("  request: %v", request)
		} else if want != 200 {
			// we got the expected error
			jsonOptions.Restore()
			continue
		}

		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			jsonOptions.Restore()
			log.Fatal(err)
		}
		if got, want := string(body), testCase.want; noSpace(got) != noSpace(want) {
			t.Errorf("testcase %2d: body: got `%s`, want %s", idx, got, want)
			t.Errorf("  request: %v", request)
		}
		jsonOptions.Restore()
	}

}

// allowFullJSON ensures that resttools JSONMarshaler uses the compact representation until
// explicitly restored; this makes some tests shorter to configure and easier to understand.
func allowCompactJSON() *resttools.JSONMarshalOptions {
	resttools.JSONMarshaler.Replace(&protojson.MarshalOptions{
		Multiline:       false,
		AllowPartial:    false,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	})
	return &resttools.JSONMarshaler
}

// allowFullJSON ensures that resttools JSONMarshaler uses the production configuration until
// explicitly restored.
func allowFullJSON() *resttools.JSONMarshalOptions {
	resttools.JSONMarshaler.Replace(nil)
	return &resttools.JSONMarshaler
}

// noSpace removes whitespace from src. This is useful for processing formatted responses or
// expected values without having to worry about whitespace matches.
func noSpace(src string) string {
	str := strings.ReplaceAll(src, "\n", "")
	return strings.ReplaceAll(str, " ", "")
}
