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

var PagedExpandLegacyMappedInput genprotopb.PagedExpandRequest

var PagedExpandLegacyMappedFromFile string

func init() {
	EchoServiceCmd.AddCommand(PagedExpandLegacyMappedCmd)

	PagedExpandLegacyMappedCmd.Flags().StringVar(&PagedExpandLegacyMappedInput.Content, "content", "", "Required. The string to expand.")

	PagedExpandLegacyMappedCmd.Flags().Int32Var(&PagedExpandLegacyMappedInput.PageSize, "page_size", 10, "Default is 10. The number of words to returned in each page.")

	PagedExpandLegacyMappedCmd.Flags().StringVar(&PagedExpandLegacyMappedInput.PageToken, "page_token", "", "The position of the page to be returned.")

	PagedExpandLegacyMappedCmd.Flags().StringVar(&PagedExpandLegacyMappedFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var PagedExpandLegacyMappedCmd = &cobra.Command{
	Use:   "paged-expand-legacy-mapped",
	Short: "This method returns a map containing lists of...",
	Long:  "This method returns a map containing lists of words that appear in the input, keyed by their  initial character. The only words returned are the ones...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if PagedExpandLegacyMappedFromFile == "" {

			cmd.MarkFlagRequired("content")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if PagedExpandLegacyMappedFromFile != "" {
			in, err = os.Open(PagedExpandLegacyMappedFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &PagedExpandLegacyMappedInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Echo", "PagedExpandLegacyMapped", &PagedExpandLegacyMappedInput)
		}
		iter := EchoClient.PagedExpandLegacyMapped(ctx, &PagedExpandLegacyMappedInput)

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
