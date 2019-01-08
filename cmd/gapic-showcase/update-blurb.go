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

var UpdateBlurbInput genprotopb.UpdateBlurbRequest

var UpdateBlurbFromFile string

var UpdateBlurbInputBlurbContent string

var UpdateBlurbInputBlurbContentImage genprotopb.Blurb_Image

var UpdateBlurbInputBlurbContentText genprotopb.Blurb_Text

func init() {
	MessagingServiceCmd.AddCommand(UpdateBlurbCmd)

	UpdateBlurbInput.Blurb = new(genprotopb.Blurb)

	UpdateBlurbInput.Blurb.CreateTime = new(timestamppb.Timestamp)

	UpdateBlurbInput.Blurb.UpdateTime = new(timestamppb.Timestamp)

	UpdateBlurbInput.UpdateMask = new(field_maskpb.FieldMask)

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInput.Blurb.Name, "blurb.name", "", "")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInput.Blurb.User, "blurb.user", "", "")

	UpdateBlurbCmd.Flags().StringVar(&UpdateBlurbInputBlurbContentText.Text, "blurb.content.text", "", "")

	UpdateBlurbCmd.Flags().BytesHexVar(&UpdateBlurbInputBlurbContentImage.Image, "blurb.content.image", []byte{}, "")

	UpdateBlurbCmd.Flags().Int64Var(&UpdateBlurbInput.Blurb.CreateTime.Seconds, "blurb.create_time.seconds", 0, "")

	UpdateBlurbCmd.Flags().Int32Var(&UpdateBlurbInput.Blurb.CreateTime.Nanos, "blurb.create_time.nanos", 0, "")

	UpdateBlurbCmd.Flags().Int64Var(&UpdateBlurbInput.Blurb.UpdateTime.Seconds, "blurb.update_time.seconds", 0, "")

	UpdateBlurbCmd.Flags().Int32Var(&UpdateBlurbInput.Blurb.UpdateTime.Nanos, "blurb.update_time.nanos", 0, "")

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
