// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	field_maskpb "google.golang.org/genproto/protobuf/field_mask"

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

func init() {
	MessagingServiceCmd.AddCommand(UpdateBlurbCmd)

	UpdateBlurbInput.Blurb = new(genprotopb.Blurb)

	UpdateBlurbInput.UpdateMask = new(field_maskpb.FieldMask)

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInput.Blurb.Name, "blurb.name", "", "The resource name of the chat room.")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInput.Blurb.User, "blurb.user", "", "Required. The resource name of the blurb's author.")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbContentText.Text, "blurb.content.text", "", "The textual content of this blurb.")

	UpdateBlurbCmd.Flags().BytesHexVar(&UpdateBlurbInputBlurbContentImage.Image, "blurb.content.image", []byte{}, "The image content of this blurb.")

	UpdateBlurbCmd.Flags().StringSliceVar(&UpdateBlurbInput.UpdateMask.Paths, "update_mask.paths", []string{}, "")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbContent, "blurb.content", "", "")

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

		}

		if Verbose {
			printVerboseInput("Messaging", "UpdateBlurb", &UpdateBlurbInput)
		}
		resp, err := MessagingClient.UpdateBlurb(ctx, &UpdateBlurbInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
