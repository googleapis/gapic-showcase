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

var ListRoomsInput genprotopb.ListRoomsRequest

var ListRoomsFromFile string

func init() {
	MessagingServiceCmd.AddCommand(ListRoomsCmd)

	ListRoomsCmd.Flags().Int32Var(&ListRoomsInput.PageSize, "page_size", 0, "The maximum number of rooms return. Server may return fewer rooms  than requested. If unspecified, server will pick an appropriate default.")

	ListRoomsCmd.Flags().StringVar(&ListRoomsInput.PageToken, "page_token", "", "The value of google.showcase.v1beta1.ListRoomsResponse.next_page_token  returned from the previous call to  `google.showcase.v1beta1.Messaging\\ListRooms` method.")

	ListRoomsCmd.Flags().StringVar(&ListRoomsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ListRoomsCmd = &cobra.Command{
	Use:   "list-rooms",
	Short: "Lists all chat rooms.",
	Long:  "Lists all chat rooms.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ListRoomsFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ListRoomsFromFile != "" {
			in, err = os.Open(ListRoomsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ListRoomsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "ListRooms", &ListRoomsInput)
		}
		iter := MessagingClient.ListRooms(ctx, &ListRoomsInput)

		// get requested page
		var items []interface{}
		data := make(map[string]interface{})

		// PageSize could be an integer with a specific precision.
		// Doing standard i := 0; i < PageSize; i++ creates i as
		// an int, creating a potential type mismatch.
		for i := ListRoomsInput.PageSize; i > 0; i-- {
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
