// Copyright 2026 Google LLC
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
// for service #7: "IAMPolicy" (.google.iam.v1.IAMPolicy).

package genrest

import (
	"bytes"
	iampbpb "cloud.google.com/go/iam/apiv1/iampb"
	"context"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
	gmux "github.com/gorilla/mux"
	"io"
	"net/http"
)

// HandleSetIamPolicy translates REST requests/responses on the wire to internal proto messages for SetIamPolicy
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=users/*}:setIamPolicy"
func (backend *RESTBackend) HandleSetIamPolicy(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=users/*}:setIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.SetIamPolicyRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=users/*}:setIamPolicy")
	response, err := backend.IAMPolicyServer.SetIamPolicy(ctx, request)
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

// HandleSetIamPolicy_1 translates REST requests/responses on the wire to internal proto messages for SetIamPolicy
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=rooms/*}:setIamPolicy"
func (backend *RESTBackend) HandleSetIamPolicy_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=rooms/*}:setIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.SetIamPolicyRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=rooms/*}:setIamPolicy")
	response, err := backend.IAMPolicyServer.SetIamPolicy(ctx, request)
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

// HandleSetIamPolicy_2 translates REST requests/responses on the wire to internal proto messages for SetIamPolicy
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=rooms/*/blurbs/*}:setIamPolicy"
func (backend *RESTBackend) HandleSetIamPolicy_2(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=rooms/*/blurbs/*}:setIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.SetIamPolicyRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=rooms/*/blurbs/*}:setIamPolicy")
	response, err := backend.IAMPolicyServer.SetIamPolicy(ctx, request)
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

// HandleSetIamPolicy_3 translates REST requests/responses on the wire to internal proto messages for SetIamPolicy
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=sequences/*}:setIamPolicy"
func (backend *RESTBackend) HandleSetIamPolicy_3(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=sequences/*}:setIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.SetIamPolicyRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=sequences/*}:setIamPolicy")
	response, err := backend.IAMPolicyServer.SetIamPolicy(ctx, request)
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

// HandleGetIamPolicy translates REST requests/responses on the wire to internal proto messages for GetIamPolicy
//
//	Generated for HTTP binding pattern: GET "/v1beta1/{resource=users/*}:getIamPolicy"
func (backend *RESTBackend) HandleGetIamPolicy(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=users/*}:getIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.GetIamPolicyRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"resource"}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=users/*}:getIamPolicy")
	response, err := backend.IAMPolicyServer.GetIamPolicy(ctx, request)
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

// HandleGetIamPolicy_1 translates REST requests/responses on the wire to internal proto messages for GetIamPolicy
//
//	Generated for HTTP binding pattern: GET "/v1beta1/{resource=rooms/*}:getIamPolicy"
func (backend *RESTBackend) HandleGetIamPolicy_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=rooms/*}:getIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.GetIamPolicyRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"resource"}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=rooms/*}:getIamPolicy")
	response, err := backend.IAMPolicyServer.GetIamPolicy(ctx, request)
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

// HandleGetIamPolicy_2 translates REST requests/responses on the wire to internal proto messages for GetIamPolicy
//
//	Generated for HTTP binding pattern: GET "/v1beta1/{resource=rooms/*/blurbs/*}:getIamPolicy"
func (backend *RESTBackend) HandleGetIamPolicy_2(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=rooms/*/blurbs/*}:getIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.GetIamPolicyRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"resource"}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=rooms/*/blurbs/*}:getIamPolicy")
	response, err := backend.IAMPolicyServer.GetIamPolicy(ctx, request)
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

// HandleGetIamPolicy_3 translates REST requests/responses on the wire to internal proto messages for GetIamPolicy
//
//	Generated for HTTP binding pattern: GET "/v1beta1/{resource=sequences/*}:getIamPolicy"
func (backend *RESTBackend) HandleGetIamPolicy_3(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=sequences/*}:getIamPolicy': %q", r.Method, r.URL)
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

	request := &iampbpb.GetIamPolicyRequest{}
	if err := resttools.CheckRequestFormat(nil, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"resource"}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=sequences/*}:getIamPolicy")
	response, err := backend.IAMPolicyServer.GetIamPolicy(ctx, request)
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

// HandleTestIamPermissions translates REST requests/responses on the wire to internal proto messages for TestIamPermissions
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=users/*}:testIamPermissions"
func (backend *RESTBackend) HandleTestIamPermissions(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=users/*}:testIamPermissions': %q", r.Method, r.URL)
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

	request := &iampbpb.TestIamPermissionsRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=users/*}:testIamPermissions")
	response, err := backend.IAMPolicyServer.TestIamPermissions(ctx, request)
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

// HandleTestIamPermissions_1 translates REST requests/responses on the wire to internal proto messages for TestIamPermissions
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=rooms/*}:testIamPermissions"
func (backend *RESTBackend) HandleTestIamPermissions_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=rooms/*}:testIamPermissions': %q", r.Method, r.URL)
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

	request := &iampbpb.TestIamPermissionsRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=rooms/*}:testIamPermissions")
	response, err := backend.IAMPolicyServer.TestIamPermissions(ctx, request)
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

// HandleTestIamPermissions_2 translates REST requests/responses on the wire to internal proto messages for TestIamPermissions
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=rooms/*/blurbs/*}:testIamPermissions"
func (backend *RESTBackend) HandleTestIamPermissions_2(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=rooms/*/blurbs/*}:testIamPermissions': %q", r.Method, r.URL)
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

	request := &iampbpb.TestIamPermissionsRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=rooms/*/blurbs/*}:testIamPermissions")
	response, err := backend.IAMPolicyServer.TestIamPermissions(ctx, request)
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

// HandleTestIamPermissions_3 translates REST requests/responses on the wire to internal proto messages for TestIamPermissions
//
//	Generated for HTTP binding pattern: POST "/v1beta1/{resource=sequences/*}:testIamPermissions"
func (backend *RESTBackend) HandleTestIamPermissions_3(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{resource=sequences/*}:testIamPermissions': %q", r.Method, r.URL)
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

	request := &iampbpb.TestIamPermissionsRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, request); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body params '*': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}

	if len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	marshaler.UseEnumNumbers = systemParameters.EnumEncodingAsInt
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{resource=sequences/*}:testIamPermissions")
	response, err := backend.IAMPolicyServer.TestIamPermissions(ctx, request)
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
