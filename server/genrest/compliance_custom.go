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
	"bytes"
	"context"
	"net/http"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
)

// customRepeatWithUnknownEnum provides REST-specific handling for a RepeatWithUnknownEnum request. It returns a JSON response with an unknown enum symbol string in an enum field.
func (backend *RESTBackend) customRepeatWithUnknownEnum(w http.ResponseWriter, r *http.Request, request *genprotopb.RepeatRequest) {
	marshaler := resttools.ToJSON()

	response, err := backend.ComplianceServer.RepeatWithUnknownEnum(context.Background(), request)
	if err != nil {
		// TODO: Properly handle error. Is StatusInternalServerError (500) the right response?
		backend.Error(w, http.StatusInternalServerError, "server error: %s", err.Error())
		return
	}

	// Make sure we have at least one sentinel value before serializing properly.
	sentinelValue := genprotopb.ComplianceData_LIFE_KINGDOM_UNSPECIFIED
	sentinelString := genprotopb.ComplianceData_LifeKingdom_name[int32(sentinelValue)]
	response.Request.Info.FKingdom = sentinelValue

	json, err := marshaler.Marshal(response)
	if err != nil {
		backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %s", err.Error())
		return
	}

	// Change the sentinel string to an unknown value.
	json = bytes.ReplaceAll(json, []byte(sentinelString), []byte("LIFE_KINGDOM_NEW"))

	w.Write(json)
}
