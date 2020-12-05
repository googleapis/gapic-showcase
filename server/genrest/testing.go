// Copyright 2020 Google LLC
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

// DO NOT EDIT. This is an auto-generated file containing the REST handlers
// for service #4: "Testing" (.google.showcase.v1beta1.Testing).

package genrest

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
)

// HandleCreateSession translates REST requests/responses on the wire to internal proto messages for CreateSession
//    Generated for HTTP binding pattern: /v1beta1/sessions
//         This matches URIs of the form: /v1beta1/sessions
func (backend *RESTBackend) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/sessions': %q", r.URL)

	var request *genprotopb.CreateSessionRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.CreateSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleGetSession translates REST requests/responses on the wire to internal proto messages for GetSession
//    Generated for HTTP binding pattern: /v1beta1/{name=sessions/*}
//         This matches URIs of the form: /v1beta1/{name:sessions/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=sessions/*}': %q", r.URL)

	var request *genprotopb.GetSessionRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.GetSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleListSessions translates REST requests/responses on the wire to internal proto messages for ListSessions
//    Generated for HTTP binding pattern: /v1beta1/sessions
//         This matches URIs of the form: /v1beta1/sessions
func (backend *RESTBackend) HandleListSessions(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/sessions': %q", r.URL)

	var request *genprotopb.ListSessionsRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.ListSessions(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleDeleteSession translates REST requests/responses on the wire to internal proto messages for DeleteSession
//    Generated for HTTP binding pattern: /v1beta1/{name=sessions/*}
//         This matches URIs of the form: /v1beta1/{name:sessions/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=sessions/*}': %q", r.URL)

	var request *genprotopb.DeleteSessionRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.DeleteSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleReportSession translates REST requests/responses on the wire to internal proto messages for ReportSession
//    Generated for HTTP binding pattern: /v1beta1/{name=sessions/*}:report
//         This matches URIs of the form: /v1beta1/{name:sessions/[a-zA-Z_%\-]+}:report
func (backend *RESTBackend) HandleReportSession(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=sessions/*}:report': %q", r.URL)

	var request *genprotopb.ReportSessionRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.ReportSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleListTests translates REST requests/responses on the wire to internal proto messages for ListTests
//    Generated for HTTP binding pattern: /v1beta1/{parent=sessions/*}/tests
//         This matches URIs of the form: /v1beta1/{parent:sessions/[a-zA-Z_%\-]+}/tests
func (backend *RESTBackend) HandleListTests(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{parent=sessions/*}/tests': %q", r.URL)

	var request *genprotopb.ListTestsRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.ListTests(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleDeleteTest translates REST requests/responses on the wire to internal proto messages for DeleteTest
//    Generated for HTTP binding pattern: /v1beta1/{name=sessions/*/tests/*}
//         This matches URIs of the form: /v1beta1/{name:sessions/[a-zA-Z_%\-]+/tests/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteTest(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=sessions/*/tests/*}': %q", r.URL)

	var request *genprotopb.DeleteTestRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.DeleteTest(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleVerifyTest translates REST requests/responses on the wire to internal proto messages for VerifyTest
//    Generated for HTTP binding pattern: /v1beta1/{name=sessions/*/tests/*}:check
//         This matches URIs of the form: /v1beta1/{name:sessions/[a-zA-Z_%\-]+/tests/[a-zA-Z_%\-]+}:check
func (backend *RESTBackend) HandleVerifyTest(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=sessions/*/tests/*}:check': %q", r.URL)

	var request *genprotopb.VerifyTestRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.TestingServer.VerifyTest(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}
