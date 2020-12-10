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
// for service #2: "Messaging" (.google.showcase.v1beta1.Messaging).

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

// HandleCreateRoom translates REST requests/responses on the wire to internal proto messages for CreateRoom
//    Generated for HTTP binding pattern: /v1beta1/rooms
//         This matches URIs of the form: /v1beta1/rooms
func (backend *RESTBackend) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/rooms': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.CreateRoomRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.CreateRoom(context.Background(), request)
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

// HandleGetRoom translates REST requests/responses on the wire to internal proto messages for GetRoom
//    Generated for HTTP binding pattern: /v1beta1/{name=rooms/*}
//         This matches URIs of the form: /v1beta1/{name:rooms/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=rooms/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.GetRoomRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.GetRoom(context.Background(), request)
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

// HandleUpdateRoom translates REST requests/responses on the wire to internal proto messages for UpdateRoom
//    Generated for HTTP binding pattern: /v1beta1/{room.name=rooms/*}
//         This matches URIs of the form: /v1beta1/{room.name:rooms/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleUpdateRoom(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{room.name=rooms/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.UpdateRoomRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.UpdateRoom(context.Background(), request)
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

// HandleDeleteRoom translates REST requests/responses on the wire to internal proto messages for DeleteRoom
//    Generated for HTTP binding pattern: /v1beta1/{name=rooms/*}
//         This matches URIs of the form: /v1beta1/{name:rooms/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=rooms/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.DeleteRoomRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.DeleteRoom(context.Background(), request)
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

