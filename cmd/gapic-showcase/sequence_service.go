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

var SequenceConfig *viper.Viper
var SequenceClient *gapic.SequenceClient
var SequenceSubCommands []string = []string{
	"create-sequence",
	"create-streaming-sequence",
	"get-sequence-report",
	"get-streaming-sequence-report",
	"attempt-sequence",
	"attempt-streaming-sequence",
}

func init() {
	rootCmd.AddCommand(SequenceServiceCmd)

	SequenceConfig = viper.New()
	SequenceConfig.SetEnvPrefix("GAPIC-SHOWCASE_SEQUENCE")
	SequenceConfig.AutomaticEnv()

	SequenceServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use GAPIC-SHOWCASE_SEQUENCE_INSECURE. Must be used with \"address\" option")
	SequenceConfig.BindPFlag("insecure", SequenceServiceCmd.PersistentFlags().Lookup("insecure"))
	SequenceConfig.BindEnv("insecure")

	SequenceServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use GAPIC-SHOWCASE_SEQUENCE_ADDRESS.")
	SequenceConfig.BindPFlag("address", SequenceServiceCmd.PersistentFlags().Lookup("address"))
	SequenceConfig.BindEnv("address")

	SequenceServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use GAPIC-SHOWCASE_SEQUENCE_TOKEN.")
	SequenceConfig.BindPFlag("token", SequenceServiceCmd.PersistentFlags().Lookup("token"))
	SequenceConfig.BindEnv("token")

	SequenceServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use GAPIC-SHOWCASE_SEQUENCE_API_KEY.")
	SequenceConfig.BindPFlag("api_key", SequenceServiceCmd.PersistentFlags().Lookup("api_key"))
	SequenceConfig.BindEnv("api_key")
}

var SequenceServiceCmd = &cobra.Command{
	Use:   "sequence",
	Short: "Sub-command for Service: Sequence",

	ValidArgs: SequenceSubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := SequenceConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if SequenceConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := SequenceConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := SequenceConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		SequenceClient, err = gapic.NewSequenceClient(ctx, opts...)
		return
	},
}
