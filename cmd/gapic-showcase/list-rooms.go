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

	ListRoomsCmd.Flags().Int32Var(&ListRoomsInput.PageSize, "page_size", 10, "Default is 10. The maximum number of rooms return. Server may...")

	ListRoomsCmd.Flags().StringVar(&ListRoomsInput.PageToken, "page_token", "", "The value of...")

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
