// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var GetUserInput genprotopb.GetUserRequest

var GetUserFromFile string

func init() {
	IdentityServiceCmd.AddCommand(GetUserCmd)

	GetUserCmd.Flags().StringVar(&GetUserInput.Name, "name", "", "Required. The resource name of the requested user.")

	GetUserCmd.Flags().StringVar(&GetUserFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetUserCmd = &cobra.Command{
	Use:   "get-user",
	Short: "Retrieves the User with the given uri.",
	Long:  "Retrieves the User with the given uri.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetUserFromFile == "" {

			cmd.MarkFlagRequired("name")

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetUserFromFile != "" {
			in, err = os.Open(GetUserFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetUserInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Identity", "GetUser", &GetUserInput)
		}
		resp, err := IdentityClient.GetUser(ctx, &GetUserInput)
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
