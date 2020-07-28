// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var AttemptSequenceInput genprotopb.AttemptSequenceRequest

var AttemptSequenceFromFile string

func init() {
	SequenceServiceCmd.AddCommand(AttemptSequenceCmd)

	AttemptSequenceCmd.Flags().StringVar(&AttemptSequenceInput.Name, "name", "", "Required. ")

	AttemptSequenceCmd.Flags().StringVar(&AttemptSequenceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var AttemptSequenceCmd = &cobra.Command{
	Use: "attempt-sequence",

	PreRun: func(cmd *cobra.Command, args []string) {

		if AttemptSequenceFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if AttemptSequenceFromFile != "" {
			in, err = os.Open(AttemptSequenceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &AttemptSequenceInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Sequence", "AttemptSequence", &AttemptSequenceInput)
		}
		err = SequenceClient.AttemptSequence(ctx, &AttemptSequenceInput)

		return err
	},
}
