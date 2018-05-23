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

package main

import (
	"log"
	"net"

  "github.com/googleapis/gapic-showcase/server"
	showcasepb "github.com/googleapis/gapic-showcase/server/genproto"
	lropb "google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	defer s.GracefulStop()

	opStore := server.NewOpertionStore()
	showcasepb.RegisterShowcaseServer(s, server.NewShowcaseServer(opStore))
	lropb.RegisterOperationsServer(s, server.NewOperationsServer(opStore))

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
