// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	anypb "google.golang.org/protobuf/types/known/anypb"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"

	"strings"
)

var EchoInput genprotopb.EchoRequest

var EchoFromFile string

var EchoInputResponse string

var EchoInputResponseContent genprotopb.EchoRequest_Content

var EchoInputResponseError genprotopb.EchoRequest_Error

var EchoInputResponseErrorDetails []string

var EchoInputSeverity string

func init() {
	EchoServiceCmd.AddCommand(EchoCmd)

	EchoInputResponseError.Error = new(statuspb.Status)

	EchoCmd.Flags().StringVar(&EchoInputResponseContent.Content, "response.content", "", "The content to be echoed by the server.")

	EchoCmd.Flags().Int32Var(&EchoInputResponseError.Error.Code, "response.error.code", 0, "The status code, which should be an enum value of...")

	EchoCmd.Flags().StringVar(&EchoInputResponseError.Error.Message, "response.error.message", "", "A developer-facing error message, which should be...")

	EchoCmd.Flags().StringArrayVar(&EchoInputResponseErrorDetails, "response.error.details", []string{}, "A list of messages that carry the error details. ...")

	EchoCmd.Flags().StringVar(&EchoInputSeverity, "severity", "", "The severity to be echoed by the server.")

	EchoCmd.Flags().StringVar(&EchoInput.Header, "header", "", "Optional. This field can be set to test the...")

	EchoCmd.Flags().StringVar(&EchoInput.OtherHeader, "other_header", "", "Optional. This field can be set to test the...")

	EchoCmd.Flags().StringVar(&EchoInputResponse, "response", "", "Choices: content, error")

	EchoCmd.Flags().StringVar(&EchoFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var EchoCmd = &cobra.Command{
	Use:   "echo",
	Short: "This method simply echoes the request. This...",
	Long:  "This method simply echoes the request. This method showcases unary RPCs.",
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

			EchoInput.Severity = genprotopb.Severity(genprotopb.Severity_value[strings.ToUpper(EchoInputSeverity)])

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
