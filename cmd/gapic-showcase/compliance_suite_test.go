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
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-showcase/server/genproto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestComplianceSuite ensures the REST test suite that we require GAPIC generators to pass works
// correctly. GAPIC generators should generate GAPICs for the Showcase API and issue the unary calls
// defined in the test suite using the GAPIC surface. The generators' test should follow the
// high-level logic below, as described in the comments.
func TestComplianceSuite(t *testing.T) {
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

	// Set handlers for each test case. When GAPIC generator tests do this, they should have
	// each of their handlers invoking the correct GAPIC library method for the Showcase API.
	restRPCs := map[string]prepRepeatDataTestFunc{
		"Compliance.RepeatDataBody":       prepRepeatDataBodyTest,
		"Compliance.RepeatDataQuery":      prepRepeatDataQueryTest,
		"Compliance.RepeatDataSimplePath": prepRepeatDataSimplePathTest,
	}

	for _, group := range suite.GetGroup() {
		rpcsToTest := group.GetRpcs()
		for requestIdx, requestProto := range group.GetRequests() {
			for rpcIdx, rpcName := range rpcsToTest {
				errorPrefix := fmt.Sprintf("[request %d/%q: rpc %q/%d/%q]",
					requestIdx, requestProto.GetName(), group.Name, rpcIdx, rpcName)

				// Ensure that we issue only the RPCs the test suite is expecting.
				rpcPrep, ok := restRPCs[rpcName]
				if !ok {
					t.Errorf("%s could not find prep function for this RPC", errorPrefix)
					continue
				}

				prepName, verb, path, requestBody, err := rpcPrep(requestProto)
				if err != nil {
					t.Errorf("%s error: %s", errorPrefix, err)
				}
				if got, want := prepName, rpcName; got != want {
					t.Errorf("%s retrieved mismatched prep function: got %q, want %q", errorPrefix, got, want)
				}

				// Issue the request. When GAPIC generator tests do this, they should simply
				// invoke the correct GAPIC library method for the Showcase API.
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

				// Check for successful response.
				if got, want := httpResponse.StatusCode, http.StatusOK; got != want {
					t.Errorf("%s response code: got %d, want %d\n   %s %s",
						errorPrefix, got, want, verb, server.URL+path)
				}

				// Unmarshal httpResponse body, interpreted as JSON.
				// should do this.
				responseBody, err := ioutil.ReadAll(httpResponse.Body)
				httpResponse.Body.Close()
				if err != nil {
					t.Errorf("%s could not read httpResponse body: %s", errorPrefix, err)
					continue
				}
				var response genproto.RepeatResponse
				if err := protojson.Unmarshal(responseBody, &response); err != nil {
					t.Errorf("%s could not unmarshal httpResponse body: %s\n   response body: %s",
						errorPrefix, err, string(responseBody))
					continue
				}

				// Check for expected response.
				if diff := cmp.Diff(response.GetInfo(), requestProto.GetInfo(), cmp.Comparer(proto.Equal)); diff != "" {
					t.Errorf("%s unexpected response: got=-, want=+:%s", errorPrefix, diff)
				}
			}
		}
	}
}

// The following are helpers for TestComplianceSuite, since Showcase doesn't intrinsically define a
// REST client. Each GAPIC generator should instead use the GAPIC it generated for the Showcase
// API.
type prepRepeatDataTestFunc func(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error)

func prepRepeatDataBodyTest(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataBody"
	bodyBytes, err := protojson.Marshal(request)
	return name, "POST", "/v1beta1/repeat:body", string(bodyBytes), err
}

func prepRepeatDataQueryTest(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataQuery"
	queryString := prepRepeatDataTestsQueryString(request, nil)
	return name, "GET", "/v1beta1/repeat:query" + queryString, body, err
}

func prepRepeatDataSimplePathTest(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataSimplePath"
	info := request.GetInfo()

	// TODO: Determine behavior for a string field path param whose value is empty. This should be
	// a failure, probably, in which case we need to augment the ComplianceGroup to allow
	// specifying expected errors.
	// TODO: Add to compliance_suite cases with near-maximal values.
	path = fmt.Sprintf("/v1beta1/repeat/%s/%d/%.20g/%t:simplepath",
		url.PathEscape(info.GetFString()), info.GetFInt32(), info.GetFDouble(), info.GetFBool())

	// exclude the path fields from the query params
	exclude := map[string]bool{
		"f_string": true,
		"f_int32":  true,
		"f_double": true,
		"f_bool":   true,
	}
	queryString := prepRepeatDataTestsQueryString(request, exclude)
	return name, "GET", path + queryString, body, err
}

// prepRepeatDataTestsQueryString returns the query string containing all fields in `request.info`
// except for those whose proto name (relative to request.info) are present in the `exclude` map
// with a value of `true`.
func prepRepeatDataTestsQueryString(request *genproto.RepeatRequest, exclude map[string]bool) string {
	info := request.GetInfo()
	queryParams := []string{}
	addParam := func(condition bool, key, value string) {
		if !condition {
			return
		}
		queryParams = append(queryParams, fmt.Sprintf("info.%s=%s", key, value))
	}

	if !exclude["f_string"] {
		addParam(len(info.GetFString()) > 0, "f_string", strings.ReplaceAll(strings.ReplaceAll(info.GetFString(), " ", "+"), "%", "%%"))
	}
	// TODO: Add additional data fields

	var queryString string
	if len(queryParams) > 0 {
		queryString = fmt.Sprintf("?%s", strings.Join(queryParams, "&"))
	}

	return queryString
}
