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

	"github.com/spf13/cobra"
	lropb "google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/grpc"
)

func init() {
	var addr, port, name string
	var opClient lropb.OperationsClient
	var conn *grpc.ClientConn

	getOpCmd := &cobra.Command{
		Use:   "getOperation --name [name]",
		Short: "Returns an operation for the given name.",
		PreRun: func(cmd *cobra.Command, args []string) {
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
				errLog.Fatalf("Cloud not connect: %v", err)
			}

			// Set client
			opClient = lropb.NewOperationsClient(conn)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Make the request
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// The response or error of this request will be handled by the
			// registered interceptors.
			opClient.GetOperation(ctx, &lropb.GetOperationRequest{Name: name})
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			conn.Close()
		},
	}
	getOpCmd.Flags().StringVarP(
		&name,
		"name",
		"n",
		"",
		"The name of the operation to get.")
	getOpCmd.MarkFlagRequired("name")
	getOpCmd.Flags().StringVarP(
		&addr,
		"address",
		"a",
		"localhost",
		"The service address to make this request to")
	getOpCmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		":7469",
		"The port to make this request to")
	rootCmd.AddCommand(getOpCmd)
}
