// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var GetBlurbInput genprotopb.GetBlurbRequest

var GetBlurbFromFile string

func init() {
	MessagingServiceCmd.AddCommand(GetBlurbCmd)

	GetBlurbCmd.Flags().StringVar(&GetBlurbInput.Name, "name", "", "")

	GetBlurbCmd.Flags().StringVar(&GetBlurbFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetBlurbCmd = &cobra.Command{
	Use:   "get-blurb",
	Short: "Retrieves the Blurb with the given resource name.",
	Long:  "Retrieves the Blurb with the given resource name.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetBlurbFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetBlurbFromFile != "" {
			in, err = os.Open(GetBlurbFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetBlurbInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "GetBlurb", &GetBlurbInput)
		}
		resp, err := MessagingClient.GetBlurb(ctx, &GetBlurbInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
