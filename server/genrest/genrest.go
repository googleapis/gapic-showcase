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
	router.HandleFunc("/v1beta1/echo:echo", rest.HandleEcho)
	router.HandleFunc("/v1beta1/echo:expand", rest.HandleExpand)
	router.HandleFunc("/v1beta1/echo:collect", rest.HandleCollect)
	router.HandleFunc("/v1beta1/echo:pagedExpand", rest.HandlePagedExpand)
	router.HandleFunc("/v1beta1/echo:wait", rest.HandleWait)
	router.HandleFunc("/v1beta1/echo:block", rest.HandleBlock)
	router.HandleFunc("/v1beta1/users", rest.HandleCreateUser)
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z_%\\-]+}", rest.HandleGetUser)
	router.HandleFunc("/v1beta1/{user.name:users/[a-zA-Z_%\\-]+}", rest.HandleUpdateUser)
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z_%\\-]+}", rest.HandleDeleteUser)
	router.HandleFunc("/v1beta1/users", rest.HandleListUsers)
	router.HandleFunc("/v1beta1/rooms", rest.HandleCreateRoom)
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z_%\\-]+}", rest.HandleGetRoom)
	router.HandleFunc("/v1beta1/{room.name:rooms/[a-zA-Z_%\\-]+}", rest.HandleUpdateRoom)
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z_%\\-]+}", rest.HandleDeleteRoom)
	router.HandleFunc("/v1beta1/rooms", rest.HandleListRooms)
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z_%\\-]+}/blurbs", rest.HandleCreateBlurb)
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z_%\\-]+/profile}/blurbs", rest.HandleCreateBlurb_1)
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z_%\\-]+/blurbs/[a-zA-Z_%\\-]+}", rest.HandleGetBlurb)
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z_%\\-]+/profile/blurbs/[a-zA-Z_%\\-]+}", rest.HandleGetBlurb_1)
	router.HandleFunc("/v1beta1/{blurb.name:rooms/[a-zA-Z_%\\-]+/blurbs/[a-zA-Z_%\\-]+}", rest.HandleUpdateBlurb)
	router.HandleFunc("/v1beta1/{blurb.name:users/[a-zA-Z_%\\-]+/profile/blurbs/[a-zA-Z_%\\-]+}", rest.HandleUpdateBlurb_1)
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z_%\\-]+/blurbs/[a-zA-Z_%\\-]+}", rest.HandleDeleteBlurb)
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z_%\\-]+/profile/blurbs/[a-zA-Z_%\\-]+}", rest.HandleDeleteBlurb_1)
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z_%\\-]+}/blurbs", rest.HandleListBlurbs)
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z_%\\-]+/profile}/blurbs", rest.HandleListBlurbs_1)
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z_%\\-]+}/blurbs:search", rest.HandleSearchBlurbs)
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z_%\\-]+/profile}/blurbs:search", rest.HandleSearchBlurbs_1)
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z_%\\-]+}/blurbs:stream", rest.HandleStreamBlurbs)
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z_%\\-]+/profile}/blurbs:stream", rest.HandleStreamBlurbs_1)
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z_%\\-]+}/blurbs:send", rest.HandleSendBlurbs)
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z_%\\-]+/profile}/blurbs:send", rest.HandleSendBlurbs_1)
	router.HandleFunc("/v1beta1/sequences", rest.HandleCreateSequence)
	router.HandleFunc("/v1beta1/{name:sequences/[a-zA-Z_%\\-]+/sequenceReport}", rest.HandleGetSequenceReport)
	router.HandleFunc("/v1beta1/{name:sequences/[a-zA-Z_%\\-]+}", rest.HandleAttemptSequence)
	router.HandleFunc("/v1beta1/sessions", rest.HandleCreateSession)
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z_%\\-]+}", rest.HandleGetSession)
	router.HandleFunc("/v1beta1/sessions", rest.HandleListSessions)
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z_%\\-]+}", rest.HandleDeleteSession)
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z_%\\-]+}:report", rest.HandleReportSession)
	router.HandleFunc("/v1beta1/{parent:sessions/[a-zA-Z_%\\-]+}/tests", rest.HandleListTests)
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z_%\\-]+/tests/[a-zA-Z_%\\-]+}", rest.HandleDeleteTest)
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z_%\\-]+/tests/[a-zA-Z_%\\-]+}:check", rest.HandleVerifyTest)
}
