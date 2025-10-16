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

var ListTestsInput genprotopb.ListTestsRequest

var ListTestsFromFile string

func init() {
	TestingServiceCmd.AddCommand(ListTestsCmd)

	ListTestsCmd.Flags().StringVar(&ListTestsInput.Parent, "parent", "", "The session.")

	ListTestsCmd.Flags().Int32Var(&ListTestsInput.PageSize, "page_size", 10, "Default is 10. The maximum number of tests to return per page.")

	ListTestsCmd.Flags().StringVar(&ListTestsInput.PageToken, "page_token", "", "The page token, for retrieving subsequent pages.")

	ListTestsCmd.Flags().StringVar(&ListTestsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ListTestsCmd = &cobra.Command{
	Use:   "list-tests",
	Short: "List the tests of a sessesion.",
	Long:  "List the tests of a sessesion.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ListTestsFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ListTestsFromFile != "" {
			in, err = os.Open(ListTestsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ListTestsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Testing", "ListTests", &ListTestsInput)
		}
		iter := TestingClient.ListTests(ctx, &ListTestsInput)

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
