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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/longrunning"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type getOpStore struct {
	getCalled bool
	OperationStore
}

func (m *getOpStore) Get(name string) (*longrunning.Operation, error) {
	m.getCalled = true
	return nil, nil
}

func TestServerGetOperation(t *testing.T) {
	store := &getOpStore{}
	server := NewOperationsServer(store)
	server.GetOperation(context.Background(), &longrunning.GetOperationRequest{Name: "name"})

	if !store.getCalled {
		t.Error("Expected to defer get to operation store.")
	}
}

type cancelOpStore struct {
	cancelCalled bool
	err          error
	OperationStore
}

func (m *cancelOpStore) Cancel(name string) error {
	m.cancelCalled = true
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestServerCancelOperation(t *testing.T) {
	store := &cancelOpStore{}
	server := NewOperationsServer(store)
	server.CancelOperation(context.Background(), &longrunning.CancelOperationRequest{Name: "name"})

	if !store.cancelCalled {
		t.Error("Expected to defer cancel to operation store.")
	}

	e := errors.New("Test Error")
	store = &cancelOpStore{err: e}
	server = NewOperationsServer(store)
	_, err := server.CancelOperation(context.Background(), &longrunning.CancelOperationRequest{Name: "name"})
	if err != e {
		t.Errorf("Expected to pass through cancel errors")
	}
}

func TestServerListOperation(t *testing.T) {
	server := NewOperationsServer(NewOperationStore())
	_, err := server.ListOperations(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Error("ListOperations should return an error indicating it is Unimplemented.")
	}
}

func TestServerDeleteOperation(t *testing.T) {
	server := NewOperationsServer(NewOperationStore())
	_, err := server.DeleteOperation(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Error("DeleteOperations should return an error indicating it is Unimplemented.")
	}
}

func TestStoreRegisterOp(t *testing.T) {
	completion, _ := ptypes.TimestampProto(time.Unix(200, 0))
	req := &pb.LongrunningRequest{
		CompletionTime: completion,
		Response: &pb.LongrunningRequest_Success{
			Success: &pb.LongrunningResponse{Content: "content"},
		},
	}
	store := &operationStoreImpl{
		nowF:  mockNow(time.Unix(100, 0)),
		store: map[string]*operationInfo{},
	}

	op, err := store.RegisterOp(req)
	if err != nil {
		t.Error(err)
	}
	if op == nil {
		t.Error("Expected RegisterOp to return an Operation.")
	}

	expectedName := fmt.Sprintf("lro-test-op-%d", time.Unix(100, 0).Unix())
	if val, ok := store.store[expectedName]; ok {
		if val.name != expectedName {
			t.Errorf("Expected registered op name to be %s, but was %s",
				expectedName, val.name)
		}
		if !time.Unix(100, 0).Equal(val.start) {
			t.Errorf("Expected start time to be %d, but was %d", time.Unix(100, 0).Unix(),
				val.start.Unix())
		}
		if !time.Unix(200, 0).Equal(val.end) {
			t.Errorf("Expected end time to be %d, but was %d", time.Unix(200, 0).Unix(),
				val.end.Unix())
		}
		if val.canceled {
			t.Errorf("A newly registered op should not be canceled.")
		}
		if req.GetSuccess() != val.resp {
			t.Errorf("Expected the op response to be %s, but was %s",
				req.GetSuccess().String(), val.resp)
		}
		if req.GetError() != val.err {
			t.Errorf("Expected the op err to be %s, but was %s",
				req.GetError().String(), val.resp)
		}
	} else {
		t.Errorf("Expected store to contain value with key %s", expectedName)
	}
}

func TestStoreRegisterOp_InvalidArgs(t *testing.T) {
	tests := []*timestamp.Timestamp{
		nil,
		{Nanos: -1},
	}
	for _, test := range tests {
		store := NewOperationStore()
		req := &pb.LongrunningRequest{
			CompletionTime: test,
		}
		_, err := store.RegisterOp(req)
		s, _ := status.FromError(err)
		if codes.InvalidArgument != s.Code() {
			t.Error("Expected to return InvalidArgument for invalid completion times.")
		}
	}
}

func TestStoreGet_NotFound(t *testing.T) {
	store := NewOperationStore()
	_, err := store.Get("non-existant")
	s, _ := status.FromError(err)
	if codes.NotFound != s.Code() {
		t.Error("Expect to return code NotFound if operation name is not found.")
	}
}

func TestStoreGet_Cancelled(t *testing.T) {
	store := &operationStoreImpl{
		nowF:  time.Now,
		store: map[string]*operationInfo{},
	}
	store.store["name"] = &operationInfo{name: "name", canceled: true}
	op, _ := store.Get("name")
	s := status.FromProto(op.GetError())
	if codes.Canceled != s.Code() {
		t.Error("Expected to return a canceled operation if an operation was canceled.")
	}
}

func mockNow(t time.Time) func() time.Time {
	return func() time.Time {
		return t
	}
}

func TestStoreGet_Done(t *testing.T) {
	end := time.Unix(100, 0)
	store := &operationStoreImpl{
		nowF:  mockNow(time.Unix(200, 0)),
		store: map[string]*operationInfo{},
	}
	expected := &pb.LongrunningResponse{Content: "content"}
	store.store["name"] = &operationInfo{
		name: "name",
		resp: expected,
		end:  end,
	}
	op, _ := store.Get("name")
	if !op.Done {
		t.Error("Operation marked to end should have been marked done.")
	}
	resp := &pb.LongrunningResponse{}
	ptypes.UnmarshalAny(op.GetResponse(), resp)
	if !proto.Equal(expected, resp) {
		t.Errorf("Expected result to be %s, but was %s", expected.String(), resp.String())
	}
}

func TestStoreGet_Err(t *testing.T) {
	end := time.Unix(100, 0)
	store := &operationStoreImpl{
		nowF:  mockNow(time.Unix(200, 0)),
		store: map[string]*operationInfo{},
	}
	expected := &statuspb.Status{Code: int32(codes.Aborted)}
	store.store["name"] = &operationInfo{
		name: "name",
		err:  expected,
		end:  end,
	}
	op, _ := store.Get("name")
	if !proto.Equal(expected, op.GetError()) {
		t.Errorf("Expected result to be %s, but was %s", expected.String(), op.GetError().String())
	}
}

func TestStoreGet_Pending(t *testing.T) {
	start := time.Unix(0, 0)
	end := time.Unix(100, 0)
	store := &operationStoreImpl{
		nowF:  mockNow(time.Unix(90, 0)),
		store: map[string]*operationInfo{},
	}
	store.store["name"] = &operationInfo{
		name:  "name",
		start: start,
		end:   end,
	}
	op, err := store.Get("name")
	if err != nil {
		t.Error(err)
	}
	meta := &pb.LongrunningMetadata{}
	ptypes.UnmarshalAny(op.GetMetadata(), meta)
	dur, err := ptypes.Duration(meta.GetTimeRemaining())
	if err != nil {
		t.Error(err)
	}
	expected := time.Duration(10) * time.Second
	if dur != expected {
		t.Errorf("Expected the duration to be %s, but was, %s", expected, dur)
	}
}

func TestStoreCancel(t *testing.T) {
	store := &operationStoreImpl{
		store: map[string]*operationInfo{},
	}
	store.store["name"] = &operationInfo{name: "name"}
	store.Cancel("name")
	if !store.store["name"].canceled {
		t.Error("Expected Cancel to mark an operation as canceled.")
	}

	err := store.Cancel("non-existant")
	s, _ := status.FromError(err)
	if codes.NotFound != s.Code() {
		t.Error("Expect to return code NotFound if operation name is not found.")
	}
}
