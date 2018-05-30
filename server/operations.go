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
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/googleapis/longrunning"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewOperationsServer returns a longrunning.OperationsServer which uses the
// given OperationStore.
func NewOperationsServer(opStore OperationStore) longrunning.OperationsServer {
	return &operationsServerImpl{store: opStore}
}

// OperationStore is an interface for storing longrunning requests for the
// Showcase API.
type OperationStore interface {
	RegisterOp(*pb.LongrunningRequest) (*longrunning.Operation, error)
	Get(string) (*longrunning.Operation, error)
	Cancel(string) error
}

// NewOperationStore returns an implemented OperationStore.
func NewOperationStore() OperationStore {
	return &operationStoreImpl{
		store: map[string]*operationInfo{},
	}
}

type operationsServerImpl struct {
	store OperationStore
}

func (s operationsServerImpl) GetOperation(ctx context.Context, in *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	return s.store.Get(in.GetName())
}

func (s operationsServerImpl) CancelOperation(ctx context.Context, in *longrunning.CancelOperationRequest) (*empty.Empty, error) {
	err := s.store.Cancel(in.GetName())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s operationsServerImpl) ListOperations(ctx context.Context, in *longrunning.ListOperationsRequest) (*longrunning.ListOperationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.ListOperations is unimplemented.")
}

func (s operationsServerImpl) DeleteOperation(ctx context.Context, in *longrunning.DeleteOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.DeleteOperation is unimplemented.")
}

type operationInfo struct {
	name     string
	start    time.Time
	end      time.Time
	canceled bool
	resp     *pb.LongrunningResponse
	err      *statuspb.Status
}

type operationStoreImpl struct {
	uid  uniqID
	nowF func() time.Time

	mu    sync.Mutex
	store map[string]*operationInfo
}

func (s *operationStoreImpl) RegisterOp(op *pb.LongrunningRequest) (*longrunning.Operation, error) {
	end, err := ptypes.Timestamp(op.CompletionTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Given operation completion time is invalid.")
	}
	name := fmt.Sprintf("lro-test-op-%d", s.uid.id())

	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[name] = &operationInfo{
		name:     name,
		start:    s.nowF(),
		end:      end,
		canceled: false,
		resp:     op.GetSuccess(),
		err:      op.GetError(),
	}
	return s.get(name)
}

func (s *operationStoreImpl) Get(name string) (*longrunning.Operation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.get(name)
}

func (s *operationStoreImpl) get(name string) (*longrunning.Operation, error) {
	op, ok := s.store[name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Operation %q not found.", name)
	}
	ret := &longrunning.Operation{
		Name: op.name,
	}

	now := s.nowF()

	if op.canceled {
		ret.Result = &longrunning.Operation_Error{
			Error: status.Newf(
				codes.Canceled,
				"Operation %q has been canceled.", name).Proto(),
		}
	} else if now.After(op.end) {
		if op.err != nil {
			ret.Result = &longrunning.Operation_Error{Error: op.err}
		} else {
			resp, _ := ptypes.MarshalAny(op.resp)
			ret.Result = &longrunning.Operation_Response{Response: resp}
			ret.Done = true
		}
	} else {
		meta, _ := ptypes.MarshalAny(
			&pb.LongrunningMetadata{
				TimeRemaining: ptypes.DurationProto(op.end.Sub(now))})
		ret.Metadata = meta
	}
	return ret, nil
}

func (s *operationStoreImpl) Cancel(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	op, ok := s.store[name]
	if !ok {
		return status.Errorf(codes.NotFound, "Operation %q not found.", name)
	}
	op.canceled = true
	s.store[name] = op
	return nil
}
