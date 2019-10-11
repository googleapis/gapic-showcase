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
	"sync"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/services"
	fallback "github.com/googleapis/grpc-fallback-go/server"
	"github.com/spf13/cobra"
	lropb "google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	var port string
	var fallbackPort string
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the showcase server",
		Run: func(cmd *cobra.Command, args []string) {
			RunShowcase(port, fallbackPort).Wait()
		},
	}
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		":7469",
		"The port that showcase will be served on.")
	runCmd.Flags().StringVarP(
		&fallbackPort,
		"fallback-port",
		"f",
		":1337",
		"The port that the fallback-proxy will be served on.")
}

// RunShowcase sets up and starts the showcase and fallback servers and returns pointers to
// them. They can be shutdown by showcaseServers.Shutdown() or showcaseServers.Wait().
func RunShowcase(port string, fallbackPort string) (showcaseServers *ShowcaseServers) {
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

	// Register Services to the server.
	pb.RegisterEchoServer(s, services.NewEchoServer())
	identityServer := services.NewIdentityServer()
	pb.RegisterIdentityServer(s, identityServer)
	messagingServer := services.NewMessagingServer(identityServer)
	pb.RegisterMessagingServer(s, messagingServer)
	operationsServer := services.NewOperationsServer(messagingServer)
	pb.RegisterTestingServer(s, services.NewTestingServer(observerRegistry))
	lropb.RegisterOperationsServer(s, operationsServer)

	var fb *fallback.FallbackServer
	if len(fallbackPort) > 0 {
		fb = fallback.NewServer(fallbackPort, "localhost"+port)
		fb.StartBackground()
	}

	// Register reflection service on gRPC server.
	reflection.Register(s)

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		s.Serve(lis)
		wait.Done()
	}()
	return &ShowcaseServers{s: s, fb: fb, wait: &wait}
}

// ShowcaseServers encapsulates information on running showcase servers, allowing for them to be
// shutdown immediately or when they stop serving.
type ShowcaseServers struct {
	s    *grpc.Server
	fb   *fallback.FallbackServer
	wait *sync.WaitGroup
}

// Shutdown will immediately start a graceful shutdown of the servers.
func (srv *ShowcaseServers) Shutdown() {
	if srv.fb != nil {
		srv.fb.Shutdown() // unfortunately, this always results in an en error log
	}
	if srv.s != nil {
		srv.s.GracefulStop()
	}
}

// Wait will wait until the servers stop serving and then call Shutdown() to shut them down
// gracefully.
func (srv *ShowcaseServers) Wait() {
	if srv.wait != nil {
		srv.wait.Wait()
	}
	srv.Shutdown()
}
