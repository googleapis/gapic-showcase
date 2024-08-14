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
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/googleapis/gapic-showcase/server/genproto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/services"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestComplianceSuiteErrors checks for non-spec-compliant HTTP requests. Not all of these
// conditions necessarily generate a server error in a real service, but the behavior is often
// ill-defined. We want Showcase to require the generators be strict in the transcoding format they
// use.
func TestComplianceSuiteErrors(t *testing.T) {
	masterSuite, server, err := complianceSuiteTestSetup()
	if err != nil {
		t.Fatal(err)
	}
	server.Start()
	defer server.Close()

	restRPCs := map[string][]prepRepeatDataNegativeTestFunc{
		"Compliance.RepeatDataBodyInfo": {
			prepRepeatDataBodyInfoNegativeTestInvalidFields,
			prepRepeatDataBodyInfoNegativeTestSnakeCasedFieldNames,
		},
		"Compliance.RepeatDataQuery": {
			prepRepeatDataQueryNegativeTestSnakeCasedFieldNames,
		},
	}

	for groupIdx, group := range masterSuite.GetGroup() {
		rpcsToTest := group.GetRpcs()
		for requestIdx, masterProto := range group.GetRequests() {
			for rpcIdx, rpcName := range rpcsToTest {
				errorPrefix := fmt.Sprintf("[request %d/%q: rpc %q/%d/%q]",
					requestIdx, masterProto.GetName(), group.Name, rpcIdx, rpcName)

				// Ensure that we issue only the RPCs the test suite is expecting.
				restTest, ok := restRPCs[rpcName]
				if !ok {
					// we don't have a negative test for this RPC
					continue
				}

				for _, rpcPrep := range restTest {
					// since these tests may modify the request protos, get a
					// clean request every time as the starting point, in order
					// to prevent previous modifications from affecting the
					// current test case
					suite, err := getCleanComplianceSuite()
					if err != nil {
						t.Fatal(err)
					}
					requestProto := suite.GetGroup()[groupIdx].GetRequests()[requestIdx]

					prepName, verb, path, requestBody, failure, err := rpcPrep(requestProto)
					if err != nil {
						t.Errorf("%s error: %s", errorPrefix, err)
					}
					if got, want := prepName, rpcName; !strings.HasPrefix(prepName, rpcName) {
						t.Errorf("%s retrieved mismatched prep function: got %q, want %q", errorPrefix, got, want)
					}

					checkExpectedFailure(t, verb, server.URL+path, requestBody, failure, errorPrefix, prepName)
				}
			}
		}
	}
}

// TestComplianceSuiteUnexpectedFieldPresence checks that we detect erroneous presence/absence of
// optional fields.
func TestComplianceSuiteUnexpectedFieldPresence(t *testing.T) {
	suite, server, err := complianceSuiteTestSetup()
	if err != nil {
		t.Fatal(err)
	}
	server.Start()
	defer server.Close()

	indexedSuite, err := services.IndexComplianceSuite(suite)
	if err != nil {
		t.Fatal(err)
	}

	requestsModified := map[string]bool{}
	for idx, testCase := range []struct {
		name        string
		requestName string
		modify      requestModifier
		subCases    map[string]prepRepeatDataTestFunc
		failureMode string
	}{
		{
			name:        "detecting requests not included in test suite",
			requestName: "Zero values for all fields",
			modify: func(request *pb.RepeatRequest) {
				request.Name = "modified name"
			},
			subCases: map[string]prepRepeatDataTestFunc{
				"body":  prepRepeatDataBodyTest,
				"query": prepRepeatDataQueryTest,
			},
			failureMode: "(ComplianceSuiteRequestNotFoundError)",
		},

		{
			name:        "detecting set optional bool field erroneously not sent",
			requestName: "Zero values for all fields",
			modify: func(request *pb.RepeatRequest) {
				request.GetInfo().PBool = nil
			},
			subCases: map[string]prepRepeatDataTestFunc{
				"body":  prepRepeatDataBodyTest,
				"query": prepRepeatDataQueryTest,
			},
			failureMode: "(ComplianceSuiteRequestMismatchError)",
		},
		{
			name:        "detecting unset optional bool field erroneously sent",
			requestName: "Basic types, no optional fields",
			modify: func(request *pb.RepeatRequest) {
				myFalse := false
				request.GetInfo().PBool = &myFalse
			},
			subCases: map[string]prepRepeatDataTestFunc{
				"body":  prepRepeatDataBodyTest,
				"query": prepRepeatDataQueryTest,
			},
			failureMode: "(ComplianceSuiteRequestMismatchError)",
		},

		{
			name:        "detecting set optional string field erroneously not sent",
			requestName: "Zero values for all fields",
			modify: func(request *pb.RepeatRequest) {
				request.GetInfo().PString = nil
			},
			subCases: map[string]prepRepeatDataTestFunc{
				"body":  prepRepeatDataBodyTest,
				"query": prepRepeatDataQueryTest,
			},
			failureMode: "(ComplianceSuiteRequestMismatchError)",
		},
		{
			name:        "detecting unset optional string field erroneously sent",
			requestName: "Basic types, no optional fields",
			modify: func(request *pb.RepeatRequest) {
				myEmpty := ""
				request.GetInfo().PString = &myEmpty
			},
			subCases: map[string]prepRepeatDataTestFunc{
				"body":  prepRepeatDataBodyTest,
				"query": prepRepeatDataQueryTest,
			},
			failureMode: "(ComplianceSuiteRequestMismatchError)",
		},

		{
			name:        "detecting set optional int32 field erroneously not sent",
			requestName: "Zero values for all fields",
			modify: func(request *pb.RepeatRequest) {
				request.GetInfo().PInt32 = nil
			},
			subCases: map[string]prepRepeatDataTestFunc{
				"body":  prepRepeatDataBodyTest,
				"query": prepRepeatDataQueryTest,
			},
			failureMode: "(ComplianceSuiteRequestMismatchError)",
		},
		{
			name:        "detecting unset optional int32 field erroneously sent",
			requestName: "Basic types, no optional fields",
			modify: func(request *pb.RepeatRequest) {
				myInt := int32(0)
				request.GetInfo().PInt32 = &myInt
			},
			subCases: map[string]prepRepeatDataTestFunc{
				"body":  prepRepeatDataBodyTest,
				"query": prepRepeatDataQueryTest,
			},
			failureMode: "(ComplianceSuiteRequestMismatchError)",
		},
	} {
		if _, done := requestsModified[testCase.requestName]; done {
			if suite, err = getCleanComplianceSuite(); err != nil {
				t.Fatal(err)
			}
			if indexedSuite, err = services.IndexComplianceSuite(suite); err != nil {
				t.Fatal(err)
			}
			requestsModified = map[string]bool{}
		}
		requestsModified[testCase.requestName] = true

		request, ok := indexedSuite[testCase.requestName]
		if !ok {
			t.Fatalf("could not find request by name: %q", testCase.requestName)
		}
		testCase.modify(request)

		for subCaseName, prep := range testCase.subCases {
			prefix := fmt.Sprintf("[case %d: %s: %s]", idx, testCase.name, subCaseName)
			verb, prepName, path, body, error := prep(request)
			if error != nil {
				t.Fatalf("%s could not construct request: %s", prefix, err)
			}
			checkExpectedFailure(t, verb, server.URL+path, body, testCase.failureMode, prefix, prepName)
		}
	}
}

