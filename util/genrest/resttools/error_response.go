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

package resttools

import (
	"encoding/json"
	"net/http"

	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
)

// gRPCToHTTP status code mapping derived from internal source:
// go/http-canonical-mapping.
var gRPCToHTTP map[codes.Code]int = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           499, // There isn't a Go constant ClientClosedConnection
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
}

// Google API Errors, as defined by https://cloud.google.com/apis/design/errors
// will consist of a googleapi.Error nested as the key `error` in a JSON object.
// So we must create such a structure to wrap our googleapi.Error in.
type googleAPIError struct {
	Error *googleapi.Error `json:"error"`
}

// GRPCToHTTP maps the given gRPC Code to the canonical HTTP Status code as
// defined by the internal source go/http-canonical-mapping.
func GRPCToHTTP(c codes.Code) int {
	httpStatus, ok := gRPCToHTTP[c]
	if !ok {
		httpStatus = http.StatusInternalServerError
	}

	return httpStatus
}

// ErrorResponse is a helper that formats the given response information,
// including the HTTP Status code, a message, and any error detail types, into
// a googleAPIError and writes the response as JSON.
func ErrorResponse(w http.ResponseWriter, status int, message string, details ...interface{}) {
	googleAPIError := &googleAPIError{
		Error: &googleapi.Error{
			Code:    status,
			Message: message,
			Details: details,
		},
	}
	w.WriteHeader(status)
	data, _ := json.Marshal(googleAPIError)
	w.Write(data)
}
