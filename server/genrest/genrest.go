// Copyright 2022 Google LLC
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
	router.HandleFunc("/v1beta1/repeat/{info.fString:.+}/{info.fInt32:.+}/{info.fDouble:.+}/{info.fBool:.+}/{info.fKingdom:.+}:simplepath", rest.HandleRepeatDataSimplePath).Methods("GET")
	router.HandleFunc("/v1beta1/repeat/{info.fString:first/.+}/{info.fChild.fString:second/.+}/bool/{info.fBool:.+}:pathresource", rest.HandleRepeatDataPathResource).Methods("GET")
	router.HandleFunc("/v1beta1/repeat/{info.fChild.fString:first/.+}/{info.fString:second/.+}/bool/{info.fBool:.+}:childfirstpathresource", rest.HandleRepeatDataPathResource_1).Methods("GET")
	router.HandleFunc("/v1beta1/repeat/{info.fString:first/.+}/{info.fChild.fString:second/.+}:pathtrailingresource", rest.HandleRepeatDataPathTrailingResource).Methods("GET")
	router.HandleFunc("/v1beta1/repeat:bodyput", rest.HandleRepeatDataBodyPut).Methods("PUT")
	router.HandleFunc("/v1beta1/repeat:bodypatch", rest.HandleRepeatDataBodyPatch).Methods("PATCH")
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
	router.HandleFunc("/v1beta1/{name:users/.+}", rest.HandleGetUser).Methods("GET")
	router.HandleFunc("/v1beta1/{user.name:users/.+}", rest.HandleUpdateUser).Methods("PATCH")
	router.HandleFunc("/v1beta1/{name:users/.+}", rest.HandleDeleteUser).Methods("DELETE")
	router.HandleFunc("/v1beta1/users", rest.HandleListUsers).Methods("GET")
	router.HandleFunc("/v1beta1/rooms", rest.HandleCreateRoom).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/.+}", rest.HandleGetRoom).Methods("GET")
	router.HandleFunc("/v1beta1/{room.name:rooms/.+}", rest.HandleUpdateRoom).Methods("PATCH")
	router.HandleFunc("/v1beta1/{name:rooms/.+}", rest.HandleDeleteRoom).Methods("DELETE")
	router.HandleFunc("/v1beta1/rooms", rest.HandleListRooms).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:rooms/.+}/blurbs", rest.HandleCreateBlurb).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/.+/profile}/blurbs", rest.HandleCreateBlurb_1).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/.+/blurbs/.+}", rest.HandleGetBlurb).Methods("GET")
	router.HandleFunc("/v1beta1/{name:users/.+/profile/blurbs/.+}", rest.HandleGetBlurb_1).Methods("GET")
	router.HandleFunc("/v1beta1/{blurb.name:rooms/.+/blurbs/.+}", rest.HandleUpdateBlurb).Methods("PATCH")
	router.HandleFunc("/v1beta1/{blurb.name:users/.+/profile/blurbs/.+}", rest.HandleUpdateBlurb_1).Methods("PATCH")
	router.HandleFunc("/v1beta1/{name:rooms/.+/blurbs/.+}", rest.HandleDeleteBlurb).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:users/.+/profile/blurbs/.+}", rest.HandleDeleteBlurb_1).Methods("DELETE")
	router.HandleFunc("/v1beta1/{parent:rooms/.+}/blurbs", rest.HandleListBlurbs).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:users/.+/profile}/blurbs", rest.HandleListBlurbs_1).Methods("GET")
	router.HandleFunc("/v1beta1/{parent:rooms/.+}/blurbs:search", rest.HandleSearchBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/.+/profile}/blurbs:search", rest.HandleSearchBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/{name:rooms/.+}/blurbs:stream", rest.HandleStreamBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{name:users/.+/profile}/blurbs:stream", rest.HandleStreamBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:rooms/.+}/blurbs:send", rest.HandleSendBlurbs).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:users/.+/profile}/blurbs:send", rest.HandleSendBlurbs_1).Methods("POST")
	router.HandleFunc("/v1beta1/sequences", rest.HandleCreateSequence).Methods("POST")
	router.HandleFunc("/v1beta1/{name:sequences/.+/sequenceReport}", rest.HandleGetSequenceReport).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sequences/.+}", rest.HandleAttemptSequence).Methods("POST")
	router.HandleFunc("/v1beta1/sessions", rest.HandleCreateSession).Methods("POST")
	router.HandleFunc("/v1beta1/{name:sessions/.+}", rest.HandleGetSession).Methods("GET")
	router.HandleFunc("/v1beta1/sessions", rest.HandleListSessions).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sessions/.+}", rest.HandleDeleteSession).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:sessions/.+}:report", rest.HandleReportSession).Methods("POST")
	router.HandleFunc("/v1beta1/{parent:sessions/.+}/tests", rest.HandleListTests).Methods("GET")
	router.HandleFunc("/v1beta1/{name:sessions/.+/tests/.+}", rest.HandleDeleteTest).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:sessions/.+/tests/.+}:check", rest.HandleVerifyTest).Methods("POST")
	router.HandleFunc("/v1beta1/operations", rest.HandleListOperations).Methods("GET")
	router.HandleFunc("/v1beta1/{name:operations/.+}", rest.HandleGetOperation).Methods("GET")
	router.HandleFunc("/v1beta1/{name:operations/.+}", rest.HandleDeleteOperation).Methods("DELETE")
	router.HandleFunc("/v1beta1/{name:operations/.+}:cancel", rest.HandleCancelOperation).Methods("POST")
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
