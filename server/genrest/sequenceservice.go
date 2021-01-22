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
// for service #4: "SequenceService" (.google.showcase.v1beta1.SequenceService).

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

// HandleCreateSequence translates REST requests/responses on the wire to internal proto messages for CreateSequence
//    Generated for HTTP binding pattern: "/v1beta1/sequences"
//         This matches URIs of the form: "/v1beta1/sequences"
func (backend *RESTBackend) HandleCreateSequence(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/sequences': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.CreateSequenceRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var bodyField genprotopb.Sequence
	if err := jsonpb.Unmarshal(r.Body, &bodyField); err != nil {
		backend.StdLog.Printf(`  error reading body into request field "sequence": %s`, err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}
	request.Sequence = &bodyField

	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
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

	response, err := backend.SequenceServiceServer.CreateSequence(context.Background(), request)
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

// HandleGetSequenceReport translates REST requests/responses on the wire to internal proto messages for GetSequenceReport
//    Generated for HTTP binding pattern: "/v1beta1/{name=sequences/*/sequenceReport}"
//         This matches URIs of the form: "/v1beta1/{name:sequences/.+/sequenceReport}"
func (backend *RESTBackend) HandleGetSequenceReport(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=sequences/*/sequenceReport}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.GetSequenceReportRequest{}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
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

	response, err := backend.SequenceServiceServer.GetSequenceReport(context.Background(), request)
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

// HandleAttemptSequence translates REST requests/responses on the wire to internal proto messages for AttemptSequence
//    Generated for HTTP binding pattern: "/v1beta1/{name=sequences/*}"
//         This matches URIs of the form: "/v1beta1/{name:sequences/.+}"
func (backend *RESTBackend) HandleAttemptSequence(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=sequences/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.AttemptSequenceRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	if err := jsonpb.Unmarshal(r.Body, request); err != nil {
		backend.StdLog.Printf(`  error reading body params "*": %s`, err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error reading URL path params: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.SequenceServiceServer.AttemptSequence(context.Background(), request)
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
