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

package services

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewOperationsServer returns a new OperationsServer for the Showcase API.
func NewOperationsServer(messagingServer MessagingServer) lropb.OperationsServer {
	return &operationsServerImpl{waiter: server.GetWaiterInstance(), messagingServer: messagingServer}
}

type operationsServerImpl struct {
	messagingServer MessagingServer
	waiter          server.Waiter
}

func (s *operationsServerImpl) GetOperation(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	if op, err := s.handleWait(in); op != nil || err != nil {
		return op, err
	}
	if op, err := s.handleSearchBlurbs(in); op != nil || err != nil {
		return op, err
	}
	return nil, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
}

func (s *operationsServerImpl) handleWait(in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	prefix := "operations/google.showcase.v1alpha3.Echo/Wait/"
	if !strings.HasPrefix(in.Name, prefix) {
		return nil, nil
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

func (s *operationsServerImpl) handleSearchBlurbs(in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	prefix := "operations/google.showcase.v1alpha3.Messaging/SearchBlurbs/"
	if !strings.HasPrefix(in.GetName(), prefix) {
		return nil, nil
	}

	req := &pb.SearchBlurbsRequest{}
	encodedBytes := strings.TrimPrefix(in.GetName(), prefix)
	waitReqBytes, err := base64.StdEncoding.DecodeString(encodedBytes)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
	}

	err = proto.Unmarshal(waitReqBytes, req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
	}

	// TODO(landrito): add some randomization here so that the search blurbs
	// operation could take multiple get calls to complete.
	listResp, err := s.messagingServer.FilteredListBlurbs(
		context.Background(),
		&pb.ListBlurbsRequest{
			Parent:    req.GetParent(),
			PageSize:  req.GetPageSize(),
			PageToken: req.GetPageToken(),
		},
		searchFilterFunc(req.GetQuery()))

	answer := &lropb.Operation{
		Name: in.GetName(),
		Done: true,
	}

	resp, _ := ptypes.MarshalAny(
		&pb.SearchBlurbsResponse{
			Blurbs:        listResp.GetBlurbs(),
			NextPageToken: listResp.GetNextPageToken(),
		})
	answer.Result = &lropb.Operation_Response{Response: resp}

	return answer, nil
}

func searchFilterFunc(query string) func(b *pb.Blurb) bool {
	return func(b *pb.Blurb) bool {
		for _, s := range strings.Fields(query) {
			if strings.Index(b.GetText(), s) >= 0 {
				return true
			}
		}
		return false
	}
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

func (s operationsServerImpl) WaitOperation(ctx context.Context, in *lropb.WaitOperationRequest) (*lropb.Operation, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.WaitOperation is unimplemented.")
}
