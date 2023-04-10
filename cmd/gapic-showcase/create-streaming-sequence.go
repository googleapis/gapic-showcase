// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var CreateStreamingSequenceInput genprotopb.CreateStreamingSequenceRequest

var CreateStreamingSequenceFromFile string

var CreateStreamingSequenceInputStreamingSequenceResponses []string

func init() {
	SequenceServiceCmd.AddCommand(CreateStreamingSequenceCmd)

	CreateStreamingSequenceInput.StreamingSequence = new(genprotopb.StreamingSequence)

	CreateStreamingSequenceCmd.Flags().StringVar(&CreateStreamingSequenceInput.StreamingSequence.Content, "StreamingSequence.content", "", "The Content that the stream will send")

	CreateStreamingSequenceCmd.Flags().StringArrayVar(&CreateStreamingSequenceInputStreamingSequenceResponses, "StreamingSequence.responses", []string{}, "Sequence of responses to return in order for each...")

	CreateStreamingSequenceCmd.Flags().StringVar(&CreateStreamingSequenceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateStreamingSequenceCmd = &cobra.Command{
	Use:   "create-streaming-sequence",
	Short: "Creates a sequence.",
	Long:  "Creates a sequence.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateStreamingSequenceFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateStreamingSequenceFromFile != "" {
			in, err = os.Open(CreateStreamingSequenceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateStreamingSequenceInput)
			if err != nil {
				return err
			}

		}

		// unmarshal JSON strings into slice of structs
		for _, item := range CreateStreamingSequenceInputStreamingSequenceResponses {
			tmp := genprotopb.StreamingSequence_Response{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			CreateStreamingSequenceInput.StreamingSequence.Responses = append(CreateStreamingSequenceInput.StreamingSequence.Responses, &tmp)
		}

		if Verbose {
			printVerboseInput("Sequence", "CreateStreamingSequence", &CreateStreamingSequenceInput)
		}
		resp, err := SequenceClient.CreateStreamingSequence(ctx, &CreateStreamingSequenceInput)
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
