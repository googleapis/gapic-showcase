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

var RepeatDataBodyInfoInput genprotopb.RepeatRequest

var RepeatDataBodyInfoFromFile string

var RepeatDataBodyInfoInputInfoFKingdom string

var RepeatDataBodyInfoInputInfoFChildFContinent string

var repeatDataBodyInfoInputInfoFChildPString string

var repeatDataBodyInfoInputInfoFChildPFloat float32

var repeatDataBodyInfoInputInfoFChildPDouble float64

var repeatDataBodyInfoInputInfoFChildPBool bool

var RepeatDataBodyInfoInputInfoFChildPContinent string

var repeatDataBodyInfoInputInfoPString string

var repeatDataBodyInfoInputInfoPInt32 int32

var repeatDataBodyInfoInputInfoPDouble float64

var repeatDataBodyInfoInputInfoPBool bool

var RepeatDataBodyInfoInputInfoPKingdom string

var RepeatDataBodyInfoInputInfoPChildFContinent string

var repeatDataBodyInfoInputInfoPChildPString string

var repeatDataBodyInfoInputInfoPChildPFloat float32

var repeatDataBodyInfoInputInfoPChildPDouble float64

var repeatDataBodyInfoInputInfoPChildPBool bool

var RepeatDataBodyInfoInputInfoPChildPContinent string

var repeatDataBodyInfoInputPInt32 int32

var repeatDataBodyInfoInputPInt64 int64

var repeatDataBodyInfoInputPDouble float64

var repeatDataBodyInfoInputIntendedBindingUri string

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataBodyInfoCmd)

	RepeatDataBodyInfoInput.Info = new(genprotopb.ComplianceData)

	RepeatDataBodyInfoInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyInfoInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyInfoInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyInfoInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyInfoInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyInfoInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Name, "name", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Info.FString, "info.f_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Int32Var(&RepeatDataBodyInfoInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataBodyInfoCmd.Flags().Int32Var(&RepeatDataBodyInfoInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataBodyInfoCmd.Flags().Int32Var(&RepeatDataBodyInfoInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataBodyInfoCmd.Flags().Uint32Var(&RepeatDataBodyInfoInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataBodyInfoCmd.Flags().Uint32Var(&RepeatDataBodyInfoInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataBodyInfoCmd.Flags().Int64Var(&RepeatDataBodyInfoInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataBodyInfoCmd.Flags().Int64Var(&RepeatDataBodyInfoInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataBodyInfoCmd.Flags().Int64Var(&RepeatDataBodyInfoInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataBodyInfoCmd.Flags().Uint64Var(&RepeatDataBodyInfoInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataBodyInfoCmd.Flags().Uint64Var(&RepeatDataBodyInfoInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().Float32Var(&RepeatDataBodyInfoInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().BytesHexVar(&RepeatDataBodyInfoInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float32Var(&RepeatDataBodyInfoInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&repeatDataBodyInfoInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float32Var(&repeatDataBodyInfoInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&repeatDataBodyInfoInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&repeatDataBodyInfoInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&repeatDataBodyInfoInputInfoPString, "info.p_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Int32Var(&repeatDataBodyInfoInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&repeatDataBodyInfoInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&repeatDataBodyInfoInputInfoPBool, "info.p_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float32Var(&RepeatDataBodyInfoInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&repeatDataBodyInfoInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float32Var(&repeatDataBodyInfoInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&repeatDataBodyInfoInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&repeatDataBodyInfoInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataBodyInfoCmd.Flags().BoolVar(&RepeatDataBodyInfoInput.ServerVerify, "server_verify", false, "If true, the server will verify that the received...")

	RepeatDataBodyInfoCmd.Flags().Int32Var(&RepeatDataBodyInfoInput.FInt32, "f_int32", 0, "Some top level fields, to test that these are...")

	RepeatDataBodyInfoCmd.Flags().Int64Var(&RepeatDataBodyInfoInput.FInt64, "f_int64", 0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&RepeatDataBodyInfoInput.FDouble, "f_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().Int32Var(&repeatDataBodyInfoInputPInt32, "p_int32", 0, "")

	RepeatDataBodyInfoCmd.Flags().Int64Var(&repeatDataBodyInfoInputPInt64, "p_int64", 0, "")

	RepeatDataBodyInfoCmd.Flags().Float64Var(&repeatDataBodyInfoInputPDouble, "p_double", 0.0, "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&repeatDataBodyInfoInputIntendedBindingUri, "intended_binding_uri", "", "")

	RepeatDataBodyInfoCmd.Flags().StringVar(&RepeatDataBodyInfoFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataBodyInfoCmd = &cobra.Command{
	Use:   "repeat-data-body-info",
	Short: "This method echoes the ComplianceData request....",
	Long:  "This method echoes the ComplianceData request. This method exercises  sending the a message-type field in the REST body. Per AIP-127, only ...",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataBodyInfoFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataBodyInfoFromFile != "" {
			in, err = os.Open(RepeatDataBodyInfoFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataBodyInfoInput)
			if err != nil {
				return err
			}

		} else {

			RepeatDataBodyInfoInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyInfoInputInfoFKingdom)])

			RepeatDataBodyInfoInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInfoInputInfoFChildFContinent)])

			RepeatDataBodyInfoInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInfoInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyInfoInputInfoPKingdom)])
				RepeatDataBodyInfoInput.Info.PKingdom = &e
			}

			RepeatDataBodyInfoInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInfoInputInfoPChildFContinent)])

			RepeatDataBodyInfoInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInfoInputInfoPChildPContinent)])

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataBodyInfoInput.Info.FChild.PString = &repeatDataBodyInfoInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataBodyInfoInput.Info.FChild.PFloat = &repeatDataBodyInfoInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataBodyInfoInput.Info.FChild.PDouble = &repeatDataBodyInfoInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataBodyInfoInput.Info.FChild.PBool = &repeatDataBodyInfoInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataBodyInfoInput.Info.PString = &repeatDataBodyInfoInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataBodyInfoInput.Info.PInt32 = &repeatDataBodyInfoInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataBodyInfoInput.Info.PDouble = &repeatDataBodyInfoInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataBodyInfoInput.Info.PBool = &repeatDataBodyInfoInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataBodyInfoInput.Info.PChild.PString = &repeatDataBodyInfoInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataBodyInfoInput.Info.PChild.PFloat = &repeatDataBodyInfoInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataBodyInfoInput.Info.PChild.PDouble = &repeatDataBodyInfoInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataBodyInfoInput.Info.PChild.PBool = &repeatDataBodyInfoInputInfoPChildPBool
			}

			if cmd.Flags().Changed("p_int32") {
				RepeatDataBodyInfoInput.PInt32 = &repeatDataBodyInfoInputPInt32
			}

			if cmd.Flags().Changed("p_int64") {
				RepeatDataBodyInfoInput.PInt64 = &repeatDataBodyInfoInputPInt64
			}

			if cmd.Flags().Changed("p_double") {
				RepeatDataBodyInfoInput.PDouble = &repeatDataBodyInfoInputPDouble
			}

			if cmd.Flags().Changed("intended_binding_uri") {
				RepeatDataBodyInfoInput.IntendedBindingUri = &repeatDataBodyInfoInputIntendedBindingUri
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataBodyInfo", &RepeatDataBodyInfoInput)
		}
		resp, err := ComplianceClient.RepeatDataBodyInfo(ctx, &RepeatDataBodyInfoInput)
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
