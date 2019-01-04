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

var ConnectFromFile string

var ConnectOutFile string

func init() {
	MessagingServiceCmd.AddCommand(ConnectCmd)

	ConnectCmd.Flags().StringVar(&ConnectFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

	ConnectCmd.Flags().StringVar(&ConnectOutFile, "out_file", "", "Absolute path to a file to pipe output to")
	ConnectCmd.MarkFlagRequired("out_file")

}

var ConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "This method starts a bidirectional stream that...",
	Long:  "This method starts a bidirectional stream that receives all blurbs that  are being created after the stream has started and sends requests to create  blurbs. If an invalid blurb is requested to be created, the stream will  close with an error.",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if ConnectFromFile != "" {
			in, err = os.Open(ConnectFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

		}

		stream, err := MessagingClient.Connect(ctx)

		out, err := os.OpenFile(ConnectOutFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}

		// start background stream receive
		go func() {
			var res *genprotopb.Blurb
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

		var ConnectInput genprotopb.ConnectRequest
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				continue
			}
			err = jsonpb.UnmarshalString(input, &ConnectInput)
			if err != nil {
				return err
			}

			err = stream.Send(&ConnectInput)
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
