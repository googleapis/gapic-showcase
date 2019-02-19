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
	"strconv"
	"strings"

	"github.com/googleapis/gapic-showcase/util"
	"github.com/spf13/cobra"
)

const CURRENT_API = "v1alpha3"
const CURRENT_RELEASE = "0.0.11"

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
				if !oneof(bumpMajor, bumpMinor, bumpPatch) {
					log.Fatalf("Expected only one of --major, --minor, and --patch.")
				}
				versions := strings.Split(CURRENT_RELEASE, ".")

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
					major += 1
					minor = 0
					patch = 0
				}
				if bumpMinor {
					minor += 1
					patch = 0
				}
				if bumpPatch {
					patch += 1
				}
				replace(CURRENT_RELEASE, fmt.Sprintf("%d.%d.%d", major, minor, patch))
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

	filetypes := []string{".go", ".md", ".yml"}
	err := filepath.Walk(showcaseDir, replacer(filetypes, old, new))
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
			t += 1
		}
	}
	return t == 1
}
