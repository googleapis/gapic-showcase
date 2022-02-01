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
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/googleapis/gapic-showcase/util"
)

var version string

func init() {
	flag.StringVar(&version, "version", "", "the version tag [required]")
}

// This script is ran in CI when a new version tag is pushed to master. This script
// places the compiled proto descriptor set, a tarball of showcase-protos alongside it's
// dependencies, and the compiled executables of the gapic-showcase cli tool inside the
// directory "dist"
//
// This script must be ran from the root directory of the gapic-showcase repository.
//
// Usage: go run ./util/cmd/release -version vX.Y.Z
func main() {
	flag.Parse()
	if version == "" {
		log.Fatalln("Missing required flag: -version")
	}

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
	util.Execute("cp", filepath.Join(apiPath, "routing.proto"), tmpAPIPath)

	longrunningPath := filepath.Join("schema", "api-common-protos", "google", "longrunning")
	tmpLongrunningPath := filepath.Join(tmpProtoPath, "google", "longrunning")
	os.MkdirAll(tmpLongrunningPath, 0755)
	util.Execute("cp", filepath.Join(longrunningPath, "operations.proto"), tmpLongrunningPath)

	rpcPath := filepath.Join("schema", "api-common-protos", "google", "rpc")
	tmpRPCPath := filepath.Join(tmpProtoPath, "google", "rpc")
	os.MkdirAll(tmpRPCPath, 0755)
	util.Execute("cp", filepath.Join(rpcPath, "status.proto"), tmpRPCPath)
	util.Execute("cp", filepath.Join(rpcPath, "error_details.proto"), tmpRPCPath)

	// Copy gRPC ServiceConfig as the source of retry config.
	retrySrc := filepath.Join("schema", "google", "showcase", "v1beta1", "showcase_grpc_service_config.json")
	util.Execute("cp", retrySrc, "dist")

	// Copy compliance suite for easy access by generators.
	complianceSrc := filepath.Join("schema", "google", "showcase", "v1beta1", "compliance_suite.json")
	util.Execute("cp", complianceSrc, "dist")

	// Copy the API Service config for easy access by generators.
	serviceYamlSrc := filepath.Join("schema", "google", "showcase", "v1beta1", "showcase_v1beta1.yaml")
	util.Execute("cp", serviceYamlSrc, "dist")

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
		"--experimental_allow_proto3_optional",
		"--proto_path=" + tmpProtoPath,
		"--include_imports",
		"--include_source_info",
		"-o",
		filepath.Join("dist", fmt.Sprintf("gapic-showcase-%s.desc", version)),
	}
	util.Execute(append(command, files...)...)
}
