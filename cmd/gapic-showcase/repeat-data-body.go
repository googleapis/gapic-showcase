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

var RepeatDataBodyInput genprotopb.RepeatRequest

var RepeatDataBodyFromFile string

var RepeatDataBodyInputInfoFKingdom string

var RepeatDataBodyInputInfoFChildFContinent string

var repeatDataBodyInputInfoFChildPString string

var repeatDataBodyInputInfoFChildPFloat float32

var repeatDataBodyInputInfoFChildPDouble float64

var repeatDataBodyInputInfoFChildPBool bool

var RepeatDataBodyInputInfoFChildPContinent string

var repeatDataBodyInputInfoPString string

var repeatDataBodyInputInfoPInt32 int32

var repeatDataBodyInputInfoPDouble float64

var repeatDataBodyInputInfoPBool bool

var RepeatDataBodyInputInfoPKingdom string

var RepeatDataBodyInputInfoPChildFContinent string

var repeatDataBodyInputInfoPChildPString string

var repeatDataBodyInputInfoPChildPFloat float32

var repeatDataBodyInputInfoPChildPDouble float64

var repeatDataBodyInputInfoPChildPBool bool

var RepeatDataBodyInputInfoPChildPContinent string

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataBodyCmd)

	RepeatDataBodyInput.Info = new(genprotopb.ComplianceData)

	RepeatDataBodyInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataBodyInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Name, "name", "", "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Info.FString, "info.f_string", "", "")

	RepeatDataBodyCmd.Flags().Int32Var(&RepeatDataBodyInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataBodyCmd.Flags().Int32Var(&RepeatDataBodyInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataBodyCmd.Flags().Int32Var(&RepeatDataBodyInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataBodyCmd.Flags().Uint32Var(&RepeatDataBodyInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataBodyCmd.Flags().Uint32Var(&RepeatDataBodyInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataBodyCmd.Flags().Int64Var(&RepeatDataBodyInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataBodyCmd.Flags().Int64Var(&RepeatDataBodyInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataBodyCmd.Flags().Int64Var(&RepeatDataBodyInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataBodyCmd.Flags().Uint64Var(&RepeatDataBodyInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataBodyCmd.Flags().Uint64Var(&RepeatDataBodyInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataBodyCmd.Flags().Float64Var(&RepeatDataBodyInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataBodyCmd.Flags().Float32Var(&RepeatDataBodyInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&RepeatDataBodyInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataBodyCmd.Flags().BytesHexVar(&RepeatDataBodyInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataBodyCmd.Flags().Float32Var(&RepeatDataBodyInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataBodyCmd.Flags().Float64Var(&RepeatDataBodyInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&RepeatDataBodyInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataBodyCmd.Flags().Float64Var(&RepeatDataBodyInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&RepeatDataBodyInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&repeatDataBodyInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataBodyCmd.Flags().Float32Var(&repeatDataBodyInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataBodyCmd.Flags().Float64Var(&repeatDataBodyInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&repeatDataBodyInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataBodyCmd.Flags().Float64Var(&RepeatDataBodyInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&RepeatDataBodyInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&repeatDataBodyInputInfoPString, "info.p_string", "", "")

	RepeatDataBodyCmd.Flags().Int32Var(&repeatDataBodyInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataBodyCmd.Flags().Float64Var(&repeatDataBodyInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&repeatDataBodyInputInfoPBool, "info.p_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataBodyCmd.Flags().Float32Var(&RepeatDataBodyInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataBodyCmd.Flags().Float64Var(&RepeatDataBodyInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&RepeatDataBodyInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataBodyCmd.Flags().Float64Var(&RepeatDataBodyInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&RepeatDataBodyInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&repeatDataBodyInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataBodyCmd.Flags().Float32Var(&repeatDataBodyInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataBodyCmd.Flags().Float64Var(&repeatDataBodyInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&repeatDataBodyInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataBodyCmd.Flags().Float64Var(&RepeatDataBodyInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataBodyCmd.Flags().BoolVar(&RepeatDataBodyInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataBodyCmd.Flags().StringVar(&RepeatDataBodyFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataBodyCmd = &cobra.Command{
	Use:   "repeat-data-body",
	Short: "This method echoes the ComplianceData request....",
	Long:  "This method echoes the ComplianceData request. This method exercises  sending the entire request object in the REST body.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataBodyFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataBodyFromFile != "" {
			in, err = os.Open(RepeatDataBodyFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataBodyInput)
			if err != nil {
				return err
			}

		} else {

			RepeatDataBodyInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyInputInfoFKingdom)])

			RepeatDataBodyInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInputInfoFChildFContinent)])

			RepeatDataBodyInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataBodyInputInfoPKingdom)])
				RepeatDataBodyInput.Info.PKingdom = &e
			}

			RepeatDataBodyInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInputInfoPChildFContinent)])

			RepeatDataBodyInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataBodyInputInfoPChildPContinent)])

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataBodyInput.Info.FChild.PString = &repeatDataBodyInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataBodyInput.Info.FChild.PFloat = &repeatDataBodyInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataBodyInput.Info.FChild.PDouble = &repeatDataBodyInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataBodyInput.Info.FChild.PBool = &repeatDataBodyInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataBodyInput.Info.PString = &repeatDataBodyInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataBodyInput.Info.PInt32 = &repeatDataBodyInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataBodyInput.Info.PDouble = &repeatDataBodyInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataBodyInput.Info.PBool = &repeatDataBodyInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataBodyInput.Info.PChild.PString = &repeatDataBodyInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataBodyInput.Info.PChild.PFloat = &repeatDataBodyInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataBodyInput.Info.PChild.PDouble = &repeatDataBodyInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataBodyInput.Info.PChild.PBool = &repeatDataBodyInputInfoPChildPBool
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataBody", &RepeatDataBodyInput)
		}
		resp, err := ComplianceClient.RepeatDataBody(ctx, &RepeatDataBodyInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
