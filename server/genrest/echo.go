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
// for service #1: "Echo" (.google.showcase.v1beta1.Echo).

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

// HandleEcho translates REST requests/responses on the wire to internal proto messages for Echo
//    Generated for HTTP binding pattern: "/v1beta1/echo:echo"
func (backend *RESTBackend) HandleEcho(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:echo': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.EchoRequest{}
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

	if queryParams := r.URL.Query(); len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.EchoServer.Echo(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}

// HandleExpand translates REST requests/responses on the wire to internal proto messages for Expand
//    Generated for HTTP binding pattern: "/v1beta1/echo:expand"
func (backend *RESTBackend) HandleExpand(w http.ResponseWriter, r *http.Request) {
	backend.Error(w, http.StatusNotImplemented, "streaming methods not implemented yet (request matched '/v1beta1/echo:expand': %q)", r.URL)
}

// HandleCollect translates REST requests/responses on the wire to internal proto messages for Collect
//    Generated for HTTP binding pattern: "/v1beta1/echo:collect"
func (backend *RESTBackend) HandleCollect(w http.ResponseWriter, r *http.Request) {
	backend.Error(w, http.StatusNotImplemented, "streaming methods not implemented yet (request matched '/v1beta1/echo:collect': %q)", r.URL)
}

// HandlePagedExpand translates REST requests/responses on the wire to internal proto messages for PagedExpand
//    Generated for HTTP binding pattern: "/v1beta1/echo:pagedExpand"
func (backend *RESTBackend) HandlePagedExpand(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:pagedExpand': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.PagedExpandRequest{}
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

	if queryParams := r.URL.Query(); len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.EchoServer.PagedExpand(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}

// HandlePagedExpandLegacy translates REST requests/responses on the wire to internal proto messages for PagedExpandLegacy
//    Generated for HTTP binding pattern: "/v1beta1/echo:pagedExpandLegacy"
func (backend *RESTBackend) HandlePagedExpandLegacy(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:pagedExpandLegacy': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.PagedExpandLegacyRequest{}
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

	if queryParams := r.URL.Query(); len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.EchoServer.PagedExpandLegacy(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}

// HandleWait translates REST requests/responses on the wire to internal proto messages for Wait
//    Generated for HTTP binding pattern: "/v1beta1/echo:wait"
func (backend *RESTBackend) HandleWait(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:wait': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.WaitRequest{}
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

	if queryParams := r.URL.Query(); len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.EchoServer.Wait(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}

// HandleBlock translates REST requests/responses on the wire to internal proto messages for Block
//    Generated for HTTP binding pattern: "/v1beta1/echo:block"
func (backend *RESTBackend) HandleBlock(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:block': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)
		return
	}

	request := &genprotopb.BlockRequest{}
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

	if queryParams := r.URL.Query(); len(queryParams) > 0 {
		backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %v", queryParams)
		return
	}
	if err := resttools.PopulateSingularFields(request, urlPathParams); err != nil {
		backend.Error(w, http.StatusBadRequest, "error reading URL path params: %s", err)
		return
	}

	marshaler := resttools.ToJSON()
	requestJSON, _ := marshaler.Marshal(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.EchoServer.Block(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	w.Write(json)
}
