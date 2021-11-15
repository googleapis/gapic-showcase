// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var GetSessionInput genprotopb.GetSessionRequest

var GetSessionFromFile string

func init() {
	TestingServiceCmd.AddCommand(GetSessionCmd)

	GetSessionCmd.Flags().StringVar(&GetSessionInput.Name, "name", "", "The session to be retrieved.")

	GetSessionCmd.Flags().StringVar(&GetSessionFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetSessionCmd = &cobra.Command{
	Use:   "get-session",
	Short: "Gets a testing session.",
	Long:  "Gets a testing session.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetSessionFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetSessionFromFile != "" {
			in, err = os.Open(GetSessionFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetSessionInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "GetSession", &GetSessionInput)
		}
		resp, err := TestingClient.GetSession(ctx, &GetSessionInput)
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
