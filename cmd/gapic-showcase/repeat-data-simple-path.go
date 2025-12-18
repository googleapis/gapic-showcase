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

var RepeatDataSimplePathInput genprotopb.RepeatRequest

var RepeatDataSimplePathFromFile string

var RepeatDataSimplePathInputInfoFKingdom string

var RepeatDataSimplePathInputInfoFChildFContinent string

var repeatDataSimplePathInputInfoFChildPString string

var repeatDataSimplePathInputInfoFChildPFloat float32

var repeatDataSimplePathInputInfoFChildPDouble float64

var repeatDataSimplePathInputInfoFChildPBool bool

var RepeatDataSimplePathInputInfoFChildPContinent string

var repeatDataSimplePathInputInfoPString string

var repeatDataSimplePathInputInfoPInt32 int32

var repeatDataSimplePathInputInfoPSint32 int32

var repeatDataSimplePathInputInfoPSfixed32 int32

var repeatDataSimplePathInputInfoPUint32 uint32

var repeatDataSimplePathInputInfoPFixed32 uint32

var repeatDataSimplePathInputInfoPInt64 int64

var repeatDataSimplePathInputInfoPSint64 int64

var repeatDataSimplePathInputInfoPSfixed64 int64

var repeatDataSimplePathInputInfoPUint64 uint64

var repeatDataSimplePathInputInfoPFixed64 uint64

var repeatDataSimplePathInputInfoPFloat float32

var repeatDataSimplePathInputInfoPDouble float64

var repeatDataSimplePathInputInfoPBool bool

var RepeatDataSimplePathInputInfoPKingdom string

var RepeatDataSimplePathInputInfoPChildFContinent string

var repeatDataSimplePathInputInfoPChildPString string

var repeatDataSimplePathInputInfoPChildPFloat float32

var repeatDataSimplePathInputInfoPChildPDouble float64

var repeatDataSimplePathInputInfoPChildPBool bool

var RepeatDataSimplePathInputInfoPChildPContinent string

var repeatDataSimplePathInputIntendedBindingUri string

var repeatDataSimplePathInputPInt32 int32

var repeatDataSimplePathInputPInt64 int64

