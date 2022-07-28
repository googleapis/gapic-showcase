// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var GetEnumInput genprotopb.EnumRequest

var GetEnumFromFile string

func init() {
	ComplianceServiceCmd.AddCommand(GetEnumCmd)

	GetEnumCmd.Flags().BoolVar(&GetEnumInput.UnknownEnum, "unknown_enum", false, "Whether the client is requesting a new, unknown...")

	GetEnumCmd.Flags().StringVar(&GetEnumFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var GetEnumCmd = &cobra.Command{
	Use:   "get-enum",
	Short: "This method requests an enum value from the...",
	Long:  "This method requests an enum value from the server. Depending on the contents of EnumRequest, the enum value returned will be a known enum declared...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if GetEnumFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if GetEnumFromFile != "" {
			in, err = os.Open(GetEnumFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &GetEnumInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "GetEnum", &GetEnumInput)
		}
		resp, err := ComplianceClient.GetEnum(ctx, &GetEnumInput)
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
