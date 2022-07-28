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

var VerifyEnumInput genprotopb.EnumResponse

var VerifyEnumFromFile string

var VerifyEnumInputContinent string

func init() {
	ComplianceServiceCmd.AddCommand(VerifyEnumCmd)

	VerifyEnumInput.Request = new(genprotopb.EnumRequest)

	VerifyEnumCmd.Flags().BoolVar(&VerifyEnumInput.Request.UnknownEnum, "request.unknown_enum", false, "Whether the client is requesting a new, unknown...")

	VerifyEnumCmd.Flags().StringVar(&VerifyEnumInputContinent, "continent", "", "The actual enum the server provided.")

	VerifyEnumCmd.Flags().StringVar(&VerifyEnumFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var VerifyEnumCmd = &cobra.Command{
	Use:   "verify-enum",
	Short: "This method is used to verify that clients can...",
	Long:  "This method is used to verify that clients can round-trip enum values, which is particularly important for unknown enum values over REST....",
	PreRun: func(cmd *cobra.Command, args []string) {

		if VerifyEnumFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if VerifyEnumFromFile != "" {
			in, err = os.Open(VerifyEnumFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &VerifyEnumInput)
			if err != nil {
				return err
			}

		} else {

			VerifyEnumInput.Continent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(VerifyEnumInputContinent)])

		}

		if Verbose {
			printVerboseInput("Compliance", "VerifyEnum", &VerifyEnumInput)
		}
		resp, err := ComplianceClient.VerifyEnum(ctx, &VerifyEnumInput)
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
