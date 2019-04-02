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
	"strings"
)

// CompileProtos regenerates all of the generated source code for the Showcase
// API including the generated messages, gRPC services, go gapic clients,
// and the generated CLI. This must be ran from the root directory
// of the gapic-showcase repository.
func CompileProtos(version string) {
	// Check if protoc is installed.
	if err := exec.Command("protoc", "--version").Run(); err != nil {
		log.Fatal("Error: 'protoc' is expected to be installed on the path.")
	}

	// Setup paths
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: unable to get working dir: %+v", err)
	}
	gosrc, err := filepath.Abs(
		filepath.Join(pwd, "..", "..", ".."),
	)
	if err != nil || !strings.HasSuffix(gosrc, filepath.Join("go", "src")) {
		log.Fatal("Error: could not find the go src path from the current working dir.")
	}

	protos := filepath.Join("schema", "google", "showcase", version, "*.proto")
	files, err := filepath.Glob(protos)
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protos)
	}

	// Run protoc
	command := []string{
		"protoc",
		"--proto_path=schema/api-common-protos",
		"--proto_path=schema",
		"--go_cli_out=" + filepath.Join("cmd", "gapic-showcase"),
		"--go_cli_opt=root=gapic-showcase",
		"--go_cli_opt=gapic=github.com/googleapis/gapic-showcase/client",
		"--go_cli_opt=fmt=false",
		"--go_gapic_out=" + gosrc,
		"--go_gapic_opt=go-gapic-package=github.com/googleapis/gapic-showcase/client;client",
		"--go_out=plugins=grpc:" + gosrc,
	}
	Execute(append(command, files...)...)
}
