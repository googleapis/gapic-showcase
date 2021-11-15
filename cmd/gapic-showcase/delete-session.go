// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var DeleteSessionInput genprotopb.DeleteSessionRequest

var DeleteSessionFromFile string

func init() {
	TestingServiceCmd.AddCommand(DeleteSessionCmd)

	DeleteSessionCmd.Flags().StringVar(&DeleteSessionInput.Name, "name", "", "The session to be deleted.")

	DeleteSessionCmd.Flags().StringVar(&DeleteSessionFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteSessionCmd = &cobra.Command{
	Use:   "delete-session",
	Short: "Delete a test session.",
	Long:  "Delete a test session.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteSessionFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteSessionFromFile != "" {
			in, err = os.Open(DeleteSessionFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteSessionInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "DeleteSession", &DeleteSessionInput)
		}
		err = TestingClient.DeleteSession(ctx, &DeleteSessionInput)
		if err != nil {
			return err
		}

		return err
	},
}
