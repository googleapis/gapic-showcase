// Copyright 2021 Google LLC
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
	_ "embed"
	"fmt"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

//go:embed compliance_suite.json
var complianceSuiteBytes []byte

// NewComplianceServer returns a new ComplianceServer for the Showcase API.
func NewComplianceServer() pb.ComplianceServer {
	// read embedded data set
	// process into protos?
	// store data internally (global? see how done in messaging)

	fmt.Printf("loaded compliance suite: size %d, starts with %q\n", len(complianceSuiteBytes), string(complianceSuiteBytes[0:100]))

	return &complianceServerImpl{
		waiter:               server.GetWaiterInstance(),
		complianceSuiteBytes: complianceSuiteBytes,
	}
}

type complianceServerImpl struct {
	waiter               server.Waiter
	complianceSuiteBytes []byte
}

func (s *complianceServerImpl) Repeat(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	echoTrailers(ctx)
	return &pb.RepeatResponse{Info: in.GetInfo()}, nil
}

func (s *complianceServerImpl) RepeatDataBody(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return s.Repeat(ctx, in)
}

func (s *complianceServerImpl) RepeatDataBodyInfo(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return s.Repeat(ctx, in)
}

func (s *complianceServerImpl) RepeatDataQuery(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return s.Repeat(ctx, in)
}

func (s *complianceServerImpl) RepeatDataSimplePath(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return s.Repeat(ctx, in)
}

func (s *complianceServerImpl) RepeatDataPathResource(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return s.Repeat(ctx, in)
}

func (s *complianceServerImpl) RepeatDataPathTrailingResource(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return s.Repeat(ctx, in)
}
