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

		// get requested page
		var items []interface{}
		data := make(map[string]interface{})

		// PageSize could be an integer with a specific precision.
		// Doing standard i := 0; i < PageSize; i++ creates i as
		// an int, creating a potential type mismatch.
		for i := ListTestsInput.PageSize; i > 0; i-- {
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
