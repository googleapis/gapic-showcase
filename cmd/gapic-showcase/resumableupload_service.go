// Code generated. DO NOT EDIT.

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gapic "github.com/googleapis/gapic-showcase/client"
)

var ResumableUploadConfig *viper.Viper
var ResumableUploadClient *gapic.ResumableUploadClient
var ResumableUploadSubCommands []string = []string{
	"upload-media",
}

func init() {
	rootCmd.AddCommand(ResumableUploadServiceCmd)

	ResumableUploadConfig = viper.New()
	ResumableUploadConfig.SetEnvPrefix("GAPIC-SHOWCASE_RESUMABLEUPLOAD")
	ResumableUploadConfig.AutomaticEnv()

	ResumableUploadServiceCmd.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use GAPIC-SHOWCASE_RESUMABLEUPLOAD_INSECURE. Must be used with \"address\" option")
	ResumableUploadConfig.BindPFlag("insecure", ResumableUploadServiceCmd.PersistentFlags().Lookup("insecure"))
	ResumableUploadConfig.BindEnv("insecure")

	ResumableUploadServiceCmd.PersistentFlags().String("address", "", "Set API address used by client. Or use GAPIC-SHOWCASE_RESUMABLEUPLOAD_ADDRESS.")
	ResumableUploadConfig.BindPFlag("address", ResumableUploadServiceCmd.PersistentFlags().Lookup("address"))
	ResumableUploadConfig.BindEnv("address")

	ResumableUploadServiceCmd.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use GAPIC-SHOWCASE_RESUMABLEUPLOAD_TOKEN.")
	ResumableUploadConfig.BindPFlag("token", ResumableUploadServiceCmd.PersistentFlags().Lookup("token"))
	ResumableUploadConfig.BindEnv("token")

	ResumableUploadServiceCmd.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use GAPIC-SHOWCASE_RESUMABLEUPLOAD_API_KEY.")
	ResumableUploadConfig.BindPFlag("api_key", ResumableUploadServiceCmd.PersistentFlags().Lookup("api_key"))
	ResumableUploadConfig.BindEnv("api_key")
}

var ResumableUploadServiceCmd = &cobra.Command{
	Use:       "resumableupload",
	Short:     "A service showcasing universal resumable upload...",
	Long:      "A service showcasing universal resumable upload protocol support.",
	ValidArgs: ResumableUploadSubCommands,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := ResumableUploadConfig.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if ResumableUploadConfig.GetBool("insecure") {
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := ResumableUploadConfig.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := ResumableUploadConfig.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		ResumableUploadClient, err = gapic.NewResumableUploadClient(ctx, opts...)
		return
	},
}
