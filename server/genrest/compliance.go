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
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gmux "github.com/gorilla/mux"

	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
)

// HandleRepeatDataBody translates REST requests/responses on the wire to internal proto messages for RepeatDataBody
//    Generated for HTTP binding pattern: "/v1beta1/repeat:body"
func (backend *RESTBackend) HandleRepeatDataBody(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat:body': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.RepeatRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	if err := jsonpb.Unmarshal(r.Body, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if queryParams := r.URL.Query(); len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.ComplianceServer.RepeatDataBody(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.MarshalToString(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write([]byte(json))
}

// HandleRepeatDataBodyInfo translates REST requests/responses on the wire to internal proto messages for RepeatDataBodyInfo
//    Generated for HTTP binding pattern: "/v1beta1/repeat:bodyinfo"
func (backend *RESTBackend) HandleRepeatDataBodyInfo(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat:bodyinfo': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.RepeatRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var bodyField genprotopb.ComplianceData
	if err := jsonpb.Unmarshal(r.Body, &bodyField); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body into request field 'info': %s", err)
		return
	}
	request.Info = &bodyField

	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	queryParams := map[string][]string(r.URL.Query())
	excludedQueryParams := []string{"info"}
	if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {
		backend.Error(w, http.StatusBadRequest, " found keys that should not appear in query params: %v", duplicates)
		return
	}
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.ComplianceServer.RepeatDataBodyInfo(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.MarshalToString(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write([]byte(json))
}

// HandleRepeatDataQuery translates REST requests/responses on the wire to internal proto messages for RepeatDataQuery
//    Generated for HTTP binding pattern: "/v1beta1/repeat:query"
func (backend *RESTBackend) HandleRepeatDataQuery(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat:query': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.RepeatRequest{}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	queryParams := map[string][]string(r.URL.Query())
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.ComplianceServer.RepeatDataQuery(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.MarshalToString(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write([]byte(json))
}

// HandleRepeatDataSimplePath translates REST requests/responses on the wire to internal proto messages for RepeatDataSimplePath
//    Generated for HTTP binding pattern: "/v1beta1/repeat/{info.f_string}/{info.f_int32}/{info.f_double}/{info.f_bool}:simplepath"
func (backend *RESTBackend) HandleRepeatDataSimplePath(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat/{info.f_string}/{info.f_int32}/{info.f_double}/{info.f_bool}:simplepath': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 4, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 4 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 4, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.RepeatRequest{}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	queryParams := map[string][]string(r.URL.Query())
	excludedQueryParams := []string{"info.f_string", "info.f_int32", "info.f_double", "info.f_bool"}
	if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {
		backend.Error(w, http.StatusBadRequest, " found keys that should not appear in query params: %v", duplicates)
		return
	}
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.ComplianceServer.RepeatDataSimplePath(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.MarshalToString(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write([]byte(json))
}

// HandleRepeatDataPathResource translates REST requests/responses on the wire to internal proto messages for RepeatDataPathResource
//    Generated for HTTP binding pattern: "/v1beta1/repeat/{info.f_string=first/*}/{info.f_child.f_string=second/*}/bool/{info.f_bool}:pathresource"
func (backend *RESTBackend) HandleRepeatDataPathResource(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat/{info.f_string=first/*}/{info.f_child.f_string=second/*}/bool/{info.f_bool}:pathresource': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 3, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 3 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 3, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.RepeatRequest{}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	queryParams := map[string][]string(r.URL.Query())
	excludedQueryParams := []string{"info.f_string", "info.f_child.f_string", "info.f_bool"}
	if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {
		backend.Error(w, http.StatusBadRequest, " found keys that should not appear in query params: %v", duplicates)
		return
	}
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.ComplianceServer.RepeatDataPathResource(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.MarshalToString(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write([]byte(json))
}

// HandleRepeatDataPathTrailingResource translates REST requests/responses on the wire to internal proto messages for RepeatDataPathTrailingResource
//    Generated for HTTP binding pattern: "/v1beta1/repeat/{info.f_string=first/*}/{info.f_child.f_string=second/**}:pathtrailingresource"
func (backend *RESTBackend) HandleRepeatDataPathTrailingResource(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/repeat/{info.f_string=first/*}/{info.f_child.f_string=second/**}:pathtrailingresource': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 2, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 2 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 2, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.RepeatRequest{}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	queryParams := map[string][]string(r.URL.Query())
	excludedQueryParams := []string{"info.f_string", "info.f_child.f_string"}
	if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {
		backend.Error(w, http.StatusBadRequest, " found keys that should not appear in query params: %v", duplicates)
		return
	}
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.ComplianceServer.RepeatDataPathTrailingResource(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.MarshalToString(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write([]byte(json))
}
