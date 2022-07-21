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

var RepeatDataPathTrailingResourceInput genprotopb.RepeatRequest

var RepeatDataPathTrailingResourceFromFile string

var RepeatDataPathTrailingResourceInputInfoFKingdom string

var RepeatDataPathTrailingResourceInputInfoFChildFContinent string

var repeatDataPathTrailingResourceInputInfoFChildPString string

var repeatDataPathTrailingResourceInputInfoFChildPFloat float32

var repeatDataPathTrailingResourceInputInfoFChildPDouble float64

var repeatDataPathTrailingResourceInputInfoFChildPBool bool

var RepeatDataPathTrailingResourceInputInfoFChildPContinent string

var repeatDataPathTrailingResourceInputInfoPString string

var repeatDataPathTrailingResourceInputInfoPInt32 int32

var repeatDataPathTrailingResourceInputInfoPDouble float64

var repeatDataPathTrailingResourceInputInfoPBool bool

var RepeatDataPathTrailingResourceInputInfoPKingdom string

var RepeatDataPathTrailingResourceInputInfoPChildFContinent string

var repeatDataPathTrailingResourceInputInfoPChildPString string

var repeatDataPathTrailingResourceInputInfoPChildPFloat float32

var repeatDataPathTrailingResourceInputInfoPChildPDouble float64

var repeatDataPathTrailingResourceInputInfoPChildPBool bool

var RepeatDataPathTrailingResourceInputInfoPChildPContinent string

var repeatDataPathTrailingResourceInputPInt32 int32

var repeatDataPathTrailingResourceInputPInt64 int64

var repeatDataPathTrailingResourceInputPDouble float64

