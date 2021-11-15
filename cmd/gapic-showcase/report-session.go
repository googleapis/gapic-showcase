// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var ReportSessionInput genprotopb.ReportSessionRequest

var ReportSessionFromFile string

func init() {
	TestingServiceCmd.AddCommand(ReportSessionCmd)

	ReportSessionCmd.Flags().StringVar(&ReportSessionInput.Name, "name", "", "The session to be reported on.")

	ReportSessionCmd.Flags().StringVar(&ReportSessionFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ReportSessionCmd = &cobra.Command{
	Use:   "report-session",
	Short: "Report on the status of a session.  This...",
	Long:  "Report on the status of a session.  This generates a report detailing which tests have been completed,  and an overall rollup.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ReportSessionFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ReportSessionFromFile != "" {
			in, err = os.Open(ReportSessionFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ReportSessionInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "ReportSession", &ReportSessionInput)
		}
		resp, err := TestingClient.ReportSession(ctx, &ReportSessionInput)
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
