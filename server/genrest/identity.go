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
// for service #1: "Identity" (.google.showcase.v1beta1.Identity).

package genrest

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
)

// HandleCreateUser translates REST requests/responses on the wire to internal proto messages for CreateUser
//    Generated for HTTP binding pattern: /v1beta1/users
//         This matches URIs of the form: /v1beta1/users
func (backend *RESTBackend) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/users': %q", r.URL)

	var request *genprotopb.CreateUserRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.IdentityServer.CreateUser(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleGetUser translates REST requests/responses on the wire to internal proto messages for GetUser
//    Generated for HTTP binding pattern: /v1beta1/{name=users/*}
//         This matches URIs of the form: /v1beta1/{name:users/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=users/*}': %q", r.URL)

	var request *genprotopb.GetUserRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.IdentityServer.GetUser(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleUpdateUser translates REST requests/responses on the wire to internal proto messages for UpdateUser
//    Generated for HTTP binding pattern: /v1beta1/{user.name=users/*}
//         This matches URIs of the form: /v1beta1/{user.name:users/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{user.name=users/*}': %q", r.URL)

	var request *genprotopb.UpdateUserRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.IdentityServer.UpdateUser(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleDeleteUser translates REST requests/responses on the wire to internal proto messages for DeleteUser
//    Generated for HTTP binding pattern: /v1beta1/{name=users/*}
//         This matches URIs of the form: /v1beta1/{name:users/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=users/*}': %q", r.URL)

	var request *genprotopb.DeleteUserRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.IdentityServer.DeleteUser(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}

// HandleListUsers translates REST requests/responses on the wire to internal proto messages for ListUsers
//    Generated for HTTP binding pattern: /v1beta1/users
//         This matches URIs of the form: /v1beta1/users
func (backend *RESTBackend) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/users': %q", r.URL)

	var request *genprotopb.ListUsersRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.IdentityServer.ListUsers(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(response)
	if err != nil {
		// TODO: Properly handle error
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(json))
}
