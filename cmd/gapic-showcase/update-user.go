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

var updateUserInputUserAge int32

var updateUserInputUserHeightFeet float64

var updateUserInputUserNickname string

var updateUserInputUserEnableNotifications bool

func init() {
	IdentityServiceCmd.AddCommand(UpdateUserCmd)

	UpdateUserInput.User = new(genprotopb.User)

	UpdateUserInput.UpdateMask = new(field_maskpb.FieldMask)

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.Name, "user.name", "", "The resource name of the user.")

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.DisplayName, "user.display_name", "", "Required. The display_name of the user.")

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.Email, "user.email", "", "Required. The email address of the user.")

	UpdateUserCmd.Flags().Int32Var(&updateUserInputUserAge, "user.age", 0, "The age of the use in years.")

	UpdateUserCmd.Flags().Float64Var(&updateUserInputUserHeightFeet, "user.height_feet", 0.0, "The height of the user in feet.")

	UpdateUserCmd.Flags().StringVar(&updateUserInputUserNickname, "user.nickname", "", "The nickname of the user.   (--...")

	UpdateUserCmd.Flags().BoolVar(&updateUserInputUserEnableNotifications, "user.enable_notifications", false, "Enables the receiving of notifications. The...")

	UpdateUserCmd.Flags().StringSliceVar(&UpdateUserInput.UpdateMask.Paths, "update_mask.paths", []string{}, "The set of field mask paths.")

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

		} else {

			if cmd.Flags().Changed("user.age") {
				UpdateUserInput.User.Age = &updateUserInputUserAge
			}

			if cmd.Flags().Changed("user.height_feet") {
				UpdateUserInput.User.HeightFeet = &updateUserInputUserHeightFeet
			}

			if cmd.Flags().Changed("user.nickname") {
				UpdateUserInput.User.Nickname = &updateUserInputUserNickname
			}

			if cmd.Flags().Changed("user.enable_notifications") {
				UpdateUserInput.User.EnableNotifications = &updateUserInputUserEnableNotifications
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
