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
// for service #3: "SequenceService" (.google.showcase.v1beta1.SequenceService).

package genrest

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
)

// HandleCreateSequence translates REST requests/responses on the wire to internal proto messages for CreateSequence
//    Generated for HTTP binding pattern: /v1beta1/sequences
//         This matches URIs of the form: /v1beta1/sequences
func (backend *RESTBackend) HandleCreateSequence(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/sequences': %q", r.URL)

	var request *genprotopb.CreateSequenceRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.SequenceServiceServer.CreateSequence(context.Background(), request)
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

// HandleGetSequenceReport translates REST requests/responses on the wire to internal proto messages for GetSequenceReport
//    Generated for HTTP binding pattern: /v1beta1/{name=sequences/*/sequenceReport}
//         This matches URIs of the form: /v1beta1/{name:sequences/[a-zA-Z_%\-]+/sequenceReport}
func (backend *RESTBackend) HandleGetSequenceReport(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=sequences/*/sequenceReport}': %q", r.URL)

	var request *genprotopb.GetSequenceReportRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.SequenceServiceServer.GetSequenceReport(context.Background(), request)
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

// HandleAttemptSequence translates REST requests/responses on the wire to internal proto messages for AttemptSequence
//    Generated for HTTP binding pattern: /v1beta1/{name=sequences/*}
//         This matches URIs of the form: /v1beta1/{name:sequences/[a-zA-Z_%\-]+}
func (backend *RESTBackend) HandleAttemptSequence(w http.ResponseWriter, r *http.Request) {
	backend.StdLog.Printf("Received request matching '/v1beta1/{name=sequences/*}': %q", r.URL)

	var request *genprotopb.AttemptSequenceRequest
	// TODO: Populate request with parameters from HTTP request

	response, err := backend.SequenceServiceServer.AttemptSequence(context.Background(), request)
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
