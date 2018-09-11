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
	"io"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/spf13/cobra"

	"google.golang.org/grpc"
)

func init() {
	var addr, port string
	var echoClient pb.EchoClient
	var conn *grpc.ClientConn

	initClient := func(cmd *cobra.Command, args []string) {
		// Set start listening.
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
		var err error
		conn, err = grpc.Dial(addr+port, grpc.WithInsecure())
		if err != nil {
			errLog.Fatalf("did not connect: %v", err)
		}

		// Set client
		echoClient = pb.NewEchoClient(conn)
	}

	closeConnection := func(cmd *cobra.Command, args []string) {
		conn.Close()
	}

	commands := []*cobra.Command{
		&cobra.Command{
			Use:    "echo [content to echo]",
			Short:  "Sends an echo request",
			Args:   cobra.MinimumNArgs(1),
			PreRun: initClient,
			Run: func(cmd *cobra.Command, args []string) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				req := &pb.EchoRequest{
					Response: &pb.EchoRequest_Content{Content: strings.Join(args, " ")}}
				resp, err := echoClient.Echo(ctx, req)
				if err != nil {
					errLog.Fatalf("%+v", err)
				}
				stdLog.Printf("Sent Request: %s", proto.MarshalTextString(req))
				stdLog.Printf("Got Response: %s", proto.MarshalTextString(resp))
			},
			PostRun: closeConnection,
		},
		&cobra.Command{
			Use:    "expand",
			Short:  "Starts a server-side stream using the streaming rpc 'expand'.",
			Args:   cobra.MinimumNArgs(1),
			PreRun: initClient,
			Run: func(cmd *cobra.Command, args []string) {
				// Make the request
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				req := &pb.ExpandRequest{Content: strings.Join(args, " ")}
				stream, err := echoClient.Expand(ctx, req)
				if err != nil {
					errLog.Fatalf("%+v", err)
				}
				stdLog.Printf("Sent Request: %s", proto.MarshalTextString(req))

				// Log the responses
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						return
					}
					if err != nil {
						stdLog.Printf("Error: %v", err)
						return
					}
					if resp.Content != "" {
						stdLog.Printf("Got Response: %s", proto.MarshalTextString(resp))
					}
				}
			},
			PostRun: closeConnection,
		},
		&cobra.Command{
			Use:    "collect",
			Short:  "Starts a client stream using the streaming rpc 'collect'.",
			PreRun: initClient,
			Run: func(cmd *cobra.Command, args []string) {
				// Start the stream
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				stream, _ := echoClient.Collect(ctx)

				// Create requests from user input
				stdLog.Print("Enter request content [empty line ends the stream]: ")
				for {
					var input string
					fmt.Scanln(&input)
					if input == "" {
						break
					}
					req := &pb.EchoRequest{
						Response: &pb.EchoRequest_Content{Content: input}}
					err := stream.Send(req)
					if err != nil {
						errLog.Fatalf("%+v", err)
					}
					stdLog.Printf("Sent Request: %s", proto.MarshalTextString(req))
				}

				resp, err := stream.CloseAndRecv()
				if err != nil {
					errLog.Fatalf("%+v", err)
				}
				stdLog.Printf("Got Response: %s", proto.MarshalTextString(resp))
			},
			PostRun: closeConnection,
		},
		&cobra.Command{
			Use:    "chat",
			Short:  "Starts a bidirectional stream using the streaming rpc 'chat'.",
			PreRun: initClient,
			Run: func(cmd *cobra.Command, args []string) {
				// Start the stream
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				stream, _ := echoClient.Chat(ctx)

				// Log responses
				go func() {
					for {
						resp, err := stream.Recv()
						if err == io.EOF {
							return
						}
						if err != nil {
							stdLog.Printf("Error: %v", err)
							return
						}
						if resp.Content != "" {
							stdLog.Printf("Got Response: %s", proto.MarshalTextString(resp))
						}
					}
				}()

				// Create requests from user input
				stdLog.Print("Enter request content [empty line ends the stream]: ")
				for {
					var input string
					fmt.Scanln(&input)
					if input == "" {
						break
					}
					req := &pb.EchoRequest{
						Response: &pb.EchoRequest_Content{Content: input}}
					err := stream.Send(req)
					if err != nil {
						errLog.Fatalf("%+v", err)
					}
					stdLog.Printf("Sent Request: %s", proto.MarshalTextString(req))
				}
			},
			PostRun: closeConnection,
		},
	}

	rootCmd.AddCommand(commands...)
	for _, command := range commands {
		command.Flags().StringVarP(
			&addr,
			"address",
			"a",
			"localhost",
			"The service address to make this request to")
		command.Flags().StringVarP(
			&port,
			"port",
			"p",
			":7469",
			"The port to make this request to")
	}
}
