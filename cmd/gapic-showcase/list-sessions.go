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

var ListSessionsInput genprotopb.ListSessionsRequest

var ListSessionsFromFile string

func init() {
	TestingServiceCmd.AddCommand(ListSessionsCmd)

	ListSessionsCmd.Flags().Int32Var(&ListSessionsInput.PageSize, "page_size", 0, "")

	ListSessionsCmd.Flags().StringVar(&ListSessionsInput.PageToken, "page_token", "", "")

	ListSessionsCmd.Flags().StringVar(&ListSessionsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ListSessionsCmd = &cobra.Command{
	Use:   "list-sessions",
	Short: "Lists the current test sessions.",
	Long:  "Lists the current test sessions.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ListSessionsFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ListSessionsFromFile != "" {
			in, err = os.Open(ListSessionsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ListSessionsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "ListSessions", &ListSessionsInput)
		}
		iter := TestingClient.ListSessions(ctx, &ListSessionsInput)

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
