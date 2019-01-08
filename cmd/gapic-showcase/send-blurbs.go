// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"bufio"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var SendBlurbsFromFile string

func init() {
	MessagingServiceCmd.AddCommand(SendBlurbsCmd)

	SendBlurbsCmd.Flags().StringVar(&SendBlurbsFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var SendBlurbsCmd = &cobra.Command{
	Use:   "send-blurbs",
	Short: "This is a stream to create multiple blurbs. If an...",
	Long:  "This is a stream to create multiple blurbs. If an invalid blurb is  requested to be created, the stream will close with an error.",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if SendBlurbsFromFile != "" {
			in, err = os.Open(SendBlurbsFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

		}

		stream, err := MessagingClient.SendBlurbs(ctx)

		if Verbose {
			fmt.Println("Client stream open. Close with ctrl+D.")
		}

		var SendBlurbsInput genprotopb.CreateBlurbRequest
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				continue
			}
			err = jsonpb.UnmarshalString(input, &SendBlurbsInput)
			if err != nil {
				return err
			}

			err = stream.Send(&SendBlurbsInput)
			if err != nil {
				return err
			}
		}
		if err = scanner.Err(); err != nil {
			return err
		}

		resp, err := stream.CloseAndRecv()
		if err != nil {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
