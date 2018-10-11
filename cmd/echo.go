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

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/spf13/cobra"

	"google.golang.org/grpc"
)

func init() {
	var addr, port string
	var pageToken string
	var pageSize int
	var echoClient pb.EchoClient
	var conn *grpc.ClientConn

	initClient := func(cmd *cobra.Command, args []string) {
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
		var err error
		conn, err = grpc.Dial(
			addr+port,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(logClientUnary),
			grpc.WithStreamInterceptor(logClientStreaming),
		)
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
			Use:    "echo [content]",
			Short:  "Sends an echo request",
			Args:   cobra.MinimumNArgs(1),
			PreRun: initClient,
			Run: func(cmd *cobra.Command, args []string) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				req := &pb.EchoRequest{
					Response: &pb.EchoRequest_Content{Content: strings.Join(args, " ")}}
				// The response or error of this request will be handled by the
				// registered interceptors.
				echoClient.Echo(ctx, req)
			},
			PostRun: closeConnection,
		},
		&cobra.Command{
			Use:    "expand [content]",
			Short:  "Starts a server-side stream using the streaming rpc 'expand'.",
			Args:   cobra.MinimumNArgs(1),
			PreRun: initClient,
			Run: func(cmd *cobra.Command, args []string) {
				// Make the request
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				req := &pb.ExpandRequest{Content: strings.Join(args, " ")}
				stream, _ := echoClient.Expand(ctx, req)
				for {
					// The response or error of this request will be handled by the
					// registered interceptors.
					_, err := stream.Recv()
					if err == io.EOF {
						return
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
					stream.Send(req)
				}
				// The response or error of this request will be handled by the
				// registered interceptors.
				stream.CloseAndRecv()
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

				// Poll for responses
				go func() {
					for {
						// The response or error of this request will be handled by the
						// registered interceptors.
						_, err := stream.Recv()
						if err == io.EOF {
							return
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
					stream.Send(req)
				}
			},
			PostRun: closeConnection,
		},
	}

	pagedExpandCmd := &cobra.Command{
		Use:    "pagedExpand [content]",
		Short:  "Expands the given content and returns the expansion as a paged list.",
		Args:   cobra.MinimumNArgs(1),
		PreRun: initClient,
		Run: func(cmd *cobra.Command, args []string) {
			// Make the request
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			req := &pb.PagedExpandRequest{
				Content:   strings.Join(args, " "),
				PageToken: pageToken,
				PageSize:  int32(pageSize),
			}

			// The response or error of this request will be handled by the
			// registered interceptors.
			echoClient.PagedExpand(ctx, req)
		},
		PostRun: closeConnection,
	}
	pagedExpandCmd.Flags().StringVarP(
		&pageToken,
		"page_token",
		"t",
		"",
		"The page token to send with this request.")
	pagedExpandCmd.Flags().IntVarP(
		&pageSize,
		"page_size",
		"s",
		0,
		"The page size to send with this request")
	commands = append(commands, pagedExpandCmd)

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
