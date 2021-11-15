// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var DeleteBlurbInput genprotopb.DeleteBlurbRequest

var DeleteBlurbFromFile string

func init() {
	MessagingServiceCmd.AddCommand(DeleteBlurbCmd)

	DeleteBlurbCmd.Flags().StringVar(&DeleteBlurbInput.Name, "name", "", "Required. The resource name of the requested blurb.")

	DeleteBlurbCmd.Flags().StringVar(&DeleteBlurbFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteBlurbCmd = &cobra.Command{
	Use:   "delete-blurb",
	Short: "Deletes a blurb.",
	Long:  "Deletes a blurb.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteBlurbFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteBlurbFromFile != "" {
			in, err = os.Open(DeleteBlurbFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteBlurbInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "DeleteBlurb", &DeleteBlurbInput)
		}
		err = MessagingClient.DeleteBlurb(ctx, &DeleteBlurbInput)
		if err != nil {
			return err
		}

		return err
	},
}
