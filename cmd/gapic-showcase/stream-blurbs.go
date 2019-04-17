// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"io"

	"os"

	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
)

var StreamBlurbsInput genprotopb.StreamBlurbsRequest

var StreamBlurbsFromFile string

func init() {
	MessagingServiceCmd.AddCommand(StreamBlurbsCmd)

	StreamBlurbsInput.ExpireTime = new(timestamppb.Timestamp)

	StreamBlurbsCmd.Flags().StringVar(&StreamBlurbsInput.Name, "name", "", "Required. The resource name of a chat room or user profile whose blurbs to stream.")

	StreamBlurbsCmd.Flags().Int64Var(&StreamBlurbsInput.ExpireTime.Seconds, "expire_time.seconds", 0, "")

	StreamBlurbsCmd.Flags().Int32Var(&StreamBlurbsInput.ExpireTime.Nanos, "expire_time.nanos", 0, "")

	StreamBlurbsCmd.Flags().StringVar(&StreamBlurbsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var StreamBlurbsCmd = &cobra.Command{
	Use:   "stream-blurbs",
	Short: "This returns a stream that emits the blurbs that...",
	Long:  "This returns a stream that emits the blurbs that are created for a  particular chat room or user profile.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if StreamBlurbsFromFile == "" {

			cmd.MarkFlagRequired("name")

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

		var item *genprotopb.StreamBlurbsResponse
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
