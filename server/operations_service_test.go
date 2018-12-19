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

func TestGetOperation_wait(t *testing.T) {
	endTime, _ := ptypes.TimestampProto(time.Now())
	waitReq := &pb.WaitRequest{EndTime: endTime}
	nameBytes, _ := proto.Marshal(waitReq)
	req := &lropb.GetOperationRequest{
		Name: fmt.Sprintf(
			"operations/google.showcase.v1alpha3.Echo/Wait/%s",
			base64.StdEncoding.EncodeToString(nameBytes)),
	}

	waiter := &mockWaiter{}
	server := &operationsServerImpl{waiter: waiter}
	server.GetOperation(context.Background(), req)
	if !proto.Equal(waiter.req, waitReq) {
		t.Error("Expected echo.Wait to defer to waiter.")
	}
}

type messagingServerWrapper struct {
	listReq *pb.ListBlurbsRequest

	MessagingServer
}

func (m *messagingServerWrapper) FilteredListBlurbs(ctx context.Context, r *pb.ListBlurbsRequest, f func(*pb.Blurb) bool) (*pb.ListBlurbsResponse, error) {
	m.listReq = r
	return m.MessagingServer.FilteredListBlurbs(ctx, r, f)
}

func TestGetOperation_searchBlurbs(t *testing.T) {
	expected := []*pb.Blurb{
		&pb.Blurb{
			User:    "users/rumble",
			Content: &pb.Blurb_Text{Text: "woof"},
		},
		&pb.Blurb{
			User:    "users/ekko",
			Content: &pb.Blurb_Text{Text: "bark"},
		},
	}
	wrapped := &messagingServerWrapper{
		MessagingServer: &messagingServerImpl{
			identityServer: &mockIdentityServer{},
			roomKeys:       map[string]int{},
			parentUids:     map[string]*uniqID{},
			token:          NewTokenGenerator(),
			blurbKeys: map[string]blurbIndex{
				"users/rumble/profile/messages/1": blurbIndex{
					row: "users/rumble/profile",
					col: 1},
				"users/rumble/profile/messages/2": blurbIndex{
					row: "users/rumble/profile",
					col: 2},
				"users/rumble/profile/messages/3": blurbIndex{
					row: "users/rumble/profile",
					col: 3},
				"users/rumble/profile/messages/4": blurbIndex{
					row: "users/rumble/profile",
					col: 4},
			},
			blurbs: map[string][]blurbEntry{
				"users/rumble/profile": []blurbEntry{
					blurbEntry{blurb: expected[0]},
					blurbEntry{
						blurb: &pb.Blurb{
							User:    "users/musubi",
							Content: &pb.Blurb_Text{Text: "meow"},
						},
					},
					blurbEntry{blurb: expected[1]},
					blurbEntry{blurb: expected[1], deleted: true},
					blurbEntry{
						blurb: &pb.Blurb{
							User:    "users/musubi",
							Content: &pb.Blurb_Text{Text: "meow"},
						},
					},
				},
			},
		},
	}

	server := NewOperationsServer(wrapped)

	searchReq := &pb.SearchBlurbsRequest{
		Query:    "woof bark",
		Parent:   "users/rumble/profile",
		PageSize: 2,
	}
	reqBytes, _ := proto.Marshal(searchReq)
	req := &lropb.GetOperationRequest{
		Name: fmt.Sprintf(
			"operations/google.showcase.v1alpha3.Messaging/SearchBlurbs/%s",
			base64.StdEncoding.EncodeToString(reqBytes)),
	}
	op, err := server.GetOperation(context.Background(), req)
	if err != nil {
		t.Errorf("GetOperation: unexpected err %+v", err)
	}

	listReq := wrapped.listReq
	if listReq.Parent != searchReq.GetParent() {
		t.Errorf(
			"GetOperation searchBlurbs: list request parent expected %s got %s",
			searchReq.GetParent(),
			listReq.Parent)
	}
	if listReq.PageSize != searchReq.GetPageSize() {
		t.Errorf(
			"GetOperation searchBlurbs: list request page size expected %d got %d",
			searchReq.GetPageSize(),
			listReq.PageSize)
	}
	if listReq.PageToken != searchReq.GetPageToken() {
		t.Errorf(
			"GetOperation searchBlurbs: list request page token expected %s got %s",
			searchReq.GetPageToken(),
			listReq.PageToken)
	}

	if !op.Done {
		t.Errorf("SearchBlurbs() for %q expected done=true got done=false", req)
	}
	if op.Metadata != nil {
		t.Errorf("SearchBlurbs() for %q expected nil metadata, got %q", req, op.Metadata)
	}
	if op.GetError() != nil {
		t.Errorf("SearchBlurbs() expected op.Error=nil, got %q", op.GetError())
	}
	if op.GetResponse() == nil {
		t.Error("SearchBlurbs() expected op.Response!=nil")
	}
	resp := &pb.SearchBlurbsResponse{}
	ptypes.UnmarshalAny(op.GetResponse(), resp)
	if len(resp.GetBlurbs()) != len(expected) {
		t.Errorf(
			"SearchBlurbs() expected blurbs size %d, got %d",
			len(expected),
			len(resp.GetBlurbs()))
	}
	for i, b := range expected {
		if !proto.Equal(b, resp.GetBlurbs()[i]) {
			t.Errorf(
				"SearchBlurbs().blurbs[%d] want %+v, got %+v",
				i,
				b,
				resp.GetBlurbs()[i],
			)
		}
	}
}