type prepRepeatDataNegativeTestFunc func(request *genproto.RepeatRequest) (verb string, name string, path string, body string, expect string, err error)

func prepRepeatDataBodyInfoNegativeTestInvalidFields(request *genproto.RepeatRequest) (verb string, name string, path string, body string, expect string, err error) {
	resttools.JSONMarshaler.Replace(nil)
	defer resttools.JSONMarshaler.Restore()

	name = "Compliance.RepeatDataBodyInfo#NegativeTestInvalidFields"
	bodyBytes, err := resttools.ToJSON().Marshal(request.Info)
	queryString := prepRepeatDataTestsQueryString(request, nil) // purposefully repeats query params, which should cause an error
	return name, "POST", "/v1beta1/repeat:bodyinfo" + queryString, string(bodyBytes), "(QueryParamsInvalidFieldError)", err
}

func prepRepeatDataBodyInfoNegativeTestSnakeCasedFieldNames(request *genproto.RepeatRequest) (verb string, name string, path string, body string, expect string, err error) {
	resttools.JSONMarshaler.Replace(&protojson.MarshalOptions{
		Multiline:       true,
		AllowPartial:    false,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true, // this should cause an error
	})
	defer resttools.JSONMarshaler.Restore()

	name = "Compliance.RepeatDataBodyInfo#NegativeTestSnakeCasedFieldNames"
	request.Info.FString += name
	bodyBytes, err := resttools.ToJSON().Marshal(request.Info)
	return name, "POST", "/v1beta1/repeat:bodyinfo", string(bodyBytes), "(BodyFieldNameIncorrectlyCasedError)", err
}

func prepRepeatDataQueryNegativeTestSnakeCasedFieldNames(request *genproto.RepeatRequest) (verb string, name string, path string, body string, expect string, err error) {
	name = "Compliance.RepeatDataQuery#NegativeTestSnakeCasedFieldNames"
	queryParams := prepRepeatDataTestsQueryParams(request, nil, queryStringSnakeCaser) // this should cause an error
	queryString := prepQueryString(queryParams)
	return name, "GET", "/v1beta1/repeat:query" + queryString, body, "(QueryParameterNameIncorrectlyCasedError)", err
}

// checkExpectedFailure issues a request using the specified verb, URL, and request body. It expects
// a failing HTTP code and a response message containing the substring in `failure`. Test errors are
// reported using the given errorPrefix and the name prepName of the prepping function.
func checkExpectedFailure(t *testing.T, verb, url, requestBody, failure, errorPrefix, prepName string) {
	// Issue the request
	httpRequest, err := http.NewRequest(verb, url, strings.NewReader(requestBody))
	if err != nil {
		t.Errorf("%s error creating request: %s", errorPrefix, err)
		return
	}
	resttools.PopulateRequestHeaders(httpRequest)
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		t.Errorf("%s error issuing call: %s", errorPrefix, err)
		return
	}

	// Check for unsuccessful response.
	if got, notWant := httpResponse.StatusCode, http.StatusOK; got == notWant {
		t.Errorf("%s response code: got %d, notWant %d  name:%q\n   %s %s\nrequest body: %s\n----------------------------------------\n",
			errorPrefix, got, notWant, prepName, verb, url, requestBody)
		return
	}

	body, err := io.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		t.Fatalf("%s could not read response body: %s", errorPrefix, err)
	}
	if got, want := string(body), failure; !strings.Contains(got, want) {
		t.Errorf("%s response body: wanted response to include %q, but instead got: %q   (status %d) header: %v name:%q\n   %s %s\nrequest body: %s\n----------------------------------------\n",
			errorPrefix, want, got, httpResponse.StatusCode, httpResponse.Header, prepName, verb, url, requestBody)
	}

}

// requestModifer is a function that modifies a request in-place.
type requestModifier func(*pb.RepeatRequest)
