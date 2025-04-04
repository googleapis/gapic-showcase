// Copyright 2022 Google LLC
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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
)

func TestErrorResponseFromHTTP(t *testing.T) {

	for _, tst := range []struct {
		details       []interface{}
		message, name string
		httpStatus    int
		wantResponse  string
	}{
		{
			name:         "internal_server",
			message:      "Had an issue",
			httpStatus:   http.StatusInternalServerError,
			details:      []interface{}{&errdetails.ErrorInfo{Reason: "foo"}},
			wantResponse: `{"error":{"code":500, "message":"Had an issue", "status":"INTERNAL", "details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo", "reason":"foo"}]}}`,
		},
		{
			name:         "bad_request",
			message:      "The request was bad",
			httpStatus:   http.StatusBadRequest,
			wantResponse: `{"error":{"code":400, "message":"The request was bad", "status":"INVALID_ARGUMENT"}}`,
		},
		{
			name:       "multiple_issues",
			message:    "Had multiple issues",
			httpStatus: http.StatusInternalServerError,
			details: []interface{}{
				&errdetails.ErrorInfo{Reason: "foo"},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequest_FieldViolation{
						{
							Field:       "an offending field",
							Description: "a description",
							Reason:      "a reason",
						},
					},
				},
			},
			wantResponse: `{"error":{"code":500, "message":"Had multiple issues", "status":"INTERNAL", "details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo", "reason":"foo"}, {"@type":"type.googleapis.com/google.rpc.BadRequest", "fieldViolations":[{"field":"an offending field", "description":"a description", "reason":"a reason"}]}]}}`,
		},
	} {
		got := httptest.NewRecorder()
		ErrorResponse(got, tst.httpStatus, NoCodeGRPC, tst.message, tst.details...)
		if got.Code != tst.httpStatus {
			t.Errorf("%s: Expected code %d, but got %d", tst.name, tst.httpStatus, got.Code)
		}
		gotResponse, err := io.ReadAll(got.Result().Body)
		if err != nil {
			t.Fatalf("%s: Error reading response body: %+v", tst.name, err)
		}
		var gotJSON interface{}
		err = json.Unmarshal([]byte(gotResponse), &gotJSON)
		if err != nil {
			t.Fatalf("%s: Error parsing actual response body as JSON: %+v", tst.name, err)
		}

		var wantJSON interface{}
		err = json.Unmarshal([]byte(tst.wantResponse), &wantJSON)
		if err != nil {
			t.Fatalf("%s: Error parsing expected response body as JSON: %+v", tst.name, err)
		}

		if diff := cmp.Diff(gotJSON, wantJSON); diff != "" {
			t.Errorf("%s: error body: got(-),want(+):%s\n\n---------- Raw JSON: got\n%s\n---------- Raw JSON: want\n%s",
				tst.name, diff, gotResponse, tst.wantResponse)
		}
	}
}

func TestErrorResponseFromGRPC(t *testing.T) {

	for _, tst := range []struct {
		details        []interface{}
		message, name  string
		grpcCode       codes.Code
		wantHTTPStatus int
		wantResponse   string
	}{
		{
			name:           "internal_server",
			message:        "Had an issue",
			grpcCode:       codes.Internal,
			details:        []interface{}{&errdetails.ErrorInfo{Reason: "foo"}},
			wantHTTPStatus: http.StatusInternalServerError,
			wantResponse:   `{"error":{"code":500, "message":"Had an issue", "status":"INTERNAL", "details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo", "reason":"foo"}]}}`,
		},
		{
			name:           "bad_request",
			message:        "The request was bad",
			grpcCode:       codes.InvalidArgument,
			wantHTTPStatus: http.StatusBadRequest,
			wantResponse:   `{"error":{"code":400, "message":"The request was bad", "status":"INVALID_ARGUMENT"}}`,
		},
		{
			name:           "multiple_issues",
			message:        "Had multiple issues",
			grpcCode:       codes.Internal,
			wantHTTPStatus: http.StatusInternalServerError,
			details: []interface{}{
				&errdetails.ErrorInfo{Reason: "foo"},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequest_FieldViolation{
						{
							Field:       "an offending field",
							Description: "a description",
							Reason:      "a reason",
						},
					},
				},
			},
			wantResponse: `{"error":{"code":500, "message":"Had multiple issues", "status":"INTERNAL", "details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo", "reason":"foo"}, {"@type":"type.googleapis.com/google.rpc.BadRequest", "fieldViolations":[{"field":"an offending field", "description":"a description", "reason":"a reason"}]}]}}`,
		},
	} {
		got := httptest.NewRecorder()
		ErrorResponse(got, NoCodeHTTP, tst.grpcCode, tst.message, tst.details...)
		if got.Code != tst.wantHTTPStatus {
			t.Errorf("%s: Expected code %d, but got %d", tst.name, tst.wantHTTPStatus, got.Code)
		}
		gotResponse, err := io.ReadAll(got.Result().Body)
		if err != nil {
			t.Fatalf("%s: Error reading response body: %+v", tst.name, err)
		}
		var gotJSON interface{}
		err = json.Unmarshal([]byte(gotResponse), &gotJSON)
		if err != nil {
			t.Fatalf("%s: Error parsing actual response body as JSON: %+v", tst.name, err)
		}

		var wantJSON interface{}
		err = json.Unmarshal([]byte(tst.wantResponse), &wantJSON)
		if err != nil {
			t.Fatalf("%s: Error parsing expected response body as JSON: %+v", tst.name, err)
		}

		if diff := cmp.Diff(gotJSON, wantJSON); diff != "" {
			t.Errorf("%s: error body: got(-),want(+):%s\n\n---------- Raw JSON: got\n%s\n---------- Raw JSON: want\n%s",
				tst.name, diff, gotResponse, tst.wantResponse)
		}
	}
}

func TestGRPCToHTTP(t *testing.T) {
	for _, tst := range []struct {
		code codes.Code
		want int
	}{
		{
			codes.Aborted,
			http.StatusConflict,
		},
		{
			100,
			http.StatusInternalServerError,
		},
	} {
		if got := GRPCToHTTP(tst.code); got != tst.want {
			t.Errorf("got %d, but expected %d", got, tst.want)
		}
	}
}
