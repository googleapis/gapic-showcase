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
	"strings"
)

func main() {
	// Check if protoc is installed.
	if err := exec.Command("protoc", "--version").Run(); err != nil {
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
	tmpDir := filepath.Join(showcaseDir, "tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatalf("Error, could not remove path %s: %v", tmpDir, err)
	}

	execute(
		"git",
		"clone",
		"-b",
		"input-contract",
		"https://github.com/googleapis/api-common-protos.git",
		filepath.Join(showcaseDir, "tmp", "api-common-protos"),
	)

	// Move showcase protos alongside it's dependencies.
	protoDest := filepath.Join(
		showcaseDir,
		"tmp",
		"api-common-protos",
		"google",
		"showcase",
		"v1alpha2")
	if err := os.MkdirAll(protoDest, 0755); err != nil {
		log.Fatalf("Error, could not make path %s: %v", protoDest, err)
	}

	files, err := filepath.Glob(filepath.Join(showcaseDir, "schema", "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + showcaseDir)
	}

	for _, f := range files {
		execute("cp", f, protoDest)
	}

	// Compile protos
	files, err = filepath.Glob(filepath.Join(protoDest, "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDest)
	}
	command := []string{
		"protoc",
		"--include_imports",
		"--include_source_info",
		"--go_out=plugins=grpc:" + gopath + "/src",
		"--proto_path=" + filepath.Join(showcaseDir, "tmp", "api-common-protos"),
	}
	execute(append(command, files...)...)
	os.Exit(0)
}

func execute(args ...string) {
	log.Print("Executing: ", strings.Join(args, " "))
	if output, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
		log.Fatalf("%s", output)
	}
}
