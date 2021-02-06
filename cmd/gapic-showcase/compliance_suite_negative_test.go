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
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/googleapis/gapic-showcase/server/genproto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestComplianceSuiteErrors checks for non-spec-compliant HTTP requests. Not all of these
// conditions necessarily generate a server error in a real service, but the behavior is often
// ill-defined. We want Showcase to require the generators be strict in the transcoding format they
// use.
func TestComplianceSuiteErrors(t *testing.T) {
	// return
	// Run the Showcase REST server locally.
	server := httptest.NewUnstartedServer(nil)
	backend := createBackends()
	restServer := newEndpointREST(nil, backend)
	server.Config = restServer.server
	server.Start()
	defer server.Close()

	// Locate, load, and unmarshal the compliance suite.
	_, thisFile, _, _ := runtime.Caller(0)
	suiteFile := filepath.Join(filepath.Dir(thisFile), "../../schema/google/showcase/v1beta1/compliance_suite.json")
	jsonProto, err := ioutil.ReadFile(suiteFile)
	if err != nil {
		t.Fatalf("could not open suite file %q", suiteFile)
	}
	var suite pb.ComplianceSuite

	if err := protojson.Unmarshal(jsonProto, &suite); err != nil {
		t.Fatalf("error unmarshalling from json %s:\n   file: %s\n   input was: %s", err, suiteFile, jsonProto)
	}

	restRPCs := map[string][]prepRepeatDataTestFunc{
		"Compliance.RepeatDataBodyInfo": {prepRepeatDataBodyInfoNegativeTestRepeatedFields},
	}

	for _, group := range suite.GetGroup() {
		rpcsToTest := group.GetRpcs()
		for requestIdx, requestProto := range group.GetRequests() {
			for rpcIdx, rpcName := range rpcsToTest {
				errorPrefix := fmt.Sprintf("[request %d/%q: rpc %q/%d/%q]",
					requestIdx, requestProto.GetName(), group.Name, rpcIdx, rpcName)

				// Ensure that we issue only the RPCs the test suite is expecting.
				restTest, ok := restRPCs[rpcName]
				if !ok {
					// we don't have a negative test for this RPC
					continue
				}

				for _, rpcPrep := range restTest {

					prepName, verb, path, requestBody, err := rpcPrep(requestProto)
					if err != nil {
						t.Errorf("%s error: %s", errorPrefix, err)
					}
					if got, want := prepName, rpcName; got != want {
						t.Errorf("%s retrieved mismatched prep function: got %q, want %q", errorPrefix, got, want)
					}

					// Issue the request
					httpRequest, err := http.NewRequest(verb, server.URL+path, strings.NewReader(requestBody))
					if err != nil {
						t.Errorf("%s error creating request: %s", errorPrefix, err)
						continue
					}
					httpResponse, err := http.DefaultClient.Do(httpRequest)
					if err != nil {
						t.Errorf("%s error issuing call: %s", errorPrefix, err)
						continue
					}

					// Check for unsuccessful response.
					if got, notWant := httpResponse.StatusCode, http.StatusOK; got == notWant {
						t.Errorf("%s response code: got %d, notWant %d\n   %s %s\nbody: %s\n----------------------------------------\n",
							errorPrefix, got, notWant, verb, server.URL+path, requestBody)
					}
				}
			}
		}
	}
}

func prepRepeatDataBodyInfoNegativeTestRepeatedFields(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataBodyInfo"
	bodyBytes, err := protojson.Marshal(request.Info)
	queryString := prepRepeatDataTestsQueryString(request, nil) // purposefully repeats query params, which should cause an error
	_ = bodyBytes
	return name, "POST", "/v1beta1/repeat:bodyinfo" + queryString, string(bodyBytes), err
}
