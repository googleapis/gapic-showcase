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
	"fmt"
	"net/http"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	code "google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// gRPCToHTTP is the status code mapping derived from the internal source go/http-canonical-mapping.
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

// httpToGRPC is the status code mapping derived from the internal source go/http-canonical-mapping.
// This is not merely the inverse of gRPCToHTTP (which, at any rate, is not injective). The
// canonical mapping also specifies codes for some HTTP status ranges, so it is imperative to use
// HTTPToGRPC()
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

// Constants to make it obvious when we're not supplying either a gRPC or HTTP status code to
// ErrorResponse(). The code we're not supplying will be obtained from the one we do supply.
const (
	NoCodeGRPC codes.Code = 9999
	NoCodeHTTP int        = -1
)

// ErrorResponse is a helper that formats the given response information,
// including the HTTP Status code, a message, and any error detail types, into
// a googleAPIError and writes the response as JSON.
func ErrorResponse(w http.ResponseWriter, httpResponseCode int, grpcStatus codes.Code, message string, details ...interface{}) {
	if httpResponseCode == NoCodeHTTP {
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

	jsonError := &pb.RestError{
		Error: &pb.RestError_Status{
			Code:    int32(httpResponseCode),
			Message: message,
			Status:  code.Code(grpcStatus),
		},
	}

	for _, oneDetail := range details {
		detailAsAny, err := anypb.New(oneDetail.(proto.Message))
		if err != nil {
			WriteShowcaseRESTImplementationError(w, fmt.Sprintf("could not interpret error detail as a protobuf Any: %+v", oneDetail))
			return
		}
		jsonError.Error.Details = append(jsonError.Error.Details, detailAsAny)
	}

	w.WriteHeader(httpResponseCode)
	data, _ := protojson.Marshal(jsonError)
	w.Write(data)
}

func WriteShowcaseRESTImplementationError(w http.ResponseWriter, message string) {
	ErrorResponse(w, 500, NoCodeGRPC, fmt.Sprintf("Showcase consistency error: %s", message))
}
