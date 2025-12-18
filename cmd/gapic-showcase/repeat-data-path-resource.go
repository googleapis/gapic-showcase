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

var RepeatDataPathResourceInput genprotopb.RepeatRequest

var RepeatDataPathResourceFromFile string

var RepeatDataPathResourceInputInfoFKingdom string

var RepeatDataPathResourceInputInfoFChildFContinent string

var repeatDataPathResourceInputInfoFChildPString string

var repeatDataPathResourceInputInfoFChildPFloat float32

var repeatDataPathResourceInputInfoFChildPDouble float64

var repeatDataPathResourceInputInfoFChildPBool bool

var RepeatDataPathResourceInputInfoFChildPContinent string

var repeatDataPathResourceInputInfoPString string

var repeatDataPathResourceInputInfoPInt32 int32

var repeatDataPathResourceInputInfoPSint32 int32

var repeatDataPathResourceInputInfoPSfixed32 int32

var repeatDataPathResourceInputInfoPUint32 uint32

var repeatDataPathResourceInputInfoPFixed32 uint32

var repeatDataPathResourceInputInfoPInt64 int64

var repeatDataPathResourceInputInfoPSint64 int64

var repeatDataPathResourceInputInfoPSfixed64 int64

var repeatDataPathResourceInputInfoPUint64 uint64

var repeatDataPathResourceInputInfoPFixed64 uint64

var repeatDataPathResourceInputInfoPFloat float32

var repeatDataPathResourceInputInfoPDouble float64

var repeatDataPathResourceInputInfoPBool bool

var RepeatDataPathResourceInputInfoPKingdom string

var RepeatDataPathResourceInputInfoPChildFContinent string

var repeatDataPathResourceInputInfoPChildPString string

var repeatDataPathResourceInputInfoPChildPFloat float32

var repeatDataPathResourceInputInfoPChildPDouble float64

var repeatDataPathResourceInputInfoPChildPBool bool

var RepeatDataPathResourceInputInfoPChildPContinent string

var repeatDataPathResourceInputIntendedBindingUri string

var repeatDataPathResourceInputPInt32 int32

var repeatDataPathResourceInputPInt64 int64

