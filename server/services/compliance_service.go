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
	_ "embed" // for storing compliance suite data, used to verify  incoming requests
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NewComplianceServer returns a new ComplianceServer for the Showcase API.
func NewComplianceServer() pb.ComplianceServer {
	return &complianceServerImpl{waiter: server.GetWaiterInstance()}
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
	return &pb.RepeatResponse{Request: in}, nil
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

func (csi *complianceServerImpl) RepeatDataBodyPut(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatDataBodyPatch(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatWithUnknownEnum(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

func (csi *complianceServerImpl) RepeatWithUnknownOptionalEnum(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error) {
	return csi.Repeat(ctx, in)
}

// complianceSuiteBytes contains the contents of the compliance suite JSON file. This requires Go
// 1.16. Note that embedding can only be applied to global variables at package scope.
//go:embed compliance_suite.json
var complianceSuiteBytes []byte

// ComplianceSuiteInitStatus contains the status result of loading the compliance test suite
type ComplianceSuiteInitStatus int

const (
	// ComplianceSuiteUninitialized means we have not attempted to parse the compliance suite data into services.ComplianceSuite.
	ComplianceSuiteUninitialized ComplianceSuiteInitStatus = iota

	// ComplianceSuiteLoaded means we have successfully parsed the compliance suite data into services.ComplianceSuite.
	ComplianceSuiteLoaded

	// ComplianceSuiteError means we failed parsing the compliance suite data into services.ComplianceSuite.
	ComplianceSuiteError
)

var (
	// ComplianceSuite holds the protocol buffer representation of the compliance suite data.
	ComplianceSuite *pb.ComplianceSuite

	// ComplianceSuiteRequests holds all the requests in ComplianceSuite, indexed by the `name` field of the request.
	ComplianceSuiteRequests map[string]*pb.RepeatRequest

	// ComplianceSuiteStatus reports the status of loading the compliance suite data into services.ComplianceSuite.
	ComplianceSuiteStatus ComplianceSuiteInitStatus

	// ComplianceSuiteStatusMessage holds a message explaining ComplianceSuiteStatus. This is
	// typically used to provide more information in the case
	// ComplianceSuiteStatus==ComplianceSuiteError.
	ComplianceSuiteStatusMessage string
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

// indexTestingRequests creates a map by request name of the requests in suiteBytes (a
// JSON-formatted encoding of pb.ComplianceSuite), for easy retrieval later.
func indexTestingRequests(suiteBytes []byte) (err error) {
	if ComplianceSuiteStatus == ComplianceSuiteLoaded {
		return nil
	}

	ComplianceSuite = &pb.ComplianceSuite{}
	if err := protojson.Unmarshal(suiteBytes, ComplianceSuite); err != nil {
		ComplianceSuiteStatus = ComplianceSuiteError
		ComplianceSuiteStatusMessage = fmt.Sprintf("(ComplianceServiceReadError) could not read compliance suite file: %s", err)
		return fmt.Errorf(ComplianceSuiteStatusMessage)

	}

	indexedSuite, err := IndexComplianceSuite(ComplianceSuite)
	if err != nil {
		ComplianceSuiteStatus = ComplianceSuiteError
		ComplianceSuiteStatusMessage = fmt.Sprintf("(ComplianceServiceSetupError) %s", err)
		return fmt.Errorf(ComplianceSuiteStatusMessage)
	}
	ComplianceSuiteRequests = indexedSuite
	ComplianceSuiteStatus = ComplianceSuiteLoaded
	ComplianceSuiteStatusMessage = "OK"
	return nil
}

func init() {
	indexTestingRequests(complianceSuiteBytes)
}
