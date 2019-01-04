// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"io"

	"os"
)

var StreamBlurbsInput genprotopb.StreamBlurbsRequest

var StreamBlurbsFromFile string

func init() {
	MessagingServiceCmd.AddCommand(StreamBlurbsCmd)

	StreamBlurbsCmd.Flags().StringVar(&StreamBlurbsInput.Name, "name", "", "")

	StreamBlurbsCmd.Flags().StringVar(&StreamBlurbsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var StreamBlurbsCmd = &cobra.Command{
	Use:   "stream-blurbs",
	Short: "This returns a stream that emits the blurbs that...",
	Long:  "This returns a stream that emits the blurbs that are created for a  particular chat room or user profile.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if StreamBlurbsFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if StreamBlurbsFromFile != "" {
			in, err = os.Open(StreamBlurbsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &StreamBlurbsInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "StreamBlurbs", &StreamBlurbsInput)
		}
		resp, err := MessagingClient.StreamBlurbs(ctx, &StreamBlurbsInput)

		var item *genprotopb.Blurb
		for {
			item, err = resp.Recv()
			if err != nil {
				break
			}

			if Verbose {
				fmt.Print("Output: ")
			}
			printMessage(item)
		}

		if err == io.EOF {
			return nil
		}

		return err
	},
}
