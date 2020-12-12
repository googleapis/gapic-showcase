// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	anypb "github.com/golang/protobuf/ptypes/any"

	durationpb "github.com/golang/protobuf/ptypes/duration"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"
)

var BlockInput genprotopb.BlockRequest

var BlockFromFile string

var BlockInputResponse string

var BlockInputResponseError genprotopb.BlockRequest_Error

var BlockInputResponseSuccess genprotopb.BlockRequest_Success

var BlockInputResponseErrorDetails []string

func init() {
	EchoServiceCmd.AddCommand(BlockCmd)

	BlockInput.ResponseDelay = new(durationpb.Duration)

	BlockInputResponseError.Error = new(statuspb.Status)

	BlockInputResponseSuccess.Success = new(genprotopb.BlockResponse)

	BlockCmd.Flags().Int64Var(&BlockInput.ResponseDelay.Seconds, "response_delay.seconds", 0, "Signed seconds of the span of time. Must be from...")

	BlockCmd.Flags().Int32Var(&BlockInput.ResponseDelay.Nanos, "response_delay.nanos", 0, "Signed fractions of a second at nanosecond...")

	BlockCmd.Flags().Int32Var(&BlockInputResponseError.Error.Code, "response.error.code", 0, "The status code, which should be an enum value of...")

	BlockCmd.Flags().StringVar(&BlockInputResponseError.Error.Message, "response.error.message", "", "A developer-facing error message, which should be...")

	BlockCmd.Flags().StringArrayVar(&BlockInputResponseErrorDetails, "response.error.details", []string{}, "A list of messages that carry the error details. ...")

	BlockCmd.Flags().StringVar(&BlockInputResponseSuccess.Success.Content, "response.success.content", "", "This content can contain anything, the server...")

	BlockCmd.Flags().StringVar(&BlockInputResponse, "response", "", "Choices: error, success")

	BlockCmd.Flags().StringVar(&BlockFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var BlockCmd = &cobra.Command{
	Use:   "block",
	Short: "This method will block (wait) for the requested...",
	Long:  "This method will block (wait) for the requested amount of time  and then return the response or error.  This method showcases how a client handles...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if BlockFromFile == "" {

			cmd.MarkFlagRequired("response")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if BlockFromFile != "" {
			in, err = os.Open(BlockFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &BlockInput)
			if err != nil {
				return err
			}

		} else {

			switch BlockInputResponse {

			case "error":
				BlockInput.Response = &BlockInputResponseError

			case "success":
				BlockInput.Response = &BlockInputResponseSuccess

			default:
				return fmt.Errorf("Missing oneof choice for response")
			}

		}

		// unmarshal JSON strings into slice of structs
		for _, item := range BlockInputResponseErrorDetails {
			tmp := anypb.Any{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			BlockInputResponseError.Error.Details = append(BlockInputResponseError.Error.Details, &tmp)
		}

		if Verbose {
			printVerboseInput("Echo", "Block", &BlockInput)
		}
		resp, err := EchoClient.Block(ctx, &BlockInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
