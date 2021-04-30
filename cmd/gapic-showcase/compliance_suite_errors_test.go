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
	"net/url"
	"strings"
	"testing"

	"github.com/googleapis/gapic-showcase/server/genproto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/services"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestComplianceSuiteVerifyErrors(t *testing.T) {
	if services.ComplianceSuiteStatus != services.ComplianceSuiteLoaded {
		t.Fatalf("compliance suite was not loaded: status %#v %s", services.ComplianceSuiteStatus, services.ComplianceSuiteStatusMessage)
	}
	// 	suite, server, err := complianceSuiteTestSetup()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	server.Start()
	// 	defer server.Close()

}

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

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		t.Fatalf("%s could not read response body: %s", errorPrefix, err)
	}
	if got, want := string(body), failure; !strings.Contains(got, want) {
		t.Errorf("%s response body: wanted response to include %q, but instead got: %q   name:%q\n   %s %s\nrequest body: %s\n----------------------------------------\n",
			errorPrefix, want, got, prepName, verb, url, requestBody)
	}

}

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
			prepRepeatDataQueryNegativeTestNumericEnums,
			prepRepeatDataQueryNegativeTestNumericOptionalEnums,
			prepRepeatDataQueryNegativeTestSnakeCasedFieldNames,
		},
		"Compliance.RepeatDataSimplePath": {prepRepeatDataSimplePathNegativeTestEnum},
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

func prepRepeatDataQueryNegativeTestNumericEnums(request *genproto.RepeatRequest) (verb string, name string, path string, body string, expect string, err error) {
	name = "Compliance.RepeatDataQuery#NegativeTestNumericEnums"
	info := request.GetInfo()
	badQueryParam := fmt.Sprintf("info.fKingdom=%d", info.GetFKingdom()) // purposefully use a number, which should cause an error

	// We clear the field so we don't set the same query param correctly below. This change
	// modifies the request, but since these tests only check that calls fail, we never need to
	// refer back to the request proto after constructing the REST query.
	info.FKingdom = pb.ComplianceData_LIFE_KINGDOM_UNSPECIFIED
	queryParams := append(prepRepeatDataTestsQueryParams(request, nil, queryStringLowerCamelCaser), badQueryParam)

	queryString := prepQueryString(queryParams)
	return name, "GET", "/v1beta1/repeat:query" + queryString, body, "(EnumValueNotStringError)", err
}

func prepRepeatDataQueryNegativeTestNumericOptionalEnums(request *genproto.RepeatRequest) (verb string, name string, path string, body string, expect string, err error) {
	name = "Compliance.RepeatDataQuery#NegativeTestNumericOptionalEnums"
	info := request.GetInfo()
	badQueryParam := fmt.Sprintf("info.pKingdom=%d", info.GetPKingdom()) // purposefully use a number, which should cause an error

	// We clear the field so we don't set the same query param correctly below. This change
	// modifies the request, but since these tests only check that calls fail, we never need to
	// refer back to the request proto after constructing the REST query.
	info.PKingdom = nil
	queryParams := append(prepRepeatDataTestsQueryParams(request, nil, queryStringLowerCamelCaser), badQueryParam)

	queryString := prepQueryString(queryParams)
	return name, "GET", "/v1beta1/repeat:query" + queryString, body, "(EnumValueNotStringError)", err
}

func prepRepeatDataSimplePathNegativeTestEnum(request *genproto.RepeatRequest) (verb string, name string, path string, body string, expect string, err error) {
	name = "Compliance.RepeatDataSimplePath#NegativeTestNumericEnum"
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
	return name, "GET", path + queryString, body, "(EnumValueNotStringError)", err
}
