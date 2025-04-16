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

func TestErrorResponse(t *testing.T) {

	for _, tst := range []struct {
		details         []interface{}
		message, name   string
		inputHTTPStatus int
		inputGRPCCode   codes.Code
		wantHTTPStatus  int
		wantResponse    string
	}{
		// HTTP → gRPC codes
		{
			name:            "internal_server",
			message:         "Had an issue",
			inputGRPCCode:   NoCodeGRPC,
			inputHTTPStatus: http.StatusInternalServerError,
			wantHTTPStatus:  http.StatusInternalServerError,
			details:         []interface{}{&errdetails.ErrorInfo{Reason: "foo"}},
			wantResponse:    `{"error":{"code":500, "message":"Had an issue", "status":"INTERNAL", "details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo", "reason":"foo"}]}}`,
		},
		{
			name:            "bad_request",
			message:         "The request was bad",
			inputGRPCCode:   NoCodeGRPC,
			inputHTTPStatus: http.StatusBadRequest,
			wantHTTPStatus:  http.StatusBadRequest,
			wantResponse:    `{"error":{"code":400, "message":"The request was bad", "status":"INVALID_ARGUMENT"}}`,
		},
		{
			name:            "multiple_issues",
			message:         "Had multiple issues",
			inputGRPCCode:   NoCodeGRPC,
			inputHTTPStatus: http.StatusInternalServerError,
			wantHTTPStatus:  http.StatusInternalServerError,
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

		// gRPC → HTTP codes
		{
			name:            "internal_server",
			message:         "Had an issue",
			inputGRPCCode:   codes.Internal,
			inputHTTPStatus: NoCodeHTTP,
			wantHTTPStatus:  http.StatusInternalServerError,
			details:         []interface{}{&errdetails.ErrorInfo{Reason: "foo"}},
			wantResponse:    `{"error":{"code":500, "message":"Had an issue", "status":"INTERNAL", "details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo", "reason":"foo"}]}}`,
		},
		{
			name:            "bad_request",
			message:         "The request was bad",
			inputGRPCCode:   codes.InvalidArgument,
			inputHTTPStatus: NoCodeHTTP,
			wantHTTPStatus:  http.StatusBadRequest,
			wantResponse:    `{"error":{"code":400, "message":"The request was bad", "status":"INVALID_ARGUMENT"}}`,
		},
		{
			name:            "multiple_issues",
			message:         "Had multiple issues",
			inputGRPCCode:   codes.Internal,
			inputHTTPStatus: NoCodeHTTP,
			wantHTTPStatus:  http.StatusInternalServerError,
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

		// error inputs
		{
			name:            "test when neither gRPC nor HTTP code is specified",
			message:         "Had an issue",
			inputGRPCCode:   NoCodeGRPC,
			inputHTTPStatus: NoCodeHTTP,
			wantHTTPStatus:  http.StatusInternalServerError,
			details:         []interface{}{&errdetails.ErrorInfo{Reason: "foo"}},
			wantResponse:    `{"error":{"code":500, "message":"Showcase consistency error: neither HTTP code or gRPC status are provided for ErrorResponse. Exactly one must be provided.", "status":"INTERNAL"}}`,
		},
		{
			name:            "test when both gRPC and HTTP codes are specified",
			message:         "Had an issue",
			inputGRPCCode:   codes.Internal,
			inputHTTPStatus: http.StatusBadRequest,
			wantHTTPStatus:  http.StatusInternalServerError,
			details:         []interface{}{&errdetails.ErrorInfo{Reason: "foo"}},
			wantResponse:    `{"error":{"code":500, "message":"Showcase consistency error: both HTTP code and gRPC status are provided for ErrorResponse. Exactly one must be provided.", "status":"INTERNAL"}}`,
		},
	} {
		got := httptest.NewRecorder()
		ErrorResponse(got, tst.inputHTTPStatus, tst.inputGRPCCode, tst.message, tst.details...)
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
			t.Errorf("converting %s: got %d, but expected %d", tst.code, got, tst.want)
		}
	}
}

func TestHTTPToGRPC(t *testing.T) {
	// This test focuses on the ranges of HTTP codes that map to a single gRPC status, as per go/http-canonical-mapping.
	for _, tst := range []struct {
		code int
		want codes.Code
	}{
		{200, codes.OK},
		{299, codes.OK},
		{350, codes.Unknown},
		{403, codes.PermissionDenied},
		{499, codes.Canceled},
		{498, codes.FailedPrecondition},
		{503, codes.Unavailable},
		{599, codes.Internal},
		{149, codes.Unknown},
	} {
		if got := HTTPToGRPC(tst.code); got != tst.want {
			t.Errorf("converting %d: got %d, but expected %d", tst.code, got, tst.want)
		}
	}
}
