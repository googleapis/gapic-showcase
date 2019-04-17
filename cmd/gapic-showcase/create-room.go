// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var CreateRoomInput genprotopb.CreateRoomRequest

var CreateRoomFromFile string

func init() {
	MessagingServiceCmd.AddCommand(CreateRoomCmd)

	CreateRoomInput.Room = new(genprotopb.Room)

	CreateRoomCmd.Flags().StringVar(&CreateRoomInput.Room.Name, "room.name", "", "The resource name of the chat room.")

	CreateRoomCmd.Flags().StringVar(&CreateRoomInput.Room.DisplayName, "room.display_name", "", "Required. The human readable name of the chat room.")

	CreateRoomCmd.Flags().StringVar(&CreateRoomInput.Room.Description, "room.description", "", "The description of the chat room.")

	CreateRoomCmd.Flags().StringVar(&CreateRoomFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateRoomCmd = &cobra.Command{
	Use:   "create-room",
	Short: "Creates a room.",
	Long:  "Creates a room.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateRoomFromFile == "" {

			cmd.MarkFlagRequired("room.display_name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateRoomFromFile != "" {
			in, err = os.Open(CreateRoomFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateRoomInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "CreateRoom", &CreateRoomInput)
		}
		resp, err := MessagingClient.CreateRoom(ctx, &CreateRoomInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
