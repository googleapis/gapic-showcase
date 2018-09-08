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
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/spf13/cobra"

	"google.golang.org/grpc"
)

var Addr string

// startCmd represents the start command
var echoCmd = &cobra.Command{
	Use:   "echo [content to echo]",
	Short: "Sends an echo request",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := Port
		// Set start listening.
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
		conn, err := grpc.Dial(Addr+Port, grpc.WithInsecure())
		if err != nil {
			ErrLog.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewEchoClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		req := &pb.EchoRequest{
			Response: &pb.EchoRequest_Content{Content: strings.Join(args, " ")}}
		StdLog.Printf("Request: %s", proto.MarshalTextString(req))
		resp, _ := c.Echo(ctx, req)
		StdLog.Printf("Response: %s", proto.MarshalTextString(resp))

	},
}

func init() {
	rootCmd.AddCommand(echoCmd)
	echoCmd.Flags().StringVarP(
		&Addr,
		"address",
		"a",
		"localhost",
		"The service address to make this request to")
	echoCmd.Flags().StringVarP(
		&Port,
		"port",
		"p",
		":7469",
		"The port to make this request to")
}
