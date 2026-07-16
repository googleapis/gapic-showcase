// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var UploadMediaInput genprotopb.UploadMediaRequest

var UploadMediaFromFile string

func init() {
	ResumableUploadServiceCmd.AddCommand(UploadMediaCmd)

	UploadMediaCmd.Flags().StringVar(&UploadMediaInput.Name, "name", "", "")

	UploadMediaCmd.Flags().StringVar(&UploadMediaFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UploadMediaCmd = &cobra.Command{
	Use:   "upload-media",
	Short: "A method with media_upload annotation enabled.",
	Long:  "A method with media_upload annotation enabled.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UploadMediaFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if UploadMediaFromFile != "" {
			in, err = os.Open(UploadMediaFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &UploadMediaInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("ResumableUpload", "UploadMedia", &UploadMediaInput)
		}
		resp, err := ResumableUploadClient.UploadMedia(ctx, &UploadMediaInput)
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
