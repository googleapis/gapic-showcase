// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"

	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
)

var CreateRoomInput genprotopb.CreateRoomRequest

var CreateRoomFromFile string

func init() {
	MessagingServiceCmd.AddCommand(CreateRoomCmd)

	CreateRoomInput.Room = new(genprotopb.Room)

	CreateRoomInput.Room.CreateTime = new(timestamppb.Timestamp)

	CreateRoomInput.Room.UpdateTime = new(timestamppb.Timestamp)

	CreateRoomCmd.Flags().StringVar(&CreateRoomInput.Room.Name, "room.name", "", "")

	CreateRoomCmd.Flags().StringVar(&CreateRoomInput.Room.DisplayName, "room.display_name", "", "")

	CreateRoomCmd.Flags().StringVar(&CreateRoomInput.Room.Description, "room.description", "", "")

	CreateRoomCmd.Flags().Int64Var(&CreateRoomInput.Room.CreateTime.Seconds, "room.create_time.seconds", 0, "")

	CreateRoomCmd.Flags().Int32Var(&CreateRoomInput.Room.CreateTime.Nanos, "room.create_time.nanos", 0, "")

	CreateRoomCmd.Flags().Int64Var(&CreateRoomInput.Room.UpdateTime.Seconds, "room.update_time.seconds", 0, "")

	CreateRoomCmd.Flags().Int32Var(&CreateRoomInput.Room.UpdateTime.Nanos, "room.update_time.nanos", 0, "")

	CreateRoomCmd.Flags().StringVar(&CreateRoomFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateRoomCmd = &cobra.Command{
	Use:   "create-room",
	Short: "Creates a room.",
	Long:  "Creates a room.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateRoomFromFile == "" {

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
