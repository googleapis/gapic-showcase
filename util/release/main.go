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
)

func main() {
	setup()
	stageProtos()
	parallelize(compileDescriptors, tarProtos, compileBinaries)
	teardown()
}

func setup() {
	os.Chdir(showcaseDir())

	distDir := filepath.Join(showcaseDir(), "dist")
	os.RemoveAll(distDir)
	os.MkdirAll(distDir, 0755)

	tmpDir := filepath.Join(showcaseDir(), "tmp")
	os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)
}

func teardown() {
	os.RemoveAll(filepath.Join(showcaseDir(), "tmp"))
}

func stageProtos() {
	// Get proto dependencies
	execute(
		"git",
		"clone",
		"-b",
		"input-contract",
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
	os.MkdirAll(protoDest, 0755)

	files, err := filepath.Glob(filepath.Join(showcaseDir(), "schema", "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in " + showcaseDir())
	}

	for _, f := range files {
		execute("cp", f, protoDest)
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
		"-o" + filepath.Join(showcaseDir(), "dist", fmt.Sprintf("gapic-showcase-%s.desc", version())),
	}
	execute(append(command, files...)...)
}

func tarProtos() {
	tmpDir := filepath.Join(showcaseDir(), "tmp")
	stagingDir := filepath.Join(tmpDir, fmt.Sprintf("gapic-showcase-v1alpha2-%s-protos", version()))
	output := filepath.Join(showcaseDir(), "dist", fmt.Sprintf("gapic-showcase-v1alpha2-%s-protos.tar.gz", version()))

	os.MkdirAll(stagingDir, 0755)
	execute("cp", "-r", filepath.Join(tmpDir, "api-common-protos", "google"), stagingDir)

	os.Chdir(tmpDir)
	defer os.Chdir(showcaseDir())
	execute("tar", "-zcf", output, fmt.Sprintf("gapic-showcase-v1alpha2-%s-protos", version()))
}

func compileBinaries() {
	stagingDir := filepath.Join(showcaseDir(), "tmp", "x")
	outputDir := filepath.Join(showcaseDir(), "dist")
	execute("go", "get", "github.com/mitchellh/gox")

	// This is a windows dependency that is not implicitly got since we only get
	// the linux dependencies
	execute("go", "get", "github.com/inconshreveable/mousetrap")

	osArchs := []string{
		"windows/amd64",
		"linux/amd64",
		"darwin/amd64",
		"linux/arm",
	}
	for _, osArch := range osArchs {
		execute(
			"gox",
			fmt.Sprintf("-osarch=%s", osArch),
			"-output",
			filepath.Join(stagingDir, fmt.Sprintf("gapic-showcase-v1alpha2-%s-{{.OS}}-{{.Arch}}/gapic-showcase", version())),
			"github.com/googleapis/gapic-showcase")
	}

	dirs, _ := filepath.Glob(filepath.Join(stagingDir, "*"))
	for _, dir := range dirs {
		os.Chdir(dir)
		// The windows binaries are suffixed with '.exe'. This allows us to create
		// tarballs of the executables whether or not they contain a suffix.
		files, _ := filepath.Glob(filepath.Join(dir, "gapic-showcase*"))
		execute(
			"tar",
			"-zcf",
			filepath.Join(outputDir, filepath.Base(dir)+".tar.gz"),
			filepath.Base(files[0]))
	}
	os.Chdir(showcaseDir())
}

func execute(args ...string) {
	log.Print("Executing: ", strings.Join(args, " "))
	if output, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
		log.Fatalf("%s", output)
	}
}

var sDirMu sync.Mutex = sync.Mutex{}
var showcaseDirMemo string

func showcaseDir() string {
	sDirMu.Lock()
	defer sDirMu.Unlock()
	if showcaseDirMemo != "" {
		return showcaseDirMemo
	}
	gopath := os.Getenv("GOPATH")
	showcaseDir := filepath.Join(
		gopath,
		"src",
		"github.com",
		"googleapis",
		"gapic-showcase")
	showcaseDirMemo = showcaseDir
	return showcaseDirMemo
}

var verMu sync.Mutex = sync.Mutex{}
var versionMemo string

func version() string {
	verMu.Lock()
	defer verMu.Unlock()
	if versionMemo != "" {
		return versionMemo
	}
	os.Chdir(showcaseDir())
	execute("go", "get")
	execute("go", "install")
	version, err := exec.Command("gapic-showcase", "--version").Output()
	if err != nil {
		log.Fatal("Failed getting showcase version")
	}

	versionMemo = strings.TrimSpace(string(version[:]))
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
