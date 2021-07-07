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

var RepeatDataQueryInput genprotopb.RepeatRequest

var RepeatDataQueryFromFile string

var RepeatDataQueryInputInfoFKingdom string

var RepeatDataQueryInputInfoFChildFContinent string

var repeatDataQueryInputInfoFChildPString string

var repeatDataQueryInputInfoFChildPFloat float32

var repeatDataQueryInputInfoFChildPDouble float64

var repeatDataQueryInputInfoFChildPBool bool

var RepeatDataQueryInputInfoFChildPContinent string

var repeatDataQueryInputInfoPString string

var repeatDataQueryInputInfoPInt32 int32

var repeatDataQueryInputInfoPDouble float64

var repeatDataQueryInputInfoPBool bool

var RepeatDataQueryInputInfoPKingdom string

var RepeatDataQueryInputInfoPChildFContinent string

var repeatDataQueryInputInfoPChildPString string

var repeatDataQueryInputInfoPChildPFloat float32

var repeatDataQueryInputInfoPChildPDouble float64

var repeatDataQueryInputInfoPChildPBool bool

var RepeatDataQueryInputInfoPChildPContinent string

var repeatDataQueryInputPInt32 int32

var repeatDataQueryInputPInt64 int64

var repeatDataQueryInputPDouble float64

