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

var IdentityConfig *viper.Viper
var IdentityClient *gapic.IdentityClient
var IdentitySubCommands []string = []string{
	"create-user",
	"get-user",
	"update-user",
	"delete-user",
	"list-users",
}

func init() {
	rootCmd.AddCommand(IdentityServiceCmd)

	IdentityConfig = viper.New()
	IdentityConfig.SetEnvPrefix("SHOWCASE_IDENTITY")
	IdentityConfig.AutomaticEnv()

	IdentityServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use SHOWCASE_IDENTITY_INSECURE. Must be used with \"address\" option")
	IdentityConfig.BindPFlag("insecure", IdentityServiceCmd.PersistentFlags().Lookup("insecure"))
	IdentityConfig.BindEnv("insecure")

	IdentityServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use SHOWCASE_IDENTITY_ADDRESS.")
	IdentityConfig.BindPFlag("address", IdentityServiceCmd.PersistentFlags().Lookup("address"))
	IdentityConfig.BindEnv("address")

	IdentityServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use SHOWCASE_IDENTITY_TOKEN.")
	IdentityConfig.BindPFlag("token", IdentityServiceCmd.PersistentFlags().Lookup("token"))
	IdentityConfig.BindEnv("token")

	IdentityServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use SHOWCASE_IDENTITY_API_KEY.")
	IdentityConfig.BindPFlag("api_key", IdentityServiceCmd.PersistentFlags().Lookup("api_key"))
	IdentityConfig.BindEnv("api_key")
}

var IdentityServiceCmd = &cobra.Command{
	Use:       "identity",
	Short:     "A simple identity service.",
	Long:      "A simple identity service.",
	ValidArgs: IdentitySubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := IdentityConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if IdentityConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := IdentityConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := IdentityConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		IdentityClient, err = gapic.NewIdentityClient(ctx, opts...)
		return
	},
}
