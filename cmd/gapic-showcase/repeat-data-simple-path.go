// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	"github.com/golang/protobuf/jsonpb"

	"os"
)

var RepeatDataSimplePathInput genprotopb.RepeatRequest

var RepeatDataSimplePathFromFile string

var repeatDataSimplePathInputInfoFChildPString string

var repeatDataSimplePathInputInfoFChildPFloat float32

var repeatDataSimplePathInputInfoFChildPDouble float64

var repeatDataSimplePathInputInfoFChildPBool bool

var repeatDataSimplePathInputInfoPString string

var repeatDataSimplePathInputInfoPInt32 int32

var repeatDataSimplePathInputInfoPDouble float64

var repeatDataSimplePathInputInfoPBool bool

var repeatDataSimplePathInputInfoPChildPString string

var repeatDataSimplePathInputInfoPChildPFloat float32

var repeatDataSimplePathInputInfoPChildPDouble float64

var repeatDataSimplePathInputInfoPChildPBool bool

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataSimplePathCmd)

	RepeatDataSimplePathInput.Info = new(genprotopb.ComplianceData)

	RepeatDataSimplePathInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataSimplePathInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataSimplePathInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataSimplePathInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataSimplePathInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataSimplePathInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Name, "name", "", "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.FString, "info.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&RepeatDataSimplePathInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&RepeatDataSimplePathInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&RepeatDataSimplePathInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint32Var(&RepeatDataSimplePathInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint32Var(&RepeatDataSimplePathInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataSimplePathCmd.Flags().Int64Var(&RepeatDataSimplePathInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataSimplePathCmd.Flags().Int64Var(&RepeatDataSimplePathInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataSimplePathCmd.Flags().Int64Var(&RepeatDataSimplePathInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint64Var(&RepeatDataSimplePathInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint64Var(&RepeatDataSimplePathInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&RepeatDataSimplePathInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().BytesHexVar(&RepeatDataSimplePathInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&RepeatDataSimplePathInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&repeatDataSimplePathInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&repeatDataSimplePathInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&repeatDataSimplePathInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&repeatDataSimplePathInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&repeatDataSimplePathInputInfoPString, "info.p_string", "", "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&repeatDataSimplePathInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&repeatDataSimplePathInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&repeatDataSimplePathInputInfoPBool, "info.p_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&RepeatDataSimplePathInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&repeatDataSimplePathInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&repeatDataSimplePathInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&repeatDataSimplePathInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&repeatDataSimplePathInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataSimplePathCmd = &cobra.Command{
	Use:   "repeat-data-simple-path",
	Short: "This method echoes the ComplianceData request....",
	Long:  "This method echoes the ComplianceData request. This method exercises  sending some parameters as 'simple' path variables (i.e., of the form ...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataSimplePathFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataSimplePathFromFile != "" {
			in, err = os.Open(RepeatDataSimplePathFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataSimplePathInput)
			if err != nil {
				return err
			}

		} else {

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataSimplePathInput.Info.FChild.PString = &repeatDataSimplePathInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataSimplePathInput.Info.FChild.PFloat = &repeatDataSimplePathInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataSimplePathInput.Info.FChild.PDouble = &repeatDataSimplePathInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataSimplePathInput.Info.FChild.PBool = &repeatDataSimplePathInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataSimplePathInput.Info.PString = &repeatDataSimplePathInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataSimplePathInput.Info.PInt32 = &repeatDataSimplePathInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataSimplePathInput.Info.PDouble = &repeatDataSimplePathInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataSimplePathInput.Info.PBool = &repeatDataSimplePathInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataSimplePathInput.Info.PChild.PString = &repeatDataSimplePathInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataSimplePathInput.Info.PChild.PFloat = &repeatDataSimplePathInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataSimplePathInput.Info.PChild.PDouble = &repeatDataSimplePathInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataSimplePathInput.Info.PChild.PBool = &repeatDataSimplePathInputInfoPChildPBool
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataSimplePath", &RepeatDataSimplePathInput)
		}
		resp, err := ComplianceClient.RepeatDataSimplePath(ctx, &RepeatDataSimplePathInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
