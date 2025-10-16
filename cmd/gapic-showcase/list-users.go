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

	ListUsersCmd.Flags().Int32Var(&ListUsersInput.PageSize, "page_size", 10, "Default is 10. The maximum number of users to return. Server may...")

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

		// populate iterator with a page
		_, err = iter.Next()
		if err != nil && err != iterator.Done {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(iter.Response)

		if err == iterator.Done {
			return nil
		}

		return err
	},
}
