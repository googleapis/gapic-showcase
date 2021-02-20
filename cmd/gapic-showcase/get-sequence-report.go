// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var GetSequenceReportInput genprotopb.GetSequenceReportRequest

var GetSequenceReportFromFile string

func init() {
	SequenceServiceCmd.AddCommand(GetSequenceReportCmd)

	GetSequenceReportCmd.Flags().StringVar(&GetSequenceReportInput.Name, "name", "", "Required. ")

	GetSequenceReportCmd.Flags().StringVar(&GetSequenceReportFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetSequenceReportCmd = &cobra.Command{
	Use:   "get-sequence-report",
	Short: "Retrieves a sequence.",
	Long:  "Retrieves a sequence.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetSequenceReportFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetSequenceReportFromFile != "" {
			in, err = os.Open(GetSequenceReportFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetSequenceReportInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Sequence", "GetSequenceReport", &GetSequenceReportInput)
		}
		resp, err := SequenceClient.GetSequenceReport(ctx, &GetSequenceReportInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
