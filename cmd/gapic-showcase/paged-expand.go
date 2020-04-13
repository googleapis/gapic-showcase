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

	PagedExpandCmd.Flags().StringVar(&PagedExpandInput.Content, "content", "", "Required. The string to expand.")

	PagedExpandCmd.Flags().Int32Var(&PagedExpandInput.PageSize, "page_size", 10, "Default is 10. The amount of words to returned in each page.")

	PagedExpandCmd.Flags().StringVar(&PagedExpandInput.PageToken, "page_token", "", "The position of the page to be returned.")

	PagedExpandCmd.Flags().StringVar(&PagedExpandFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var PagedExpandCmd = &cobra.Command{
	Use:   "paged-expand",
	Short: "This is similar to the Expand method but instead...",
	Long:  "This is similar to the Expand method but instead of returning a stream of  expanded words, this method returns a paged list of expanded words.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if PagedExpandFromFile == "" {

			cmd.MarkFlagRequired("content")

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
