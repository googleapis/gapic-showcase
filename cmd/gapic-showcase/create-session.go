// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"

	"strings"
)

var CreateSessionInput genprotopb.CreateSessionRequest

var CreateSessionFromFile string

var CreateSessionInputSessionVersion string

func init() {
	TestingServiceCmd.AddCommand(CreateSessionCmd)

	CreateSessionInput.Session = new(genprotopb.Session)

	CreateSessionCmd.Flags().StringVar(&CreateSessionInput.Session.Name, "session.name", "", "The name of the session. The ID must conform to ^[a-z]+$  If this is not provided, Showcase chooses one at random.")

	CreateSessionCmd.Flags().StringVar(&CreateSessionInputSessionVersion, "session.version", "", "Required. The version this session is using.")

	CreateSessionCmd.Flags().StringVar(&CreateSessionFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var CreateSessionCmd = &cobra.Command{
	Use:   "create-session",
	Short: "Creates a new testing session.",
	Long:  "Creates a new testing session.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if CreateSessionFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if CreateSessionFromFile != "" {
			in, err = os.Open(CreateSessionFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &CreateSessionInput)
			if err != nil {
				return err
			}

		} else {

			CreateSessionInput.Session.Version = genprotopb.Session_Version(genprotopb.Session_Version_value[strings.ToUpper(CreateSessionInputSessionVersion)])

		}

		if Verbose {
			printVerboseInput("Testing", "CreateSession", &CreateSessionInput)
		}
		resp, err := TestingClient.CreateSession(ctx, &CreateSessionInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
