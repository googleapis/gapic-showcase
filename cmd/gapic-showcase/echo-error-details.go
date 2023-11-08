// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var EchoErrorDetailsInput genprotopb.EchoErrorDetailsRequest

var EchoErrorDetailsFromFile string

func init() {
	EchoServiceCmd.AddCommand(EchoErrorDetailsCmd)

	EchoErrorDetailsCmd.Flags().StringVar(&EchoErrorDetailsInput.SingleDetailText, "single_detail_text", "", "Content to return in a singular `*.error.details`...")

	EchoErrorDetailsCmd.Flags().StringSliceVar(&EchoErrorDetailsInput.MultiDetailText, "multi_detail_text", []string{}, "Content to return in a repeated `*.error.details`...")

	EchoErrorDetailsCmd.Flags().StringVar(&EchoErrorDetailsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var EchoErrorDetailsCmd = &cobra.Command{
	Use:   "echo-error-details",
	Short: "This method returns error details in a repeated...",
	Long:  "This method returns error details in a repeated 'google.protobuf.Any'  field. This method showcases handling errors thus encoded, particularly  over...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if EchoErrorDetailsFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if EchoErrorDetailsFromFile != "" {
			in, err = os.Open(EchoErrorDetailsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &EchoErrorDetailsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Echo", "EchoErrorDetails", &EchoErrorDetailsInput)
		}
		resp, err := EchoClient.EchoErrorDetails(ctx, &EchoErrorDetailsInput)
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
