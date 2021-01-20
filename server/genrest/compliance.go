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

// DO NOT EDIT. This is an auto-generated file containing the REST handlers
// for service #0: "Compliance" (.google.showcase.v1beta1.Compliance).

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

// HandleRepeatDataBody translates REST requests/responses on the wire to internal proto messages for RepeatDataBody
//    Generated for HTTP binding pattern: /v1beta1/repeat:body
//         This matches URIs of the form: /v1beta1/repeat:body
func (backend *RESTBackend) HandleRepeatDataBody(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat:body': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.RepeatRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	if err := jsonpb.Unmarshal(r.Body, request); err != nil {
		backend.StdLog.Printf(`  error reading body params "*": %s`, err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}
	// TODO: Ensure we handle URL-encoded values in path variables
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.ComplianceServer.RepeatDataBody(context.Background(), request)
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

// HandleRepeatDataQuery translates REST requests/responses on the wire to internal proto messages for RepeatDataQuery
//    Generated for HTTP binding pattern: /v1beta1/repeat:query
//         This matches URIs of the form: /v1beta1/repeat:query
func (backend *RESTBackend) HandleRepeatDataQuery(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat:query': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.RepeatRequest{}
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

	response, err := backend.ComplianceServer.RepeatDataQuery(context.Background(), request)
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

// HandleRepeatDataSimplePath translates REST requests/responses on the wire to internal proto messages for RepeatDataSimplePath
//    Generated for HTTP binding pattern: /v1beta1/repeat/{info.f_string}/{info.f_int32}/{info.f_double}/{info.f_bool}:simplepath
//         This matches URIs of the form: /v1beta1/repeat/{info.f_string:[0-9a-zA-Z_% "\-]+}/{info.f_int32:[0-9a-zA-Z_% "\-]+}/{info.f_double:[0-9a-zA-Z_% "\-]+}/{info.f_bool:[0-9a-zA-Z_% "\-]+}:simplepath
func (backend *RESTBackend) HandleRepeatDataSimplePath(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat/{info.f_string}/{info.f_int32}/{info.f_double}/{info.f_bool}:simplepath': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 4, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 4 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 4, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.RepeatRequest{}
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

	response, err := backend.ComplianceServer.RepeatDataSimplePath(context.Background(), request)
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
