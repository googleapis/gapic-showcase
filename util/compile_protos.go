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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	outDir, err := ioutil.TempDir(os.TempDir(), "gapic-showcase")
	if err != nil {
		log.Fatalf("Error: unable to create a temporary dir: %+v\n", err)
	}
	defer os.RemoveAll(outDir)

	protos := filepath.Join("schema", "google", "showcase", version, "*.proto")
	files, err := filepath.Glob(protos)
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protos)
	}

	// Run protoc
	command := []string{
		"protoc",
		"--experimental_allow_proto3_optional",
		"--proto_path=schema/api-common-protos",
		"--proto_path=schema",
		// Disable go_cli generation until the generator supports it.
		//
		// TODO(ndietz): reenable once issue is resolved:
		// https://github.com/googleapis/gapic-generator-go/issues/378
		// "--go_cli_out=" + filepath.Join("cmd", "gapic-showcase"),
		// "--go_cli_opt=root=gapic-showcase",
		// "--go_cli_opt=gapic=github.com/googleapis/gapic-showcase/client",
		// "--go_cli_opt=fmt=false",
		"--go_gapic_out=" + outDir,
		"--go_gapic_opt=go-gapic-package=github.com/googleapis/gapic-showcase/client;client",
		"--go_gapic_opt=grpc-service-config=schema/google/showcase/v1beta1/showcase_grpc_service_config.json",
		"--go_out=plugins=grpc:" + outDir,
	}
	Execute(append(command, files...)...)

	// Copy generated code back into repo.
	tempClient := filepath.Join(outDir, "github.com", "googleapis", "gapic-showcase", "client")
	tempServer := filepath.Join(outDir, "github.com", "googleapis", "gapic-showcase", "server")
	command = []string{
		"cp",
		"-r",
		tempClient,
		tempServer,
		pwd,
	}
	Execute(command...)

	// Fix some generated errors.
	fixes := []struct {
		file string
		fix  string
	}{
		{
			"cmd/gapic-showcase/verify-test.go",
			"/ByteSliceVar/d",
		},
		{
			"cmd/gapic-showcase/wait.go",
			"s/EndEnd_time/EndEndTime/g",
		},
	}
	command = []string{
		"sed",
		"-i.bak",
	}
	for _, f := range fixes {
		Execute(append(command, f.fix, f.file)...)

		// Remove the backup file.
		Execute("rm", fmt.Sprintf("%s.bak", f.file))
	}

	// Format generated output
	Execute("go", "fmt", "./...")
}
