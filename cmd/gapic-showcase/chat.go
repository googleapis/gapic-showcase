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

var ChatFromFile string

var ChatOutFile string

func init() {
	EchoServiceCmd.AddCommand(ChatCmd)

	ChatCmd.Flags().StringVar(&ChatFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

	ChatCmd.Flags().StringVar(&ChatOutFile, "out_file", "", "Absolute path to a file to pipe output to")
	ChatCmd.MarkFlagRequired("out_file")

}

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "This method, upon receiving a request on the...",
	Long:  "This method, upon receiving a request on the stream, the same content will  be passed  back on the stream. This method showcases bidirectional ...",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ChatFromFile != "" {
			in, err = os.Open(ChatFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

		}

		stream, err := EchoClient.Chat(ctx)

		out, err := os.OpenFile(ChatOutFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}

		// start background stream receive
		go func() {
			var res *genprotopb.EchoResponse
			for {
				res, err = stream.Recv()
				if err != nil {
					return
				}

				str := res.String()
				if OutputJSON {
					str, _ = marshaler.MarshalToString(res)
				}
				fmt.Fprintln(out, str)
			}
		}()

		if Verbose {
			fmt.Println("Client stream open. Close with ctrl+D.")
		}

		var ChatInput genprotopb.EchoRequest
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				continue
			}
			err = jsonpb.UnmarshalString(input, &ChatInput)
			if err != nil {
				return err
			}

			err = stream.Send(&ChatInput)
			if err != nil {
				return err
			}
		}
		if err = scanner.Err(); err != nil {
			return err
		}

		err = stream.CloseSend()

		return err
	},
}
