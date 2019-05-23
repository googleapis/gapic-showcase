// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var DeleteTestInput genprotopb.DeleteTestRequest

var DeleteTestFromFile string

func init() {
	TestingServiceCmd.AddCommand(DeleteTestCmd)

	DeleteTestCmd.Flags().StringVar(&DeleteTestInput.Name, "name", "", "The test to be deleted.")

	DeleteTestCmd.Flags().StringVar(&DeleteTestFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteTestCmd = &cobra.Command{
	Use:   "delete-test",
	Short: "Explicitly decline to implement a test.   This...",
	Long:  "Explicitly decline to implement a test.   This removes the test from subsequent `ListTests` calls, and  attempting to do the test will error.   This...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteTestFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteTestFromFile != "" {
			in, err = os.Open(DeleteTestFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteTestInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "DeleteTest", &DeleteTestInput)
		}
		err = TestingClient.DeleteTest(ctx, &DeleteTestInput)

		return err
	},
}
