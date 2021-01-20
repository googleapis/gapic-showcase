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
	"context"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gmux "github.com/gorilla/mux"

	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
)

// HandleEcho translates REST requests/responses on the wire to internal proto messages for Echo
//    Generated for HTTP binding pattern: /v1beta1/echo:echo
//         This matches URIs of the form: /v1beta1/echo:echo
func (backend *RESTBackend) HandleEcho(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:echo': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.EchoRequest{}
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

	response, err := backend.EchoServer.Echo(context.Background(), request)
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

// HandleExpand translates REST requests/responses on the wire to internal proto messages for Expand
//    Generated for HTTP binding pattern: /v1beta1/echo:expand
//         This matches URIs of the form: /v1beta1/echo:expand
func (backend *RESTBackend) HandleExpand(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/echo:expand': %q", r.URL)
	w.Write([]byte("ERROR: not implementing streaming methods yet"))
}

// HandleCollect translates REST requests/responses on the wire to internal proto messages for Collect
//    Generated for HTTP binding pattern: /v1beta1/echo:collect
//         This matches URIs of the form: /v1beta1/echo:collect
func (backend *RESTBackend) HandleCollect(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/echo:collect': %q", r.URL)
	w.Write([]byte("ERROR: not implementing streaming methods yet"))
}

// HandlePagedExpand translates REST requests/responses on the wire to internal proto messages for PagedExpand
//    Generated for HTTP binding pattern: /v1beta1/echo:pagedExpand
//         This matches URIs of the form: /v1beta1/echo:pagedExpand
func (backend *RESTBackend) HandlePagedExpand(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:pagedExpand': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.PagedExpandRequest{}
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

	response, err := backend.EchoServer.PagedExpand(context.Background(), request)
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

// HandleWait translates REST requests/responses on the wire to internal proto messages for Wait
//    Generated for HTTP binding pattern: /v1beta1/echo:wait
//         This matches URIs of the form: /v1beta1/echo:wait
func (backend *RESTBackend) HandleWait(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:wait': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.WaitRequest{}
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

	response, err := backend.EchoServer.Wait(context.Background(), request)
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

// HandleBlock translates REST requests/responses on the wire to internal proto messages for Block
//    Generated for HTTP binding pattern: /v1beta1/echo:block
//         This matches URIs of the form: /v1beta1/echo:block
func (backend *RESTBackend) HandleBlock(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/echo:block': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.BlockRequest{}
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

	response, err := backend.EchoServer.Block(context.Background(), request)
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