var repeatDataPathResourceInputPDouble float64

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

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&RepeatDataPathResourceInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&repeatDataPathResourceInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&repeatDataPathResourceInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&repeatDataPathResourceInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&repeatDataPathResourceInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&repeatDataPathResourceInputInfoPString, "info.p_string", "", "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&repeatDataPathResourceInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&repeatDataPathResourceInputInfoPSint32, "info.p_sint32", 0, "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&repeatDataPathResourceInputInfoPSfixed32, "info.p_sfixed32", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint32Var(&repeatDataPathResourceInputInfoPUint32, "info.p_uint32", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint32Var(&repeatDataPathResourceInputInfoPFixed32, "info.p_fixed32", 0, "")

	RepeatDataPathResourceCmd.Flags().Int64Var(&repeatDataPathResourceInputInfoPInt64, "info.p_int64", 0, "")

	RepeatDataPathResourceCmd.Flags().Int64Var(&repeatDataPathResourceInputInfoPSint64, "info.p_sint64", 0, "")

	RepeatDataPathResourceCmd.Flags().Int64Var(&repeatDataPathResourceInputInfoPSfixed64, "info.p_sfixed64", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint64Var(&repeatDataPathResourceInputInfoPUint64, "info.p_uint64", 0, "")

	RepeatDataPathResourceCmd.Flags().Uint64Var(&repeatDataPathResourceInputInfoPFixed64, "info.p_fixed64", 0, "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&repeatDataPathResourceInputInfoPFloat, "info.p_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&repeatDataPathResourceInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&repeatDataPathResourceInputInfoPBool, "info.p_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&RepeatDataPathResourceInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&repeatDataPathResourceInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float32Var(&repeatDataPathResourceInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&repeatDataPathResourceInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&repeatDataPathResourceInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataPathResourceCmd.Flags().StringVar(&RepeatDataPathResourceInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataPathResourceCmd.Flags().BoolVar(&RepeatDataPathResourceInput.ServerVerify, "server_verify", false, "If true, the server will verify that the received...")

	RepeatDataPathResourceCmd.Flags().StringVar(&repeatDataPathResourceInputIntendedBindingUri, "intended_binding_uri", "", "The URI template this request is expected to be...")

	RepeatDataPathResourceCmd.Flags().Int32Var(&RepeatDataPathResourceInput.FInt32, "f_int32", 0, "Some top level fields, to test that these are...")

	RepeatDataPathResourceCmd.Flags().Int64Var(&RepeatDataPathResourceInput.FInt64, "f_int64", 0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&RepeatDataPathResourceInput.FDouble, "f_double", 0.0, "")

	RepeatDataPathResourceCmd.Flags().Int32Var(&repeatDataPathResourceInputPInt32, "p_int32", 0, "")

	RepeatDataPathResourceCmd.Flags().Int64Var(&repeatDataPathResourceInputPInt64, "p_int64", 0, "")

	RepeatDataPathResourceCmd.Flags().Float64Var(&repeatDataPathResourceInputPDouble, "p_double", 0.0, "")

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

			RepeatDataPathResourceInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataPathResourceInputInfoFKingdom)])

			RepeatDataPathResourceInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathResourceInputInfoFChildFContinent)])

			RepeatDataPathResourceInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathResourceInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataPathResourceInputInfoPKingdom)])
				RepeatDataPathResourceInput.Info.PKingdom = &e
			}

			RepeatDataPathResourceInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathResourceInputInfoPChildFContinent)])

			RepeatDataPathResourceInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataPathResourceInputInfoPChildPContinent)])

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

			if cmd.Flags().Changed("info.p_sint32") {
				RepeatDataPathResourceInput.Info.PSint32 = &repeatDataPathResourceInputInfoPSint32
			}

			if cmd.Flags().Changed("info.p_sfixed32") {
				RepeatDataPathResourceInput.Info.PSfixed32 = &repeatDataPathResourceInputInfoPSfixed32
			}

			if cmd.Flags().Changed("info.p_uint32") {
				RepeatDataPathResourceInput.Info.PUint32 = &repeatDataPathResourceInputInfoPUint32
			}

			if cmd.Flags().Changed("info.p_fixed32") {
				RepeatDataPathResourceInput.Info.PFixed32 = &repeatDataPathResourceInputInfoPFixed32
			}

			if cmd.Flags().Changed("info.p_int64") {
				RepeatDataPathResourceInput.Info.PInt64 = &repeatDataPathResourceInputInfoPInt64
			}

			if cmd.Flags().Changed("info.p_sint64") {
				RepeatDataPathResourceInput.Info.PSint64 = &repeatDataPathResourceInputInfoPSint64
			}

			if cmd.Flags().Changed("info.p_sfixed64") {
				RepeatDataPathResourceInput.Info.PSfixed64 = &repeatDataPathResourceInputInfoPSfixed64
			}

			if cmd.Flags().Changed("info.p_uint64") {
				RepeatDataPathResourceInput.Info.PUint64 = &repeatDataPathResourceInputInfoPUint64
			}

			if cmd.Flags().Changed("info.p_fixed64") {
				RepeatDataPathResourceInput.Info.PFixed64 = &repeatDataPathResourceInputInfoPFixed64
			}

			if cmd.Flags().Changed("info.p_float") {
				RepeatDataPathResourceInput.Info.PFloat = &repeatDataPathResourceInputInfoPFloat
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

			if cmd.Flags().Changed("intended_binding_uri") {
				RepeatDataPathResourceInput.IntendedBindingUri = &repeatDataPathResourceInputIntendedBindingUri
			}

			if cmd.Flags().Changed("p_int32") {
				RepeatDataPathResourceInput.PInt32 = &repeatDataPathResourceInputPInt32
			}

			if cmd.Flags().Changed("p_int64") {
				RepeatDataPathResourceInput.PInt64 = &repeatDataPathResourceInputPInt64
			}

			if cmd.Flags().Changed("p_double") {
				RepeatDataPathResourceInput.PDouble = &repeatDataPathResourceInputPDouble
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataPathResource", &RepeatDataPathResourceInput)
		}
		resp, err := ComplianceClient.RepeatDataPathResource(ctx, &RepeatDataPathResourceInput)
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
