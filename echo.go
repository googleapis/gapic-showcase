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

package feature_testing

import (
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	"io"
	"strings"
)

type EchoServer struct{}

func (s *EchoServer) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	err := status.ErrorProto(in.GetError())
	if err != nil {
		return nil, err
	}
	return &pb.EchoResponse{Content: in.GetContent()}, nil
}

func (s *EchoServer) Expand(in *pb.ExpandRequest, stream pb.EchoService_ExpandServer) error {
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

func (s *EchoServer) Collect(stream pb.EchoService_CollectServer) error {
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

func (s *EchoServer) Chat(stream pb.EchoService_ChatServer) error {
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
