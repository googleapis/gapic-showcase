// Copyright 2024 Google LLC
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
// for service #8: "Operations" (.google.longrunning.Operations).

package genrest

import (
	longrunningpbpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"context"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
	gmux "github.com/gorilla/mux"
	"net/http"
)

// HandleListOperations translates REST requests/responses on the wire to internal proto messages for ListOperations
//
//	Generated for HTTP binding pattern: GET "/v1beta1/operations"
func (backend *RESTBackend) HandleListOperations(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/operations': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)
	backend.StdLog.Printf("  urlRequestHeaders:\n%s", resttools.PrettyPrintHeaders(r, "    "))

	resttools.IncludeRequestHeadersInResponse(w, r)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &longrunningpbpb.ListOperationsRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/operations")
	response, err := backend.OperationsServer.ListOperations(ctx, request)
	if err != nil {
		backend.ReportGRPCError(w, err)
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}

// HandleGetOperation translates REST requests/responses on the wire to internal proto messages for GetOperation
//
//	Generated for HTTP binding pattern: GET "/v1beta1/{name=operations/**}"
func (backend *RESTBackend) HandleGetOperation(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=operations/**}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)
	backend.StdLog.Printf("  urlRequestHeaders:\n%s", resttools.PrettyPrintHeaders(r, "    "))

	resttools.IncludeRequestHeadersInResponse(w, r)

	if numUrlPathParams != 1 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &longrunningpbpb.GetOperationRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"name"}
	if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {
		backend.Error(w, http.StatusBadRequest, "(QueryParamsInvalidFieldError) found keys that should not appear in query params: %v", duplicates)
		return
	}
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{name=operations/**}")
	response, err := backend.OperationsServer.GetOperation(ctx, request)
	if err != nil {
		backend.ReportGRPCError(w, err)
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}

// HandleDeleteOperation translates REST requests/responses on the wire to internal proto messages for DeleteOperation
//
//	Generated for HTTP binding pattern: DELETE "/v1beta1/{name=operations/**}"
func (backend *RESTBackend) HandleDeleteOperation(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=operations/**}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)
	backend.StdLog.Printf("  urlRequestHeaders:\n%s", resttools.PrettyPrintHeaders(r, "    "))

	resttools.IncludeRequestHeadersInResponse(w, r)

	if numUrlPathParams != 1 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &longrunningpbpb.DeleteOperationRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"name"}
	if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {
		backend.Error(w, http.StatusBadRequest, "(QueryParamsInvalidFieldError) found keys that should not appear in query params: %v", duplicates)
		return
	}
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{name=operations/**}")
	response, err := backend.OperationsServer.DeleteOperation(ctx, request)
	if err != nil {
		backend.ReportGRPCError(w, err)
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}

// HandleCancelOperation translates REST requests/responses on the wire to internal proto messages for CancelOperation
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{name=operations/**}:cancel"
func (backend *RESTBackend) HandleCancelOperation(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=operations/**}:cancel': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)
	backend.StdLog.Printf("  urlRequestHeaders:\n%s", resttools.PrettyPrintHeaders(r, "    "))

	resttools.IncludeRequestHeadersInResponse(w, r)

	if numUrlPathParams != 1 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &longrunningpbpb.CancelOperationRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"name"}
	if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {
		backend.Error(w, http.StatusBadRequest, "(QueryParamsInvalidFieldError) found keys that should not appear in query params: %v", duplicates)
		return
	}
	if err := resttools.PopulateFields(request, queryParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading query params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{name=operations/**}:cancel")
	response, err := backend.OperationsServer.CancelOperation(ctx, request)
	if err != nil {
		backend.ReportGRPCError(w, err)
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}
