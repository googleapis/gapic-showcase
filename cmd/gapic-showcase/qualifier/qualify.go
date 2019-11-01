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

package qualifier

import (
	"fmt"
	"os"

	trace "github.com/google/go-trace"
)

// Settings contains the settings for a run of the qualification suite.
type Settings struct {
	Generator
	ShowcasePort int
	Timestamp    string // for tracing and diagnostics
	Verbose      bool
}

// Run runs the qualification suite using the values in `settings`.
func Run(settings *Settings) {
	// TODO: Return an error rather than exiting. We can return an error type ErrorCode that
	// wraps the current errors and includes the appropriate return code, and change main() so
	// that if the error it gets is of type ErrorCode, it exits with that code.
	const (
		retCodeSuccess = iota
		retCodeInternalError
		retCodeFailedDependencies
		retCodeUsageError
		retCodeScenarioFailure
	)

	if err := getAssets(); err != nil {
		os.Exit(retCodeInternalError)
	}

	if err := checkDependencies(); err != nil {
		os.Exit(retCodeFailedDependencies)
	}

	allScenarios := getTestScenarios(settings)

	success := true
	for _, scenario := range allScenarios {
		if err := scenario.Run(); err != nil {
			os.Exit(retCodeInternalError)
		}
		status := "SUCCESS"
		if !scenario.Success() {
			success = false
			status = "FAILURE"
		}
		fmt.Printf("scenario %q: %s    %s\n", scenario.name, status, scenario.sandbox)
	}
	if !success {
		os.Exit(retCodeScenarioFailure)
	}
}

// getTestScenarios returns a list of Scenario as found in the acceptanceSuite.
func getTestScenarios(settings *Settings) []*Scenario {
	scenarios := []*Scenario{}
	configDirs := getFilesByDir(acceptanceSuite)
	for _, config := range configDirs {
		scenarios = append(scenarios, &Scenario{
			name:         config.directory,
			timestamp:    settings.Timestamp,
			showcasePort: settings.ShowcasePort,
			generator:    &settings.Generator,
			files:        config.files,
			fileBox:      acceptanceSuite,
			schemaBox:    schemaSuite,
		})

	}
	trace.Trace("adding scenarios %#v", scenarios)
	return scenarios
}