var repeatDataPathTrailingResourceInputIntendedBindingUri string

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataPathTrailingResourceCmd)

	RepeatDataPathTrailingResourceInput.Info = new(genprotopb.ComplianceData)

	RepeatDataPathTrailingResourceInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataPathTrailingResourceInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathTrailingResourceInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathTrailingResourceInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataPathTrailingResourceInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathTrailingResourceInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Name, "name", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Info.FString, "info.f_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Int32Var(&RepeatDataPathTrailingResourceInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Int32Var(&RepeatDataPathTrailingResourceInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Int32Var(&RepeatDataPathTrailingResourceInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Uint32Var(&RepeatDataPathTrailingResourceInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Uint32Var(&RepeatDataPathTrailingResourceInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Int64Var(&RepeatDataPathTrailingResourceInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Int64Var(&RepeatDataPathTrailingResourceInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Int64Var(&RepeatDataPathTrailingResourceInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Uint64Var(&RepeatDataPathTrailingResourceInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Uint64Var(&RepeatDataPathTrailingResourceInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float32Var(&RepeatDataPathTrailingResourceInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().BytesHexVar(&RepeatDataPathTrailingResourceInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float32Var(&RepeatDataPathTrailingResourceInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&repeatDataPathTrailingResourceInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float32Var(&repeatDataPathTrailingResourceInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&repeatDataPathTrailingResourceInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&repeatDataPathTrailingResourceInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&repeatDataPathTrailingResourceInputInfoPString, "info.p_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Int32Var(&repeatDataPathTrailingResourceInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&repeatDataPathTrailingResourceInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&repeatDataPathTrailingResourceInputInfoPBool, "info.p_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float32Var(&RepeatDataPathTrailingResourceInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&repeatDataPathTrailingResourceInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float32Var(&repeatDataPathTrailingResourceInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&repeatDataPathTrailingResourceInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&repeatDataPathTrailingResourceInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataPathTrailingResourceCmd.Flags().BoolVar(&RepeatDataPathTrailingResourceInput.ServerVerify, "server_verify", false, "If true, the server will verify that the received...")

	RepeatDataPathTrailingResourceCmd.Flags().Int32Var(&RepeatDataPathTrailingResourceInput.FInt32, "f_int32", 0, "Some top level fields, to test that these are...")

	RepeatDataPathTrailingResourceCmd.Flags().Int64Var(&RepeatDataPathTrailingResourceInput.FInt64, "f_int64", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&RepeatDataPathTrailingResourceInput.FDouble, "f_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Int32Var(&repeatDataPathTrailingResourceInputPInt32, "p_int32", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Int64Var(&repeatDataPathTrailingResourceInputPInt64, "p_int64", 0, "")

	RepeatDataPathTrailingResourceCmd.Flags().Float64Var(&repeatDataPathTrailingResourceInputPDouble, "p_double", 0.0, "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&repeatDataPathTrailingResourceInputIntendedBindingUri, "intended_binding_uri", "", "")

	RepeatDataPathTrailingResourceCmd.Flags().StringVar(&RepeatDataPathTrailingResourceFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataPathTrailingResourceCmd = &cobra.Command{
	Use:   "repeat-data-path-trailing-resource",
	Short: "Same as RepeatDataSimplePath, but with a trailing...",
	Long:  "Same as RepeatDataSimplePath, but with a trailing resource.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataPathTrailingResourceFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataPathTrailingResourceFromFile != "" {
			in, err = os.Open(RepeatDataPathTrailingResourceFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataPathTrailingResourceInput)
			if err != nil {
				return err
			}

		} else {

			RepeatDataPathTrailingResourceInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataPathTrailingResourceInputInfoFKingdom)])

			RepeatDataPathTrailingResourceInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathTrailingResourceInputInfoFChildFContinent)])

			RepeatDataPathTrailingResourceInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathTrailingResourceInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataPathTrailingResourceInputInfoPKingdom)])
				RepeatDataPathTrailingResourceInput.Info.PKingdom = &e
			}

			RepeatDataPathTrailingResourceInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathTrailingResourceInputInfoPChildFContinent)])

			RepeatDataPathTrailingResourceInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathTrailingResourceInputInfoPChildPContinent)])

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataPathTrailingResourceInput.Info.FChild.PString = &repeatDataPathTrailingResourceInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataPathTrailingResourceInput.Info.FChild.PFloat = &repeatDataPathTrailingResourceInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataPathTrailingResourceInput.Info.FChild.PDouble = &repeatDataPathTrailingResourceInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataPathTrailingResourceInput.Info.FChild.PBool = &repeatDataPathTrailingResourceInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataPathTrailingResourceInput.Info.PString = &repeatDataPathTrailingResourceInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataPathTrailingResourceInput.Info.PInt32 = &repeatDataPathTrailingResourceInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataPathTrailingResourceInput.Info.PDouble = &repeatDataPathTrailingResourceInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataPathTrailingResourceInput.Info.PBool = &repeatDataPathTrailingResourceInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataPathTrailingResourceInput.Info.PChild.PString = &repeatDataPathTrailingResourceInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataPathTrailingResourceInput.Info.PChild.PFloat = &repeatDataPathTrailingResourceInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataPathTrailingResourceInput.Info.PChild.PDouble = &repeatDataPathTrailingResourceInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataPathTrailingResourceInput.Info.PChild.PBool = &repeatDataPathTrailingResourceInputInfoPChildPBool
			}

			if cmd.Flags().Changed("p_int32") {
				RepeatDataPathTrailingResourceInput.PInt32 = &repeatDataPathTrailingResourceInputPInt32
			}

			if cmd.Flags().Changed("p_int64") {
				RepeatDataPathTrailingResourceInput.PInt64 = &repeatDataPathTrailingResourceInputPInt64
			}

			if cmd.Flags().Changed("p_double") {
				RepeatDataPathTrailingResourceInput.PDouble = &repeatDataPathTrailingResourceInputPDouble
			}

			if cmd.Flags().Changed("intended_binding_uri") {
				RepeatDataPathTrailingResourceInput.IntendedBindingUri = &repeatDataPathTrailingResourceInputIntendedBindingUri
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataPathTrailingResource", &RepeatDataPathTrailingResourceInput)
		}
		resp, err := ComplianceClient.RepeatDataPathTrailingResource(ctx, &RepeatDataPathTrailingResourceInput)
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
