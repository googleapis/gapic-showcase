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
	"net/http"
	"net/url"
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
	suite, server, err := complianceSuiteTestSetup()
	if err != nil {
		t.Fatal(err)
	}
	server.Start()
	defer server.Close()

	restRPCs := map[string][]prepRepeatDataTestFunc{
		"Compliance.RepeatDataBodyInfo":   {prepRepeatDataBodyInfoNegativeTestRepeatedFields},
		"Compliance.RepeatDataQuery":      {prepRepeatDataQueryNegativeTestNumericEnums, prepRepeatDataQueryNegativeTestNumericOptionalEnums},
		"Compliance.RepeatDataSimplePath": {prepRepeatDataSimplePathNegativeTestEnum},
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

func prepRepeatDataQueryNegativeTestNumericEnums(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataQuery"
	info := request.GetInfo()
	badQueryParam := fmt.Sprintf("f_kingdom=%d", info.GetFKingdom()) // purposefully use a number, which should cause an error

	// We clear the field so we don't set the same query param correctly below. This change
	// modifies the request, but since these tests only check that calls fail, we never need to
	// refer back to the request proto after constructing the REST query.
	info.FKingdom = pb.ComplianceData_LIFE_KINGDOM_UNSPECIFIED
	queryParams := append(prepRepeatDataTestsQueryParams(request, nil), badQueryParam)

	queryString := prepQueryString(queryParams)
	return name, "GET", "/v1beta1/repeat:query" + queryString, body, err
}

func prepRepeatDataQueryNegativeTestNumericOptionalEnums(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataQuery"
	info := request.GetInfo()
	badQueryParam := fmt.Sprintf("p_kingdom=%d", info.GetPKingdom()) // purposefully use a number, which should cause an error

	// We clear the field so we don't set the same query param correctly below. This change
	// modifies the request, but since these tests only check that calls fail, we never need to
	// refer back to the request proto after constructing the REST query.
	info.PKingdom = nil
	queryParams := append(prepRepeatDataTestsQueryParams(request, nil), badQueryParam)

	queryString := prepQueryString(queryParams)
	return name, "GET", "/v1beta1/repeat:query" + queryString, body, err
}

func prepRepeatDataSimplePathNegativeTestEnum(request *genproto.RepeatRequest) (verb string, name string, path string, body string, err error) {
	name = "Compliance.RepeatDataSimplePath"
	info := request.GetInfo()

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
		{"f_kingdom", "%d", info.GetFKingdom()}, // purposefully use a number, which should cause an error
	} {
		pathParts = append(pathParts, url.PathEscape(fmt.Sprintf(part.format, part.value)))
		nonQueryParamNames["info."+part.name] = true
	}
	path = fmt.Sprintf("/v1beta1/repeat/%s:simplepath", strings.Join(pathParts, "/"))

	queryString := prepRepeatDataTestsQueryString(request, nonQueryParamNames)
	return name, "GET", path + queryString, body, err
}
