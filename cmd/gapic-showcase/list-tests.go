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

	ListTestsCmd.Flags().StringVar(&ListTestsInput.Parent, "parent", "", "")

	ListTestsCmd.Flags().Int32Var(&ListTestsInput.PageSize, "page_size", 0, "")

	ListTestsCmd.Flags().StringVar(&ListTestsInput.PageToken, "page_token", "", "")

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