func init() {
	ComplianceServiceCmd.AddCommand(RepeatDataQueryCmd)

	RepeatDataQueryInput.Info = new(genprotopb.ComplianceData)

	RepeatDataQueryInput.Info.FChild = new(genprotopb.ComplianceDataChild)

	RepeatDataQueryInput.Info.FChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataQueryInput.Info.FChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataQueryInput.Info.PChild = new(genprotopb.ComplianceDataChild)

	RepeatDataQueryInput.Info.PChild.FChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataQueryInput.Info.PChild.PChild = new(genprotopb.ComplianceDataGrandchild)

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Name, "name", "", "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Info.FString, "info.f_string", "", "")

	RepeatDataQueryCmd.Flags().Int32Var(&RepeatDataQueryInput.Info.FInt32, "info.f_int32", 0, "")

	RepeatDataQueryCmd.Flags().Int32Var(&RepeatDataQueryInput.Info.FSint32, "info.f_sint32", 0, "")

	RepeatDataQueryCmd.Flags().Int32Var(&RepeatDataQueryInput.Info.FSfixed32, "info.f_sfixed32", 0, "")

	RepeatDataQueryCmd.Flags().Uint32Var(&RepeatDataQueryInput.Info.FUint32, "info.f_uint32", 0, "")

	RepeatDataQueryCmd.Flags().Uint32Var(&RepeatDataQueryInput.Info.FFixed32, "info.f_fixed32", 0, "")

	RepeatDataQueryCmd.Flags().Int64Var(&RepeatDataQueryInput.Info.FInt64, "info.f_int64", 0, "")

	RepeatDataQueryCmd.Flags().Int64Var(&RepeatDataQueryInput.Info.FSint64, "info.f_sint64", 0, "")

	RepeatDataQueryCmd.Flags().Int64Var(&RepeatDataQueryInput.Info.FSfixed64, "info.f_sfixed64", 0, "")

	RepeatDataQueryCmd.Flags().Uint64Var(&RepeatDataQueryInput.Info.FUint64, "info.f_uint64", 0, "")

	RepeatDataQueryCmd.Flags().Uint64Var(&RepeatDataQueryInput.Info.FFixed64, "info.f_fixed64", 0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.Info.FDouble, "info.f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().Float32Var(&RepeatDataQueryInput.Info.FFloat, "info.f_float", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.Info.FBool, "info.f_bool", false, "")

	RepeatDataQueryCmd.Flags().BytesHexVar(&RepeatDataQueryInput.Info.FBytes, "info.f_bytes", []byte{}, "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInputInfoFKingdom, "info.f_kingdom", "", "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Info.FChild.FString, "info.f_child.f_string", "", "")

	RepeatDataQueryCmd.Flags().Float32Var(&RepeatDataQueryInput.Info.FChild.FFloat, "info.f_child.f_float", 0.0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.Info.FChild.FDouble, "info.f_child.f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.Info.FChild.FBool, "info.f_child.f_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInputInfoFChildFContinent, "info.f_child.f_continent", "", "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Info.FChild.FChild.FString, "info.f_child.f_child.f_string", "", "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.Info.FChild.FChild.FDouble, "info.f_child.f_child.f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.Info.FChild.FChild.FBool, "info.f_child.f_child.f_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&repeatDataQueryInputInfoFChildPString, "info.f_child.p_string", "", "")

	RepeatDataQueryCmd.Flags().Float32Var(&repeatDataQueryInputInfoFChildPFloat, "info.f_child.p_float", 0.0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&repeatDataQueryInputInfoFChildPDouble, "info.f_child.p_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&repeatDataQueryInputInfoFChildPBool, "info.f_child.p_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInputInfoFChildPContinent, "info.f_child.p_continent", "", "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Info.FChild.PChild.FString, "info.f_child.p_child.f_string", "", "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.Info.FChild.PChild.FDouble, "info.f_child.p_child.f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.Info.FChild.PChild.FBool, "info.f_child.p_child.f_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&repeatDataQueryInputInfoPString, "info.p_string", "", "")

	RepeatDataQueryCmd.Flags().Int32Var(&repeatDataQueryInputInfoPInt32, "info.p_int32", 0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&repeatDataQueryInputInfoPDouble, "info.p_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&repeatDataQueryInputInfoPBool, "info.p_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInputInfoPKingdom, "info.p_kingdom", "", "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Info.PChild.FString, "info.p_child.f_string", "", "")

	RepeatDataQueryCmd.Flags().Float32Var(&RepeatDataQueryInput.Info.PChild.FFloat, "info.p_child.f_float", 0.0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.Info.PChild.FDouble, "info.p_child.f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.Info.PChild.FBool, "info.p_child.f_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInputInfoPChildFContinent, "info.p_child.f_continent", "", "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Info.PChild.FChild.FString, "info.p_child.f_child.f_string", "", "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.Info.PChild.FChild.FDouble, "info.p_child.f_child.f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.Info.PChild.FChild.FBool, "info.p_child.f_child.f_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&repeatDataQueryInputInfoPChildPString, "info.p_child.p_string", "", "")

	RepeatDataQueryCmd.Flags().Float32Var(&repeatDataQueryInputInfoPChildPFloat, "info.p_child.p_float", 0.0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&repeatDataQueryInputInfoPChildPDouble, "info.p_child.p_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&repeatDataQueryInputInfoPChildPBool, "info.p_child.p_bool", false, "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInputInfoPChildPContinent, "info.p_child.p_continent", "", "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryInput.Info.PChild.PChild.FString, "info.p_child.p_child.f_string", "", "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.Info.PChild.PChild.FDouble, "info.p_child.p_child.f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.Info.PChild.PChild.FBool, "info.p_child.p_child.f_bool", false, "")

	RepeatDataQueryCmd.Flags().BoolVar(&RepeatDataQueryInput.ServerVerify, "server_verify", false, "If true, the server will verify that the received...")

	RepeatDataQueryCmd.Flags().Int32Var(&RepeatDataQueryInput.FInt32, "f_int32", 0, "Some top level fields, to test that these are...")

	RepeatDataQueryCmd.Flags().Int64Var(&RepeatDataQueryInput.FInt64, "f_int64", 0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&RepeatDataQueryInput.FDouble, "f_double", 0.0, "")

	RepeatDataQueryCmd.Flags().Int32Var(&repeatDataQueryInputPInt32, "p_int32", 0, "")

	RepeatDataQueryCmd.Flags().Int64Var(&repeatDataQueryInputPInt64, "p_int64", 0, "")

	RepeatDataQueryCmd.Flags().Float64Var(&repeatDataQueryInputPDouble, "p_double", 0.0, "")

	RepeatDataQueryCmd.Flags().StringVar(&RepeatDataQueryFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

}

var RepeatDataQueryCmd = &cobra.Command{
	Use:   "repeat-data-query",
	Short: "This method echoes the ComplianceData request....",
	Long:  "This method echoes the ComplianceData request. This method exercises  sending all request fields as query parameters.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if RepeatDataQueryFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if RepeatDataQueryFromFile != "" {
			in, err = os.Open(RepeatDataQueryFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &RepeatDataQueryInput)
			if err != nil {
				return err
			}

		} else {

			RepeatDataQueryInput.Info.FKingdom = genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataQueryInputInfoFKingdom)])

			RepeatDataQueryInput.Info.FChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataQueryInputInfoFChildFContinent)])

			RepeatDataQueryInput.Info.FChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataQueryInputInfoFChildPContinent)])

			if cmd.Flags().Changed("info.p_kingdom") {
				e := genprotopb.ComplianceData_LifeKingdom(genprotopb.ComplianceData_LifeKingdom_value[strings.ToUpper(RepeatDataQueryInputInfoPKingdom)])
				RepeatDataQueryInput.Info.PKingdom = &e
			}

			RepeatDataQueryInput.Info.PChild.FContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataQueryInputInfoPChildFContinent)])

			RepeatDataQueryInput.Info.PChild.PContinent = genprotopb.Continent(genprotopb.Continent_value[strings.ToUpper(RepeatDataQueryInputInfoPChildPContinent)])

			if cmd.Flags().Changed("info.f_child.p_string") {
				RepeatDataQueryInput.Info.FChild.PString = &repeatDataQueryInputInfoFChildPString
			}

			if cmd.Flags().Changed("info.f_child.p_float") {
				RepeatDataQueryInput.Info.FChild.PFloat = &repeatDataQueryInputInfoFChildPFloat
			}

			if cmd.Flags().Changed("info.f_child.p_double") {
				RepeatDataQueryInput.Info.FChild.PDouble = &repeatDataQueryInputInfoFChildPDouble
			}

			if cmd.Flags().Changed("info.f_child.p_bool") {
				RepeatDataQueryInput.Info.FChild.PBool = &repeatDataQueryInputInfoFChildPBool
			}

			if cmd.Flags().Changed("info.p_string") {
				RepeatDataQueryInput.Info.PString = &repeatDataQueryInputInfoPString
			}

			if cmd.Flags().Changed("info.p_int32") {
				RepeatDataQueryInput.Info.PInt32 = &repeatDataQueryInputInfoPInt32
			}

			if cmd.Flags().Changed("info.p_double") {
				RepeatDataQueryInput.Info.PDouble = &repeatDataQueryInputInfoPDouble
			}

			if cmd.Flags().Changed("info.p_bool") {
				RepeatDataQueryInput.Info.PBool = &repeatDataQueryInputInfoPBool
			}

			if cmd.Flags().Changed("info.p_child.p_string") {
				RepeatDataQueryInput.Info.PChild.PString = &repeatDataQueryInputInfoPChildPString
			}

			if cmd.Flags().Changed("info.p_child.p_float") {
				RepeatDataQueryInput.Info.PChild.PFloat = &repeatDataQueryInputInfoPChildPFloat
			}

			if cmd.Flags().Changed("info.p_child.p_double") {
				RepeatDataQueryInput.Info.PChild.PDouble = &repeatDataQueryInputInfoPChildPDouble
			}

			if cmd.Flags().Changed("info.p_child.p_bool") {
				RepeatDataQueryInput.Info.PChild.PBool = &repeatDataQueryInputInfoPChildPBool
			}

			if cmd.Flags().Changed("p_int32") {
				RepeatDataQueryInput.PInt32 = &repeatDataQueryInputPInt32
			}

			if cmd.Flags().Changed("p_int64") {
				RepeatDataQueryInput.PInt64 = &repeatDataQueryInputPInt64
			}

			if cmd.Flags().Changed("p_double") {
				RepeatDataQueryInput.PDouble = &repeatDataQueryInputPDouble
			}

		}

		if Verbose {
			printVerboseInput("Compliance", "RepeatDataQuery", &RepeatDataQueryInput)
		}
		resp, err := ComplianceClient.RepeatDataQuery(ctx, &RepeatDataQueryInput)

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)

		return err
	},
}
