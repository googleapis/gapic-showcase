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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: `print_changes ${VERSION}`")
	}
	version := os.Args[1]

	file, err := os.Open("CHANGELOG.md")
	if err != nil {
		log.Fatalf("Failed to open CHANGELOG.md: %+v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var line string
	versionFound := false
	changes := []string{}
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		if strings.Contains(line, version) {
			versionFound = true
			continue
		}

		if versionFound && strings.HasPrefix(line, "#") {
			break
		}
		if versionFound && line != "" {
			changes = append(changes, line)
		}
	}
	if !versionFound {
		log.Fatalf("Could not find version: %s", version)
	}
	fmt.Print(strings.Join(changes, "\n"))
}
