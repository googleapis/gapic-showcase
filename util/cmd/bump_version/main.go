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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/googleapis/gapic-showcase/util"
	"github.com/spf13/cobra"
)

const currentAPIVersion = "v1beta1"
const currentReleaseVersion = "0.5.0"

// This script updates the release version or API version of files in gapic-showcase.
// This script is used on API and release version bumps. This script must be ran in
// the root directory of gapic-showcase.
//
// Usage: go run ./util/cmd/bump_version/main.go -h
func main() {
	var bumpMajor, bumpMinor, bumpPatch bool
	var newAPI string

	cmd := &cobra.Command{
		Short: "Utility script to bump the API and realease versions in all relevant files.",
		Run: func(cmd *cobra.Command, args []string) {
			if newAPI == "" && !bumpMajor && !bumpMinor && !bumpPatch {
				return
			}

			if bumpMajor || bumpMinor || bumpPatch {
				if !oneof(bumpMajor, bumpMinor, bumpPatch) {
					log.Fatalf("Expected only one of --major, --minor, and --patch.")
				}
				versions := strings.Split(currentReleaseVersion, ".")

				atoi := func(s string) int {
					i, err := strconv.Atoi(s)
					if err != nil {
						log.Fatal(err)
					}
					return i
				}

				major := atoi(versions[0])
				minor := atoi(versions[1])
				patch := atoi(versions[2])

				if bumpMajor {
					major++
					minor = 0
					patch = 0
				}
				if bumpMinor {
					minor++
					patch = 0
				}
				if bumpPatch {
					patch++
				}
				replace(currentReleaseVersion, fmt.Sprintf("%d.%d.%d", major, minor, patch))
			}

			if newAPI != "" && currentAPIVersion != newAPI {
				versionRegexStr := "^([v]\\d+)([p_]\\d+)?((alpha|beta)\\d*)?"
				versionRegex, err := regexp.Compile(versionRegexStr)
				if err != nil {
					log.Fatalf("Unexpected Error compiling version regex: %+v", err)
				}
				if !versionRegex.Match([]byte(newAPI)) {
					log.Fatalf("The API version must conform to the regex: %s", versionRegexStr)
				}

				schemaDir := filepath.Join("schema", "google", "showcase")
				err = os.Rename(filepath.Join(schemaDir, currentAPIVersion), filepath.Join(schemaDir, newAPI))
				if err != nil {
					log.Fatalf("Failed to change proto directory: %+v", err)
				}

				replace(currentAPIVersion, newAPI)
				util.CompileProtos(newAPI)
			}
		},
	}

	cmd.Flags().BoolVarP(
		&bumpMajor,
		"major",
		"",
		false,
		"Pass this flag to bump the major version")
	cmd.Flags().BoolVarP(
		&bumpMinor,
		"minor",
		"",
		false,
		"Pass this flag to bump the minor version")
	cmd.Flags().BoolVarP(
		&bumpPatch,
		"patch",
		"",
		false,
		"Pass this flag to bump the patch version")
	cmd.Flags().StringVarP(
		&newAPI,
		"api",
		"a",
		"",
		"The new API version to set.")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func replace(old, new string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: unable to get working dir: %+v", err)
	}

	filetypes := []string{".go", ".md", ".yml", ".proto"}
	err = filepath.Walk(pwd, replacer(filetypes, old, new))
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func replacer(filetypes []string, old, new string) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(fi.Name(), "CHANGELOG.md") {
			return nil
		}

		matched := false
		for _, t := range filetypes {
			matched = matched || strings.HasSuffix(path, t)
		}
		if !matched {
			return nil
		}

		oldBytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("%v", err)
		}

		newBytes := bytes.Replace(oldBytes, []byte(old), []byte(new), -1)
		if !bytes.Equal(oldBytes, newBytes) {
			if err = ioutil.WriteFile(path, newBytes, 0); err != nil {
				log.Fatalf("%v", err)
			}
		}
		return nil
	}
}

func oneof(bs ...bool) bool {
	t := 0
	for _, b := range bs {
		if b {
			t++
		}
	}
	return t == 1
}
