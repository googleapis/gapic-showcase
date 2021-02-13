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
	suite, server, err := complianceSuiteTestSetup()
	if err != nil {
		t.Fatal(err)
	}
	server.Start()
	defer server.Close()

	// Set handlers for each test case. When GAPIC generator tests do this, they should have
	// each of their handlers invoking the correct GAPIC library method for the Showcase API.
	restRPCs := map[string]prepRepeatDataTestFunc{
		"Compliance.RepeatDataBody":                 prepRepeatDataBodyTest,
		"Compliance.RepeatDataBodyInfo":             prepRepeatDataBodyInfoTest,
		"Compliance.RepeatDataQuery":                prepRepeatDataQueryTest,
		"Compliance.RepeatDataSimplePath":           prepRepeatDataSimplePathTest,
		"Compliance.RepeatDataPathResource":         prepRepeatDataPathResourceTest,
		"Compliance.RepeatDataPathTrailingResource": prepRepeatDataPathTrailingResourceTest,
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
					t.Errorf("%s response code: got %d, want %d\n   %s %s\n\n",
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
					t.Errorf("%s could not unmarshal httpResponse body: %s\n   response body: %s\n   request: %s\n",
						errorPrefix, err, string(responseBody), requestBody)
					continue
				}

				// Check for expected response.
				if diff := cmp.Diff(response.GetInfo(), requestProto.GetInfo(), cmp.Comparer(proto.Equal)); diff != "" {
					t.Errorf("%s unexpected response: got=-, want=+:%s\n   %s %s\n------------------------------\n",
						errorPrefix, diff, verb, server.URL+path)
				}
			}
		}
	}
}

func complianceSuiteTestSetup() (suite *pb.ComplianceSuite, server *httptest.Server, err error) {
	// Run the Showcase REST server locally.
	server = httptest.NewUnstartedServer(nil)
	backend := createBackends()
	restServer := newEndpointREST(nil, backend)
	server.Config = restServer.server

	// Locate, load, and unmarshal the compliance suite.
	_, thisFile, _, _ := runtime.Caller(0)
	suiteFile := filepath.Join(filepath.Dir(thisFile), "../../schema/google/showcase/v1beta1/compliance_suite.json")
	jsonProto, err := ioutil.ReadFile(suiteFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open suite file %q", suiteFile)
	}

	suite = &pb.ComplianceSuite{}
	if err := protojson.Unmarshal(jsonProto, suite); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling from json %s:\n   file: %s\n   input was: %s", err, suiteFile, jsonProto)
	}

	return suite, server, nil
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

func prepRepeatDataBodyInfoTest(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataBodyInfo"
	bodyBytes, err := protojson.Marshal(request.Info)
	queryString := prepRepeatDataTestsQueryString(request, map[string]bool{"info": true})
	_ = bodyBytes
	return name, "POST", "/v1beta1/repeat:bodyinfo" + queryString, string(bodyBytes), err
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

	pathParts := []string{}
	nonQueryParamNames := map[string]bool{}

	for _, part := range []struct {
		name   string
		format string
		value  interface{}
	}{
		{"f_string", "%s", info.GetFString()},
		{"f_int32", "%d", info.GetFInt32()},
		{"f_double", "%g", info.GetFDouble()},
		{"f_bool", "%t", info.GetFBool()},
		{"f_kingdom", "%s", info.GetFKingdom()},
	} {
		pathParts = append(pathParts, url.PathEscape(fmt.Sprintf(part.format, part.value)))
		nonQueryParamNames["info."+part.name] = true
	}
	path = fmt.Sprintf("/v1beta1/repeat/%s:simplepath", strings.Join(pathParts, "/"))

	queryString := prepRepeatDataTestsQueryString(request, nonQueryParamNames)
	return name, "GET", path + queryString, body, err
}

func prepRepeatDataPathResourceTest(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataPathResource"
	info := request.GetInfo()

	pathParts := []string{}
	nonQueryParamNames := map[string]bool{}

	for _, part := range []struct {
		name           string
		format         string
		value          interface{}
		requiredPrefix string
	}{
		{"f_string", "%s", info.GetFString(), "first/"},
		{"f_child.f_string", "%s", info.GetFChild().GetFString(), "second/"},
		{"f_bool", "bool/%t", info.GetFBool(), ""},
	} {
		if len(part.requiredPrefix) > 0 && !strings.HasPrefix(part.value.(string), part.requiredPrefix) {
			err = fmt.Errorf("expected value of %q to begin with %q; got %q", part.name, part.requiredPrefix, part.value)
			return
		}
		pathParts = append(pathParts, url.PathEscape(fmt.Sprintf(part.format, part.value)))
		nonQueryParamNames["info."+part.name] = true
	}
	path = fmt.Sprintf("/v1beta1/repeat/%s:pathresource", strings.Join(pathParts, "/"))

	queryString := prepRepeatDataTestsQueryString(request, nonQueryParamNames)
	return name, "GET", path + queryString, body, err
}

