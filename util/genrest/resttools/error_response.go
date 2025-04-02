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
	"fmt"
	"net/http"

	"github.com/iancoleman/strcase"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
)

// gRPCToHTTP is the status code mapping derived from internal source:
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

// httpToGRPC is the status code mapping derived from internal source:
// go/http-canonical-mapping. This is not merely the inverse of gRPCToHTTP (which, at any rate, is
// not injective). The canonical mapping also specifies codes for some HTTP status ranges, so it is
// imperative to use HTTPToGRPC()
var httpToGRPC = map[int]codes.Code{
	http.StatusOK:                           codes.OK,
	http.StatusMultipleChoices:              codes.Unknown,
	http.StatusBadRequest:                   codes.InvalidArgument,
	http.StatusUnauthorized:                 codes.Unauthenticated,
	http.StatusForbidden:                    codes.PermissionDenied,
	http.StatusNotFound:                     codes.NotFound,
	http.StatusConflict:                     codes.Aborted,
	http.StatusRequestedRangeNotSatisfiable: codes.OutOfRange,
	http.StatusTooManyRequests:              codes.ResourceExhausted,
	http.StatusNotImplemented:               codes.Unimplemented,
	http.StatusServiceUnavailable:           codes.Unavailable,
	http.StatusGatewayTimeout:               codes.DeadlineExceeded,
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

// HTTPToGRPC maps the given HTTP status code to the canonical gRPC
// Code as defined by the internal source go/http-canonical-mapping.
func HTTPToGRPC(httpStatus int) codes.Code {
	var gRPCCode codes.Code
	switch {
	case httpStatus >= 200 && httpStatus < 300:
		gRPCCode = codes.OK
	case httpStatus >= 300 && httpStatus < 400:
		gRPCCode = codes.Unknown
	case httpStatus >= 400 && httpStatus < 500:
		gRPCCode = codes.FailedPrecondition
	case httpStatus >= 500 && httpStatus < 600:
		gRPCCode = codes.Internal
	default:
		gRPCCode = codes.Unknown
	}

	if codeOverride, ok := httpToGRPC[httpStatus]; ok {
		gRPCCode = codeOverride
	}
	return gRPCCode
}

type JSONErrorInternals struct {
	// The HTTP status code that corresponds to `google.rpc.Status.code`
	Code int `json:"code"`
	// This corresponds to `google.rpc.Status.message`.
	Message string `json:"message"`
	// This is the enum version for `google.rpc.Status.code`, cased as in
	// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto.
	Status string `json:"status"`
	// This corresponds to `google.rpc.Status.details`.
	Details []any `json:"details"`
}

type JSONErrorResponseBody struct {
	Error JSONErrorInternals `json:"error"`
}

// Google API Errors, as defined by https://cloud.google.com/apis/design/errors
// will consist of a googleapi.Error nested as the key `error` in a JSON object.
// So we must create such a structure to wrap our googleapi.Error in.
type googleAPIError struct {
	Error *googleapi.Error `json:"error"`
}

const (
	NoCodeGRPC        codes.Code = 9999
	CalculateCodeHTTP int        = -1
)

// ErrorResponse is a helper that formats the given response information,
// including the HTTP Status code, a message, and any error detail types, into
// a googleAPIError and writes the response as JSON.
func ErrorResponse(w http.ResponseWriter, httpResponseCode int, grpcStatus codes.Code, message string, details ...interface{}) {
	if httpResponseCode == CalculateCodeHTTP {
		if grpcStatus == NoCodeGRPC {
			WriteShowcaseRESTImplementationError(w, "neither HTTP code or gRPC status are provided for ErrorResponse. Exactly one must be provided.")
			return
		}
		httpResponseCode = GRPCToHTTP(grpcStatus)
	} else {
		if grpcStatus != NoCodeGRPC {
			WriteShowcaseRESTImplementationError(w, "both HTTP code and gRPC status are provided for ErrorResponse. Exactly one must be provided.")
			return
		}
		grpcStatus = HTTPToGRPC(httpResponseCode)
	}

	jsonErrorBody := &JSONErrorResponseBody{
		Error: JSONErrorInternals{
			Code:    httpResponseCode,
			Message: message,
			Status:  strcase.ToScreamingSnake(grpcStatus.String()),
			Details: details,
		},
	}

	w.WriteHeader(httpResponseCode)
	data, _ := json.Marshal(jsonErrorBody)
	w.Write(data)
}

func WriteShowcaseRESTImplementationError(w http.ResponseWriter, message string) {
	ErrorResponse(w, 500, NoCodeGRPC, fmt.Sprintf("Showcase consistency error: %s", message))
}