var repeatDataSimplePathInputPDouble float64

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

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&RepeatDataSimplePathInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&repeatDataSimplePathInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&repeatDataSimplePathInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&repeatDataSimplePathInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&repeatDataSimplePathInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&repeatDataSimplePathInputInfoPString, "info.p_string", "", "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&repeatDataSimplePathInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&repeatDataSimplePathInputInfoPSint32, "info.p_sint32", 0, "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&repeatDataSimplePathInputInfoPSfixed32, "info.p_sfixed32", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint32Var(&repeatDataSimplePathInputInfoPUint32, "info.p_uint32", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint32Var(&repeatDataSimplePathInputInfoPFixed32, "info.p_fixed32", 0, "")

	RepeatDataSimplePathCmd.Flags().Int64Var(&repeatDataSimplePathInputInfoPInt64, "info.p_int64", 0, "")

	RepeatDataSimplePathCmd.Flags().Int64Var(&repeatDataSimplePathInputInfoPSint64, "info.p_sint64", 0, "")

	RepeatDataSimplePathCmd.Flags().Int64Var(&repeatDataSimplePathInputInfoPSfixed64, "info.p_sfixed64", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint64Var(&repeatDataSimplePathInputInfoPUint64, "info.p_uint64", 0, "")

	RepeatDataSimplePathCmd.Flags().Uint64Var(&repeatDataSimplePathInputInfoPFixed64, "info.p_fixed64", 0, "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&repeatDataSimplePathInputInfoPFloat, "info.p_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&repeatDataSimplePathInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&repeatDataSimplePathInputInfoPBool, "info.p_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&RepeatDataSimplePathInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&repeatDataSimplePathInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float32Var(&repeatDataSimplePathInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&repeatDataSimplePathInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&repeatDataSimplePathInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataSimplePathCmd.Flags().StringVar(&RepeatDataSimplePathInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataSimplePathCmd.Flags().BoolVar(&RepeatDataSimplePathInput.ServerVerify, "server_verify", false, "If true, the server will verify that the received...")

	RepeatDataSimplePathCmd.Flags().StringVar(&repeatDataSimplePathInputIntendedBindingUri, "intended_binding_uri", "", "The URI template this request is expected to be...")

	RepeatDataSimplePathCmd.Flags().Int32Var(&RepeatDataSimplePathInput.FInt32, "f_int32", 0, "Some top level fields, to test that these are...")

	RepeatDataSimplePathCmd.Flags().Int64Var(&RepeatDataSimplePathInput.FInt64, "f_int64", 0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&RepeatDataSimplePathInput.FDouble, "f_double", 0.0, "")

	RepeatDataSimplePathCmd.Flags().Int32Var(&repeatDataSimplePathInputPInt32, "p_int32", 0, "")

	RepeatDataSimplePathCmd.Flags().Int64Var(&repeatDataSimplePathInputPInt64, "p_int64", 0, "")

	RepeatDataSimplePathCmd.Flags().Float64Var(&repeatDataSimplePathInputPDouble, "p_double", 0.0, "")

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

			RepeatDataSimplePathInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataSimplePathInputInfoFKingdom)])

			RepeatDataSimplePathInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataSimplePathInputInfoFChildFContinent)])

			RepeatDataSimplePathInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataSimplePathInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataSimplePathInputInfoPKingdom)])
				RepeatDataSimplePathInput.Info.PKingdom = &e
			}

			RepeatDataSimplePathInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataSimplePathInputInfoPChildFContinent)])

			RepeatDataSimplePathInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataSimplePathInputInfoPChildPContinent)])

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

			if cmd.Flags().Changed("info.p_sint32") {
				RepeatDataSimplePathInput.Info.PSint32 = &repeatDataSimplePathInputInfoPSint32
			}

			if cmd.Flags().Changed("info.p_sfixed32") {
				RepeatDataSimplePathInput.Info.PSfixed32 = &repeatDataSimplePathInputInfoPSfixed32
			}

			if cmd.Flags().Changed("info.p_uint32") {
				RepeatDataSimplePathInput.Info.PUint32 = &repeatDataSimplePathInputInfoPUint32
			}

			if cmd.Flags().Changed("info.p_fixed32") {
				RepeatDataSimplePathInput.Info.PFixed32 = &repeatDataSimplePathInputInfoPFixed32
			}

			if cmd.Flags().Changed("info.p_int64") {
				RepeatDataSimplePathInput.Info.PInt64 = &repeatDataSimplePathInputInfoPInt64
			}

			if cmd.Flags().Changed("info.p_sint64") {
				RepeatDataSimplePathInput.Info.PSint64 = &repeatDataSimplePathInputInfoPSint64
			}

			if cmd.Flags().Changed("info.p_sfixed64") {
				RepeatDataSimplePathInput.Info.PSfixed64 = &repeatDataSimplePathInputInfoPSfixed64
			}

			if cmd.Flags().Changed("info.p_uint64") {
				RepeatDataSimplePathInput.Info.PUint64 = &repeatDataSimplePathInputInfoPUint64
			}

			if cmd.Flags().Changed("info.p_fixed64") {
				RepeatDataSimplePathInput.Info.PFixed64 = &repeatDataSimplePathInputInfoPFixed64
			}

			if cmd.Flags().Changed("info.p_float") {
				RepeatDataSimplePathInput.Info.PFloat = &repeatDataSimplePathInputInfoPFloat
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

			if cmd.Flags().Changed("intended_binding_uri") {
				RepeatDataSimplePathInput.IntendedBindingUri = &repeatDataSimplePathInputIntendedBindingUri
			}

			if cmd.Flags().Changed("p_int32") {
				RepeatDataSimplePathInput.PInt32 = &repeatDataSimplePathInputPInt32
			}

			if cmd.Flags().Changed("p_int64") {
				RepeatDataSimplePathInput.PInt64 = &repeatDataSimplePathInputPInt64
			}

			if cmd.Flags().Changed("p_double") {
				RepeatDataSimplePathInput.PDouble = &repeatDataSimplePathInputPDouble
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataSimplePath", &RepeatDataSimplePathInput)
		}
		resp, err := ComplianceClient.RepeatDataSimplePath(ctx, &RepeatDataSimplePathInput)
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
