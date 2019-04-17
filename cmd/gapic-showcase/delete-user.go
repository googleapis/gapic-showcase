// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var DeleteUserInput genprotopb.DeleteUserRequest

var DeleteUserFromFile string

func init() {
	IdentityServiceCmd.AddCommand(DeleteUserCmd)

	DeleteUserCmd.Flags().StringVar(&DeleteUserInput.Name, "name", "", "Required. The resource name of the user to delete.")

	DeleteUserCmd.Flags().StringVar(&DeleteUserFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var DeleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Deletes a user, their profile, and all of their...",
	Long:  "Deletes a user, their profile, and all of their authored messages.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if DeleteUserFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if DeleteUserFromFile != "" {
			in, err = os.Open(DeleteUserFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &DeleteUserInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Identity", "DeleteUser", &DeleteUserInput)
		}
		err = IdentityClient.DeleteUser(ctx, &DeleteUserInput)

		return err
	},
}
