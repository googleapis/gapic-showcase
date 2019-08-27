// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	anypb "github.com/golang/protobuf/ptypes/any"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"
)

var EchoInput genprotopb.EchoRequest

var EchoFromFile string

var EchoInputResponse string

var EchoInputResponseContent genprotopb.EchoRequest_Content

var EchoInputResponseError genprotopb.EchoRequest_Error

var EchoInputResponseErrorDetails []string

func init() {
	EchoServiceCmd.AddCommand(EchoCmd)

	EchoInputResponseError.Error = new(statuspb.Status)

	EchoCmd.Flags().StringVar(&EchoInputResponseContent.Content, "response.content", "", "The content to be echoed by the server.")

	EchoCmd.Flags().Int32Var(&EchoInputResponseError.Error.Code, "response.error.code", 0, "The status code, which should be an enum value of...")

	EchoCmd.Flags().StringVar(&EchoInputResponseError.Error.Message, "response.error.message", "", "A developer-facing error message, which should be...")

	EchoCmd.Flags().StringArrayVar(&EchoInputResponseErrorDetails, "response.error.details", []string{}, "A list of messages that carry the error details. ...")

	EchoCmd.Flags().StringVar(&EchoInputResponse, "response", "", "Choices: content, error")

	EchoCmd.Flags().StringVar(&EchoFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var EchoCmd = &cobra.Command{
	Use:   "echo",
	Short: "This method simply echos the request. This method...",
	Long:  "This method simply echos the request. This method is showcases unary rpcs.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if EchoFromFile == "" {

			cmd.MarkFlagRequired("response")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if EchoFromFile != "" {
			in, err = os.Open(EchoFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &EchoInput)
			if err != nil {
				return err
			}

		} else {

			switch EchoInputResponse {

			case "content":
				EchoInput.Response = &EchoInputResponseContent

			case "error":
				EchoInput.Response = &EchoInputResponseError

			default:
				return fmt.Errorf("Missing oneof choice for response")
			}

		}

		// unmarshal JSON strings into slice of structs
		for _, item := range EchoInputResponseErrorDetails {
			tmp := anypb.Any{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			EchoInputResponseError.Error.Details = append(EchoInputResponseError.Error.Details, &tmp)
		}

		if Verbose {
			printVerboseInput("Echo", "Echo", &EchoInput)
		}
		resp, err := EchoClient.Echo(ctx, &EchoInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
