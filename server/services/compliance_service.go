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

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NewComplianceServer returns a new ComplianceServer for the Showcase API.
func NewComplianceServer() pb.ComplianceServer {
	server := &complianceServerImpl{
		waiter: server.GetWaiterInstance(),
	}

	return server
}

type complianceServerImpl struct {
	waiter server.Waiter
}

// requestMatchesExpectations returns an error iff the received request asks for server verification and its
// contents do not match a known suite testing request with the same name.
func (csi *complianceServerImpl) requestMatchesExpectation(received *pb.RepeatRequest) error {
	if !received.GetServerVerify() {
		return nil
	}
	if ComplianceSuiteStatus == ComplianceSuiteError {
		return fmt.Errorf(ComplianceSuiteStatusMessage)
	}

	name := received.GetName()
	expectedRequest, ok := ComplianceSuiteRequests[name]
	if !ok {
		return fmt.Errorf("(ComplianceSuiteRequestNotFoundError) compliance suite does not contain a request %q", name)
	}

	if diff := cmp.Diff(received.GetInfo(), expectedRequest.GetInfo(), cmp.Comparer(proto.Equal)); diff != "" {
		return fmt.Errorf("(ComplianceSuiteRequestMismatchError) contents of request %q do not match test suite", name)
	}

	return nil
}

func (csi *complianceServerImpl) Repeat(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	echoTrailers(ctx)
	if err := csi.requestMatchesExpectation(in); err != nil {
		return nil, err
	}
	return &pb.RepeatResponse{Info: in.GetInfo()}, nil
}

func (csi *complianceServerImpl) RepeatDataBody(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatDataBodyInfo(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatDataQuery(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatDataSimplePath(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatDataPathResource(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatDataPathTrailingResource(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

//go:embed compliance_suite.json
var complianceSuiteBytes []byte

// complianceSuiteStatus contains the status result of loading the compliance test suite
type complianceSuiteStatus int

const (
	ComplianceSuiteUninitialized complianceSuiteStatus = iota
	ComplianceSuiteLoaded
	ComplianceSuiteError
)

var (
	ComplianceSuiteRequests      map[string]*pb.RepeatRequest // all requests, indexed by name
	ComplianceSuite              *pb.ComplianceSuite
	ComplianceSuiteStatus        complianceSuiteStatus
	ComplianceSuiteStatusMessage string // message explaining the status
)

// IndexComplianceSuite creates a map by request name of the the requests in the
// suite, for easy retrieval later.
func IndexComplianceSuite(suite *pb.ComplianceSuite) (map[string]*pb.RepeatRequest, error) {
	indexedSuite := make(map[string]*pb.RepeatRequest)
	for _, group := range suite.GetGroup() {
		for _, requestProto := range group.GetRequests() {
			name := requestProto.GetName()
			if _, exists := indexedSuite[name]; exists {
				return nil, fmt.Errorf("multiple requests in compliance suite have name %q", name)
			}
			indexedSuite[name] = requestProto
		}
	}
	return indexedSuite, nil
}

// indexTestingRequests creates a map by request name of the the requests in the
// complianceSuiteBytes, for easy retrieval later.
func indexTestingRequests() {
	if ComplianceSuiteStatus == ComplianceSuiteLoaded {
		return
	}

	ComplianceSuite = &pb.ComplianceSuite{}
	if err := protojson.Unmarshal(complianceSuiteBytes, ComplianceSuite); err != nil {
		ComplianceSuiteStatus = ComplianceSuiteError
		ComplianceSuiteStatusMessage = fmt.Sprintf("(ComplianceServiceReadError) could not read compliance suite file: %s", err)
		return

	}

	indexedSuite, err := IndexComplianceSuite(ComplianceSuite)
	if err != nil {
		ComplianceSuiteStatus = ComplianceSuiteError
		ComplianceSuiteStatusMessage = fmt.Sprintf("(ComplianceServiceSetupError) %s", err)
		return
	}
	ComplianceSuiteRequests = indexedSuite
	ComplianceSuiteStatus = ComplianceSuiteLoaded
	ComplianceSuiteStatusMessage = "OK"
}

func init() {
	indexTestingRequests()
}
