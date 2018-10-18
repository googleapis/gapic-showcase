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
	"strings"
)

func Execute(args ...string) {
	log.Print("Executing: ", strings.Join(args, " "))
	if output, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
		log.Fatalf("%s", output)
	}
}

func ExecuteInDir(dir string, args ...string) {
	log.Printf("Executing in %s: %s", dir, strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("%s", output)
	}
}
