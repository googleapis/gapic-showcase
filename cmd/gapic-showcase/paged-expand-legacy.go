// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var PagedExpandLegacyInput genprotopb.PagedExpandLegacyRequest

var PagedExpandLegacyFromFile string

func init() {
	EchoServiceCmd.AddCommand(PagedExpandLegacyCmd)

	PagedExpandLegacyCmd.Flags().StringVar(&PagedExpandLegacyInput.Content, "content", "", "Required. The string to expand.")

	PagedExpandLegacyCmd.Flags().Int32Var(&PagedExpandLegacyInput.MaxResults, "max_results", 0, "The number of words to returned in each page. ...")

	PagedExpandLegacyCmd.Flags().StringVar(&PagedExpandLegacyInput.PageToken, "page_token", "", "The position of the page to be returned.")

	PagedExpandLegacyCmd.Flags().StringVar(&PagedExpandLegacyFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var PagedExpandLegacyCmd = &cobra.Command{
	Use:   "paged-expand-legacy",
	Short: "This is similar to the PagedExpand except that it...",
	Long:  "This is similar to the PagedExpand except that it uses  max_results instead of page_size, as some legacy APIs still  do. New APIs should NOT use this...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if PagedExpandLegacyFromFile == "" {

			cmd.MarkFlagRequired("content")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if PagedExpandLegacyFromFile != "" {
			in, err = os.Open(PagedExpandLegacyFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &PagedExpandLegacyInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Echo", "PagedExpandLegacy", &PagedExpandLegacyInput)
		}
		resp, err := EchoClient.PagedExpandLegacy(ctx, &PagedExpandLegacyInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
