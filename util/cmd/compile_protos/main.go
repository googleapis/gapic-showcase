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
	"path/filepath"
	"strings"

	"github.com/googleapis/gapic-showcase/util"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Stage the protos
	protoDir := filepath.Join(pwd, "schema")
	protoDest := filepath.Join(pwd, "tmp", "protos")
	util.StageProtos("v1alpha3", protoDir, protoDest)

	// Assumes that this is being ran at the root of the gapic-showcase project.
	gopath := filepath.Join(pwd, "..", "..", "..")
	compileProtos(protoDest, gopath)
	generateCli(protoDest, filepath.Join(pwd, "cmd", "gapic-showcase"))
	generateClient(protoDest, gopath)

	// Remove the staged protos.
	if err := os.RemoveAll(filepath.Join(pwd, "tmp")); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", filepath.Join(pwd, "tmp"), err)
	}
}

func compileProtos(protoDir, out string) {
	// Compile protos
	files, err := filepath.Glob(filepath.Join(protoDir, "google/showcase/**/*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in: " + protoDir)
	}
	command := []string{
		"protoc",
		"--go_out=plugins=grpc:" + out,
		"--proto_path=" + protoDir,
	}
	util.Execute(append(command, files...)...)
}

func generateCli(protoDir, out string) {
	// Delete the files that were generated prior. This can't be done wholesale since
	// there are some hand-written go files in the output path.
	files, err := filepath.Glob(filepath.Join(out, "*.go"))
	if err != nil {
		log.Fatal("Error: failed to clean up cli path")
	}
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}

		firstLine := bufio.NewScanner(file).Text()
		file.Close()
		if strings.Contains(firstLine, "// Code generated. DO NOT EDIT.") {
			os.Remove(f)
		}
	}

	// Generate the CLI.
	files, err = filepath.Glob(filepath.Join(protoDir, "google/showcase/**/*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDir)
	}
	command := []string{
		"protoc",
		"--proto_path=" + protoDir,
		"--go_cli_out=" + out,
		"--go_cli_opt=root=gapic-showcase",
		"--go_cli_opt=gapic=github.com/googleapis/gapic-showcase/client",
		"--go_cli_opt=fmt=false",
	}
	util.Execute(append(command, files...)...)
}

func generateClient(protoDir, out string) {
	files, err := filepath.Glob(filepath.Join(protoDir, "google/showcase/**/*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDir)
	}
	command := []string{
		"protoc",
		"--proto_path=" + protoDir,
		"--go_gapic_out=" + out,
		"--go_gapic_opt=go-gapic-package=github.com/googleapis/gapic-showcase/client;client",
	}
	util.Execute(append(command, files...)...)
}
