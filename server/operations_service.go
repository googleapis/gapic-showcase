// Copyright 2018 Google LLC
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

package server

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewEchoServer returns a new EchoServer for the Showcase API.
func NewOperationsServer() lropb.OperationsServer {
	return &operationsServerImpl{waiter: waiterSingleton}
}

type operationsServerImpl struct {
	waiter Waiter
}

func (s *operationsServerImpl) GetOperation(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	prefix := "operations/google.showcase.v1alpha2.Echo/Wait/"
	if !strings.HasPrefix(in.Name, prefix) {
		return nil, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
	}

	waitReq := &pb.WaitRequest{}
	encodedBytes := strings.TrimPrefix(in.Name, prefix)
	waitReqBytes, err := base64.StdEncoding.DecodeString(encodedBytes)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
	}

	err = proto.Unmarshal(waitReqBytes, waitReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
	}

	return s.waiter.Wait(waitReq), nil
}

func (s operationsServerImpl) CancelOperation(ctx context.Context, in *lropb.CancelOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.CancelOperation is unimplemented.")
}

func (s operationsServerImpl) ListOperations(ctx context.Context, in *lropb.ListOperationsRequest) (*lropb.ListOperationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.ListOperations is unimplemented.")
}

func (s operationsServerImpl) DeleteOperation(ctx context.Context, in *lropb.DeleteOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.DeleteOperation is unimplemented.")
}
