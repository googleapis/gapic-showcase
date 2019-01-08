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

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsInput.Parent, "parent", "", "")

	ListBlurbsCmd.Flags().Int32Var(&ListBlurbsInput.PageSize, "page_size", 0, "")

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsInput.PageToken, "page_token", "", "")

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ListBlurbsCmd = &cobra.Command{
	Use:   "list-blurbs",
	Short: "Lists blurbs for a specific chat room or user...",
	Long:  "Lists blurbs for a specific chat room or user profile depending on the  parent resource name.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ListBlurbsFromFile == "" {

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

		// get requested page
		page, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				fmt.Println("No more results")
				return nil
			}

			return err
		}

		data := make(map[string]interface{})
		data["page"] = page

		//get next page token
		_, err = iter.Next()
		if err != nil && err != iterator.Done {
			return err
		}
		data["nextToken"] = iter.PageInfo().Token

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(data)

		return err
	},
}
