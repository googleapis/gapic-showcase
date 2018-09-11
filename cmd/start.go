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

package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/spf13/cobra"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func init() {
	var port string
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the showcase server",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(port)
		},
	}
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		":7469",
		"The port that showcase will be served on.")
}

func startServer(port string) {
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
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(logStreamRequests),
		grpc.UnaryInterceptor(logUnaryRequests),
	}
	s := grpc.NewServer(opts...)
	defer s.GracefulStop()
	pb.RegisterEchoServer(s, server.NewEchoServer())

	// Register reflection service on gRPC server.
	reflection.Register(s)
	s.Serve(lis)
}

func logUnaryRequests(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	stdLog.Printf("Received Unary Request for Method: %s\n", info.FullMethod)
	stdLog.Printf("    Request:  %+v\n", req)
	resp, err := handler(ctx, req)
	if err == nil {
		stdLog.Printf("    Returning Response: %+v\n", resp)
	} else {
		stdLog.Printf("    Returning Error: %+v\n", err)
	}
	stdLog.Println("")
	return resp, err
}

type loggingServerStream struct {
	s    grpc.ServerStream
	info *grpc.StreamServerInfo
}

func (s *loggingServerStream) SetHeader(m metadata.MD) error {
	return s.s.SetHeader(m)
}

func (s *loggingServerStream) SendHeader(m metadata.MD) error {
	return s.s.SendHeader(m)
}

func (s *loggingServerStream) SetTrailer(m metadata.MD) {
	s.s.SetTrailer(m)
}

func (s *loggingServerStream) Context() context.Context {
	return s.s.Context()
}

func (s *loggingServerStream) SendMsg(m interface{}) error {
	stdLog.Printf("%s Stream for Method: %s\n", s.streamType(), s.info.FullMethod)
	stdLog.Printf("    Sending Message:  %+v\n", m)
	stdLog.Println("")

	return s.s.SendMsg(m)
}

func (s *loggingServerStream) RecvMsg(m interface{}) error {
	err := s.s.RecvMsg(m)
	if fmt.Sprintf("%v", m) != "" {
		stdLog.Printf("%s Stream for Method: %s\n", s.streamType(), s.info.FullMethod)
		stdLog.Printf("    Recieving Message:  %v\n", m)
		stdLog.Println("")
	}

	return err
}

func (s *loggingServerStream) streamType() string {
	if s.info.IsClientStream && s.info.IsServerStream {
		return "Bi-directional"
	} else if s.info.IsClientStream {
		return "Client"
	}
	return "Server"
}

func logStreamRequests(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	loggingStream := &loggingServerStream{s: ss, info: info}
	return handler(srv, loggingStream)
}
