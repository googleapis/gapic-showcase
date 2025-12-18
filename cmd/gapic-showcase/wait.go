// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	anypb "google.golang.org/protobuf/types/known/anypb"

	durationpb "google.golang.org/protobuf/types/known/durationpb"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var WaitInput genprotopb.WaitRequest

var WaitFromFile string

var WaitFollow bool

var WaitPollOperation string

var WaitInputEnd string

var WaitInputEndEndTime genprotopb.WaitRequest_EndTime

var WaitInputEndTtl genprotopb.WaitRequest_Ttl

var WaitInputResponse string

var WaitInputResponseError genprotopb.WaitRequest_Error

var WaitInputResponseSuccess genprotopb.WaitRequest_Success

var WaitInputResponseErrorDetails []string

func init() {
	EchoServiceCmd.AddCommand(WaitCmd)

	WaitInputEndEndTime.EndTime = new(timestamppb.Timestamp)

	WaitInputEndTtl.Ttl = new(durationpb.Duration)

	WaitInputResponseError.Error = new(statuspb.Status)

	WaitInputResponseSuccess.Success = new(genprotopb.WaitResponse)

	WaitCmd.Flags().Int64Var(&WaitInputEndEndTime.EndTime.Seconds, "end.end_time.seconds", 0, "Represents seconds of UTC time since Unix epoch...")

	WaitCmd.Flags().Int32Var(&WaitInputEndEndTime.EndTime.Nanos, "end.end_time.nanos", 0, "Non-negative fractions of a second at nanosecond...")

	WaitCmd.Flags().Int64Var(&WaitInputEndTtl.Ttl.Seconds, "end.ttl.seconds", 0, "Signed seconds of the span of time. Must be from...")

	WaitCmd.Flags().Int32Var(&WaitInputEndTtl.Ttl.Nanos, "end.ttl.nanos", 0, "Signed fractions of a second at nanosecond...")

	WaitCmd.Flags().Int32Var(&WaitInputResponseError.Error.Code, "response.error.code", 0, "The status code, which should be an enum value of...")

	WaitCmd.Flags().StringVar(&WaitInputResponseError.Error.Message, "response.error.message", "", "A developer-facing error message, which should be...")

	WaitCmd.Flags().StringArrayVar(&WaitInputResponseErrorDetails, "response.error.details", []string{}, "A list of messages that carry the error details. ...")

	WaitCmd.Flags().StringVar(&WaitInputResponseSuccess.Success.Content, "response.success.content", "", "This content of the result.")

	WaitCmd.Flags().StringVar(&WaitInputEnd, "end", "", "Choices: end_time, ttl")

	WaitCmd.Flags().StringVar(&WaitInputResponse, "response", "", "Choices: error, success")

	WaitCmd.Flags().StringVar(&WaitFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

	WaitCmd.Flags().BoolVar(&WaitFollow, "follow", false, "Block until the long running operation completes")

	EchoServiceCmd.AddCommand(WaitPollCmd)

	WaitPollCmd.Flags().BoolVar(&WaitFollow, "follow", false, "Block until the long running operation completes")

	WaitPollCmd.Flags().StringVar(&WaitPollOperation, "operation", "", "Required. Operation name to poll for")

	WaitPollCmd.MarkFlagRequired("operation")

}

var WaitCmd = &cobra.Command{
	Use:   "wait",
	Short: "This method will wait for the requested amount of...",
	Long:  "This method will wait for the requested amount of time and then return.  This method showcases how a client handles a request timeout.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if WaitFromFile == "" {

			cmd.MarkFlagRequired("end")

			cmd.MarkFlagRequired("response")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if WaitFromFile != "" {
			in, err = os.Open(WaitFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &WaitInput)
			if err != nil {
				return err
			}

		} else {

			switch WaitInputEnd {

			case "end_time":
				WaitInput.End = &WaitInputEndEndTime

			case "ttl":
				WaitInput.End = &WaitInputEndTtl

			default:
				return fmt.Errorf("Missing oneof choice for end")
			}

			switch WaitInputResponse {

			case "error":
				WaitInput.Response = &WaitInputResponseError

			case "success":
				WaitInput.Response = &WaitInputResponseSuccess

			default:
				return fmt.Errorf("Missing oneof choice for response")
			}

		}

		// unmarshal JSON strings into slice of structs
		for _, item := range WaitInputResponseErrorDetails {
			tmp := anypb.Any{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			WaitInputResponseError.Error.Details = append(WaitInputResponseError.Error.Details, &tmp)
		}

		if Verbose {
			printVerboseInput("Echo", "Wait", &WaitInput)
		}
		resp, err := EchoClient.Wait(ctx, &WaitInput)
		if err != nil {
			return err
		}

		if !WaitFollow {
			var s interface{}
			s = resp.Name()

			if OutputJSON {
				d := make(map[string]string)
				d["operation"] = resp.Name()
				s = d
			}

			printMessage(s)
			return err
		}

		result, err := resp.Wait(ctx)
		if err != nil {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(result)

		return err
	},
}

var WaitPollCmd = &cobra.Command{
	Use:   "poll-wait",
	Short: "Poll the status of a WaitOperation by name",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		op := EchoClient.WaitOperation(WaitPollOperation)

		if WaitFollow {
			resp, err := op.Wait(ctx)
			if err != nil {
				return err
			}

			if Verbose {
				fmt.Print("Output: ")
			}
			printMessage(resp)
			return err
		}

		resp, err := op.Poll(ctx)
		if err != nil {
			return err
		} else if resp != nil {
			if Verbose {
				fmt.Print("Output: ")
			}

			printMessage(resp)
			return
		}

		if op.Done() {
			fmt.Println(fmt.Sprintf("Operation %s is done", op.Name()))
		} else {
			fmt.Println(fmt.Sprintf("Operation %s not done", op.Name()))
		}

		return err
	},
}
