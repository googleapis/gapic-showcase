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

package services

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/proto"
)

func TestComplianceRepeats(t *testing.T) {
	// Note that additional thorough test cases are exercised in
	// cmd/gapic-showcase/compliance_suite_test.go.
	server := NewComplianceServer()
	info := &pb.ComplianceData{
		FString:   "Terra Incognita",
		FInt32:    1,
		FSint32:   -2,
		FSfixed32: -300000000,
		FUint32:   5,
		FFixed32:  700000000,
		FInt64:    9,
		FSint64:   -1100000000,
		FSfixed64: -1300000000,
		FUint64:   1700000000000000000,
		FFixed64:  1900000000000000000,

		FDouble: 6.02e23,
		FFloat:  3.1415,
		FBool:   true,
		FBytes:  []byte("Lorem ipsum"),
	}
	request := &pb.RepeatRequest{Info: info}

	for idx, rpc := range [](func(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error)){
		server.RepeatDataQuery,
		server.RepeatDataBody,
		server.RepeatDataBodyInfo,
		server.RepeatDataSimplePath,
		server.RepeatDataPathResource,
		server.RepeatDataPathTrailingResource,
		server.RepeatDataBodyPut,
		server.RepeatDataBodyPatch,
		server.RepeatWithUnknownEnum,
		server.RepeatWithUnknownOptionalEnum,
	} {
		response, err := rpc(context.Background(), request)
		if err != nil {
			t.Errorf("call %d: error: %s", idx, err)
		}
		if diff := cmp.Diff(response.GetRequest(), request, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("call %d: got=-, want=+:%s", idx, diff)
		}
	}
}

func TestMatchingComplianceSuiteRequests(t *testing.T) {
	server := &complianceServerImpl{}

	info := &pb.ComplianceData{
		FString:   "Terra Incognita",
		FInt32:    1,
		FSint32:   -2,
		FSfixed32: -300000000,
	}
	request := &pb.RepeatRequest{
		Name: "Basic data types", // matches a name in compliance_suite.json
		Info: info,
	}

	if got := server.requestMatchesExpectation(request); got != nil {
		t.Errorf("expected request to trivially match when serverVerify unset. Got error: %s", got)
	}

	request.ServerVerify = true
	if err := server.requestMatchesExpectation(request); err == nil {
		t.Errorf("expected verified request with differing data to not match")
	} else {
		if got, want := err.Error(), "(ComplianceSuiteRequestMismatchError)"; !strings.Contains(got, want) {
			t.Errorf("error message does not contain expected substring: want: %q  got %q", want, got)
		}
		if _, got := server.Repeat(context.Background(), request); got == nil {
			t.Errorf("expected Repeat() to error with unverified request, but it didn't")
		}
	}

	request.Name = "non-existent case"
	if err := server.requestMatchesExpectation(request); err == nil {
		t.Errorf("expected verified request with unmatched name to cause an error")
	} else {
		if got, want := err.Error(), "(ComplianceSuiteRequestNotFoundError)"; !strings.Contains(got, want) {
			t.Errorf("error message does not contain expected substring: want: %q  got %q", want, got)
		}
		if _, got := server.Repeat(context.Background(), request); got == nil {
			t.Errorf("expected Repeat() to error with unverified request, but it didn't")
		}
	}

	request = ComplianceSuiteRequests["Basic data types"] // matches a name in compliance_suite.json
	if got := server.requestMatchesExpectation(request); got != nil {
		t.Errorf("expected test suite case to match. Got error: %s", got)
	}
	if _, got := server.Repeat(context.Background(), request); got != nil {
		t.Errorf("expected Repeat() to succeed with verified request, but got error: %s", got)
	}
}

func TestIndexingComplianceSuite(t *testing.T) {
	// set up
	ComplianceSuiteStatus = ComplianceSuiteUninitialized

	suiteBytes := []byte("nonexistent_field: 5 ")

	if err := indexTestingRequests(suiteBytes); err == nil {
		t.Errorf("expected JSON unmarshaling to fail, but it succeeded")
	} else {
		if got, want := err.Error(), "(ComplianceServiceReadError)"; !strings.Contains(got, want) {
			t.Errorf("error message does not contain expected substring: want: %q  got %q", want, got)
		}
	}

	suiteBytes = []byte(`
            {
              "group": [
                 {
                  "name": "sample suite",
                  "requests": [
                     { "name": "Alpha"},
                     { "name": "Beta"},
                     { "name": "Alpha"}
                  ]
                 }
                ]
             }
               `)

	if err := indexTestingRequests(suiteBytes); err == nil {
		t.Errorf("expected JSON unmarshaling to fail, but it succeeded")
	} else {
		if got, want := err.Error(), "(ComplianceServiceSetupError)"; !strings.Contains(got, want) {
			t.Errorf("error message does not contain expected substring: want: %q  got %q", want, got)
		}
	}

	// test that the indexing error gets properly propagated
	server := &complianceServerImpl{}
	request := &pb.RepeatRequest{
		Name:         "Basic data types", // matches a name in compliance_suite.json
		ServerVerify: true,
	}
	if err := server.requestMatchesExpectation(request); err == nil {
		t.Errorf("expected verified request with differing data to not match")
	} else {
		if got, want := err.Error(), "(ComplianceServiceSetupError)"; !strings.Contains(got, want) {
			t.Errorf("error message does not contain expected substring: want: %q  got %q", want, got)
		}
		if _, got := server.Repeat(context.Background(), request); got == nil {
			t.Errorf("expected Repeat() to error with unverified request, but it didn't")
		}
	}

	// clean up
	ComplianceSuiteStatus = ComplianceSuiteUninitialized
	if err := indexTestingRequests(complianceSuiteBytes); err != nil {
		t.Errorf("initializing ComplianceSuite with real suite data should not have errored, but got: %s", err)
	}

	if err := indexTestingRequests(complianceSuiteBytes); err != nil {
		t.Errorf("initializing ComplianceSuite a second time should not have errored, but got: %s", err)
	}
}
