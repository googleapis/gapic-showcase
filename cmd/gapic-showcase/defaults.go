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
	"fmt"
	"os"
)

func init() {
	// Since the showcase server is locally ran, the default address needs to be set
	// since the go gapic assumes port :443.
	services := []string{"ECHO", "IDENTITY", "MESSAGING", "TESTING"}
	envVars := map[string]string{
		"ADDRESS":  "localhost:7469",
		"INSECURE": "true",
	}
	for _, service := range services {
		for key, value := range envVars {
			envName := fmt.Sprintf("GAPIC-SHOWCASE_%s_%s", service, key)
			if os.Getenv(envName) == "" {
				os.Setenv(envName, value)
			}
		}
	}
}
