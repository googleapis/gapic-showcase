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

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsInput.Parent, "parent", "", "Required. The resource name of the requested room or profile whos blurbs to list.")

	ListBlurbsCmd.Flags().Int32Var(&ListBlurbsInput.PageSize, "page_size", 0, "The maximum number of blurbs to return. Server may return fewer  blurbs than requested. If unspecified, server will pick an appropriate  default.")

	ListBlurbsCmd.Flags().StringVar(&ListBlurbsInput.PageToken, "page_token", "", "The value of google.showcase.v1beta1.ListBlurbsResponse.next_page_token  returned from the previous call to  `google.showcase.v1beta1.Messaging\\ListBlurbs` method.")

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

		// get requested page
		var items []interface{}
		data := make(map[string]interface{})

		// PageSize could be an integer with a specific precision.
		// Doing standard i := 0; i < PageSize; i++ creates i as
		// an int, creating a potential type mismatch.
		for i := ListBlurbsInput.PageSize; i > 0; i-- {
			item, err := iter.Next()
			if err == iterator.Done {
				err = nil
				break
			} else if err != nil {
				return err
			}

			items = append(items, item)
		}

		data["page"] = items
		data["nextToken"] = iter.PageInfo().Token

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(data)

		return err
	},
}
