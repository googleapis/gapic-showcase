// Copyright 2019 Google LLC
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

package qualifier

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	packr "github.com/gobuffalo/packr/v2"
	trace "github.com/google/go-trace"
)

var (
	acceptanceSuite *packr.Box // check suite and their files: protos, configs, tests
	schemaSuite     *packr.Box // the schema files that may be needed by various suite checks
)

const (
	// various file types for classifying suite files
	fileTypeSampleConfig = "com.google.api.codegen.SampleConfigProto"
	fileTypeProtobuf     = "proto"

	// external command to run
	cmdSampleTester = "sample-tester"
)

// getAssets loads the assets needed by the qualifier suite. These are Packr boxes with the files
// obtained from the source tree.
func getAssets() error {
	acceptanceSuite = packr.New("acceptance suite", "acceptance_suite")
	schemaSuite = packr.New("schema", "../../../schema")

	if len(acceptanceSuite.List()) == 0 || len(schemaSuite.List()) == 0 {
		return fmt.Errorf("release error: some of the asset boxes are empty")
	}

	traceBox(acceptanceSuite)
	traceBox(schemaSuite)
	return nil
}

// checkDependencies checks that any external dependencies (typically programs that the qualifier
// suite calls) are accessible in this invocation.
func checkDependencies() error {
	notFound := []string{}
	trace.Trace("")

	sampleTesterPath, err := exec.LookPath(cmdSampleTester)
	if err != nil {
		notFound = append(notFound, cmdSampleTester)
	}

	if len(notFound) > 0 {
		msg := fmt.Sprintf("could not find dependencies in $PATH: %q", notFound)
		log.Printf(msg)
		return fmt.Errorf(msg)
	}
	trace.Trace("found %q: %s", cmdSampleTester, sampleTesterPath)
	return nil
}

func traceBox(box *packr.Box) {
	trace.Trace("suite %q has %d entries", box.Name, len(box.List()))
}

// filesByDir contains a list of all the files under the given directory. Each entry in `Files`
// includes the path components for `Directory`.
type filesByDir struct {
	directory string
	files     []string // filenames are paths relative to `Directory`
}

// getFilesByDir returns a list of `filesByDir` objects, one per top-level directory in `box`.
func getFilesByDir(box *packr.Box) []*filesByDir {
	filesInAllDirs := []*filesByDir{}
	var filesInThisDir *filesByDir
	commitThisDir := func() {
		if filesInThisDir != nil {
			filesInAllDirs = append(filesInAllDirs, filesInThisDir)
			filesInThisDir = nil
		}
	}

	// We assume allFiles is returned in a top-down order, so that
	// all files under each first level directory appear
	// contiguously.
	allFiles := box.List()
	previousDir := ""
	for _, oneFile := range allFiles {
		dir := strings.Split(oneFile, string(os.PathSeparator))[0]
		if dir != previousDir {
			commitThisDir()
			filesInThisDir = &filesByDir{directory: dir}
		}
		previousDir = dir
		filesInThisDir.files = append(filesInThisDir.files, oneFile[len(dir)+len(string(os.PathSeparator)):])
	}
	commitThisDir()

	return filesInAllDirs
}

// GetMatchingFiles returns the a list of files in `box` matching
// `srcPath`. `srcPath` should end with `os.PathSeparator` iff it refers
// to a directory.
//
// As a convenience, this also returns an additional value, useful
// for copying files within matched directories:
//
//  * `replacePath`: a copy of `dstPath` (which is assumed to have no
//     trailing separator), with `os.PathSeparator` appended if
//     `srcPath` refers to a directory
func GetMatchingFiles(box *packr.Box, dstPath string, srcPath string) (files []string, replacePath string, err error) {

	trace.Trace("reading %q", srcPath)

	replacePath = dstPath

	// If `srcPath` specifies a single file, match just that.
	if !strings.HasSuffix(srcPath, string(os.PathSeparator)) {
		if !box.Has(srcPath) {
			err = fmt.Errorf("file box %q has no file %q", box.Name, srcPath)
			return
		}
		files = []string{srcPath}
		return
	}

	// Let the caller replace the directory part of the path with
	// the separator included.
	replacePath = dstPath + string(os.PathSeparator)

	for _, entry := range box.List() {
		if strings.HasPrefix(entry, srcPath) {
			files = append(files, entry)
		}
	}
	if len(files) == 0 {
		err = fmt.Errorf("file box %q has no files matching %q", box.Name, srcPath)
		return
	}

	trace.Trace("obtained %d files", len(files))
	return
}
