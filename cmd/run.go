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
	info *grpc.StreamServerInfo

	grpc.ServerStream
}

func (s *loggingServerStream) SendMsg(m interface{}) error {
	stdLog.Printf("%s Stream for Method: %s\n", s.streamType(), s.info.FullMethod)
	stdLog.Printf("    Sending Message:  %+v\n", m)
	stdLog.Println("")

	return s.ServerStream.SendMsg(m)
}

func (s *loggingServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
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
	loggingStream := &loggingServerStream{info, ss}
	return handler(srv, loggingStream)
}
