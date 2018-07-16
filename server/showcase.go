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
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type uniqID struct {
	i int64
}

func (u *uniqID) id() int64 {
	return atomic.AddInt64(&u.i, 1) - 1
}

// NewShowcaseServer returns a new ShowcaseServer for the Showcase API.
func NewShowcaseServer(opStore OperationStore) pb.ShowcaseServer {
	return &showcaseServerImpl{
		operationStore: opStore,
		sleepF:         time.Sleep,
	}
}

type showcaseServerImpl struct {
	uid uniqID

	mu             sync.Mutex
	retryStore     map[string][]*statuspb.Status
	operationStore OperationStore
	sleepF         func(time.Duration)
}

func (s *showcaseServerImpl) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	err := status.ErrorProto(in.GetError())
	if err != nil {
		return nil, err
	}
	return &pb.EchoResponse{Content: in.GetContent()}, nil
}

func (s *showcaseServerImpl) Expand(in *pb.ExpandRequest, stream pb.Showcase_ExpandServer) error {
	for _, word := range strings.Fields(in.GetContent()) {
		err := stream.Send(&pb.EchoResponse{Content: word})
		if err != nil {
			return err
		}
	}
	if in.GetError() != nil {
		return status.ErrorProto(in.GetError())
	}
	return nil
}

func (s *showcaseServerImpl) Collect(stream pb.Showcase_CollectServer) error {
	var resp []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.EchoResponse{Content: strings.Join(resp, " ")})
		}
		if err != nil {
			return err
		}
		s := status.ErrorProto(req.GetError())
		if s != nil {
			return s
		}
		if req.GetContent() != "" {
			resp = append(resp, req.GetContent())
		}
	}
}

func (s *showcaseServerImpl) Chat(stream pb.Showcase_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s := status.ErrorProto(req.GetError())
		if s != nil {
			return s
		}
		stream.Send(&pb.EchoResponse{Content: req.GetContent()})
	}
}

func (s *showcaseServerImpl) Timeout(ctx context.Context, in *pb.TimeoutRequest) (*pb.TimeoutResponse, error) {
	d, _ := ptypes.Duration(in.GetResponseDelay())
	s.sleepF(d)
	if in.GetError() != nil {
		return nil, status.ErrorProto(in.GetError())
	}
	return in.GetSuccess(), nil
}

func (s *showcaseServerImpl) SetupRetry(ctx context.Context, in *pb.SetupRetryRequest) (*pb.RetryId, error) {
	if in.GetResponses() == nil || len(in.GetResponses()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "A list of responses must be specified.")
	}
	id := fmt.Sprintf("retry-test-%d", s.uid.id())

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.retryStore == nil {
		s.retryStore = map[string][]*statuspb.Status{}
	}
	s.retryStore[id] = in.GetResponses()
	return &pb.RetryId{Id: id}, nil
}

func (s *showcaseServerImpl) Retry(ctx context.Context, in *pb.RetryId) (*empty.Empty, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "An Id must be specified.")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	resps, ok := s.retryStore[in.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Retry with Id: %s was not found.", in.GetId())
	}
	resp, resps := resps[0], resps[1:]
	if status.FromProto(resp).Code() == codes.OK {
		delete(s.retryStore, in.GetId())
		return &empty.Empty{}, nil
	}
	if len(resps) == 0 {
		delete(s.retryStore, in.GetId())
	} else {
		s.retryStore[in.GetId()] = resps
	}
	return nil, status.ErrorProto(resp)
}

func (s *showcaseServerImpl) Longrunning(ctx context.Context, in *pb.LongrunningRequest) (*lropb.Operation, error) {
	return s.operationStore.RegisterOp(in)
}

func (s *showcaseServerImpl) Pagination(ctx context.Context, in *pb.PaginationRequest) (*pb.PaginationResponse, error) {
	if in.GetPageSize() < 0 || in.GetPageSizeOverride() < 0 {
		return nil, status.Error(codes.InvalidArgument, "The page size provided must not be negative.")
	}

	if in.GetMaxResponse() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "The maximum response provided must be positive.")
	}

	start := int32(0)
	if in.GetPageToken() != "" {
		token, err := strconv.Atoi(in.GetPageToken())
		token32 := int32(token)
		if err != nil || token32 < 0 || token32 > in.GetMaxResponse() {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid page token: %s. Token must be within the range [0, request.MaxResponse]", in.GetPageToken())
		}
		start = token32
	}

	actualSize := in.GetPageSize()
	if in.GetPageSizeOverride() > 0 {
		actualSize = in.GetPageSizeOverride()
	}

	end := start + actualSize
	if actualSize == 0 {
		end = in.GetMaxResponse()
	}
	if end > in.GetMaxResponse() {
		end = in.GetMaxResponse()
	}

	nextToken := ""
	if end < in.GetMaxResponse() {
		nextToken = strconv.Itoa(int(end))
	}

	page := []int32{}
	for i := start; i < end; i++ {
		page = append(page, i)
	}

	return &pb.PaginationResponse{
		Responses:     page,
		NextPageToken: nextToken,
	}, nil
}

func (s *showcaseServerImpl) ParameterFlattening(ctx context.Context, in *pb.ParameterFlatteningMessage) (*pb.ParameterFlatteningMessage, error) {
	return in, nil
}

func (s *showcaseServerImpl) ResourceName(ctx context.Context, in *pb.ResourceNameMessage) (*pb.ResourceNameMessage, error) {
	return in, nil
}

func (s *showcaseServerImpl) ResourceSet(ctx context.Context, in *pb.ResourceSetMessage) (*pb.ResourceSetMessage, error) {
	return in, nil
}
