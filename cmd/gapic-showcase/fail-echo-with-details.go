// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var FailEchoWithDetailsInput genprotopb.FailEchoWithDetailsRequest

var FailEchoWithDetailsFromFile string

func init() {
	EchoServiceCmd.AddCommand(FailEchoWithDetailsCmd)

	FailEchoWithDetailsCmd.Flags().StringVar(&FailEchoWithDetailsInput.Message, "message", "", "Optional message to echo back in the PoetryError....")

	FailEchoWithDetailsCmd.Flags().StringVar(&FailEchoWithDetailsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var FailEchoWithDetailsCmd = &cobra.Command{
	Use:   "fail-echo-with-details",
	Short: "This method always fails with a gRPC 'Aborted'...",
	Long:  "This method always fails with a gRPC 'Aborted' error status that contains  multiple error details.  These include one instance of each of the...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if FailEchoWithDetailsFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if FailEchoWithDetailsFromFile != "" {
			in, err = os.Open(FailEchoWithDetailsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &FailEchoWithDetailsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Echo", "FailEchoWithDetails", &FailEchoWithDetailsInput)
		}
		resp, err := EchoClient.FailEchoWithDetails(ctx, &FailEchoWithDetailsInput)
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
