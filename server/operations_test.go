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
  "testing"
  "time"

  "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/grpc/grpc-go/status"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"golang.org/x/net/context"
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
  err error
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
  server := NewOperationsServer(NewOpertionStore())
  _, err := server.ListOperations(context.Background(), nil)
  s, _ := status.FromError(err)
  if codes.Unimplemented != s.Code() {
    t.Error("ListOperations should return an error indicating it is Unimplemented.")
  }
}

func TestServerDeleteOperation(t *testing.T) {
  server := NewOperationsServer(NewOpertionStore())
  _, err := server.DeleteOperation(context.Background(), nil)
  s, _ := status.FromError(err)
  if codes.Unimplemented != s.Code() {
    t.Error("DeleteOperations should return an error indicating it is Unimplemented.")
  }
}

func TestStoreGet_NotFound(t *testing.T) {
    store := NewOpertionStore()
    _, err := store.Get("non-existant");
    s, _ := status.FromError(err)
    if codes.NotFound != s.Code() {
      t.Error("Expect to return code NotFound if operation name is not found.")
    }
}

func TestStoreGet_Cancelled(t *testing.T) {
  store := &operationStoreImpl{
    nowF: time.Now,
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
    nowF: mockNow(time.Unix(200, 0)),
    store: map[string]*operationInfo{},
  }
  expected := &pb.LongrunningResponse{Content: "content"}
  store.store["name"] = &operationInfo{
    name: "name",
    resp: expected,
    end: end,
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
