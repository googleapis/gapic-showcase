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

	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
)

var WaitInput genprotopb.WaitRequest

var WaitFromFile string

var WaitFollow bool

var WaitPollOperation string

var WaitInputResponse string

var WaitInputResponseError genprotopb.WaitRequest_Error

var WaitInputResponseSuccess genprotopb.WaitRequest_Success

var WaitInputResponseErrorDetails []string

func init() {
	EchoServiceCmd.AddCommand(WaitCmd)

	WaitInput.EndTime = new(timestamppb.Timestamp)

	WaitInputResponseError.Error = new(statuspb.Status)

	WaitInputResponseSuccess.Success = new(genprotopb.WaitResponse)

	WaitCmd.Flags().Int64Var(&WaitInput.EndTime.Seconds, "end_time.seconds", 0, "")

	WaitCmd.Flags().Int32Var(&WaitInput.EndTime.Nanos, "end_time.nanos", 0, "")

	WaitCmd.Flags().Int32Var(&WaitInputResponseError.Error.Code, "response.error.code", 0, "")

	WaitCmd.Flags().StringVar(&WaitInputResponseError.Error.Message, "response.error.message", "", "")

	WaitCmd.Flags().StringArrayVar(&WaitInputResponseErrorDetails, "response.error.details", []string{}, "")

	WaitCmd.Flags().StringVar(&WaitInputResponseSuccess.Success.Content, "response.success.content", "", "")

	WaitCmd.Flags().StringVar(&WaitInputResponse, "response", "", "")

	WaitCmd.Flags().StringVar(&WaitFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

	WaitCmd.Flags().BoolVar(&WaitFollow, "follow", false, "Block until the long running operation completes")

	EchoServiceCmd.AddCommand(WaitPollCmd)

	WaitPollCmd.Flags().BoolVar(&WaitFollow, "follow", false, "Block until the long running operation completes")

	WaitPollCmd.Flags().StringVar(&WaitPollOperation, "operation", "", "Required. Operation name to poll for")

	WaitPollCmd.MarkFlagRequired("operation")

}

var WaitCmd = &cobra.Command{
	Use:   "wait",
	Short: "This method will wait the requested amount of and...",
	Long:  "This method will wait the requested amount of and then return.  This method showcases how a client handles a request timing out.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if WaitFromFile == "" {

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

		fmt.Println(fmt.Sprintf("Operation %s not done", op.Name()))

		return err
	},
}