// HandleListRooms translates REST requests/responses on the wire to internal proto messages for ListRooms
//    Generated for HTTP binding pattern: /v1beta1/rooms
//         This matches URIs of the form: /v1beta1/rooms
func (backend *RESTBackend) HandleListRooms(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/rooms': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 0, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 0 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 0, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.ListRoomsRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.ListRooms(context.Background(), request)
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

// HandleCreateBlurb translates REST requests/responses on the wire to internal proto messages for CreateBlurb
//    Generated for HTTP binding pattern: /v1beta1/{parent=rooms/*}/blurbs
//         This matches URIs of the form: /v1beta1/{parent:rooms/[a-zA-Z_%\-]+}/blurbs
func (backend *RESTBackend) HandleCreateBlurb(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{parent=rooms/*}/blurbs': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.CreateBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.CreateBlurb(context.Background(), request)
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

// HandleCreateBlurb_1 translates REST requests/responses on the wire to internal proto messages for CreateBlurb
//    Generated for HTTP binding pattern: /v1beta1/{parent=users/*/profile}/blurbs
//         This matches URIs of the form: /v1beta1/{parent:users/[a-zA-Z_%\-]+/profile}/blurbs
func (backend *RESTBackend) HandleCreateBlurb_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{parent=users/*/profile}/blurbs': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.CreateBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.CreateBlurb(context.Background(), request)
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

// HandleGetBlurb translates REST requests/responses on the wire to internal proto messages for GetBlurb
//    Generated for HTTP binding pattern: /v1beta1/{name=rooms/*/blurbs/*}
//         This matches URIs of the form: /v1beta1/{name:rooms/[a-zA-Z_%\-]+/blurbs/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleGetBlurb(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=rooms/*/blurbs/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.GetBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.GetBlurb(context.Background(), request)
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

// HandleGetBlurb_1 translates REST requests/responses on the wire to internal proto messages for GetBlurb
//    Generated for HTTP binding pattern: /v1beta1/{name=users/*/profile/blurbs/*}
//         This matches URIs of the form: /v1beta1/{name:users/[a-zA-Z_%\-]+/profile/blurbs/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleGetBlurb_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=users/*/profile/blurbs/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.GetBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.GetBlurb(context.Background(), request)
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

// HandleUpdateBlurb translates REST requests/responses on the wire to internal proto messages for UpdateBlurb
//    Generated for HTTP binding pattern: /v1beta1/{blurb.name=rooms/*/blurbs/*}
//         This matches URIs of the form: /v1beta1/{blurb.name:rooms/[a-zA-Z_%\-]+/blurbs/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleUpdateBlurb(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{blurb.name=rooms/*/blurbs/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.UpdateBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.UpdateBlurb(context.Background(), request)
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

// HandleUpdateBlurb_1 translates REST requests/responses on the wire to internal proto messages for UpdateBlurb
//    Generated for HTTP binding pattern: /v1beta1/{blurb.name=users/*/profile/blurbs/*}
//         This matches URIs of the form: /v1beta1/{blurb.name:users/[a-zA-Z_%\-]+/profile/blurbs/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleUpdateBlurb_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{blurb.name=users/*/profile/blurbs/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.UpdateBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.UpdateBlurb(context.Background(), request)
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

// HandleDeleteBlurb translates REST requests/responses on the wire to internal proto messages for DeleteBlurb
//    Generated for HTTP binding pattern: /v1beta1/{name=rooms/*/blurbs/*}
//         This matches URIs of the form: /v1beta1/{name:rooms/[a-zA-Z_%\-]+/blurbs/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteBlurb(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=rooms/*/blurbs/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.DeleteBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.DeleteBlurb(context.Background(), request)
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

// HandleDeleteBlurb_1 translates REST requests/responses on the wire to internal proto messages for DeleteBlurb
//    Generated for HTTP binding pattern: /v1beta1/{name=users/*/profile/blurbs/*}
//         This matches URIs of the form: /v1beta1/{name:users/[a-zA-Z_%\-]+/profile/blurbs/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteBlurb_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{name=users/*/profile/blurbs/*}': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.DeleteBlurbRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.DeleteBlurb(context.Background(), request)
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

// HandleListBlurbs translates REST requests/responses on the wire to internal proto messages for ListBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{parent=rooms/*}/blurbs
//         This matches URIs of the form: /v1beta1/{parent:rooms/[a-zA-Z_%\-]+}/blurbs
func (backend *RESTBackend) HandleListBlurbs(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{parent=rooms/*}/blurbs': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.ListBlurbsRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.ListBlurbs(context.Background(), request)
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

// HandleListBlurbs_1 translates REST requests/responses on the wire to internal proto messages for ListBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{parent=users/*/profile}/blurbs
//         This matches URIs of the form: /v1beta1/{parent:users/[a-zA-Z_%\-]+/profile}/blurbs
func (backend *RESTBackend) HandleListBlurbs_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{parent=users/*/profile}/blurbs': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.ListBlurbsRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.ListBlurbs(context.Background(), request)
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

// HandleSearchBlurbs translates REST requests/responses on the wire to internal proto messages for SearchBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{parent=rooms/*}/blurbs:search
//         This matches URIs of the form: /v1beta1/{parent:rooms/[a-zA-Z_%\-]+}/blurbs:search
func (backend *RESTBackend) HandleSearchBlurbs(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{parent=rooms/*}/blurbs:search': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.SearchBlurbsRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.SearchBlurbs(context.Background(), request)
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

// HandleSearchBlurbs_1 translates REST requests/responses on the wire to internal proto messages for SearchBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{parent=users/*/profile}/blurbs:search
//         This matches URIs of the form: /v1beta1/{parent:users/[a-zA-Z_%\-]+/profile}/blurbs:search
func (backend *RESTBackend) HandleSearchBlurbs_1(w http.ResponseWriter, r *http.Request) {
	urlPathParams := gmux.Vars(r)
	numUrlPathParams := len(urlPathParams)

	backend.StdLog.Printf("Received %s request matching '/v1beta1/{parent=users/*/profile}/blurbs:search': %q", r.Method, r.URL)
	backend.StdLog.Printf("  urlPathParams (expect 1, have %d): %q", numUrlPathParams, urlPathParams)

	if numUrlPathParams != 1 {
		w.Write([]byte(fmt.Sprintf("unexpected number of URL variables: expected 1, have %d: %#v", numUrlPathParams, urlPathParams)))
		return
	}

	request := &genprotopb.SearchBlurbsRequest{}
	// TODO: Populate request with parameters from HTTP request
	if err := resttools.PopulateFields(request, urlPathParams); err != nil {
		backend.StdLog.Printf("  error: %s", err)
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	requestJSON, _ := marshaler.MarshalToString(request)
	backend.StdLog.Printf("  request: %s", requestJSON)

	response, err := backend.MessagingServer.SearchBlurbs(context.Background(), request)
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

// HandleStreamBlurbs translates REST requests/responses on the wire to internal proto messages for StreamBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{name=rooms/*}/blurbs:stream
//         This matches URIs of the form: /v1beta1/{name:rooms/[a-zA-Z_%\-]+}/blurbs:stream
func (backend *RESTBackend) HandleStreamBlurbs(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=rooms/*}/blurbs:stream': %q", r.URL)
	w.Write([]byte("ERROR: not implementing streaming methods yet"))
}

// HandleStreamBlurbs_1 translates REST requests/responses on the wire to internal proto messages for StreamBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{name=users/*/profile}/blurbs:stream
//         This matches URIs of the form: /v1beta1/{name:users/[a-zA-Z_%\-]+/profile}/blurbs:stream
func (backend *RESTBackend) HandleStreamBlurbs_1(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=users/*/profile}/blurbs:stream': %q", r.URL)
	w.Write([]byte("ERROR: not implementing streaming methods yet"))
}

// HandleSendBlurbs translates REST requests/responses on the wire to internal proto messages for SendBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{parent=rooms/*}/blurbs:send
//         This matches URIs of the form: /v1beta1/{parent:rooms/[a-zA-Z_%\-]+}/blurbs:send
func (backend *RESTBackend) HandleSendBlurbs(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{parent=rooms/*}/blurbs:send': %q", r.URL)
	w.Write([]byte("ERROR: not implementing streaming methods yet"))
}

// HandleSendBlurbs_1 translates REST requests/responses on the wire to internal proto messages for SendBlurbs
//    Generated for HTTP binding pattern: /v1beta1/{parent=users/*/profile}/blurbs:send
//         This matches URIs of the form: /v1beta1/{parent:users/[a-zA-Z_%\-]+/profile}/blurbs:send
func (backend *RESTBackend) HandleSendBlurbs_1(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{parent=users/*/profile}/blurbs:send': %q", r.URL)
	w.Write([]byte("ERROR: not implementing streaming methods yet"))
}
