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

var AttemptStreamingSequenceInputErrorDetails []string

func init() {
	SequenceServiceCmd.AddCommand(AttemptStreamingSequenceCmd)

	AttemptStreamingSequenceInput.Error = new(statuspb.Status)

	AttemptStreamingSequenceCmd.Flags().StringVar(&AttemptStreamingSequenceInput.Name, "name", "", "Required. ")

	AttemptStreamingSequenceCmd.Flags().StringVar(&AttemptStreamingSequenceInput.Content, "content", "", "The content that will be split into words and...")

	AttemptStreamingSequenceCmd.Flags().Int32Var(&AttemptStreamingSequenceInput.Error.Code, "error.code", 0, "The status code, which should be an enum value of...")

	AttemptStreamingSequenceCmd.Flags().StringVar(&AttemptStreamingSequenceInput.Error.Message, "error.message", "", "A developer-facing error message, which should be...")

	AttemptStreamingSequenceCmd.Flags().StringArrayVar(&AttemptStreamingSequenceInputErrorDetails, "error.details", []string{}, "A list of messages that carry the error details. ...")

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
		for _, item := range AttemptStreamingSequenceInputErrorDetails {
			tmp := anypb.Any{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			AttemptStreamingSequenceInput.Error.Details = append(AttemptStreamingSequenceInput.Error.Details, &tmp)
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
