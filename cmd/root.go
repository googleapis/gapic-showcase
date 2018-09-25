// Copyright 2018 Google LLC
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

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var stdLog, errLog *log.Logger

var rootCmd = &cobra.Command{
	Use:   "gapic-showcase",
	Short: "An API for showcasing GAPIC features",
	Long: `Showcase is an API used to showcase Generate API Client (GAPIC)
features. This tool is used start the API and as well as being used for
running any commands to manage the showcase testing infrastructure.`,
	Version: "0.0.7",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Setup Loggers
	stdLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errLog = log.New(os.Stderr, "", log.Ldate|log.Ltime)

	// Make roots version option only emit the version. This is used in circleci.
	// The template looks weird on purpose. Leaving as a single line causes the
	// output to append an extra character.
	rootCmd.SetVersionTemplate(
		`{{printf "%s" .Version}}
`)
}
