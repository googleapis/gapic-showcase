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
	"strings"
	"sync"

	"github.com/googleapis/gapic-showcase/util"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Setup Dirs.
	dist := filepath.Join(pwd, "dist")
	tmpDir := filepath.Join(pwd, "tmp")
	protoDir := filepath.Join(tmpDir, "protos")
	setup(dist, protoDir)

	// Stage the protos
	showcaseProtoDir := filepath.Join(pwd, "schema")
	util.StageProtos("v1alpha3", showcaseProtoDir, protoDir)

	// Generate the things.
	compileDescriptors(protoDir, dist)
	tarProtos(protoDir, dist)
	compileBinaries(protoDir, tmpDir, dist)

	// Create Github Release
	githubRelease(dist)

	// Push Docker Image
	publishDocker()

	// Teardown.
	teardown(tmpDir)
}

func setup(dist, tmpDir string) {
	if err := os.RemoveAll(dist); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", dist, err)
	}
	if err := os.MkdirAll(dist, 0755); err != nil {
		log.Fatalf("Failed to make the directory %s: %v", dist, err)
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", tmpDir, err)
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Fatalf("Failed to make the directory %s: %v", tmpDir, err)
	}
}

func teardown(tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", tmpDir, err)
	}
}

func compileDescriptors(protoDir, out string) {
	showcaseProtoGlob := filepath.Join(protoDir, "google", "showcase", "**", "*.proto")
	files, err := filepath.Glob(showcaseProtoGlob)
	if err != nil {
		log.Fatal("Error: failed to find protos in " + showcaseProtoGlob)
	}
	command := []string{
		"protoc",
		"--proto_path=" + protoDir,
		"--include_imports",
		"--include_source_info",
		"-o",
		filepath.Join(out, fmt.Sprintf("gapic-showcase-%s.desc", version())),
	}
	util.Execute(append(command, files...)...)
}

func tarProtos(protoDir, out string) {
	output := filepath.Join(out, fmt.Sprintf("gapic-showcase-%s-protos.tar.gz", version()))
	util.Execute("tar", "-zcf", output, "-C", protoDir, "google")
}

func compileBinaries(protoDir, tmpDir, out string) {
	// Mousetrap is a windows dependency that is not implicitly got since
	// we only get the linux dependencies.
	util.Execute("go", "get", "github.com/mitchellh/gox", "github.com/inconshreveable/mousetrap")

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
			filepath.Join(tmpDir, fmt.Sprintf("gapic-showcase-%s-{{.OS}}-{{.Arch}}", version()), "gapic-showcase"),
			"github.com/googleapis/gapic-showcase/cmd/gapic-showcase")
	}

	dirs, _ := filepath.Glob(filepath.Join(tmpDir, "*"))
	for _, dir := range dirs {
		// The windows binaries are suffixed with '.exe'. This allows us to create
		// tarballs of the executables whether or not they contain a suffix.
		files, _ := filepath.Glob(filepath.Join(dir, "gapic-showcase*"))
		util.Execute(
			"tar",
			"-zcf",
			filepath.Join(out, filepath.Base(dir)+".tar.gz"),
			"-C",
			filepath.Dir(files[0]),
			filepath.Base(files[0]))
	}
}

func githubRelease(distFolder string) {
	util.Execute("github.com/tcnksm/ghr")
	util.Execute("ghr", "-prerelease", "-draft", fmt.Sprintf("v%s", version()), distFolder)
}

func publishDocker() {
	imageName := "gcr.io/gapic-images/gapic-showcase"
	util.Execute("docker", "build", "-t", imageName, ".")
	util.Execute("docker", "tag", imageName, fmt.Sprintf("%s:%s", imageName, version()))
	util.Execute("gcloud", "auth", "configure-docker")
	util.Execute("gcloud", "docker", "--", "push", fmt.Sprintf("%s:latest", imageName))
	util.Execute("gcloud", "docker", "--", "push", fmt.Sprintf("%s:%s", imageName, version()))
}

var verOnce sync.Once
var versionMemo string

func version() string {
	verOnce.Do(func() {
		util.Execute("go", "get", "./...")
		util.Execute("go", "install", "./cmd/gapic-showcase")
		version, err := exec.Command("gapic-showcase", "--version").Output()
		if err != nil {
			log.Fatal("Failed getting showcase version")
		}

		versionMemo = strings.TrimSpace(string(version))
	})
	return versionMemo
}
