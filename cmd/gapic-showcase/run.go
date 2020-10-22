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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func message(err error) string {
	if err == nil {
		return "ok"
	}
	return err.Error()
}

func init() {
	config := RuntimeConfig{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the showcase server",
		Run: func(cmd *cobra.Command, args []string) {
			cmuxServer := CreateAllEndpoints(config)

			done := make(chan os.Signal, 2)
			signal.Notify(done, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGSTOP)
			go func() {
				sig := <-done
				stdLog.Printf("Got signal %q", sig)
				stdLog.Printf("Shutting down server: %s", message(cmuxServer.Shutdown()))

				// TODO: Delete the following line once this PR is
				// merged: https://github.com/soheilhy/cmux/pull/69. The issue is
				// that the user cannot Ctrl-C the main binary, probably due to
				// https://github.com/soheilhy/cmux/pull/69#issuecomment-712928041.
				os.Exit(1)
			}()

			stdLog.Printf("Server finished: %s", message(cmuxServer.Serve()))
		},
	}
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(
		&config.port,
		"port",
		"p",
		":7469",
		"The port that showcase will be served on.")
	runCmd.Flags().StringVarP(
		&config.fallbackPort,
		"fallback-port",
		"f",
		":1337",
		"The port that the fallback-proxy will be served on.")
	runCmd.Flags().StringVar(
		&config.tlsCaCert,
		"mtls-ca-cert",
		"",
		"The Root CA certificate path for custom mutual TLS channel.")
	runCmd.Flags().StringVar(
		&config.tlsCert,
		"mtls-cert",
		"",
		"The server certificate path for custom mutual TLS channel.")
	runCmd.Flags().StringVar(
		&config.tlsKey,
		"mtls-key",
		"",
		"The server private key path for custom mutual TLS channel.")
}
