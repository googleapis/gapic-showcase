// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var CreateUserInput genprotopb.CreateUserRequest

var CreateUserFromFile string

var createUserInputUserAge int32

var createUserInputUserHeightFeet float64

var createUserInputUserNickname string

var createUserInputUserEnableNotifications bool

func init() {
	IdentityServiceCmd.AddCommand(CreateUserCmd)

	CreateUserInput.User = new(genprotopb.User)

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.Name, "user.name", "", "The resource name of the user.")

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.DisplayName, "user.display_name", "", "Required. The display_name of the user.")

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.Email, "user.email", "", "Required. The email address of the user.")

	CreateUserCmd.Flags().Int32Var(&createUserInputUserAge, "user.age", 0, "The age of the user in years.")

	CreateUserCmd.Flags().Float64Var(&createUserInputUserHeightFeet, "user.height_feet", 0.0, "The height of the user in feet.")

	CreateUserCmd.Flags().StringVar(&createUserInputUserNickname, "user.nickname", "", "The nickname of the user.   (--...")

	CreateUserCmd.Flags().BoolVar(&createUserInputUserEnableNotifications, "user.enable_notifications", false, "Enables the receiving of notifications. The...")

	CreateUserCmd.Flags().StringVar(&CreateUserFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Creates a user.",
	Long:  "Creates a user.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateUserFromFile == "" {

			cmd.MarkFlagRequired("user.display_name")

			cmd.MarkFlagRequired("user.email")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateUserFromFile != "" {
			in, err = os.Open(CreateUserFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateUserInput)
			if err != nil {
				return err
			}

		} else {

			if cmd.Flags().Changed("user.age") {
				CreateUserInput.User.Age = &createUserInputUserAge
			}

			if cmd.Flags().Changed("user.height_feet") {
				CreateUserInput.User.HeightFeet = &createUserInputUserHeightFeet
			}

			if cmd.Flags().Changed("user.nickname") {
				CreateUserInput.User.Nickname = &createUserInputUserNickname
			}

			if cmd.Flags().Changed("user.enable_notifications") {
				CreateUserInput.User.EnableNotifications = &createUserInputUserEnableNotifications
			}

		}

		if Verbose {
			printVerboseInput("Identity", "CreateUser", &CreateUserInput)
		}
		resp, err := IdentityClient.CreateUser(ctx, &CreateUserInput)
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
