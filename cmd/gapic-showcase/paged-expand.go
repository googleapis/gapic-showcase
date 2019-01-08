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

var PagedExpandInput genprotopb.PagedExpandRequest

var PagedExpandFromFile string

func init() {
	EchoServiceCmd.AddCommand(PagedExpandCmd)

	PagedExpandCmd.Flags().StringVar(&PagedExpandInput.Content, "content", "", "")

	PagedExpandCmd.Flags().Int32Var(&PagedExpandInput.PageSize, "page_size", 0, "")

	PagedExpandCmd.Flags().StringVar(&PagedExpandInput.PageToken, "page_token", "", "")

	PagedExpandCmd.Flags().StringVar(&PagedExpandFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var PagedExpandCmd = &cobra.Command{
	Use:   "paged-expand",
	Short: "This is similar to the Expand method but instead...",
	Long:  "This is similar to the Expand method but instead of returning a stream of  expanded words, this method returns a paged list of expanded words.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if PagedExpandFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if PagedExpandFromFile != "" {
			in, err = os.Open(PagedExpandFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &PagedExpandInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Echo", "PagedExpand", &PagedExpandInput)
		}
		iter := EchoClient.PagedExpand(ctx, &PagedExpandInput)

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
