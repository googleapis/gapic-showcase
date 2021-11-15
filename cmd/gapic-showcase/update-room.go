// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var UpdateRoomInput genprotopb.UpdateRoomRequest

var UpdateRoomFromFile string

func init() {
	MessagingServiceCmd.AddCommand(UpdateRoomCmd)

	UpdateRoomInput.Room = new(genprotopb.Room)

	UpdateRoomInput.UpdateMask = new(fieldmaskpb.FieldMask)

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomInput.Room.Name, "room.name", "", "The resource name of the chat room.")

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomInput.Room.DisplayName, "room.display_name", "", "Required. The human readable name of the chat room.")

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomInput.Room.Description, "room.description", "", "The description of the chat room.")

	UpdateRoomCmd.Flags().StringSliceVar(&UpdateRoomInput.UpdateMask.Paths, "update_mask.paths", []string{}, "The set of field mask paths.")

	UpdateRoomCmd.Flags().StringVar(&UpdateRoomFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UpdateRoomCmd = &cobra.Command{
	Use:   "update-room",
	Short: "Updates a room.",
	Long:  "Updates a room.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UpdateRoomFromFile == "" {

			cmd.MarkFlagRequired("room.display_name")

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
