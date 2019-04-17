// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var CreateBlurbInput genprotopb.CreateBlurbRequest

var CreateBlurbFromFile string

var CreateBlurbInputBlurbContent string

var CreateBlurbInputBlurbContentImage genprotopb.Blurb_Image

var CreateBlurbInputBlurbContentText genprotopb.Blurb_Text

func init() {
	MessagingServiceCmd.AddCommand(CreateBlurbCmd)

	CreateBlurbInput.Blurb = new(genprotopb.Blurb)

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInput.Parent, "parent", "", "Required. The resource name of the chat room or user profile that this blurb will  be tied to.")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInput.Blurb.Name, "blurb.name", "", "The resource name of the chat room.")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInput.Blurb.User, "blurb.user", "", "Required. The resource name of the blurb's author.")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInputBlurbContentText.Text, "blurb.content.text", "", "The textual content of this blurb.")

	CreateBlurbCmd.Flags().BytesHexVar(&CreateBlurbInputBlurbContentImage.Image, "blurb.content.image", []byte{}, "The image content of this blurb.")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbInputBlurbContent, "blurb.content", "", "")

	CreateBlurbCmd.Flags().StringVar(&CreateBlurbFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateBlurbCmd = &cobra.Command{
	Use:   "create-blurb",
	Short: "Creates a blurb. If the parent is a room, the...",
	Long:  "Creates a blurb. If the parent is a room, the blurb is understood to be a  message in that room. If the parent is a profile, the blurb is understood  to be a post on the profile.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateBlurbFromFile == "" {

			cmd.MarkFlagRequired("parent")

			cmd.MarkFlagRequired("blurb.user")

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
