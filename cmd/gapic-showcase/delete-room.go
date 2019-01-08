// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var DeleteRoomInput genprotopb.DeleteRoomRequest

var DeleteRoomFromFile string

func init() {
	MessagingServiceCmd.AddCommand(DeleteRoomCmd)

	DeleteRoomCmd.Flags().StringVar(&DeleteRoomInput.Name, "name", "", "")

	DeleteRoomCmd.Flags().StringVar(&DeleteRoomFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteRoomCmd = &cobra.Command{
	Use:   "delete-room",
	Short: "Deletes a room and all of its blurbs.",
	Long:  "Deletes a room and all of its blurbs.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteRoomFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteRoomFromFile != "" {
			in, err = os.Open(DeleteRoomFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteRoomInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "DeleteRoom", &DeleteRoomInput)
		}
		err = MessagingClient.DeleteRoom(ctx, &DeleteRoomInput)

		return err
	},
}
