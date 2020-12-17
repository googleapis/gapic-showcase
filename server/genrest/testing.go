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
	"fmt"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gmux "github.com/gorilla/mux"

	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
)

// HandleCreateSession translates REST requests/responses on the wire to internal proto messages for CreateSession
//    Generated for HTTP binding pattern: /v1beta1/sessions
//         This matches URIs of the form: /v1beta1/sessions
func (backend *RESTBackend) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/sessions': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.CreateSessionRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var bodyField genprotopb.Session
	if err := jsonpb.Unmarshal(r.Body, &bodyField); err != nil {
		backend.StdLog.Printf(`  error reading body into request field "session": %s`, err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}
	request.Session = &bodyField

	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.CreateSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

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
//         This matches URIs of the form: /v1beta1/{name:sessions/[0-9a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=sessions/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.GetSessionRequest{}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.GetSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

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
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/sessions': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.ListSessionsRequest{}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.ListSessions(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

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
//         This matches URIs of the form: /v1beta1/{name:sessions/[0-9a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=sessions/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.DeleteSessionRequest{}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.DeleteSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

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
//         This matches URIs of the form: /v1beta1/{name:sessions/[0-9a-zA-Z_%\-]+}:report
func (backend *RESTBackend) HandleReportSession(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=sessions/*}:report': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.ReportSessionRequest{}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.ReportSession(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

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
//         This matches URIs of the form: /v1beta1/{parent:sessions/[0-9a-zA-Z_%\-]+}/tests
func (backend *RESTBackend) HandleListTests(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{parent=sessions/*}/tests': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.ListTestsRequest{}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.ListTests(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

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
//         This matches URIs of the form: /v1beta1/{name:sessions/[0-9a-zA-Z_%\-]+/tests/[0-9a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteTest(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=sessions/*/tests/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.DeleteTestRequest{}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.DeleteTest(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

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
//         This matches URIs of the form: /v1beta1/{name:sessions/[0-9a-zA-Z_%\-]+/tests/[0-9a-zA-Z_%\-]+}:check
func (backend *RESTBackend) HandleVerifyTest(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=sessions/*/tests/*}:check': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.VerifyTestRequest{}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	// TODO: Ensure we handle URL-encoded values in query parameters
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.StdLog.Printf("  error reading query params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.TestingServer.VerifyTest(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}
