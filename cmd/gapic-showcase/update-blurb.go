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

var UpdateBlurbInput genprotopb.UpdateBlurbRequest

var UpdateBlurbFromFile string

var UpdateBlurbInputBlurbContent string

var UpdateBlurbInputBlurbContentImage genprotopb.Blurb_Image

var UpdateBlurbInputBlurbContentText genprotopb.Blurb_Text

var UpdateBlurbInputBlurbLegacyId string

var UpdateBlurbInputBlurbLegacyIdLegacyRoomId genprotopb.Blurb_LegacyRoomId

var UpdateBlurbInputBlurbLegacyIdLegacyUserId genprotopb.Blurb_LegacyUserId

func init() {
	MessagingServiceCmd.AddCommand(UpdateBlurbCmd)

	UpdateBlurbInput.Blurb = new(genprotopb.Blurb)

	UpdateBlurbInput.UpdateMask = new(fieldmaskpb.FieldMask)

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInput.Blurb.Name, "blurb.name", "", "The resource name of the chat room.")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInput.Blurb.User, "blurb.user", "", "Required. The resource name of the blurb's author.")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbContentText.Text, "blurb.content.text", "", "The textual content of this blurb.")

	UpdateBlurbCmd.Flags().BytesHexVar(&UpdateBlurbInputBlurbContentImage.Image, "blurb.content.image", []byte{}, "The image content of this blurb.")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbLegacyIdLegacyRoomId.LegacyRoomId, "blurb.legacy_id.legacy_room_id", "", "The legacy id of the room. This field is used to...")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbLegacyIdLegacyUserId.LegacyUserId, "blurb.legacy_id.legacy_user_id", "", "The legacy id of the user. This field is used to...")

	UpdateBlurbCmd.Flags().StringSliceVar(&UpdateBlurbInput.UpdateMask.Paths, "update_mask.paths", []string{}, "The set of field mask paths.")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbContent, "blurb.content", "", "Choices: text, image")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbLegacyId, "blurb.legacy_id", "", "Choices: legacy_room_id, legacy_user_id")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UpdateBlurbCmd = &cobra.Command{
	Use:   "update-blurb",
	Short: "Updates a blurb.",
	Long:  "Updates a blurb.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UpdateBlurbFromFile == "" {

			cmd.MarkFlagRequired("blurb.user")

			cmd.MarkFlagRequired("blurb.content")

			cmd.MarkFlagRequired("blurb.legacy_id")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if UpdateBlurbFromFile != "" {
			in, err = os.Open(UpdateBlurbFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &UpdateBlurbInput)
			if err != nil {
				return err
			}

		} else {

			switch UpdateBlurbInputBlurbContent {

			case "image":
				UpdateBlurbInput.Blurb.Content = &UpdateBlurbInputBlurbContentImage

			case "text":
				UpdateBlurbInput.Blurb.Content = &UpdateBlurbInputBlurbContentText

			default:
				return fmt.Errorf("Missing oneof choice for blurb.content")
			}

			switch UpdateBlurbInputBlurbLegacyId {

			case "legacy_room_id":
				UpdateBlurbInput.Blurb.LegacyId = &UpdateBlurbInputBlurbLegacyIdLegacyRoomId

			case "legacy_user_id":
				UpdateBlurbInput.Blurb.LegacyId = &UpdateBlurbInputBlurbLegacyIdLegacyUserId

			default:
				return fmt.Errorf("Missing oneof choice for blurb.legacy_id")
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "UpdateBlurb", &UpdateBlurbInput)
		}
		resp, err := MessagingClient.UpdateBlurb(ctx, &UpdateBlurbInput)
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
