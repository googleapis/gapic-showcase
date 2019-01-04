// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var RegisterTestInput genprotopb.RegisterTestRequest

var RegisterTestFromFile string

func init() {
	TestingServiceCmd.AddCommand(RegisterTestCmd)

	RegisterTestCmd.Flags().StringVar(&RegisterTestInput.Name, "name", "", "")

	RegisterTestCmd.Flags().StringSliceVar(&RegisterTestInput.Answers, "answers", []string{}, "")

	RegisterTestCmd.Flags().StringVar(&RegisterTestFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RegisterTestCmd = &cobra.Command{
	Use:   "register-test",
	Short: "Register a response to a test.   In cases where a...",
	Long:  "Register a response to a test.   In cases where a test involves registering a final answer at the  end of the test, this method provides the means to do so.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RegisterTestFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RegisterTestFromFile != "" {
			in, err = os.Open(RegisterTestFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RegisterTestInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "RegisterTest", &RegisterTestInput)
		}
		err = TestingClient.RegisterTest(ctx, &RegisterTestInput)

		return err
	},
}
