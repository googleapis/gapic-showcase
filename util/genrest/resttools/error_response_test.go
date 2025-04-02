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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
)

func TestErrorResponse(t *testing.T) {

	for _, tst := range []struct {
		details       []interface{}
		message, name string
		status        int
	}{
		{
			name:    "internal_server",
			message: "Had an issue",
			status:  http.StatusInternalServerError,
			details: []interface{}{"foo"},
		},
		{
			name:    "bad_request",
			message: "The request was bad",
			status:  http.StatusBadRequest,
		},
	} {
		got := httptest.NewRecorder()
		ErrorResponse(got, tst.status, NoCodeGRPC, tst.message, tst.details...)
		if got.Code != tst.status {
			t.Errorf("%s: Expected %d, but got %d", tst.name, tst.status, got.Code)
		}
		err := googleapi.CheckResponse(got.Result())
		var gerr *googleapi.Error
		if !errors.As(err, &gerr) {
			t.Fatalf("%s: Expected response to be a googleapi.Error, but got %v", tst.name, err)
		}

		if diff := cmp.Diff(gerr.Message, tst.message); diff != "" {
			t.Errorf("%s: got(-),want(+):%s\n", tst.name, diff)
		}
		if diff := cmp.Diff(gerr.Details, tst.details); diff != "" {
			t.Errorf("%s: got(-),want(+):%s\n", tst.name, diff)
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
