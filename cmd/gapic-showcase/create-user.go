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

var CreateUserInput genprotopb.CreateUserRequest

var CreateUserFromFile string

func init() {
	IdentityServiceCmd.AddCommand(CreateUserCmd)

	CreateUserInput.User = new(genprotopb.User)

	CreateUserInput.User.CreateTime = new(timestamppb.Timestamp)

	CreateUserInput.User.UpdateTime = new(timestamppb.Timestamp)

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.Name, "user.name", "", "")

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.DisplayName, "user.display_name", "", "")

	CreateUserCmd.Flags().StringVar(&CreateUserInput.User.Email, "user.email", "", "")

	CreateUserCmd.Flags().Int64Var(&CreateUserInput.User.CreateTime.Seconds, "user.create_time.seconds", 0, "")

	CreateUserCmd.Flags().Int32Var(&CreateUserInput.User.CreateTime.Nanos, "user.create_time.nanos", 0, "")

	CreateUserCmd.Flags().Int64Var(&CreateUserInput.User.UpdateTime.Seconds, "user.update_time.seconds", 0, "")

	CreateUserCmd.Flags().Int32Var(&CreateUserInput.User.UpdateTime.Nanos, "user.update_time.nanos", 0, "")

	CreateUserCmd.Flags().StringVar(&CreateUserFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Creates a user.",
	Long:  "Creates a user.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateUserFromFile == "" {

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
