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

var CollectFromFile string

func init() {
	EchoServiceCmd.AddCommand(CollectCmd)

	CollectCmd.Flags().StringVar(&CollectFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CollectCmd = &cobra.Command{
	Use:   "collect",
	Short: "This method will collect the words given to it....",
	Long:  "This method will collect the words given to it. When the stream is closed  by the client, this method will return the a concatenation of the strings ...",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CollectFromFile != "" {
			in, err = os.Open(CollectFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

		}

		stream, err := EchoClient.Collect(ctx)

		if Verbose {
			fmt.Println("Client stream open. Close with ctrl+D.")
		}

		var CollectInput genprotopb.EchoRequest
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				continue
			}
			err = jsonpb.UnmarshalString(input, &CollectInput)
			if err != nil {
				return err
			}

			err = stream.Send(&CollectInput)
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
