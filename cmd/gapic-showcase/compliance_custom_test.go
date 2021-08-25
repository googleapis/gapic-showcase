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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/googleapis/gapic-showcase/server/genproto"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestRepeatWithUnknownEnum
func TestRepeatWithUnknownEnum(t *testing.T) {
	_, server, err := complianceSuiteTestSetup()
	if err != nil {
		t.Fatal(err)
	}
	server.Start()
	defer server.Close()

	resttools.JSONMarshaler.Replace(nil)
	defer resttools.JSONMarshaler.Restore()

	request := &genprotopb.RepeatRequest{
		FInt32: 23,
	}

	for idx, variant := range []string{"invalidenum", "invalidoptionalenum"} {
		errorPrefix := fmt.Sprintf("[%d %q]", idx, variant)

		// First ensure the request would be otherwise successful
		responseBody, requestBody := getJSONResponse(t, request, server.URL+"/v1beta1/repeat:body", errorPrefix)
		var response genproto.RepeatResponse
		if err := protojson.Unmarshal(responseBody, &response); err != nil {
			t.Fatalf("%s could not unmarshal valid response body: %s\n   response body: %s\n   request: %s\n",
				errorPrefix, err, string(responseBody), string(requestBody))
		}

		// Then ensure the expected error occurs
		responseBody, requestBody = getJSONResponse(t, request,
			fmt.Sprintf("%s/v1beta1/repeat:%s", server.URL, variant), errorPrefix)
		err = protojson.Unmarshal(responseBody, &response)
		if err == nil {
			t.Fatalf("%s did not receive an error:\n   response body: %s\n   request: %s\n",
				errorPrefix, string(responseBody), string(requestBody))
		}
		if !strings.Contains(err.Error(), "invalid value for enum type") {
			t.Fatalf("%s received different error than expected: %s\n   response body: %s\n   request: %s\n",
				errorPrefix, err, string(responseBody), string(requestBody))
		}
	}
}

func getJSONResponse(t *testing.T, request *genprotopb.RepeatRequest, uri, errorPrefix string) (responseBody, requestBody []byte) {
	verb := "POST"
	requestBody, err := resttools.ToJSON().Marshal(request)
	if err != nil {
		t.Fatalf("%s error encoding request: %s", errorPrefix, err)
	}

	httpRequest, err := http.NewRequest(verb, uri, strings.NewReader(string(requestBody)))
	if err != nil {
		t.Fatalf("%s error creating request: %s", errorPrefix, err)
	}
	resttools.PopulateRequestHeaders(httpRequest)

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		t.Fatalf("%s error issuing call: %s", errorPrefix, err)
	}

	// Check for successful response.
	if got, want := httpResponse.StatusCode, http.StatusOK; got != want {
		t.Errorf("%s response code: got %d, want %d\n   %s %s\n\n",
			errorPrefix, got, want, verb, uri)
	}

	responseBody, err = ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		t.Fatalf("%s could not read httpResponse body: %s", errorPrefix, err)
	}
	return responseBody, requestBody
}

// repeat wtih same message for invalidenum; check for error
