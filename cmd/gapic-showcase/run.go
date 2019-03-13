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
	"strings"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/services"
	"github.com/spf13/cobra"
	lropb "google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	var port string
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the showcase server",
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure port is of the right form.
			if !strings.HasPrefix(port, ":") {
				port = ":" + port
			}

			// Start listening.
			lis, err := net.Listen("tcp", port)
			if err != nil {
				log.Fatalf("Showcase failed to listen on port '%s': %v", port, err)
			}
			stdLog.Printf("Showcase listening on port: %s", port)

			// Setup Server.
			logger := &loggerObserver{}
			observerRegistry := server.ShowcaseObserverRegistry()
			observerRegistry.RegisterUnaryObserver(logger)
			observerRegistry.RegisterStreamRequestObserver(logger)
			observerRegistry.RegisterStreamResponseObserver(logger)

			opts := []grpc.ServerOption{
				grpc.StreamInterceptor(observerRegistry.StreamInterceptor),
				grpc.UnaryInterceptor(observerRegistry.UnaryInterceptor),
			}
			s := grpc.NewServer(opts...)
			defer s.GracefulStop()

			// Register Services to the server.
			pb.RegisterEchoServer(s, services.NewEchoServer())
			identityServer := services.NewIdentityServer()
			pb.RegisterIdentityServer(s, identityServer)
			messagingServer := services.NewMessagingServer(identityServer)
			pb.RegisterMessagingServer(s, messagingServer)
			operationsServer := services.NewOperationsServer(messagingServer)
			pb.RegisterTestingServer(s, services.NewTestingServer(observerRegistry))
			lropb.RegisterOperationsServer(s, operationsServer)

			// Register reflection service on gRPC server.
			reflection.Register(s)
			s.Serve(lis)
		},
	}
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		":7469",
		"The port that showcase will be served on.")
}
