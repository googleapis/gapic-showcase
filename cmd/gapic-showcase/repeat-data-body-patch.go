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

var RepeatDataBodyPatchInput genprotopb.RepeatRequest

var RepeatDataBodyPatchFromFile string

var RepeatDataBodyPatchInputInfoFKingdom string

var RepeatDataBodyPatchInputInfoFChildFContinent string

var repeatDataBodyPatchInputInfoFChildPString string

var repeatDataBodyPatchInputInfoFChildPFloat float32

var repeatDataBodyPatchInputInfoFChildPDouble float64

var repeatDataBodyPatchInputInfoFChildPBool bool

var RepeatDataBodyPatchInputInfoFChildPContinent string

var repeatDataBodyPatchInputInfoPString string

var repeatDataBodyPatchInputInfoPInt32 int32

var repeatDataBodyPatchInputInfoPDouble float64

var repeatDataBodyPatchInputInfoPBool bool

var RepeatDataBodyPatchInputInfoPKingdom string

var RepeatDataBodyPatchInputInfoPChildFContinent string

var repeatDataBodyPatchInputInfoPChildPString string

var repeatDataBodyPatchInputInfoPChildPFloat float32

var repeatDataBodyPatchInputInfoPChildPDouble float64

var repeatDataBodyPatchInputInfoPChildPBool bool

var RepeatDataBodyPatchInputInfoPChildPContinent string

var repeatDataBodyPatchInputPInt32 int32

var repeatDataBodyPatchInputPInt64 int64

var repeatDataBodyPatchInputPDouble float64

var repeatDataBodyPatchInputIntendedBindingUri string

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataBodyPatchCmd)

	RepeatDataBodyPatchInput.Info = new(genprotopb.ComplianceData)

	RepeatDataBodyPatchInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyPatchInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPatchInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPatchInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyPatchInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPatchInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Name, "name", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Info.FString, "info.f_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Int32Var(&RepeatDataBodyPatchInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataBodyPatchCmd.Flags().Int32Var(&RepeatDataBodyPatchInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataBodyPatchCmd.Flags().Int32Var(&RepeatDataBodyPatchInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataBodyPatchCmd.Flags().Uint32Var(&RepeatDataBodyPatchInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataBodyPatchCmd.Flags().Uint32Var(&RepeatDataBodyPatchInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataBodyPatchCmd.Flags().Int64Var(&RepeatDataBodyPatchInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataBodyPatchCmd.Flags().Int64Var(&RepeatDataBodyPatchInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataBodyPatchCmd.Flags().Int64Var(&RepeatDataBodyPatchInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataBodyPatchCmd.Flags().Uint64Var(&RepeatDataBodyPatchInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataBodyPatchCmd.Flags().Uint64Var(&RepeatDataBodyPatchInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().Float32Var(&RepeatDataBodyPatchInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().BytesHexVar(&RepeatDataBodyPatchInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float32Var(&RepeatDataBodyPatchInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&repeatDataBodyPatchInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float32Var(&repeatDataBodyPatchInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&repeatDataBodyPatchInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&repeatDataBodyPatchInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&repeatDataBodyPatchInputInfoPString, "info.p_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Int32Var(&repeatDataBodyPatchInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&repeatDataBodyPatchInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&repeatDataBodyPatchInputInfoPBool, "info.p_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float32Var(&RepeatDataBodyPatchInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&repeatDataBodyPatchInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float32Var(&repeatDataBodyPatchInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&repeatDataBodyPatchInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&repeatDataBodyPatchInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataBodyPatchCmd.Flags().BoolVar(&RepeatDataBodyPatchInput.ServerVerify, "server_verify", false, "If true, the server will verify that the received...")

	RepeatDataBodyPatchCmd.Flags().Int32Var(&RepeatDataBodyPatchInput.FInt32, "f_int32", 0, "Some top level fields, to test that these are...")

	RepeatDataBodyPatchCmd.Flags().Int64Var(&RepeatDataBodyPatchInput.FInt64, "f_int64", 0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&RepeatDataBodyPatchInput.FDouble, "f_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().Int32Var(&repeatDataBodyPatchInputPInt32, "p_int32", 0, "")

	RepeatDataBodyPatchCmd.Flags().Int64Var(&repeatDataBodyPatchInputPInt64, "p_int64", 0, "")

	RepeatDataBodyPatchCmd.Flags().Float64Var(&repeatDataBodyPatchInputPDouble, "p_double", 0.0, "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&repeatDataBodyPatchInputIntendedBindingUri, "intended_binding_uri", "", "")

	RepeatDataBodyPatchCmd.Flags().StringVar(&RepeatDataBodyPatchFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataBodyPatchCmd = &cobra.Command{
	Use:   "repeat-data-body-patch",
	Short: "This method echoes the ComplianceData request,...",
	Long:  "This method echoes the ComplianceData request, using the HTTP PATCH method.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataBodyPatchFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataBodyPatchFromFile != "" {
			in, err = os.Open(RepeatDataBodyPatchFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataBodyPatchInput)
			if err != nil {
				return err
			}

		} else {

			RepeatDataBodyPatchInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyPatchInputInfoFKingdom)])

			RepeatDataBodyPatchInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPatchInputInfoFChildFContinent)])

			RepeatDataBodyPatchInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPatchInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyPatchInputInfoPKingdom)])
				RepeatDataBodyPatchInput.Info.PKingdom = &e
			}

			RepeatDataBodyPatchInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPatchInputInfoPChildFContinent)])

			RepeatDataBodyPatchInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPatchInputInfoPChildPContinent)])

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataBodyPatchInput.Info.FChild.PString = &repeatDataBodyPatchInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataBodyPatchInput.Info.FChild.PFloat = &repeatDataBodyPatchInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataBodyPatchInput.Info.FChild.PDouble = &repeatDataBodyPatchInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataBodyPatchInput.Info.FChild.PBool = &repeatDataBodyPatchInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataBodyPatchInput.Info.PString = &repeatDataBodyPatchInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataBodyPatchInput.Info.PInt32 = &repeatDataBodyPatchInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataBodyPatchInput.Info.PDouble = &repeatDataBodyPatchInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataBodyPatchInput.Info.PBool = &repeatDataBodyPatchInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataBodyPatchInput.Info.PChild.PString = &repeatDataBodyPatchInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataBodyPatchInput.Info.PChild.PFloat = &repeatDataBodyPatchInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataBodyPatchInput.Info.PChild.PDouble = &repeatDataBodyPatchInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataBodyPatchInput.Info.PChild.PBool = &repeatDataBodyPatchInputInfoPChildPBool
			}

			if cmd.Flags().Changed("p_int32") {
				RepeatDataBodyPatchInput.PInt32 = &repeatDataBodyPatchInputPInt32
			}

			if cmd.Flags().Changed("p_int64") {
				RepeatDataBodyPatchInput.PInt64 = &repeatDataBodyPatchInputPInt64
			}

			if cmd.Flags().Changed("p_double") {
				RepeatDataBodyPatchInput.PDouble = &repeatDataBodyPatchInputPDouble
			}

			if cmd.Flags().Changed("intended_binding_uri") {
				RepeatDataBodyPatchInput.IntendedBindingUri = &repeatDataBodyPatchInputIntendedBindingUri
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataBodyPatch", &RepeatDataBodyPatchInput)
		}
		resp, err := ComplianceClient.RepeatDataBodyPatch(ctx, &RepeatDataBodyPatchInput)
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
