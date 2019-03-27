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

package util

import (
	"log"
	"os/exec"
)

// Execute runs the given strings as a command.
func Execute(args ...string) {
	if output, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
		log.Fatalf("%s", output)
	}
}

// ExecuteInDir runs the given strings as a command in the directory specified.
func ExecuteInDir(dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("%s", output)
	}
}
