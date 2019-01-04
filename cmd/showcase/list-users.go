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

	ListUsersCmd.Flags().Int32Var(&ListUsersInput.PageSize, "page_size", 0, "")

	ListUsersCmd.Flags().StringVar(&ListUsersInput.PageToken, "page_token", "", "")

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
