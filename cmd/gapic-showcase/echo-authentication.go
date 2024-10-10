// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var EchoAuthenticationInput genprotopb.EchoAuthenticationRequest

func init() {
	EchoServiceCmd.AddCommand(EchoAuthenticationCmd)

}

var EchoAuthenticationCmd = &cobra.Command{
	Use:   "echo-authentication",
	Short: "This method returns authentication details from...",
	Long:  "This method returns authentication details from the incoming request, including things like  Authorization headers for clients to verify that the...",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if Verbose {
			printVerboseInput("Echo", "EchoAuthentication", &EchoAuthenticationInput)
		}
		resp, err := EchoClient.EchoAuthentication(ctx, &EchoAuthenticationInput)
		if err != nil {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
