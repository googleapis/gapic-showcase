// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	anypb "github.com/golang/protobuf/ptypes/any"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"io"

	"os"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"
)

var ExpandInput genprotopb.ExpandRequest

var ExpandFromFile string

var ExpandInputErrorDetails []string

func init() {
	EchoServiceCmd.AddCommand(ExpandCmd)

	ExpandInput.Error = new(statuspb.Status)

	ExpandCmd.Flags().StringVar(&ExpandInput.Content, "content", "", "")

	ExpandCmd.Flags().Int32Var(&ExpandInput.Error.Code, "error.code", 0, "")

	ExpandCmd.Flags().StringVar(&ExpandInput.Error.Message, "error.message", "", "")

	ExpandCmd.Flags().StringArrayVar(&ExpandInputErrorDetails, "error.details", []string{}, "")

	ExpandCmd.Flags().StringVar(&ExpandFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var ExpandCmd = &cobra.Command{
	Use:   "expand",
	Short: "This method split the given content into words...",
	Long:  "This method split the given content into words and will pass each word back  through the stream. This method showcases server-side streaming rpcs.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if ExpandFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ExpandFromFile != "" {
			in, err = os.Open(ExpandFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &ExpandInput)
			if err != nil {
				return err
			}

		}

		// unmarshal JSON strings into slice of structs
		for _, item := range ExpandInputErrorDetails {
			tmp := anypb.Any{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			ExpandInput.Error.Details = append(ExpandInput.Error.Details, &tmp)

		}

		if Verbose {
			printVerboseInput("Echo", "Expand", &ExpandInput)
		}
		resp, err := EchoClient.Expand(ctx, &ExpandInput)

		var item *genprotopb.EchoResponse
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
