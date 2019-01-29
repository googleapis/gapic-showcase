// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var VerifyTestInput genprotopb.VerifyTestRequest

var VerifyTestFromFile string

func init() {
	TestingServiceCmd.AddCommand(VerifyTestCmd)

	VerifyTestCmd.Flags().StringVar(&VerifyTestInput.Name, "name", "", "")

	VerifyTestCmd.Flags().BytesHexVar(&VerifyTestInput.Answer, "answer", []byte{}, "")

	VerifyTestCmd.Flags().StringVar(&VerifyTestFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var VerifyTestCmd = &cobra.Command{
	Use:   "verify-test",
	Short: "Register a response to a test.   In cases where a...",
	Long:  "Register a response to a test.   In cases where a test involves registering a final answer at the  end of the test, this method provides the means to do so.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if VerifyTestFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if VerifyTestFromFile != "" {
			in, err = os.Open(VerifyTestFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &VerifyTestInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "VerifyTest", &VerifyTestInput)
		}
		resp, err := TestingClient.VerifyTest(ctx, &VerifyTestInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