func prepRepeatDataPathTrailingResourceTest(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataPathTrailingResource"
	info := request.GetInfo()

	pathParts := []string{}
	nonQueryParamNames := map[string]bool{}

	for _, part := range []struct {
		name           string
		format         string
		value          interface{}
		requiredPrefix string
	}{
		{"f_string", "%s", info.GetFString(), "first/"},
		{"f_child.f_string", "%s", info.GetFChild().GetFString(), "second/"},
	} {
		if len(part.requiredPrefix) > 0 && !strings.HasPrefix(part.value.(string), part.requiredPrefix) {
			err = fmt.Errorf("expected value of %q to begin with %q; got %q", part.name, part.requiredPrefix, part.value)
			return
		}
		pathParts = append(pathParts, url.PathEscape(fmt.Sprintf(part.format, part.value)))
		nonQueryParamNames["info."+part.name] = true
	}
	path = fmt.Sprintf("/v1beta1/repeat/%s:pathtrailingresource", strings.Join(pathParts, "/"))

	queryString := prepRepeatDataTestsQueryString(request, nonQueryParamNames)
	return name, "GET", path + queryString, body, err
}

// prepRepeatDataTestsQueryString returns the query string containing all fields in `request.info`
// except for those whose proto name (relative to request.info) are present in the `exclude` map
// with a value of `true`.
func prepRepeatDataTestsQueryString(request *genproto.RepeatRequest, exclude map[string]bool) string {
	return prepQueryString(prepRepeatDataTestsQueryParams(request, exclude))
}
func prepRepeatDataTestsQueryParams(request *genproto.RepeatRequest, exclude map[string]bool) []string {
	info := request.GetInfo()
	queryParams := []string{}
	addParam := func(key string, condition bool, value string) {
		if exclude["info"] || exclude["info."+key] || !condition {
			return
		}
		queryParams = append(queryParams, fmt.Sprintf("info.%s=%s", key, value))
	}

	addParam("f_string", len(info.GetFString()) > 0, url.QueryEscape(info.GetFString()))
	addParam("f_int32", info.GetFInt32() != 0, fmt.Sprintf("%d", info.GetFInt32()))
	addParam("f_sint32", info.GetFSint32() != 0, fmt.Sprintf("%d", info.GetFSint32()))
	addParam("f_sfixed32", info.GetFSfixed32() != 0, fmt.Sprintf("%d", info.GetFSfixed32()))
	addParam("f_uint32", info.GetFUint32() != 0, fmt.Sprintf("%d", info.GetFUint32()))
	addParam("f_fixed32", info.GetFFixed32() != 0, fmt.Sprintf("%d", info.GetFFixed32()))
	addParam("f_int64", info.GetFInt64() != 0, fmt.Sprintf("%d", info.GetFInt64()))
	addParam("f_sint64", info.GetFSint64() != 0, fmt.Sprintf("%d", info.GetFSint64()))
	addParam("f_sfixed64", info.GetFSfixed64() != 0, fmt.Sprintf("%d", info.GetFSfixed64()))
	addParam("f_uint64", info.GetFUint64() != 0, fmt.Sprintf("%d", info.GetFUint64()))
	addParam("f_fixed64", info.GetFFixed64() != 0, fmt.Sprintf("%d", info.GetFFixed64()))

	addParam("f_double", info.GetFDouble() != 0, url.QueryEscape(fmt.Sprintf("%g", info.GetFDouble())))
	addParam("f_float", info.GetFFloat() != 0, url.QueryEscape(fmt.Sprintf("%g", info.GetFFloat())))
	addParam("f_bool", info.GetFBool(), "true")
	addParam("f_bytes", len(info.GetFBytes()) > 0, url.QueryEscape(string(info.GetFBytes()))) // TODO: Check this is correct, given runes in strings
	addParam("f_kingdom", info.GetFKingdom() != pb.ComplianceData_UNASSIGNED, info.GetFKingdom().String())

	addParam("p_string", info.PString != nil, url.QueryEscape(info.GetPString()))
	addParam("p_int32", info.PInt32 != nil, fmt.Sprintf("%d", info.GetPInt32()))
	addParam("p_double", info.PDouble != nil, url.QueryEscape(fmt.Sprintf("%g", info.GetPDouble())))
	addParam("p_bool", info.PBool != nil, fmt.Sprintf("%t", info.GetPBool()))
	addParam("p_kingdom", info.PKingdom != nil, info.GetPKingdom().String())

	addParam("f_child.f_string", len(info.GetFChild().GetFString()) > 0, url.QueryEscape(info.GetFChild().GetFString()))
	addParam("f_child.f_float", info.GetFChild().GetFFloat() != 0, url.QueryEscape(fmt.Sprintf("%g", info.GetFChild().GetFFloat())))
	addParam("f_child.f_double", info.GetFChild().GetFDouble() != 0, url.QueryEscape(fmt.Sprintf("%g", info.GetFChild().GetFDouble())))
	addParam("f_child.f_bool", info.GetFChild().GetFBool(), "true")

	// If needed for test cases, we'll have to add remaining nested message fields.

	return queryParams
}

func prepQueryString(queryParams []string) string {
	var queryString string
	if len(queryParams) > 0 {
		queryString = fmt.Sprintf("?%s", strings.Join(queryParams, "&"))
	}

	return queryString
}
