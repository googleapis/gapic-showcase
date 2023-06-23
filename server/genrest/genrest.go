// Copyright 2023 Google LLC
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
	"fmt"
	"net/http"

	"github.com/googleapis/gapic-showcase/server/services"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"

	gmux "github.com/gorilla/mux"
	"google.golang.org/grpc/status"
)

type RESTBackend services.Backend

func RegisterHandlers(router *gmux.Router, backend *services.Backend) {
	rest := (*RESTBackend)(backend)
	router.HandleFunc("/v1beta1/repeat:body", rest.HandleRepeatDataBody).Methods("POST")
	router.HandleFunc("/v1beta1/repeat:bodyinfo", rest.HandleRepeatDataBodyInfo).Methods("POST")
	router.HandleFunc("/v1beta1/repeat:query", rest.HandleRepeatDataQuery).Methods("GET")
	router.HandleFunc("/v1beta1/repeat/{info.fString:[a-zA-Z0-9_\\-]+}/{info.fInt32:[a-zA-Z0-9_\\-]+}/{info.fDouble:[a-zA-Z0-9_\\-]+}/{info.fBool:[a-zA-Z0-9_\\-]+}/{info.fKingdom:[a-zA-Z0-9_\\-]+}:simplepath", rest.HandleRepeatDataSimplePath).Methods("GET")
	router.HandleFunc("/v1beta1/repeat/{info.fString:first/[a-zA-Z0-9_\\-]+}/{info.fChild.fString:second/[a-zA-Z0-9_\\-]+}/bool/{info.fBool:[a-zA-Z0-9_\\-]+}:pathresource", rest.HandleRepeatDataPathResource).Methods("GET")
	router.HandleFunc("/v1beta1/repeat/{info.fChild.fString:first/[a-zA-Z0-9_\\-]+}/{info.fString:second/[a-zA-Z0-9_\\-]+}/bool/{info.fBool:[a-zA-Z0-9_\\-]+}:childfirstpathresource", rest.HandleRepeatDataPathResource_1).Methods("GET")
	router.HandleFunc("/v1beta1/repeat/{info.fString:first/[a-zA-Z0-9_\\-]+}/{info.fChild.fString:second/[a-zA-Z0-9_\\-\\/]+}:pathtrailingresource", rest.HandleRepeatDataPathTrailingResource).Methods("GET")
	router.HandleFunc("/v1beta1/repeat:bodyput", rest.HandleRepeatDataBodyPut).Methods("PUT")
	router.HandleFunc("/v1beta1/repeat:bodypatch", rest.HandleRepeatDataBodyPatch).Methods("PATCH")
	router.HandleFunc("/v1beta1/repeat:bodypatch", rest.HandleRepeatDataBodyPatch).HeadersRegexp("X-HTTP-Method-Override", "^PATCH$").Methods("POST")
	router.HandleFunc("/v1beta1/compliance/enum", rest.HandleGetEnum).Methods("GET")
	router.HandleFunc("/v1beta1/compliance/enum", rest.HandleVerifyEnum).Methods("POST")
	router.HandleFunc("/v1beta1/echo:echo", rest.HandleEcho).Methods("POST")
	router.HandleFunc("/v1beta1/echo:expand", rest.HandleExpand).Methods("POST")
	router.HandleFunc("/v1beta1/echo:collect", rest.HandleCollect).Methods("POST")
	router.HandleFunc("/v1beta1/echo:pagedExpand", rest.HandlePagedExpand).Methods("POST")
	router.HandleFunc("/v1beta1/echo:pagedExpandLegacy", rest.HandlePagedExpandLegacy).Methods("POST")
	router.HandleFunc("/v1beta1/echo:pagedExpandLegacyMapped", rest.HandlePagedExpandLegacyMapped).Methods("POST")
	router.HandleFunc("/v1beta1/echo:wait", rest.HandleWait).Methods("POST")
	router.HandleFunc("/v1beta1/echo:block", rest.HandleBlock).Methods("POST")
	router.HandleFunc("/v1beta1/users", rest.HandleCreateUser).Methods("POST")
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z0-9_\\-]+}", rest.HandleGetUser).Methods("GET")
	router.HandleFunc("/v1beta1/{user.name:users/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateUser).Methods("PATCH")
	router.HandleFunc("/v1beta1/{user.name:users/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateUser).HeadersRegexp("X-HTTP-Method-Override", "^PATCH$").Methods("POST")
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z0-9_\\-]+}", rest.HandleDeleteUser).Methods("DELETE")
	router.HandleFunc("/v1beta1/users", rest.HandleListUsers).Methods("GET")
	router.HandleFunc("/v1beta1/rooms", rest.HandleCreateRoom).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z0-9_\\-]+}", rest.HandleGetRoom).Methods("GET")
	router.HandleFunc("/v1beta1/{room.name:rooms/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateRoom).Methods("PATCH")
	router.HandleFunc("/v1beta1/{room.name:rooms/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateRoom).HeadersRegexp("X-HTTP-Method-Override", "^PATCH$").Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z0-9_\\-]+}", rest.HandleDeleteRoom).Methods("DELETE")
	router.HandleFunc("/v1beta1/rooms", rest.HandleListRooms).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z0-9_\\-]+}/blurbs", rest.HandleCreateBlurb).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z0-9_\\-]+/profile}/blurbs", rest.HandleCreateBlurb_1).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z0-9_\\-]+/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleGetBlurb).Methods("GET")
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z0-9_\\-]+/profile/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleGetBlurb_1).Methods("GET")
	router.HandleFunc("/v1beta1/{blurb.name:rooms/[a-zA-Z0-9_\\-]+/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateBlurb).Methods("PATCH")
	router.HandleFunc("/v1beta1/{blurb.name:rooms/[a-zA-Z0-9_\\-]+/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateBlurb).HeadersRegexp("X-HTTP-Method-Override", "^PATCH$").Methods("POST")
	router.HandleFunc("/v1beta1/{blurb.name:users/[a-zA-Z0-9_\\-]+/profile/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateBlurb_1).Methods("PATCH")
	router.HandleFunc("/v1beta1/{blurb.name:users/[a-zA-Z0-9_\\-]+/profile/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleUpdateBlurb_1).HeadersRegexp("X-HTTP-Method-Override", "^PATCH$").Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z0-9_\\-]+/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleDeleteBlurb).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z0-9_\\-]+/profile/blurbs/[a-zA-Z0-9_\\-]+}", rest.HandleDeleteBlurb_1).Methods("DELETE")
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z0-9_\\-]+}/blurbs", rest.HandleListBlurbs).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z0-9_\\-]+/profile}/blurbs", rest.HandleListBlurbs_1).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z0-9_\\-]+}/blurbs:search", rest.HandleSearchBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z0-9_\\-]+/profile}/blurbs:search", rest.HandleSearchBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/[a-zA-Z0-9_\\-]+}/blurbs:stream", rest.HandleStreamBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{name:users/[a-zA-Z0-9_\\-]+/profile}/blurbs:stream", rest.HandleStreamBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:rooms/[a-zA-Z0-9_\\-]+}/blurbs:send", rest.HandleSendBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/[a-zA-Z0-9_\\-]+/profile}/blurbs:send", rest.HandleSendBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/sequences", rest.HandleCreateSequence).Methods("POST")
	router.HandleFunc("/v1beta1/streamingSequences", rest.HandleCreateStreamingSequence).Methods("POST")
	router.HandleFunc("/v1beta1/{name:sequences/[a-zA-Z0-9_\\-]+/sequenceReport}", rest.HandleGetSequenceReport).Methods("GET")
	router.HandleFunc("/v1beta1/{name:streamingSequences/[a-zA-Z0-9_\\-]+/streamingSequenceReport}", rest.HandleGetStreamingSequenceReport).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sequences/[a-zA-Z0-9_\\-]+}", rest.HandleAttemptSequence).Methods("POST")
	router.HandleFunc("/v1beta1/{name:streamingSequences/[a-zA-Z0-9_\\-]+}:stream", rest.HandleAttemptStreamingSequence).Methods("POST")
	router.HandleFunc("/v1beta1/sessions", rest.HandleCreateSession).Methods("POST")
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z0-9_\\-]+}", rest.HandleGetSession).Methods("GET")
	router.HandleFunc("/v1beta1/sessions", rest.HandleListSessions).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z0-9_\\-]+}", rest.HandleDeleteSession).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z0-9_\\-]+}:report", rest.HandleReportSession).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:sessions/[a-zA-Z0-9_\\-]+}/tests", rest.HandleListTests).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z0-9_\\-]+/tests/[a-zA-Z0-9_\\-]+}", rest.HandleDeleteTest).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:sessions/[a-zA-Z0-9_\\-]+/tests/[a-zA-Z0-9_\\-]+}:check", rest.HandleVerifyTest).Methods("POST")
	router.HandleFunc("/v1beta1/{name:projects/[a-zA-Z0-9_\\-]+}/locations", rest.HandleListLocations).Methods("GET")
	router.HandleFunc("/v1beta1/{name:projects/[a-zA-Z0-9_\\-]+/locations/[a-zA-Z0-9_\\-]+}", rest.HandleGetLocation).Methods("GET")
	router.HandleFunc("/v1beta1/{resource:users/[a-zA-Z0-9_\\-]+}:setIamPolicy", rest.HandleSetIamPolicy).Methods("POST")
	router.HandleFunc("/v1beta1/{resource:rooms/[a-zA-Z0-9_\\-]+}:setIamPolicy", rest.HandleSetIamPolicy_1).Methods("POST")
	router.HandleFunc("/v1beta1/{resource:rooms/[a-zA-Z0-9_\\-]+/blurbs/[a-zA-Z0-9_\\-]+}:setIamPolicy", rest.HandleSetIamPolicy_2).Methods("POST")
	router.HandleFunc("/v1beta1/{resource:sequences/[a-zA-Z0-9_\\-]+}:setIamPolicy", rest.HandleSetIamPolicy_3).Methods("POST")
	router.HandleFunc("/v1beta1/{resource:users/[a-zA-Z0-9_\\-]+}:getIamPolicy", rest.HandleGetIamPolicy).Methods("GET")
	router.HandleFunc("/v1beta1/{resource:rooms/[a-zA-Z0-9_\\-]+}:getIamPolicy", rest.HandleGetIamPolicy_1).Methods("GET")
	router.HandleFunc("/v1beta1/{resource:rooms/[a-zA-Z0-9_\\-]+/blurbs/[a-zA-Z0-9_\\-]+}:getIamPolicy", rest.HandleGetIamPolicy_2).Methods("GET")
	router.HandleFunc("/v1beta1/{resource:sequences/[a-zA-Z0-9_\\-]+}:getIamPolicy", rest.HandleGetIamPolicy_3).Methods("GET")
	router.HandleFunc("/v1beta1/{resource:users/[a-zA-Z0-9_\\-]+}:testIamPermissions", rest.HandleTestIamPermissions).Methods("POST")
	router.HandleFunc("/v1beta1/{resource:rooms/[a-zA-Z0-9_\\-]+}:testIamPermissions", rest.HandleTestIamPermissions_1).Methods("POST")
	router.HandleFunc("/v1beta1/{resource:rooms/[a-zA-Z0-9_\\-]+/blurbs/[a-zA-Z0-9_\\-]+}:testIamPermissions", rest.HandleTestIamPermissions_2).Methods("POST")
	router.HandleFunc("/v1beta1/{resource:sequences/[a-zA-Z0-9_\\-]+}:testIamPermissions", rest.HandleTestIamPermissions_3).Methods("POST")
	router.HandleFunc("/v1beta1/operations", rest.HandleListOperations).Methods("GET")
	router.HandleFunc("/v1beta1/{name:operations/[a-zA-Z0-9_\\-\\/]+}", rest.HandleGetOperation).Methods("GET")
	router.HandleFunc("/v1beta1/{name:operations/[a-zA-Z0-9_\\-\\/]+}", rest.HandleDeleteOperation).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:operations/[a-zA-Z0-9_\\-\\/]+}:cancel", rest.HandleCancelOperation).Methods("POST")
	router.PathPrefix("/").HandlerFunc(rest.catchAllHandler)
}

func (backend *RESTBackend) catchAllHandler(w http.ResponseWriter, r *http.Request) {
	backend.Error(w, http.StatusBadRequest, "unrecognized request: %s %q", r.Method, r.URL)
}

func (backend *RESTBackend) Error(w http.ResponseWriter, httpStatus int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	backend.ErrLog.Print(message)
	resttools.ErrorResponse(w, httpStatus, message)
}
func (backend *RESTBackend) ReportGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	backend.ErrLog.Print(st.Message())
	code := resttools.GRPCToHTTP(st.Code())
	resttools.ErrorResponse(w, code, st.Message(), st.Details()...)
}
