// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"strconv"
	"time"

	trace "github.com/google/go-trace"
	"github.com/spf13/cobra"

	"github.com/googleapis/gapic-showcase/cmd/gapic-showcase/qualifier"
)

func init() {
	var timestamp string
	timestamp = time.Now().Format("20060102.150405")
	trace.Trace("timestamp = %q", timestamp)

	settings := &qualifier.Settings{
		Timestamp:    timestamp,
		Verbose:      Verbose,
		ShowcasePort: 7469,
	}
	qualifyCmd := &cobra.Command{
		Use:   "qualify [language]",
		Short: "Tests a provided GAPIC generator against an acceptance suite",
		Long: `qualify will execute a suite of acceptance checks against the GAPIC generator for the specified language.
This confirms that the generator behaves and emits artifacts as specified in generator requirements under
 a variety of inputs for various types of RPCs. Each acceptance check typically attempts to generate client
libraries and corresponding standalone samples for the Showcase "Echo" service. The generator is invoked
 as a protoc plugin; its location  may be specified via --dir, and additional generator options
 that are needed to successfully generate the GAPIC for this API may be specified via --options,`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Consider moving this to a more central place for debugging all of
			// showcase.
			trace.On(false) // set to true for debugging

			servers := RunShowcase(strconv.Itoa(settings.ShowcasePort), "")
			defer servers.Shutdown()

			settings.Language = args[0]
			trace.Trace("settings: %v", settings)
			qualifier.Run(settings)
		},
	}
	rootCmd.AddCommand(qualifyCmd)
	qualifyCmd.Flags().StringVarP(
		&settings.Directory,
		"dir",
		"d",
		"",
		"The directory in which to find the protoc plugin implementing the given GAPIC generator")
	qualifyCmd.Flags().StringVarP(
		&settings.Options,
		"options",
		"o",
		"",
		"The options to pass to the generator in order to generate a GAPIC for the Showcase \"Echo\" service")

}
