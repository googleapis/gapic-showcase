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

var UpdateUserInput genprotopb.UpdateUserRequest

var UpdateUserFromFile string

func init() {
	IdentityServiceCmd.AddCommand(UpdateUserCmd)

	UpdateUserInput.User = new(genprotopb.User)

	UpdateUserInput.User.CreateTime = new(timestamppb.Timestamp)

	UpdateUserInput.User.UpdateTime = new(timestamppb.Timestamp)

	UpdateUserInput.UpdateMask = new(field_maskpb.FieldMask)

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.Name, "user.name", "", "")

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.DisplayName, "user.display_name", "", "")

	UpdateUserCmd.Flags().StringVar(&UpdateUserInput.User.Email, "user.email", "", "")

	UpdateUserCmd.Flags().Int64Var(&UpdateUserInput.User.CreateTime.Seconds, "user.create_time.seconds", 0, "")

	UpdateUserCmd.Flags().Int32Var(&UpdateUserInput.User.CreateTime.Nanos, "user.create_time.nanos", 0, "")

	UpdateUserCmd.Flags().Int64Var(&UpdateUserInput.User.UpdateTime.Seconds, "user.update_time.seconds", 0, "")

	UpdateUserCmd.Flags().Int32Var(&UpdateUserInput.User.UpdateTime.Nanos, "user.update_time.nanos", 0, "")

	UpdateUserCmd.Flags().StringSliceVar(&UpdateUserInput.UpdateMask.Paths, "update_mask.paths", []string{}, "")

	UpdateUserCmd.Flags().StringVar(&UpdateUserFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var UpdateUserCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Updates a user.",
	Long:  "Updates a user.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if UpdateUserFromFile == "" {

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
