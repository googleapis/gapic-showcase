// Code generated. DO NOT EDIT.

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	gapic "github.com/googleapis/gapic-showcase/client"
)

var ComplianceConfig *viper.Viper
var ComplianceClient *gapic.ComplianceClient
var ComplianceSubCommands []string = []string{
	"repeat-data-body",
	"repeat-data-body-info",
	"repeat-data-query",
	"repeat-data-simple-path",
	"repeat-data-path-resource",
	"repeat-data-path-trailing-resource",
}

func init() {
	rootCmd.AddCommand(ComplianceServiceCmd)

	ComplianceConfig = viper.New()
	ComplianceConfig.SetEnvPrefix("GAPIC-SHOWCASE_COMPLIANCE")
	ComplianceConfig.AutomaticEnv()

	ComplianceServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use GAPIC-SHOWCASE_COMPLIANCE_INSECURE. Must be used with \"address\" option")
	ComplianceConfig.BindPFlag("insecure", ComplianceServiceCmd.PersistentFlags().Lookup("insecure"))
	ComplianceConfig.BindEnv("insecure")

	ComplianceServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use GAPIC-SHOWCASE_COMPLIANCE_ADDRESS.")
	ComplianceConfig.BindPFlag("address", ComplianceServiceCmd.PersistentFlags().Lookup("address"))
	ComplianceConfig.BindEnv("address")

	ComplianceServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use GAPIC-SHOWCASE_COMPLIANCE_TOKEN.")
	ComplianceConfig.BindPFlag("token", ComplianceServiceCmd.PersistentFlags().Lookup("token"))
	ComplianceConfig.BindEnv("token")

	ComplianceServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use GAPIC-SHOWCASE_COMPLIANCE_API_KEY.")
	ComplianceConfig.BindPFlag("api_key", ComplianceServiceCmd.PersistentFlags().Lookup("api_key"))
	ComplianceConfig.BindEnv("api_key")
}

var ComplianceServiceCmd = &cobra.Command{
	Use:       "compliance",
	Short:     "This service is used to test that GAPICs can...",
	Long:      "This service is used to test that GAPICs can transcode proto3 requests to  REST format correctly for various types of HTTP annotations.",
	ValidArgs: ComplianceSubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := ComplianceConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if ComplianceConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := ComplianceConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := ComplianceConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		ComplianceClient, err = gapic.NewComplianceClient(ctx, opts...)
		return
	},
}
