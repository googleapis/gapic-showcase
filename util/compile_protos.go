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
	"os"
	"os/exec"
	"path/filepath"
)

func CompileProtos(version string) {
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

	Execute(
		"git",
		"clone",
		"-b",
		"input-contract-edits",
		"https://github.com/googleapis/api-common-protos.git",
		filepath.Join(showcaseDir, "tmp", "api-common-protos"),
	)

	// Move showcase protos alongside its dependencies.
	protoDest := filepath.Join(
		showcaseDir,
		"tmp",
		"api-common-protos",
		"google",
		"showcase",
		version)
	if err := os.MkdirAll(protoDest, 0755); err != nil {
		log.Fatalf("Error, could not make path %s: %v", protoDest, err)
	}

	files, err := filepath.Glob(filepath.Join(showcaseDir, "schema", "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + showcaseDir)
	}

	for _, f := range files {
		Execute("cp", f, protoDest)
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
	Execute(append(command, files...)...)
}
