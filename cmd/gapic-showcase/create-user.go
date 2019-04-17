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

func init() {
	IdentityServiceCmd.AddCommand(CreateUserCmd)

	CreateUserInput.User = new(genprotopb.User)

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.Name, "user.name", "", "The resource name of the user.")

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.DisplayName, "user.display_name", "", "Required. The display_name of the user.")

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.Email, "user.email", "", "Required. The email address of the user.")

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

		}

		if Verbose {
			printVerboseInput("Identity", "CreateUser", &CreateUserInput)
		}
		resp, err := IdentityClient.CreateUser(ctx, &CreateUserInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
