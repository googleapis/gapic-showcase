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
	"os"

	"github.com/spf13/cobra"
	"github.com/takama/daemon"
)

func init() {
	var port string

	daemonCmd := &cobra.Command{
		Use:   "daemon",
		Short: "Start the Showcase server as a daemon",
		Long: `This command handles the Showcase server in a daemon. The daemon has a
  life cycle in which the daemon must be installed using the install subcommand,
  then started, and finally stopped.`,
	}
	rootCmd.AddCommand(daemonCmd)

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install the Showcase server daemon",
		Run: func(cmd *cobra.Command, args []string) {
			d, _ := daemon.New("Showcase", "Gapic Showcase Service")
			status, err := d.Install("run", "--port", port)
			if err != nil {
				errLog.Println(status, "\nError: ", err)
				os.Exit(1)
			}
			stdLog.Println(status)
		},
	}
	installCmd.Flags().StringVarP(
		&port,
		"port",
		"p",
		":7469",
		"The port that showcase will be served on.")

	subcommands := []*cobra.Command{
		&cobra.Command{
			Use:   "start",
			Short: "Start the Showcase server daemon",
			Run: func(cmd *cobra.Command, args []string) {
				d, _ := daemon.New("Showcase", "Gapic Showcase Service")
				status, err := d.Start()
				if err != nil {
					errLog.Println(status, "\nError: ", err)
					os.Exit(1)
				}
				stdLog.Println(status)
			},
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the Showcase server daemon",
			Run: func(cmd *cobra.Command, args []string) {
				d, _ := daemon.New("Showcase", "Gapic Showcase Service")
				status, err := d.Stop()
				if err != nil {
					errLog.Println(status, "\nError: ", err)
					os.Exit(1)
				}
				stdLog.Println(status)
			},
		},
		&cobra.Command{
			Use:   "uninstall",
			Short: "Uninstall the Showcase server daemon",
			Run: func(cmd *cobra.Command, args []string) {
				d, _ := daemon.New("Showcase", "Gapic Showcase Service")
				status, err := d.Remove()
				if err != nil {
					errLog.Println(status, "\nError: ", err)
					os.Exit(1)
				}
				stdLog.Println(status)
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Fetch the status of the Showcase server daemon",
			Run: func(cmd *cobra.Command, args []string) {
				d, _ := daemon.New("Showcase", "Gapic Showcase Service")
				status, err := d.Status()
				if err != nil {
					errLog.Println(status, "\nError: ", err)
					os.Exit(1)
				}
				stdLog.Println(status)
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Fetch the status of the Showcase server daemon",
			Run: func(cmd *cobra.Command, args []string) {
				d, _ := daemon.New("Showcase", "Gapic Showcase Service")
				status, err := d.Status()
				if err != nil {
					errLog.Println(status, "\nError: ", err)
					os.Exit(1)
				}
				stdLog.Println(status)
			},
		},
	}

	daemonCmd.AddCommand(append(subcommands, installCmd)...)
}
