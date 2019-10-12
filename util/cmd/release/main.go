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
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/googleapis/gapic-showcase/util"
)

// This script is run in CI when a new version tag is pushed to master. This script
// places the compiled proto descriptor set, a tarball of showcase-protos alongside its
// dependencies, and the compiled executables of the gapic-showcase cli tool inside the
// directory "dist"
//
// This script must be run from the root directory of the gapic-showcase repository.
//
// Usage: go run ./util/cmd/release
func main() {
	if err := os.RemoveAll("tmp"); err != nil {
		log.Fatalf("Failed to remove the directory 'tmp': %v", err)
	}

	tmpProtoPath := filepath.Join("tmp", "protos")
	if err := os.MkdirAll(tmpProtoPath, 0755); err != nil {
		log.Fatalf("Failed to make the directory 'tmp': %v", err)
	}
	defer os.RemoveAll("tmp")

	if err := os.RemoveAll("dist"); err != nil {
		log.Fatalf("Failed to remove the directory 'dist': %v", err)
	}
	if err := os.MkdirAll("dist", 0755); err != nil {
		log.Fatalf("Failed to make the directory 'dist': %v", err)
	}

	// Move schema files alongside their dependencies.
	util.Execute("cp", "-rf", filepath.Join("schema", "google"), tmpProtoPath)

	apiPath := filepath.Join("schema", "api-common-protos", "google", "api")
	tmpAPIPath := filepath.Join(tmpProtoPath, "google", "api")
	os.MkdirAll(tmpAPIPath, 0755)
	util.Execute("cp", filepath.Join(apiPath, "annotations.proto"), tmpAPIPath)
	util.Execute("cp", filepath.Join(apiPath, "client.proto"), tmpAPIPath)
	util.Execute("cp", filepath.Join(apiPath, "field_behavior.proto"), tmpAPIPath)
	util.Execute("cp", filepath.Join(apiPath, "http.proto"), tmpAPIPath)
	util.Execute("cp", filepath.Join(apiPath, "resource.proto"), tmpAPIPath)

	longrunningPath := filepath.Join("schema", "api-common-protos", "google", "longrunning")
	tmpLongrunningPath := filepath.Join(tmpProtoPath, "google", "longrunning")
	os.MkdirAll(tmpLongrunningPath, 0755)
	util.Execute("cp", filepath.Join(longrunningPath, "operations.proto"), tmpLongrunningPath)

	rpcPath := filepath.Join("schema", "api-common-protos", "google", "rpc")
	tmpRPCPath := filepath.Join(tmpProtoPath, "google", "rpc")
	os.MkdirAll(tmpRPCPath, 0755)
	util.Execute("cp", filepath.Join(rpcPath, "status.proto"), tmpRPCPath)
	util.Execute("cp", filepath.Join(rpcPath, "error_details.proto"), tmpRPCPath)

	// Find gapic-showcase version
	version, err := exec.Command("go", "run", "./cmd/gapic-showcase", "--version").Output()
	if err != nil {
		log.Fatalf("Failed getting showcase version: %+v", err)
	}

	// Tar Protos
	output := filepath.Join("dist", fmt.Sprintf("gapic-showcase-%s-protos.tar.gz", version))
	util.Execute("tar", "-zcf", output, "-C", tmpProtoPath, "google")

	// Check if protoc is installed.
	if err := exec.Command("protoc", "--version").Run(); err != nil {
		log.Fatal("Error: 'protoc' is expected to be installed on the path.")
	}

	// Compile protos
	files, err := filepath.Glob(filepath.Join(tmpProtoPath, "google", "showcase", "v1beta1", "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + tmpProtoPath)
	}
	command := []string{
		"protoc",
		"--proto_path=" + tmpProtoPath,
		"--include_imports",
		"--include_source_info",
		"-o",
		filepath.Join("dist", fmt.Sprintf("gapic-showcase-%s.desc", version)),
	}
	util.Execute(append(command, files...)...)

	// Get Packr to generate the boxes with local files before compiling. The caller needs to
	// make sure the installation directory is in the $PATH.
	util.Execute("go", "get", "github.com/gobuffalo/packr/v2/packr2")
	util.Execute("packr2")

	// Get cross compiler. Mousetrap is a windows dependency that is not implicitly got since we
	// only get the linux dependencies.  The caller needs to make sure the installation
	// directory is in the $PATH.
	util.Execute("go", "get", "-u", "github.com/mitchellh/gox", "github.com/inconshreveable/mousetrap")

	// Compile binaries
	stagingDir := filepath.Join("tmp", "binaries")
	osArchs := []string{
		"windows/amd64",
		"linux/amd64",
		"darwin/amd64",
		"linux/arm",
	}
	for _, osArch := range osArchs {
		util.Execute(
			"gox",
			fmt.Sprintf("-osarch=%s", osArch),
			"-output",
			filepath.Join(stagingDir, fmt.Sprintf("gapic-showcase-%s-{{.OS}}-{{.Arch}}", version), "gapic-showcase"),
			"github.com/googleapis/gapic-showcase/cmd/gapic-showcase")
	}

	dirs, _ := filepath.Glob(filepath.Join(stagingDir, "*"))
	for _, dir := range dirs {
		// The windows binaries are suffixed with '.exe'. This allows us to create
		// tarballs of the executables whether or not they contain a suffix.
		files, _ := filepath.Glob(filepath.Join(dir, "gapic-showcase*"))
		util.Execute(
			"tar",
			"-zcf",
			filepath.Join("dist", filepath.Base(dir)+".tar.gz"),
			"-C",
			filepath.Dir(files[0]),
			filepath.Base(files[0]))
	}

	// Remove the generated Packr files from the source tree to avoid confusion.
	util.Execute("packr2", "clean")

}
