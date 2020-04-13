// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/api/iterator"

	"os"
)

var ListBlurbsInput genprotopb.ListBlurbsRequest

var ListBlurbsFromFile string

func init() {
	MessagingServiceCmd.AddCommand(ListBlurbsCmd)

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsInput.Parent, "parent", "", "Required. The resource name of the requested room or...")

	ListBlurbsCmd.Flags().Int32Var(&ListBlurbsInput.PageSize, "page_size", 10, "Default is 10. The maximum number of blurbs to return. Server...")

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsInput.PageToken, "page_token", "", "The value of...")

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ListBlurbsCmd = &cobra.Command{
	Use:   "list-blurbs",
	Short: "Lists blurbs for a specific chat room or user...",
	Long:  "Lists blurbs for a specific chat room or user profile depending on the  parent resource name.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ListBlurbsFromFile == "" {

			cmd.MarkFlagRequired("parent")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ListBlurbsFromFile != "" {
			in, err = os.Open(ListBlurbsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ListBlurbsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "ListBlurbs", &ListBlurbsInput)
		}
		iter := MessagingClient.ListBlurbs(ctx, &ListBlurbsInput)

		// populate iterator with a page
		_, err = iter.Next()
		if err != nil && err != iterator.Done {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(iter.Response)

		return err
	},
}
