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

var UpdateUserInput genprotopb.UpdateUserRequest

var UpdateUserFromFile string

func init() {
	IdentityServiceCmd.AddCommand(UpdateUserCmd)

	UpdateUserInput.User = new(genprotopb.User)

	UpdateUserInput.UpdateMask = new(field_maskpb.FieldMask)

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.Name, "user.name", "", "The resource name of the user.")

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.DisplayName, "user.display_name", "", "Required. The display_name of the user.")

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.Email, "user.email", "", "Required. The email address of the user.")

	UpdateUserCmd.Flags().StringSliceVar(&UpdateUserInput.UpdateMask.Paths, "update_mask.paths", []string{}, "")

	UpdateUserCmd.Flags().StringVar(&UpdateUserFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UpdateUserCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Updates a user.",
	Long:  "Updates a user.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UpdateUserFromFile == "" {

			cmd.MarkFlagRequired("user.display_name")

			cmd.MarkFlagRequired("user.email")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if UpdateUserFromFile != "" {
			in, err = os.Open(UpdateUserFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &UpdateUserInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Identity", "UpdateUser", &UpdateUserInput)
		}
		resp, err := IdentityClient.UpdateUser(ctx, &UpdateUserInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
