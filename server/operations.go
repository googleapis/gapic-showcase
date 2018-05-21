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

package showcase

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
)

type OperationsServer struct {
	store OperationStore
}

func NewOperationsServer(opStore OperationStore) *OperationsServer {
	return &OperationsServer{store: opStore}
}

func (s *OperationsServer) GetOperation(ctx context.Context, in *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	return s.store.Get(in.GetName())
}

func (s *OperationsServer) CancelOperation(ctx context.Context, in *longrunning.CancelOperationRequest) (*empty.Empty, error) {
	err := s.store.Cancel(in.GetName())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OperationsServer) ListOperations(ctx context.Context, in *longrunning.ListOperationsRequest) (*longrunning.ListOperationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.ListOperations is unimplemented.")
}

func (s *OperationsServer) DeleteOperation(ctx context.Context, in *longrunning.DeleteOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.DeleteOperation is unimplemented.")
}
