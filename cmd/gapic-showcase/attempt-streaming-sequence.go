// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"io"

	"os"
)

var AttemptStreamingSequenceInput genprotopb.AttemptStreamingSequenceRequest

var AttemptStreamingSequenceFromFile string

func init() {
	SequenceServiceCmd.AddCommand(AttemptStreamingSequenceCmd)

	AttemptStreamingSequenceCmd.Flags().StringVar(&AttemptStreamingSequenceInput.Name, "name", "", "Required. ")

	AttemptStreamingSequenceCmd.Flags().Int32Var(&AttemptStreamingSequenceInput.FailIndex, "fail_index", 0, "")

	AttemptStreamingSequenceCmd.Flags().StringVar(&AttemptStreamingSequenceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var AttemptStreamingSequenceCmd = &cobra.Command{
	Use:   "attempt-streaming-sequence",
	Short: "Attempts a streaming sequence.",
	Long:  "Attempts a streaming sequence.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if AttemptStreamingSequenceFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if AttemptStreamingSequenceFromFile != "" {
			in, err = os.Open(AttemptStreamingSequenceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &AttemptStreamingSequenceInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Sequence", "AttemptStreamingSequence", &AttemptStreamingSequenceInput)
		}
		resp, err := SequenceClient.AttemptStreamingSequence(ctx, &AttemptStreamingSequenceInput)
		if err != nil {
			return err
		}

		var item *genprotopb.AttemptStreamingSequenceResponse
		for {
			item, err = resp.Recv()
			if err != nil {
				break
			}

			if Verbose {
				fmt.Print("Output: ")
			}
			printMessage(item)
		}

		if err == io.EOF {
			return nil
		}

		return err
	},
}
