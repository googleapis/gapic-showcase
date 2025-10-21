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

	ListSessionsCmd.Flags().Int32Var(&ListSessionsInput.PageSize, "page_size", 10, "Default is 10. The maximum number of sessions to return per page.")

	ListSessionsCmd.Flags().StringVar(&ListSessionsInput.PageToken, "page_token", "", "The page token, for retrieving subsequent pages.")

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
