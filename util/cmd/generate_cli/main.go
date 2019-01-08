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
	"sync"

	"github.com/googleapis/gapic-showcase/util"
)

func main() {
	setup()
	stageProtos()
	generate_gapic()
	generate_cli()
	teardown()
}

func setup() {
	tmpDir := filepath.Join(showcaseDir(), "tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", tmpDir, err)
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Fatalf("Failed to make the directory %s: %v", tmpDir, err)
	}
}

func teardown() {
	tmpDir := filepath.Join(showcaseDir(), "tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", tmpDir, err)
	}
}

func stageProtos() {
	// Get proto dependencies
	util.Execute(
		"git",
		"clone",
		"-b",
		"input-contract",
		"https://github.com/googleapis/api-common-protos.git",
		filepath.Join(showcaseDir(), "tmp", "api-common-protos"),
	)

	// Move showcase protos alongside its dependencies.
	protoDest := filepath.Join(
		showcaseDir(),
		"tmp",
		"api-common-protos",
		"google",
		"showcase",
		"v1alpha3")
	if err := os.MkdirAll(protoDest, 0755); err != nil {
		log.Fatalf("Failed to make the dir %s: %v", protoDest, err)
	}

	files, err := filepath.Glob(filepath.Join(showcaseDir(), "schema", "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + showcaseDir())
	}

	for _, f := range files {
		util.Execute("cp", f, protoDest)
	}
}

func generate_gapic() {
	// Check if protoc is installed.
	if err := exec.Command("protoc", "--version").Run(); err != nil {
		log.Fatal("Error: 'protoc' is expected to be installed on the path.")
	}

	// Compile protos
	protoDest := filepath.Join(
		showcaseDir(),
		"tmp",
		"api-common-protos",
		"google",
		"showcase",
		"v1alpha3")

	files, err := filepath.Glob(filepath.Join(protoDest, "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDest)
	}
	command := []string{
		"protoc",
		"--proto_path=" + filepath.Join(showcaseDir(), "tmp", "api-common-protos"),
		"--go_gapic_out=" + filepath.Join(os.Getenv("GOPATH"), "src"),
		"--go_gapic_opt=go-gapic-package=github.com/googleapis/gapic-showcase/client;client",
	}
	util.Execute(append(command, files...)...)
}

func generate_cli() {
	// Check if protoc is installed.
	if err := exec.Command("protoc", "--version").Run(); err != nil {
		log.Fatal("Error: 'protoc' is expected to be installed on the path.")
	}

	// Compile protos
	protoDest := filepath.Join(
		showcaseDir(),
		"tmp",
		"api-common-protos",
		"google",
		"showcase",
		"v1alpha3")

	files, err := filepath.Glob(filepath.Join(protoDest, "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDest)
	}
	command := []string{
		"protoc",
		"--proto_path=" + filepath.Join(showcaseDir(), "tmp", "api-common-protos"),
		"--go_cli_out=" + filepath.Join(showcaseDir(), "cmd", "showcase"),
		"--go_cli_opt=root=showcase",
		"--go_cli_opt=gapic=github.com/googleapis/gapic-showcase/client",
	}
	util.Execute(append(command, files...)...)
}

var sOnce sync.Once
var showcaseDirMemo string

func showcaseDir() string {
	sOnce.Do(func() {
		gopath := os.Getenv("GOPATH")
		showcaseDir := filepath.Join(
			gopath,
			"src",
			"github.com",
			"googleapis",
			"gapic-showcase")
		showcaseDirMemo = showcaseDir
	})
	return showcaseDirMemo
}
