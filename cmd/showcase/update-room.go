// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	field_maskpb "google.golang.org/genproto/protobuf/field_mask"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"

	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
)

var UpdateRoomInput genprotopb.UpdateRoomRequest

var UpdateRoomFromFile string

func init() {
	MessagingServiceCmd.AddCommand(UpdateRoomCmd)

	UpdateRoomInput.Room = new(genprotopb.Room)

	UpdateRoomInput.Room.CreateTime = new(timestamppb.Timestamp)

	UpdateRoomInput.Room.UpdateTime = new(timestamppb.Timestamp)

	UpdateRoomInput.UpdateMask = new(field_maskpb.FieldMask)

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomInput.Room.Name, "room.name", "", "")

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomInput.Room.DisplayName, "room.display_name", "", "")

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomInput.Room.Description, "room.description", "", "")

	UpdateRoomCmd.Flags().Int64Var(&UpdateRoomInput.Room.CreateTime.Seconds, "room.create_time.seconds", 0, "")

	UpdateRoomCmd.Flags().Int32Var(&UpdateRoomInput.Room.CreateTime.Nanos, "room.create_time.nanos", 0, "")

	UpdateRoomCmd.Flags().Int64Var(&UpdateRoomInput.Room.UpdateTime.Seconds, "room.update_time.seconds", 0, "")

	UpdateRoomCmd.Flags().Int32Var(&UpdateRoomInput.Room.UpdateTime.Nanos, "room.update_time.nanos", 0, "")

	UpdateRoomCmd.Flags().StringSliceVar(&UpdateRoomInput.UpdateMask.Paths, "update_mask.paths", []string{}, "")

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UpdateRoomCmd = &cobra.Command{
	Use:   "update-room",
	Short: "Updates a room.",
	Long:  "Updates a room.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UpdateRoomFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if UpdateRoomFromFile != "" {
			in, err = os.Open(UpdateRoomFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &UpdateRoomInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "UpdateRoom", &UpdateRoomInput)
		}
		resp, err := MessagingClient.UpdateRoom(ctx, &UpdateRoomInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
