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

package genrest

import (
	"context"
	"net/http"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
)

// .google.showcase.v1beta1.Compliance.HandleRepeatDataBody
func (backend *RESTBackend) customRepeatWithUnknownEnum(w http.ResponseWriter, r *http.Request, request *genprotopb.RepeatRequest) {
	marshaler := resttools.ToJSON()

	response, err := backend.ComplianceServer.RepeatWithUnknownEnum(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	// FIXME: change response field to sentinel enum

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	// FIXME: change sentinel to unknown value

	w.Write(json)
}
