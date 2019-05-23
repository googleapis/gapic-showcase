// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var SearchBlurbsInput genprotopb.SearchBlurbsRequest

var SearchBlurbsFromFile string

var SearchBlurbsFollow bool

var SearchBlurbsPollOperation string

func init() {
	MessagingServiceCmd.AddCommand(SearchBlurbsCmd)

	SearchBlurbsCmd.Flags().StringVar(&SearchBlurbsInput.Query, "query", "", "Required. The query used to search for blurbs containing to...")

	SearchBlurbsCmd.Flags().StringVar(&SearchBlurbsInput.Parent, "parent", "", "The rooms or profiles to search. If unset,...")

	SearchBlurbsCmd.Flags().Int32Var(&SearchBlurbsInput.PageSize, "page_size", 0, "The maximum number of blurbs return. Server may...")

	SearchBlurbsCmd.Flags().StringVar(&SearchBlurbsInput.PageToken, "page_token", "", "The value of ...")

	SearchBlurbsCmd.Flags().StringVar(&SearchBlurbsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

	SearchBlurbsCmd.Flags().BoolVar(&SearchBlurbsFollow, "follow", false, "Block until the long running operation completes")

	MessagingServiceCmd.AddCommand(SearchBlurbsPollCmd)

	SearchBlurbsPollCmd.Flags().BoolVar(&SearchBlurbsFollow, "follow", false, "Block until the long running operation completes")

	SearchBlurbsPollCmd.Flags().StringVar(&SearchBlurbsPollOperation, "operation", "", "Required. Operation name to poll for")

	SearchBlurbsPollCmd.MarkFlagRequired("operation")

}

var SearchBlurbsCmd = &cobra.Command{
	Use:   "search-blurbs",
	Short: "This method searches through all blurbs across...",
	Long:  "This method searches through all blurbs across all rooms and profiles  for blurbs containing to words found in the query. Only posts that  contain an...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if SearchBlurbsFromFile == "" {

			cmd.MarkFlagRequired("query")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if SearchBlurbsFromFile != "" {
			in, err = os.Open(SearchBlurbsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &SearchBlurbsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "SearchBlurbs", &SearchBlurbsInput)
		}
		resp, err := MessagingClient.SearchBlurbs(ctx, &SearchBlurbsInput)

		if !SearchBlurbsFollow {
			var s interface{}
			s = resp.Name()

			if OutputJSON {
				d := make(map[string]string)
				d["operation"] = resp.Name()
				s = d
			}

			printMessage(s)
			return err
		}

		result, err := resp.Wait(ctx)
		if err != nil {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(result)

		return err
	},
}

var SearchBlurbsPollCmd = &cobra.Command{
	Use:   "poll-search-blurbs",
	Short: "Poll the status of a SearchBlurbsOperation by name",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		op := MessagingClient.SearchBlurbsOperation(SearchBlurbsPollOperation)

		if SearchBlurbsFollow {
			resp, err := op.Wait(ctx)
			if err != nil {
				return err
			}

			if Verbose {
				fmt.Print("Output: ")
			}
			printMessage(resp)
			return err
		}

		resp, err := op.Poll(ctx)
		if err != nil {
			return err
		} else if resp != nil {
			if Verbose {
				fmt.Print("Output: ")
			}

			printMessage(resp)
			return
		}

		fmt.Println(fmt.Sprintf("Operation %s not done", op.Name()))

		return err
	},
}
