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

const CURRENT_API = "v1alpha2"
const CURRENT_RELEASE = "0.0.7"

func main() {
	var bumpMajor, bumpMinor, bumpPatch bool
	var newApi string

	// This is definitely overkill, I am just lazy and didn't want to write the
	// code to parse the args. ¯\_(ツ)_/¯
	cmd := &cobra.Command{
		Short: "Utility script to bump the API and realease versions in all relevant files.",
		Run: func(cmd *cobra.Command, args []string) {
			if newApi == "" && !bumpMajor && !bumpMinor && !bumpPatch {
				return
			}

			if bumpMajor || bumpMinor || bumpPatch {
				re := regexp.MustCompile("[0-9]+")
				versions := re.FindAllString(CURRENT_RELEASE, 3)

				major, _ := strconv.Atoi(versions[0])
				newMajor := major + btoi(bumpMajor)

				minor, _ := strconv.Atoi(versions[1])
				newMinor := minor + btoi(bumpMinor)

				patch, _ := strconv.Atoi(versions[2])
				newPatch := patch + btoi(bumpPatch)

				newRelease := fmt.Sprintf("%d.%d.%d", newMajor, newMinor, newPatch)
				replace(CURRENT_RELEASE, newRelease)
			}

			if newApi != "" && CURRENT_API != newApi {
				replace(CURRENT_API, newApi)
				util.CompileProtos(newApi)
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
		&newApi,
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
	showcaseDir := filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"googleapis",
		"gapic-showcase")

	filePatterns := []string{"*.go", "*.md", "*.yml"}

	for _, pattern := range filePatterns {
		err := filepath.Walk(showcaseDir, replacer(pattern, old, new))
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func replacer(pattern, old, new string) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if old == new {
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, fi.Name())

		if matched {
			replaced := false
			read, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("%v", err)
			}

			if strings.Contains(string(read), old) {
				replaced = true
			}

			newContents := strings.Replace(string(read), old, new, -1)

			err = ioutil.WriteFile(path, []byte(newContents), 0)
			if err != nil {
				log.Fatalf("%v", err)
			}

			if replaced {
				log.Printf("Replaced %s with %s in %s", old, new, path)
			}
		}
		return nil
	}
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
