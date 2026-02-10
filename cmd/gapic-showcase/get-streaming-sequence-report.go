// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var GetStreamingSequenceReportInput genprotopb.GetStreamingSequenceReportRequest

var GetStreamingSequenceReportFromFile string

func init() {
	SequenceServiceCmd.AddCommand(GetStreamingSequenceReportCmd)

	GetStreamingSequenceReportCmd.Flags().StringVar(&GetStreamingSequenceReportInput.Name, "name", "", "Required. ")

	GetStreamingSequenceReportCmd.Flags().StringVar(&GetStreamingSequenceReportFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetStreamingSequenceReportCmd = &cobra.Command{
	Use:   "get-streaming-sequence-report",
	Short: "Retrieves a sequence report which can be used to...",
	Long:  "Retrieves a sequence report which can be used to retrieve information  about a sequences of responses in a server streaming call.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetStreamingSequenceReportFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetStreamingSequenceReportFromFile != "" {
			in, err = os.Open(GetStreamingSequenceReportFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetStreamingSequenceReportInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Sequence", "GetStreamingSequenceReport", &GetStreamingSequenceReportInput)
		}
		resp, err := SequenceClient.GetStreamingSequenceReport(ctx, &GetStreamingSequenceReportInput)
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
