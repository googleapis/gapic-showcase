package main

import (
	"strconv"
	"time"

	trace "github.com/google/go-trace"
	"github.com/spf13/cobra"

	"github.com/googleapis/gapic-showcase/cmd/gapic-showcase/qualifier"
)

func init() {
	trace.On(false) // set to true for debugging

	var timestamp string
	timestamp = time.Now().Format("20060102.150405")
	trace.Trace("timestamp = %q", timestamp)

	settings := &qualifier.Settings{
		Timestamp:    timestamp,
		Verbose:      Verbose,
		ShowcasePort: 7469,
	}
	qualifyCmd := &cobra.Command{
		Use:   "qualify",
		Short: "Tests a provided GAPIC generator against an acceptance suite",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			servers := RunShowcase(strconv.Itoa(settings.ShowcasePort), "")
			settings.Language = args[0]
			trace.Trace("settings: %v", settings)
			qualifier.Run(settings)
			servers.Shutdown()
		},
	}
	rootCmd.AddCommand(qualifyCmd)
	qualifyCmd.Flags().StringVarP(
		&settings.PluginDirectory,
		"dir",
		"d",
		"",
		"The directory in which to find the protoc plugin implementing the given GAPIC generator")
	qualifyCmd.Flags().StringVarP(
		&settings.PluginOptions,
		"options",
		"o",
		"",
		"The options to pass to the protoc plugin in order to generate a GAPIC for the showcase Echo service")

}