func TestGetOperation_notFoundOperation(t *testing.T) {
	req := &lropb.GetOperationRequest{
		Name: "BOGUS",
	}
	server := NewOperationsServer(nil)
	_, err := server.GetOperation(context.Background(), req)
	s, _ := status.FromError(err)
	if codes.NotFound != s.Code() {
		t.Errorf("GetOperation with invalid name expected code=%d, got %d", codes.NotFound, s.Code())
	}
}

func TestGetOperation_invalidEncodedName(t *testing.T) {
	prefixes := []string{
		"operations/google.showcase.v1alpha3.Echo/Wait",
		"operations/google.showcase.v1alpha3.Messaging/SearchBlurbs",
	}
	for _, p := range prefixes {
		req := &lropb.GetOperationRequest{
			Name: fmt.Sprintf("%s/BOGUS", p),
		}
		server := NewOperationsServer(nil)
		_, err := server.GetOperation(context.Background(), req)
		s, _ := status.FromError(err)
		if codes.NotFound != s.Code() {
			t.Errorf("GetOperation with invalid name expected code=%d, got %d", codes.NotFound, s.Code())
		}
	}
}

func TestGetOperation_invalidMarshalledProto(t *testing.T) {
	prefixes := []string{
		"operations/google.showcase.v1alpha3.Echo/Wait",
		"operations/google.showcase.v1alpha3.Messaging/SearchBlurbs",
	}
	for _, p := range prefixes {
		name := fmt.Sprintf(
			"%s/%s",
			p,
			base64.StdEncoding.EncodeToString([]byte("BOGUS")))
		req := &lropb.GetOperationRequest{
			Name: name,
		}
		server := NewOperationsServer(nil)
		_, err := server.GetOperation(context.Background(), req)
		s, _ := status.FromError(err)
		if codes.NotFound != s.Code() {
			t.Errorf("GetOperation with invalid name expected code=%d, got %d", codes.NotFound, s.Code())
		}
	}
}

func TestCancelOperation(t *testing.T) {
	server := NewOperationsServer(nil)
	_, err := server.CancelOperation(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Errorf("CancelOperation expected code=%d, got %d", codes.Unimplemented, s.Code())
	}
}

func TestServerListOperation(t *testing.T) {
	server := NewOperationsServer(nil)
	_, err := server.ListOperations(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Errorf("ListOperations expected code=%d, got %d", codes.Unimplemented, s.Code())
	}
}

func TestServerDeleteOperation(t *testing.T) {
	server := NewOperationsServer(nil)
	_, err := server.DeleteOperation(context.Background(), nil)
	s, _ := status.FromError(err)
	if codes.Unimplemented != s.Code() {
		t.Errorf("DeleteOperations expected code=%d, got %d", codes.Unimplemented, s.Code())
	}
}
