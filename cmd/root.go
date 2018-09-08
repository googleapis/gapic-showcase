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

var StdLog, ErrLog *log.Logger
var Port string

var rootCmd = &cobra.Command{
	Use:   "gapic-showcase",
	Short: "An API for showcasing GAPIC features",
	Long: `Showcase is an API used to showcase Generate API Client (GAPIC)
features. This tool is used start the API and as well as being used for
running any commands to manage the showcase testing infrastructure.`,
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
	StdLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	ErrLog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}
