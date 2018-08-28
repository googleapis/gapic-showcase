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
	showcaseDir := gopath + "/" + "src/github.com/googleapis/gapic-showcase"
	os.RemoveAll(showcaseDir + "/tmp")

	cmd = exec.Command(
		"git",
		"clone",
		"-b",
		"input-contract",
		"https://github.com/googleapis/api-common-protos.git",
		showcaseDir+"/tmp/api-common-protos")
	runAndLog(cmd)

	// Move showcase protos alongside it's dependencies.
	protoDest := showcaseDir + "/tmp/api-common-protos/google/showcase/v1alpha1"
	cmd = exec.Command("mkdir", "-p", protoDest)
	runAndLog(cmd)

	files, err := filepath.Glob(showcaseDir + "/schema/*.proto")
	if err != nil {
		log.Fatal("Error: failed to find protos in " + showcaseDir)
	}

	for _, f := range files {
		cmd = exec.Command("cp", f, protoDest)
		runAndLog(cmd)
	}

	// Compile protos
	files, err = filepath.Glob(protoDest + "/*.proto")
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDest)
	}
	params := []string{
		"--go_out=plugins=grpc:" + gopath + "/src",
		"--proto_path=" + showcaseDir + "/tmp/api-common-protos",
	}
	params = append(params, files...)
	cmd = exec.Command("protoc", params...)
	runAndLog(cmd)

	os.Exit(0)
}

func runAndLog(cmd *exec.Cmd) {
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
	}
}
