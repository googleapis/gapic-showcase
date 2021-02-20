// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var RepeatDataPathResourceInput genprotopb.RepeatRequest

var RepeatDataPathResourceFromFile string

var repeatDataPathResourceInputInfoFChildPString string

var repeatDataPathResourceInputInfoFChildPFloat float32

var repeatDataPathResourceInputInfoFChildPDouble float64

var repeatDataPathResourceInputInfoFChildPBool bool

var repeatDataPathResourceInputInfoPString string

var repeatDataPathResourceInputInfoPInt32 int32

var repeatDataPathResourceInputInfoPDouble float64

var repeatDataPathResourceInputInfoPBool bool

var repeatDataPathResourceInputInfoPChildPString string

var repeatDataPathResourceInputInfoPChildPFloat float32

var repeatDataPathResourceInputInfoPChildPDouble float64

var repeatDataPathResourceInputInfoPChildPBool bool

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataPathResourceCmd)

	RepeatDataPathResourceInput.Info = new(genprotopb.ComplianceData)

	RepeatDataPathResourceInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataPathResourceInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathResourceInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathResourceInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataPathResourceInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathResourceInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Name, "name", "", "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.FString, "info.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&RepeatDataPathResourceInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&RepeatDataPathResourceInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&RepeatDataPathResourceInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint32Var(&RepeatDataPathResourceInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint32Var(&RepeatDataPathResourceInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataPathResourceCmd.Flags().Int64Var(&RepeatDataPathResourceInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataPathResourceCmd.Flags().Int64Var(&RepeatDataPathResourceInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataPathResourceCmd.Flags().Int64Var(&RepeatDataPathResourceInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint64Var(&RepeatDataPathResourceInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint64Var(&RepeatDataPathResourceInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&RepeatDataPathResourceInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().BytesHexVar(&RepeatDataPathResourceInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&RepeatDataPathResourceInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&repeatDataPathResourceInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&repeatDataPathResourceInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&repeatDataPathResourceInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&repeatDataPathResourceInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&repeatDataPathResourceInputInfoPString, "info.p_string", "", "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&repeatDataPathResourceInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&repeatDataPathResourceInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&repeatDataPathResourceInputInfoPBool, "info.p_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&RepeatDataPathResourceInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&repeatDataPathResourceInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&repeatDataPathResourceInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&repeatDataPathResourceInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&repeatDataPathResourceInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataPathResourceCmd = &cobra.Command{
	Use:   "repeat-data-path-resource",
	Short: "Same as RepeatDataSimplePath, but with a path...",
	Long:  "Same as RepeatDataSimplePath, but with a path resource.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataPathResourceFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataPathResourceFromFile != "" {
			in, err = os.Open(RepeatDataPathResourceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataPathResourceInput)
			if err != nil {
				return err
			}

		} else {

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataPathResourceInput.Info.FChild.PString = &repeatDataPathResourceInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataPathResourceInput.Info.FChild.PFloat = &repeatDataPathResourceInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataPathResourceInput.Info.FChild.PDouble = &repeatDataPathResourceInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataPathResourceInput.Info.FChild.PBool = &repeatDataPathResourceInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataPathResourceInput.Info.PString = &repeatDataPathResourceInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataPathResourceInput.Info.PInt32 = &repeatDataPathResourceInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataPathResourceInput.Info.PDouble = &repeatDataPathResourceInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataPathResourceInput.Info.PBool = &repeatDataPathResourceInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataPathResourceInput.Info.PChild.PString = &repeatDataPathResourceInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataPathResourceInput.Info.PChild.PFloat = &repeatDataPathResourceInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataPathResourceInput.Info.PChild.PDouble = &repeatDataPathResourceInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataPathResourceInput.Info.PChild.PBool = &repeatDataPathResourceInputInfoPChildPBool
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataPathResource", &RepeatDataPathResourceInput)
		}
		resp, err := ComplianceClient.RepeatDataPathResource(ctx, &RepeatDataPathResourceInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
