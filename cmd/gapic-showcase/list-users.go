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

var ListUsersInput genprotopb.ListUsersRequest

var ListUsersFromFile string

func init() {
	IdentityServiceCmd.AddCommand(ListUsersCmd)

	ListUsersCmd.Flags().Int32Var(&ListUsersInput.PageSize, "page_size", 0, "The maximum number of users to return. Server may...")

	ListUsersCmd.Flags().StringVar(&ListUsersInput.PageToken, "page_token", "", "The value of...")

	ListUsersCmd.Flags().StringVar(&ListUsersFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ListUsersCmd = &cobra.Command{
	Use:   "list-users",
	Short: "Lists all users.",
	Long:  "Lists all users.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ListUsersFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ListUsersFromFile != "" {
			in, err = os.Open(ListUsersFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ListUsersInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Identity", "ListUsers", &ListUsersInput)
		}
		iter := IdentityClient.ListUsers(ctx, &ListUsersInput)

		// get requested page
		var items []interface{}
		data := make(map[string]interface{})

		// PageSize could be an integer with a specific precision.
		// Doing standard i := 0; i < PageSize; i++ creates i as
		// an int, creating a potential type mismatch.
		for i := ListUsersInput.PageSize; i > 0; i-- {
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
