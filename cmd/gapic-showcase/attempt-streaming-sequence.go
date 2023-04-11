// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	anypb "google.golang.org/protobuf/types/known/anypb"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"io"

	"os"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"
)

var AttemptStreamingSequenceInput genprotopb.AttemptStreamingSequenceRequest

var AttemptStreamingSequenceFromFile string

var AttemptStreamingSequenceInputAttemptStatusDetails []string

func init() {
	SequenceServiceCmd.AddCommand(AttemptStreamingSequenceCmd)

	AttemptStreamingSequenceInput.AttemptStatus = new(statuspb.Status)

	AttemptStreamingSequenceCmd.Flags().StringVar(&AttemptStreamingSequenceInput.Name, "name", "", "Required. ")

	AttemptStreamingSequenceCmd.Flags().Int32Var(&AttemptStreamingSequenceInput.AttemptStatus.Code, "attempt_status.code", 0, "The status code, which should be an enum value of...")

	AttemptStreamingSequenceCmd.Flags().StringVar(&AttemptStreamingSequenceInput.AttemptStatus.Message, "attempt_status.message", "", "A developer-facing error message, which should be...")

	AttemptStreamingSequenceCmd.Flags().StringArrayVar(&AttemptStreamingSequenceInputAttemptStatusDetails, "attempt_status.details", []string{}, "A list of messages that carry the error details. ...")

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

		// unmarshal JSON strings into slice of structs
		for _, item := range AttemptStreamingSequenceInputAttemptStatusDetails {
			tmp := anypb.Any{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			AttemptStreamingSequenceInput.AttemptStatus.Details = append(AttemptStreamingSequenceInput.AttemptStatus.Details, &tmp)
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
