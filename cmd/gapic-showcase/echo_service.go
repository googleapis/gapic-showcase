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

var EchoConfig *viper.Viper
var EchoClient *gapic.EchoClient
var EchoSubCommands []string = []string{
	"echo",
	"echo-error-details",
	"expand",
	"collect",
	"chat",
	"paged-expand",
	"paged-expand-legacy",
	"paged-expand-legacy-mapped",
	"wait",
	"poll-wait", "block",
}

func init() {
	rootCmd.AddCommand(EchoServiceCmd)

	EchoConfig = viper.New()
	EchoConfig.SetEnvPrefix("GAPIC-SHOWCASE_ECHO")
	EchoConfig.AutomaticEnv()

	EchoServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use GAPIC-SHOWCASE_ECHO_INSECURE. Must be used with \"address\" option")
	EchoConfig.BindPFlag("insecure", EchoServiceCmd.PersistentFlags().Lookup("insecure"))
	EchoConfig.BindEnv("insecure")

	EchoServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use GAPIC-SHOWCASE_ECHO_ADDRESS.")
	EchoConfig.BindPFlag("address", EchoServiceCmd.PersistentFlags().Lookup("address"))
	EchoConfig.BindEnv("address")

	EchoServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use GAPIC-SHOWCASE_ECHO_TOKEN.")
	EchoConfig.BindPFlag("token", EchoServiceCmd.PersistentFlags().Lookup("token"))
	EchoConfig.BindEnv("token")

	EchoServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use GAPIC-SHOWCASE_ECHO_API_KEY.")
	EchoConfig.BindPFlag("api_key", EchoServiceCmd.PersistentFlags().Lookup("api_key"))
	EchoConfig.BindEnv("api_key")
}

var EchoServiceCmd = &cobra.Command{
	Use:       "echo",
	Short:     "This service is used showcase the four main types...",
	Long:      "This service is used showcase the four main types of rpcs - unary, server  side streaming, client side streaming, and bidirectional streaming. This ...",
	ValidArgs: EchoSubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := EchoConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if EchoConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := EchoConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := EchoConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		EchoClient, err = gapic.NewEchoClient(ctx, opts...)
		return
	},
}
