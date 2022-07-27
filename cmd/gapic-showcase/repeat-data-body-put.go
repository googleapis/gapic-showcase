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

var RepeatDataBodyPutInput genprotopb.RepeatRequest

var RepeatDataBodyPutFromFile string

var RepeatDataBodyPutInputInfoFKingdom string

var RepeatDataBodyPutInputInfoFChildFContinent string

var repeatDataBodyPutInputInfoFChildPString string

var repeatDataBodyPutInputInfoFChildPFloat float32

var repeatDataBodyPutInputInfoFChildPDouble float64

var repeatDataBodyPutInputInfoFChildPBool bool

var RepeatDataBodyPutInputInfoFChildPContinent string

var repeatDataBodyPutInputInfoPString string

var repeatDataBodyPutInputInfoPInt32 int32

var repeatDataBodyPutInputInfoPDouble float64

var repeatDataBodyPutInputInfoPBool bool

var RepeatDataBodyPutInputInfoPKingdom string

var RepeatDataBodyPutInputInfoPChildFContinent string

var repeatDataBodyPutInputInfoPChildPString string

var repeatDataBodyPutInputInfoPChildPFloat float32

var repeatDataBodyPutInputInfoPChildPDouble float64

var repeatDataBodyPutInputInfoPChildPBool bool

var RepeatDataBodyPutInputInfoPChildPContinent string

var repeatDataBodyPutInputIntendedBindingUri string

var repeatDataBodyPutInputPInt32 int32

var repeatDataBodyPutInputPInt64 int64

var repeatDataBodyPutInputPDouble float64

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataBodyPutCmd)

	RepeatDataBodyPutInput.Info = new(genprotopb.ComplianceData)

	RepeatDataBodyPutInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyPutInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPutInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPutInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyPutInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPutInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Name, "name", "", "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Info.FString, "info.f_string", "", "")

	RepeatDataBodyPutCmd.Flags().Int32Var(&RepeatDataBodyPutInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataBodyPutCmd.Flags().Int32Var(&RepeatDataBodyPutInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataBodyPutCmd.Flags().Int32Var(&RepeatDataBodyPutInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataBodyPutCmd.Flags().Uint32Var(&RepeatDataBodyPutInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataBodyPutCmd.Flags().Uint32Var(&RepeatDataBodyPutInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataBodyPutCmd.Flags().Int64Var(&RepeatDataBodyPutInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataBodyPutCmd.Flags().Int64Var(&RepeatDataBodyPutInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataBodyPutCmd.Flags().Int64Var(&RepeatDataBodyPutInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataBodyPutCmd.Flags().Uint64Var(&RepeatDataBodyPutInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataBodyPutCmd.Flags().Uint64Var(&RepeatDataBodyPutInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().Float32Var(&RepeatDataBodyPutInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataBodyPutCmd.Flags().BytesHexVar(&RepeatDataBodyPutInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float32Var(&RepeatDataBodyPutInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&repeatDataBodyPutInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float32Var(&repeatDataBodyPutInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&repeatDataBodyPutInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&repeatDataBodyPutInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&repeatDataBodyPutInputInfoPString, "info.p_string", "", "")

	RepeatDataBodyPutCmd.Flags().Int32Var(&repeatDataBodyPutInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&repeatDataBodyPutInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&repeatDataBodyPutInputInfoPBool, "info.p_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float32Var(&RepeatDataBodyPutInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&repeatDataBodyPutInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float32Var(&repeatDataBodyPutInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&repeatDataBodyPutInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&repeatDataBodyPutInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataBodyPutCmd.Flags().BoolVar(&RepeatDataBodyPutInput.ServerVerify, "server_verify", false, "If true, the server will verify that the received...")

	RepeatDataBodyPutCmd.Flags().StringVar(&repeatDataBodyPutInputIntendedBindingUri, "intended_binding_uri", "", "The URI template this request is expected to be...")

	RepeatDataBodyPutCmd.Flags().Int32Var(&RepeatDataBodyPutInput.FInt32, "f_int32", 0, "Some top level fields, to test that these are...")

	RepeatDataBodyPutCmd.Flags().Int64Var(&RepeatDataBodyPutInput.FInt64, "f_int64", 0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&RepeatDataBodyPutInput.FDouble, "f_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().Int32Var(&repeatDataBodyPutInputPInt32, "p_int32", 0, "")

	RepeatDataBodyPutCmd.Flags().Int64Var(&repeatDataBodyPutInputPInt64, "p_int64", 0, "")

	RepeatDataBodyPutCmd.Flags().Float64Var(&repeatDataBodyPutInputPDouble, "p_double", 0.0, "")

	RepeatDataBodyPutCmd.Flags().StringVar(&RepeatDataBodyPutFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataBodyPutCmd = &cobra.Command{
	Use:   "repeat-data-body-put",
	Short: "This method echoes the ComplianceData request,...",
	Long:  "This method echoes the ComplianceData request, using the HTTP PUT method.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataBodyPutFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataBodyPutFromFile != "" {
			in, err = os.Open(RepeatDataBodyPutFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataBodyPutInput)
			if err != nil {
				return err
			}

		} else {

			RepeatDataBodyPutInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyPutInputInfoFKingdom)])

			RepeatDataBodyPutInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPutInputInfoFChildFContinent)])

			RepeatDataBodyPutInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPutInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyPutInputInfoPKingdom)])
				RepeatDataBodyPutInput.Info.PKingdom = &e
			}

			RepeatDataBodyPutInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPutInputInfoPChildFContinent)])

			RepeatDataBodyPutInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyPutInputInfoPChildPContinent)])

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataBodyPutInput.Info.FChild.PString = &repeatDataBodyPutInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataBodyPutInput.Info.FChild.PFloat = &repeatDataBodyPutInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataBodyPutInput.Info.FChild.PDouble = &repeatDataBodyPutInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataBodyPutInput.Info.FChild.PBool = &repeatDataBodyPutInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataBodyPutInput.Info.PString = &repeatDataBodyPutInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataBodyPutInput.Info.PInt32 = &repeatDataBodyPutInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataBodyPutInput.Info.PDouble = &repeatDataBodyPutInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataBodyPutInput.Info.PBool = &repeatDataBodyPutInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataBodyPutInput.Info.PChild.PString = &repeatDataBodyPutInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataBodyPutInput.Info.PChild.PFloat = &repeatDataBodyPutInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataBodyPutInput.Info.PChild.PDouble = &repeatDataBodyPutInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataBodyPutInput.Info.PChild.PBool = &repeatDataBodyPutInputInfoPChildPBool
			}

			if cmd.Flags().Changed("intended_binding_uri") {
				RepeatDataBodyPutInput.IntendedBindingUri = &repeatDataBodyPutInputIntendedBindingUri
			}

			if cmd.Flags().Changed("p_int32") {
				RepeatDataBodyPutInput.PInt32 = &repeatDataBodyPutInputPInt32
			}

			if cmd.Flags().Changed("p_int64") {
				RepeatDataBodyPutInput.PInt64 = &repeatDataBodyPutInputPInt64
			}

			if cmd.Flags().Changed("p_double") {
				RepeatDataBodyPutInput.PDouble = &repeatDataBodyPutInputPDouble
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataBodyPut", &RepeatDataBodyPutInput)
		}
		resp, err := ComplianceClient.RepeatDataBodyPut(ctx, &RepeatDataBodyPutInput)
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
