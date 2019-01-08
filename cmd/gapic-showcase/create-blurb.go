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

var CreateBlurbInput genprotopb.CreateBlurbRequest

var CreateBlurbFromFile string

var CreateBlurbInputBlurbContent string

var CreateBlurbInputBlurbContentImage genprotopb.Blurb_Image

var CreateBlurbInputBlurbContentText genprotopb.Blurb_Text

func init() {
	MessagingServiceCmd.AddCommand(CreateBlurbCmd)

	CreateBlurbInput.Blurb = new(genprotopb.Blurb)

	CreateBlurbInput.Blurb.CreateTime = new(timestamppb.Timestamp)

	CreateBlurbInput.Blurb.UpdateTime = new(timestamppb.Timestamp)

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInput.Parent, "parent", "", "")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInput.Blurb.Name, "blurb.name", "", "")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInput.Blurb.User, "blurb.user", "", "")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInputBlurbContentText.Text, "blurb.content.text", "", "")

	CreateBlurbCmd.Flags().BytesHexVar(&CreateBlurbInputBlurbContentImage.Image, "blurb.content.image", []byte{}, "")

	CreateBlurbCmd.Flags().Int64Var(&CreateBlurbInput.Blurb.CreateTime.Seconds, "blurb.create_time.seconds", 0, "")

	CreateBlurbCmd.Flags().Int32Var(&CreateBlurbInput.Blurb.CreateTime.Nanos, "blurb.create_time.nanos", 0, "")

	CreateBlurbCmd.Flags().Int64Var(&CreateBlurbInput.Blurb.UpdateTime.Seconds, "blurb.update_time.seconds", 0, "")

	CreateBlurbCmd.Flags().Int32Var(&CreateBlurbInput.Blurb.UpdateTime.Nanos, "blurb.update_time.nanos", 0, "")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInputBlurbContent, "blurb.content", "", "")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateBlurbCmd = &cobra.Command{
	Use:   "create-blurb",
	Short: "Creates a blurb. If the parent is a room, the...",
	Long:  "Creates a blurb. If the parent is a room, the blurb is understood to be a  message in that room. If the parent is a profile, the blurb is understood  to be a post on the profile.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateBlurbFromFile == "" {

			cmd.MarkFlagRequired("blurb.content")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateBlurbFromFile != "" {
			in, err = os.Open(CreateBlurbFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateBlurbInput)
			if err != nil {
				return err
			}

		} else {

			switch CreateBlurbInputBlurbContent {

			case "image":
				CreateBlurbInput.Blurb.Content = &CreateBlurbInputBlurbContentImage

			case "text":
				CreateBlurbInput.Blurb.Content = &CreateBlurbInputBlurbContentText

			default:
				return fmt.Errorf("Missing oneof choice for blurb.content")
			}

		}

		if Verbose {
			printVerboseInput("Messaging", "CreateBlurb", &CreateBlurbInput)
		}
		resp, err := MessagingClient.CreateBlurb(ctx, &CreateBlurbInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
