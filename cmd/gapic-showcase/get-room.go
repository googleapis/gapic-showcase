// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var GetRoomInput genprotopb.GetRoomRequest

var GetRoomFromFile string

func init() {
	MessagingServiceCmd.AddCommand(GetRoomCmd)

	GetRoomCmd.Flags().StringVar(&GetRoomInput.Name, "name", "", "")

	GetRoomCmd.Flags().StringVar(&GetRoomFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetRoomCmd = &cobra.Command{
	Use:   "get-room",
	Short: "Retrieves the Room with the given resource name.",
	Long:  "Retrieves the Room with the given resource name.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetRoomFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetRoomFromFile != "" {
			in, err = os.Open(GetRoomFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetRoomInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "GetRoom", &GetRoomInput)
		}
		resp, err := MessagingClient.GetRoom(ctx, &GetRoomInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
