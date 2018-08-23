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

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

// NewTestingServer returns a new TestingServer for the Showcase API.
func NewTestingServer() pb.TestingServer {
	return &testingServerImpl{}
}

type testingServerImpl struct {
}

func (s *testingServerImpl) ReportSession(ctx context.Context, in *pb.ReportSessionRequest) (*pb.ReportSessionResponse, error) {
	return GetSessionSingleton().GetReport()
}

func (s *testingServerImpl) DeleteTest(ctx context.Context, in *pb.DeleteTestRequest) (*empty.Empty, error) {
	err := GetSessionSingleton().DeleteTest(in.Name)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *testingServerImpl) RegisterTest(ctx context.Context, in *pb.RegisterTestRequest) (*empty.Empty, error) {
	return nil, nil
}
