// Copyright 2020 Google LLC
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
	"log"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	locpb "google.golang.org/genproto/googleapis/cloud/location"

	iampb "cloud.google.com/go/iam/apiv1/iampb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

// Backend contains the various service backends that will be
// accessible via one or more transport endpoints.
type Backend struct {
	// Showcase schema
	EchoServer            pb.EchoServer
	IdentityServer        pb.IdentityServer
	MessagingServer       pb.MessagingServer
	SequenceServiceServer pb.SequenceServiceServer
	ComplianceServer      pb.ComplianceServer
	TestingServer         pb.TestingServer

	// Supporting protos
	OperationsServer lropb.OperationsServer
	LocationsServer  locpb.LocationsServer
	IAMPolicyServer  iampb.IAMPolicyServer

	// Other supporting data structures
	StdLog, ErrLog   *log.Logger
	ObserverRegistry server.GrpcObserverRegistry
}
