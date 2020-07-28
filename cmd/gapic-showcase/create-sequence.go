// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var CreateSequenceInput genprotopb.CreateSequenceRequest

var CreateSequenceFromFile string

var CreateSequenceInputSequenceResponses []string

func init() {
	SequenceServiceCmd.AddCommand(CreateSequenceCmd)

	CreateSequenceInput.Sequence = new(genprotopb.Sequence)

	CreateSequenceCmd.Flags().StringArrayVar(&CreateSequenceInputSequenceResponses, "sequence.responses", []string{}, "Sequence of responses to return in order for each...")

	CreateSequenceCmd.Flags().StringVar(&CreateSequenceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateSequenceCmd = &cobra.Command{
	Use: "create-sequence",

	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateSequenceFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateSequenceFromFile != "" {
			in, err = os.Open(CreateSequenceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateSequenceInput)
			if err != nil {
				return err
			}

		}

		// unmarshal JSON strings into slice of structs
		for _, item := range CreateSequenceInputSequenceResponses {
			tmp := genprotopb.Sequence_Response{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			CreateSequenceInput.Sequence.Responses = append(CreateSequenceInput.Sequence.Responses, &tmp)
		}

		if Verbose {
			printVerboseInput("Sequence", "CreateSequence", &CreateSequenceInput)
		}
		resp, err := SequenceClient.CreateSequence(ctx, &CreateSequenceInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
