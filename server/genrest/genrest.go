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

// DO NOT EDIT. This is an auto-generated file registering the REST handlers.
// for the various Showcase services.

package genrest

import (
	"github.com/googleapis/gapic-showcase/server/services"

	gmux "github.com/gorilla/mux"
)

type RESTBackend services.Backend

func RegisterHandlers(router *gmux.Router, backend *services.Backend) {
	rest := (*RESTBackend)(backend)
	router.HandleFunc("/v1beta1/echo:echo", rest.HandleEcho).Methods("POST")
	router.HandleFunc("/v1beta1/echo:expand", rest.HandleExpand).Methods("POST")
	router.HandleFunc("/v1beta1/echo:collect", rest.HandleCollect).Methods("POST")
	router.HandleFunc("/v1beta1/echo:pagedExpand", rest.HandlePagedExpand).Methods("POST")
	router.HandleFunc("/v1beta1/echo:wait", rest.HandleWait).Methods("POST")
	router.HandleFunc("/v1beta1/echo:block", rest.HandleBlock).Methods("POST")
	router.HandleFunc("/v1beta1/users", rest.HandleCreateUser).Methods("POST")
	router.HandleFunc("/v1beta1/{name:users/[0-9a-zA-Z_%\\-]+}", rest.HandleGetUser).Methods("GET")
	router.HandleFunc("/v1beta1/{user.name:users/[0-9a-zA-Z_%\\-]+}", rest.HandleUpdateUser).Methods("PATCH")
	router.HandleFunc("/v1beta1/{name:users/[0-9a-zA-Z_%\\-]+}", rest.HandleDeleteUser).Methods("DELETE")
	router.HandleFunc("/v1beta1/users", rest.HandleListUsers).Methods("GET")
	router.HandleFunc("/v1beta1/rooms", rest.HandleCreateRoom).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[0-9a-zA-Z_%\\-]+}", rest.HandleGetRoom).Methods("GET")
	router.HandleFunc("/v1beta1/{room.name:rooms/[0-9a-zA-Z_%\\-]+}", rest.HandleUpdateRoom).Methods("PATCH")
	router.HandleFunc("/v1beta1/{name:rooms/[0-9a-zA-Z_%\\-]+}", rest.HandleDeleteRoom).Methods("DELETE")
	router.HandleFunc("/v1beta1/rooms", rest.HandleListRooms).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:rooms/[0-9a-zA-Z_%\\-]+}/blurbs", rest.HandleCreateBlurb).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/[0-9a-zA-Z_%\\-]+/profile}/blurbs", rest.HandleCreateBlurb_1).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[0-9a-zA-Z_%\\-]+/blurbs/[0-9a-zA-Z_%\\-]+}", rest.HandleGetBlurb).Methods("GET")
	router.HandleFunc("/v1beta1/{name:users/[0-9a-zA-Z_%\\-]+/profile/blurbs/[0-9a-zA-Z_%\\-]+}", rest.HandleGetBlurb_1).Methods("GET")
	router.HandleFunc("/v1beta1/{blurb.name:rooms/[0-9a-zA-Z_%\\-]+/blurbs/[0-9a-zA-Z_%\\-]+}", rest.HandleUpdateBlurb).Methods("PATCH")
	router.HandleFunc("/v1beta1/{blurb.name:users/[0-9a-zA-Z_%\\-]+/profile/blurbs/[0-9a-zA-Z_%\\-]+}", rest.HandleUpdateBlurb_1).Methods("PATCH")
	router.HandleFunc("/v1beta1/{name:rooms/[0-9a-zA-Z_%\\-]+/blurbs/[0-9a-zA-Z_%\\-]+}", rest.HandleDeleteBlurb).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:users/[0-9a-zA-Z_%\\-]+/profile/blurbs/[0-9a-zA-Z_%\\-]+}", rest.HandleDeleteBlurb_1).Methods("DELETE")
	router.HandleFunc("/v1beta1/{parent:rooms/[0-9a-zA-Z_%\\-]+}/blurbs", rest.HandleListBlurbs).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:users/[0-9a-zA-Z_%\\-]+/profile}/blurbs", rest.HandleListBlurbs_1).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:rooms/[0-9a-zA-Z_%\\-]+}/blurbs:search", rest.HandleSearchBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/[0-9a-zA-Z_%\\-]+/profile}/blurbs:search", rest.HandleSearchBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[0-9a-zA-Z_%\\-]+}/blurbs:stream", rest.HandleStreamBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{name:users/[0-9a-zA-Z_%\\-]+/profile}/blurbs:stream", rest.HandleStreamBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:rooms/[0-9a-zA-Z_%\\-]+}/blurbs:send", rest.HandleSendBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/[0-9a-zA-Z_%\\-]+/profile}/blurbs:send", rest.HandleSendBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/sequences", rest.HandleCreateSequence).Methods("POST")
	router.HandleFunc("/v1beta1/{name:sequences/[0-9a-zA-Z_%\\-]+/sequenceReport}", rest.HandleGetSequenceReport).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sequences/[0-9a-zA-Z_%\\-]+}", rest.HandleAttemptSequence).Methods("POST")
	router.HandleFunc("/v1beta1/sessions", rest.HandleCreateSession).Methods("POST")
	router.HandleFunc("/v1beta1/{name:sessions/[0-9a-zA-Z_%\\-]+}", rest.HandleGetSession).Methods("GET")
	router.HandleFunc("/v1beta1/sessions", rest.HandleListSessions).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sessions/[0-9a-zA-Z_%\\-]+}", rest.HandleDeleteSession).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:sessions/[0-9a-zA-Z_%\\-]+}:report", rest.HandleReportSession).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:sessions/[0-9a-zA-Z_%\\-]+}/tests", rest.HandleListTests).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sessions/[0-9a-zA-Z_%\\-]+/tests/[0-9a-zA-Z_%\\-]+}", rest.HandleDeleteTest).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:sessions/[0-9a-zA-Z_%\\-]+/tests/[0-9a-zA-Z_%\\-]+}:check", rest.HandleVerifyTest).Methods("POST")
}
