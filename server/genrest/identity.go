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
// for service #2: "Identity" (.google.showcase.v1beta1.Identity).

package genrest

import (
	"bytes"
	"context"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
	gmux "github.com/gorilla/mux"
	"io"
	"net/http"
)

// HandleCreateUser translates REST requests/responses on the wire to internal proto messages for CreateUser
//
//	Generated for HTTP binding pattern: POST "/v1beta1/users"
func (backend *RESTBackend) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/users': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &genprotopb.CreateUserRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes := make([]byte, r.ContentLength)
	if _, err := bodyReader.Read(rBytes); err != nil && err != io.EOF {
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/users")
	response, err := backend.IdentityServer.CreateUser(ctx, request)
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

// HandleGetUser translates REST requests/responses on the wire to internal proto messages for GetUser
//
//	Generated for HTTP binding pattern: GET "/v1beta1/{name=users/*}"
func (backend *RESTBackend) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=users/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &genprotopb.GetUserRequest{}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{name=users/*}")
	response, err := backend.IdentityServer.GetUser(ctx, request)
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

// HandleUpdateUser translates REST requests/responses on the wire to internal proto messages for UpdateUser
//
//	Generated for HTTP binding pattern: PATCH "/v1beta1/{user.name=users/*}"
func (backend *RESTBackend) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{user.name=users/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &genprotopb.UpdateUserRequest{}
	// Intentional: Field values in the URL path override those set in the body.
	var bodyField genprotopb.User
	var jsonReader bytes.Buffer
	bodyReader := io.TeeReader(r.Body, &jsonReader)
	rBytes := make([]byte, r.ContentLength)
	if _, err := bodyReader.Read(rBytes); err != nil && err != io.EOF {
		backend.Error(w, http.StatusBadRequest, "error reading body content: %s", err)
		return
	}

	if err := resttools.FromJSON().Unmarshal(rBytes, &bodyField); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading body into request field 'user': %s", err)
		return
	}

	if err := resttools.CheckRequestFormat(&jsonReader, r, request.ProtoReflect()); err != nil {
		backend.Error(w, http.StatusBadRequest, "REST request failed format check: %s", err)
		return
	}
	request.User = &bodyField

	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	// TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both
	excludedQueryParams := []string{"user", "user.name"}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{user.name=users/*}")
	response, err := backend.IdentityServer.UpdateUser(ctx, request)
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

// HandleDeleteUser translates REST requests/responses on the wire to internal proto messages for DeleteUser
//
//	Generated for HTTP binding pattern: DELETE "/v1beta1/{name=users/*}"
func (backend *RESTBackend) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=users/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &genprotopb.DeleteUserRequest{}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/{name=users/*}")
	response, err := backend.IdentityServer.DeleteUser(ctx, request)
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

// HandleListUsers translates REST requests/responses on the wire to internal proto messages for ListUsers
//
//	Generated for HTTP binding pattern: GET "/v1beta1/users"
func (backend *RESTBackend) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/users': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	systemParameters, queryParams, err := resttools.GetSystemParameters(r)
	if err != nil {
		backend.Error(w, http.StatusBadRequest, "error in query string: %s", err)
		return
	}

	request := &genprotopb.ListUsersRequest{}
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

	ctx := context.WithValue(r.Context(), resttools.BindingURIKey, "/v1beta1/users")
	response, err := backend.IdentityServer.ListUsers(ctx, request)
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
