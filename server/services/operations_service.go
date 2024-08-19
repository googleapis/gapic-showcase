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
	"math/rand"
	"strconv"
	"strings"
	"time"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
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
	prefix := "operations/google.showcase.v1beta1.Echo/Wait/"
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
	prefix := "operations/google.showcase.v1beta1.Messaging/SearchBlurbs/"
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

// CancelOperation returns a successful response if the resource name is not blank
func (s operationsServerImpl) CancelOperation(ctx context.Context, in *lropb.CancelOperationRequest) (*empty.Empty, error) {
	if in.Name == "" {
		return nil, status.Error(codes.NotFound, "cannot cancel operation without a name.")
	}
	return &empty.Empty{}, nil
}

// ListOperations returns a fixed response matching the given PageSize if the resource name is not blank
func (s operationsServerImpl) ListOperations(ctx context.Context, in *lropb.ListOperationsRequest) (*lropb.ListOperationsResponse, error) {
	if in.Name == "" {
		return nil, status.Error(codes.NotFound, "cannot list operation without a name.")
	}

	var operations []*lropb.Operation
	if in.PageSize > 0 {
		for i := 1; i <= int(in.PageSize); i++ {
			var result *lropb.Operation_Response
			if i%2 == 0 {
				result = &lropb.Operation_Response{}
			}
			operations = append(operations, &lropb.Operation{
				Name:   "the/thing/" + strconv.Itoa(i),
				Done:   result != nil,
				Result: result,
			})
		}
	} else {
		operations = append(operations, &lropb.Operation{
			Name: "a/pending/thing",
			Done: false,
		})
	}

	return &lropb.ListOperationsResponse{
		Operations: operations,
	}, nil
}

// DeleteOperation returns a successful response if the resource name is not blank
func (s operationsServerImpl) DeleteOperation(ctx context.Context, in *lropb.DeleteOperationRequest) (*empty.Empty, error) {
	if in.Name == "" {
		return nil, status.Error(codes.NotFound, "cannot delete operation without a name.")
	}
	return &empty.Empty{}, nil
}

// WaitOperation randomly waits and returns an operation with the same name
func (s operationsServerImpl) WaitOperation(ctx context.Context, in *lropb.WaitOperationRequest) (*lropb.Operation, error) {
	if in.Name == "" {
		return nil, status.Error(codes.NotFound, "cannot wait on a operation without a name.")
	}

	num := rand.Intn(500)
	time.Sleep(time.Duration(num) * time.Millisecond)

	var result *lropb.Operation_Response
	if num%2 == 0 {
		result = &lropb.Operation_Response{}
	}
	return &lropb.Operation{
		Name:   in.Name,
		Done:   result != nil,
		Result: result,
	}, nil
}
