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
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	featurepb "github.com/googleapis/feature-testing-server/server/genproto"
	"github.com/grpc/grpc-go/status"

	lropb "google.golang.org/genproto/googleapis/longrunning"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type operationInfo struct {
	name     string
	start    time.Time
	end      time.Time
	canceled bool
	resp     *featurepb.LongrunningResponse
	err      *statuspb.Status
}

type OperationStore interface {
	RegisterOp(*featurepb.LongrunningRequest) (*lropb.Operation, error)
	Get(string) (*lropb.Operation, error)
	Cancel(string) error
}

type OperationStoreImpl struct {
	nowF  func() time.Time
	store map[string]*operationInfo
}

func (s *OperationStoreImpl) WithNowF(nowFunc func() time.Time) *OperationStoreImpl {
	return &OperationStoreImpl{
		nowF:  nowFunc,
		store: s.store,
	}
}

func (s *OperationStoreImpl) RegisterOp(op *featurepb.LongrunningRequest) (*lropb.Operation, error) {
	end, err := ptypes.Timestamp(op.CompletionTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Given operation completion time is invalid.")
	}
	now := s.nowF()
	name := fmt.Sprintf("lro-test-op-%d", now.UTC().Unix())
	s.store[name] = &operationInfo{
		name:     name,
		start:    now,
		end:      end,
		canceled: false,
		resp:     op.GetSuccess(),
		err:      op.GetError(),
	}
	return s.Get(name)
}

func (s *OperationStoreImpl) Get(name string) (*lropb.Operation, error) {
	op, ok := s.store[name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Operation '%s' not found.", name)
	}
	ret := &lropb.Operation{
		Name: op.name,
	}

	now := s.now()

	if op.canceled {
		ret.Result = &lropb.Operation_Error{
			Error: status.Newf(
				codes.Canceled,
				"Operation '%s' has been canceled.", name).Proto(),
		}
	} else if now.After(op.end) {
		if op.err != nil {
			ret.Result = &lropb.Operation_Error{Error: op.err}
		} else {
			resp, err := ptypes.MarshalAny(op.resp)
			if err != nil {
				return nil, err
			}
			ret.Result = &lropb.Operation_Response{Response: resp}
			ret.Done = true
			meta, err := ptypes.MarshalAny(&featurepb.LongrunningMetadata{TimeRemaining: ptypes.DurationProto(0)})
			if err != nil {
				return nil, err
			}
			ret.Metadata = meta
		}
	} else {
		meta, err := ptypes.MarshalAny(
			&featurepb.LongrunningMetadata{
				TimeRemaining: ptypes.DurationProto(now.Sub(op.end))})
		if err != nil {
			return nil, err
		}
		ret.Metadata = meta
	}
	return ret, nil
}

func (s *OperationStoreImpl) Cancel(name string) error {
	op, ok := s.store[name]
	if !ok {
		return status.Errorf(codes.NotFound, "Operation '%s' not found.", name)
	}
	op.canceled = true
	s.store[name] = op
	return nil
}

func (s *OperationStoreImpl) now() time.Time {
	if s.nowF != nil {
		return s.nowF()
	}
	return time.Now()
}
