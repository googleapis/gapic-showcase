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

var TestingConfig *viper.Viper
var TestingClient *gapic.TestingClient
var TestingSubCommands []string = []string{
	"report-session",
	"delete-test",
	"register-test",
}

func init() {
	rootCmd.AddCommand(TestingServiceCmd)

	TestingConfig = viper.New()
	TestingConfig.SetEnvPrefix("SHOWCASE_TESTING")
	TestingConfig.AutomaticEnv()

	TestingServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use SHOWCASE_TESTING_INSECURE. Must be used with \"address\" option")
	TestingConfig.BindPFlag("insecure", TestingServiceCmd.PersistentFlags().Lookup("insecure"))
	TestingConfig.BindEnv("insecure")

	TestingServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use SHOWCASE_TESTING_ADDRESS.")
	TestingConfig.BindPFlag("address", TestingServiceCmd.PersistentFlags().Lookup("address"))
	TestingConfig.BindEnv("address")

	TestingServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use SHOWCASE_TESTING_TOKEN.")
	TestingConfig.BindPFlag("token", TestingServiceCmd.PersistentFlags().Lookup("token"))
	TestingConfig.BindEnv("token")

	TestingServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use SHOWCASE_TESTING_API_KEY.")
	TestingConfig.BindPFlag("api_key", TestingServiceCmd.PersistentFlags().Lookup("api_key"))
	TestingConfig.BindEnv("api_key")
}

var TestingServiceCmd = &cobra.Command{
	Use:       "testing",
	Short:     "A service to facilitate running discrete sets of...",
	Long:      "A service to facilitate running discrete sets of tests  against Showcase.",
	ValidArgs: TestingSubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := TestingConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if TestingConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := TestingConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := TestingConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		TestingClient, err = gapic.NewTestingClient(ctx, opts...)
		return
	},
}
