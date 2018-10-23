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
	setup()
	stageProtos()
	parallelize(compileDescriptors, tarProtos, compileBinaries)
	teardown()
}

// TODO(landrito): theres a lot of similar things happening here in and in the
// compile protos command. Consolodate the repeated code into a base util file.
func setup() {
	distDir := filepath.Join(showcaseDir(), "dist")
	if err := os.RemoveAll(distDir); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", distDir, err)
	}
	if err := os.MkdirAll(distDir, 0755); err != nil {
		log.Fatalf("Failed to make the directory %s: %v", distDir, err)
	}

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
		"input-contract-edits",
		"https://github.com/googleapis/api-common-protos.git",
		filepath.Join(showcaseDir(), "tmp", "api-common-protos"),
	)

	// Move showcase protos alongside it's dependencies.
	protoDest := filepath.Join(
		showcaseDir(),
		"tmp",
		"api-common-protos",
		"google",
		"showcase",
		"v1alpha2")
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

func compileDescriptors() {
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
		"v1alpha2")

	files, err := filepath.Glob(filepath.Join(protoDest, "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + protoDest)
	}
	command := []string{
		"protoc",
		"--proto_path=" + filepath.Join(showcaseDir(), "tmp", "api-common-protos"),
		"--include_imports",
		"--include_source_info",
		"-o",
		filepath.Join(showcaseDir(), "dist", fmt.Sprintf("gapic-showcase-%s.desc", version())),
	}
	util.Execute(append(command, files...)...)
}

func tarProtos() {
	protoDir := filepath.Join(showcaseDir(), "tmp", "api-common-protos")
	output := filepath.Join(showcaseDir(), "dist", fmt.Sprintf("gapic-showcase-%s-protos.tar.gz", version()))

	util.Execute("tar", "-zcf", output, "-C", protoDir, "google")
}

func compileBinaries() {
	stagingDir := filepath.Join(showcaseDir(), "tmp", "x")
	outputDir := filepath.Join(showcaseDir(), "dist")

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
			filepath.Join(stagingDir, fmt.Sprintf("gapic-showcase-%s-{{.OS}}-{{.Arch}}", version()), "gapic-showcase"),
			"github.com/googleapis/gapic-showcase")
	}

	dirs, _ := filepath.Glob(filepath.Join(stagingDir, "*"))
	for _, dir := range dirs {
		// The windows binaries are suffixed with '.exe'. This allows us to create
		// tarballs of the executables whether or not they contain a suffix.
		files, _ := filepath.Glob(filepath.Join(dir, "gapic-showcase*"))
		util.Execute(
			"tar",
			"-zcf",
			filepath.Join(outputDir, filepath.Base(dir)+".tar.gz"),
			"-C",
			filepath.Dir(files[0]),
			filepath.Base(files[0]))
	}
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

var verOnce sync.Once
var versionMemo string

func version() string {
	verOnce.Do(func() {
		util.ExecuteInDir(showcaseDir(), "go", "get")
		util.ExecuteInDir(showcaseDir(), "go", "install")
		version, err := exec.Command("gapic-showcase", "--version").Output()
		if err != nil {
			log.Fatal("Failed getting showcase version")
		}

		versionMemo = strings.TrimSpace(string(version))
	})
	return versionMemo
}

func parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(functions))
	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			copy()
			waitGroup.Done()
		}(function)
	}
}
