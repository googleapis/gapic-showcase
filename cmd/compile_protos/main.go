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
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Check if protoc is installed.
	cmd := exec.Command("protoc", "--version")
	if err := cmd.Run(); err != nil {
		log.Fatal("Error: 'protoc' is expected to be installed on the path.")
	}

	// Get proto dependencies
	gopath := os.Getenv("GOPATH")
	showcaseDir := filepath.Join(
		gopath,
		"src",
		"github.com",
		"googleapis",
		"gapic-showcase")
	os.RemoveAll(filepath.Join(showcaseDir, "tmp"))

	err := exec.Command(
		"git",
		"clone",
		"-b",
		"input-contract",
		"https://github.com/googleapis/api-common-protos.git",
		filepath.Join(showcaseDir, "tmp", "api-common-protos"),
	).Run()
	if err != nil {
		log.Fatal(err)
	}

	// Move showcase protos alongside it's dependencies.
	protoDest := filepath.Join(
		showcaseDir,
		"tmp",
		"api-common-protos",
		"google",
		"showcase",
		"v1alpha2")
	err = exec.Command("mkdir", "-p", protoDest).Run()
	if err != nil {
		log.Fatal(err)
	}

	files, err := filepath.Glob(filepath.Join(showcaseDir, "schema", "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + showcaseDir)
	}

	for _, f := range files {
		err = exec.Command("cp", f, protoDest).Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Compile protos
	files, err = filepath.Glob(filepath.Join(protoDest, "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDest)
	}
	params := []string{
		"--go_out=plugins=grpc:" + gopath + "/src",
		"--proto_path=" + filepath.Join(showcaseDir, "tmp", "api-common-protos"),
	}
	params = append(params, files...)
	err = exec.Command("protoc", params...).Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
