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
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetOperation(t *testing.T) {
	endTime, _ := ptypes.TimestampProto(time.Now())
	waitReq := &pb.WaitRequest{EndTime: endTime}
	nameBytes, _ := proto.Marshal(waitReq)
	req := &lropb.GetOperationRequest{
		Name: fmt.Sprintf(
			"operations/google.showcase.v1alpha2.Echo/Wait/%s",
			base64.StdEncoding.EncodeToString(nameBytes)),
	}

	waiter := &mockWaiter{}
	server := &operationsServerImpl{waiter: waiter}
	server.GetOperation(context.Background(), req)
	if !proto.Equal(waiter.req, waitReq) {
		t.Error("Expected echo.Wait to defer to waiter.")
	}
}

func TestGetOperation_notFoundOperation(t *testing.T) {
	req := &lropb.GetOperationRequest{
		Name: "BOGUS",
	}
	server := NewOperationsServer()
	_, err := server.GetOperation(context.Background(), req)
	s, _ := status.FromError(err)
	if codes.NotFound != s.Code() {
		t.Errorf("GetOperation with invalid name expected code=%d, got %d", codes.NotFound, s.Code())
	}
}

func TestGetOperation_invalidEncodedName(t *testing.T) {
	req := &lropb.GetOperationRequest{
		Name: "operations/google.showcase.v1alpha2.Echo/Wait/BOGUS",
	}
	server := NewOperationsServer()
	_, err := server.GetOperation(context.Background(), req)
	s, _ := status.FromError(err)
	if codes.NotFound != s.Code() {
		t.Errorf("GetOperation with invalid name expected code=%d, got %d", codes.NotFound, s.Code())
	}
}

func TestCancelOperation(t *testing.T) {
	server := NewOperationsServer()
	_, err := server.CancelOperation(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Errorf("CancelOperation expected code=%d, got %d", codes.Unimplemented, s.Code())
	}
}

func TestServerListOperation(t *testing.T) {
	server := NewOperationsServer()
	_, err := server.ListOperations(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Errorf("ListOperations expected code=%d, got %d", codes.Unimplemented, s.Code())
	}
}

func TestServerDeleteOperation(t *testing.T) {
	server := NewOperationsServer()
	_, err := server.DeleteOperation(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Errorf("DeleteOperations expected code=%d, got %d", codes.Unimplemented, s.Code())
	}
}
