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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	packr "github.com/gobuffalo/packr/v2"
	trace "github.com/google/go-trace"
	yaml "gopkg.in/yaml.v2"
)

// Constants pertaining to the syntax of include files. Include files are akin to symbolic links:
// they reference content that will be copied at run-time from the indicated location to the include
// file's location under the indicated names. For example, an include file called "include.boren"
// with contents "api-common-protos/google/" will be replaced at run-time with a directory called
// "boren" whose contents will be a copy of everything under "api-common-protos/google/"
//
// We need some mechanism for including common protos because the protos under test may reference
// them. This particular scheme is also useful in that it allows us to explicitly pass to the
// generator only the protos explicitly listed in the acceptance check, rather than also having to
// include all the protos that are transitively included.
const (
	// In a qualifier suite directory, any file that begins with this prefix is considered an
	// include file. The rest of the filename after the prefix is the destination file or
	// directory name to which the files referenced by the content of the include file will be
	// copied.
	includeFilePrefix = "include."

	// In the contents of an include file, this symbol will be treated as the directory
	// separator. At run time, it will be replaced by the appropriate os.PathSeparator.
	includeFileDirectorySeparator = '/'
)

type Scenario struct {
	// Every element in `files` corresponds to a file in `fileBox` under
	// a directory called `name`, which is NOT part of the string
	// element itself.
	name        string              // name of the scenario
	timestamp   string              // timestamp for debug and trace messages
	files       []string            // a flat list of files in the scenario
	fileBox     *packr.Box          // the files explicitly specified in the scenario
	schemaBox   *packr.Box          // additional files that may be referenced by include files in a scenario
	sandbox     string              // the sandbox created to run this scenario
	filesByType map[string][]string // scenario files partitioned by file type

	generator    *Generator
	showcasePort int // the port of the running showcase server against which samples will run

	generationOutput       []byte           // combined output of the generation process
	generationProcessState *os.ProcessState // info on the completed generation process
	generationPassed       bool             // whether all generation checks passed
	sampleTestsPassed      bool             // whether all sample tests specified passed
}

// Success returns true iff all the expectations for this scenario were met.
func (scenario *Scenario) Success() bool {
	return scenario.generationPassed && // the generation checks passed and...
		(!scenario.generationProcessState.Success() || // ...either generation failed, or it succeeded with ...
			(len(scenario.filesByType[fileTypeSampleTest]) == 0 || scenario.sampleTestsPassed)) // all sample tests, if any, passing.
}

// Run runs this particular acceptance scenario by creating a sandbox in which it generates
// libraries and samples and then tests that the generation and any sample tests performed as
// expected.
func (scenario *Scenario) Run() error {
	if err := scenario.classifyConfigs(); err != nil {
		return fmt.Errorf("could not classify config files for scenario %q: %w", scenario.name, err)
	}

	if err := scenario.createSandbox(); err != nil {
		return fmt.Errorf("could not create sandbox for scenario %q: %w", scenario.name, err)
	}
	trace.Trace("created sandbox: %q", scenario.sandbox)

	if err := scenario.Generate(); err != nil {
		return err
	}
	scenario.CheckGeneration()
	scenario.RunTests()
	return nil
}

// Generate generates the client libraries and samples for this scenario.
func (scenario *Scenario) Generate() error {
	var err error
	scenario.generationOutput, scenario.generationProcessState, err = scenario.generator.Run(scenario.sandbox, scenario.filesByType)
	trace.Trace("run error: %v", err)
	trace.Trace("run scenario.generationOutput: %s", scenario.generationOutput)
	return err
}

// CheckGeneration verifies that generation matched configured expectations when it succeeded or
// failed.
func (scenario *Scenario) CheckGeneration() {
	// TODO: Fill in in order to handle expected, configured failures
	trace.Trace("CheckGeneration")
	scenario.generationPassed = scenario.generationProcessState.Success()
}

// RunTests runs tests on the generated samples.
func (scenario *Scenario) RunTests() {
	// TODO: Fill in
	trace.Trace("RunTests (TODO)")
}

// sandboxPath returns `relativePath` re-rooted at `scenario.sandbox`,
// which must have been previously set.
func (scenario *Scenario) sandboxPath(relativePath string) string {
	return filepath.Join(scenario.sandbox, relativePath)
}

// fileBoxPath returns relativePath re-rooted at scenario.name, so that it references an object in
// `scenario.fileBox`.
func (scenario *Scenario) fileBoxPath(relativePath string) string {
	return filepath.Join(scenario.name, relativePath)
}

// fromFileBox returns the file named by `relativePath` from `scenario.fileBox`.
func (scenario *Scenario) fromFileBox(relativePath string) ([]byte, error) {
	return scenario.fileBox.Find(scenario.fileBoxPath(relativePath))
}

// classifyConfigs classifies all the files in the scenario, storing the results in
// `scenario.filesByType`. Each file is classified as either an "include" file or a "scenario"
// file. Scenario files are also classified by the type of data they contain. Thus, "scenario" files
// have at least two labels (ie multiple entries in `scenario.filesByType`).
func (scenario *Scenario) classifyConfigs() (err error) {
	scenario.filesByType = make(map[string][]string)
	// TODO: process `include.*` files first
	for _, file := range scenario.files {
		types, err := scenario.getFileTypes(file)
		if err != nil {
			return err
		}
		for _, fType := range types {
			similarFiles := scenario.filesByType[fType]
			scenario.filesByType[fType] = append(similarFiles, file)
		}
		trace.Trace("%s: type %q", file, types)
	}
	return nil
}

// getFileTypes gets the various type labels for `file`. The type labels are either the
// singleton list {"include"}, or a list {"scenario", ...} that includes the various types of data
// included in that file.
func (scenario *Scenario) getFileTypes(file string) (types []string, err error) {

	// Identify include files in the scenario. Include files always begin with `includeFilePrefix`
	pathParts := strings.Split(file, string(os.PathSeparator))
	if len(pathParts) > 0 && strings.HasPrefix(pathParts[len(pathParts)-1], includeFilePrefix) {
		return []string{"include"}, nil
	}

	extension := filepath.Ext(file)
	switch extension {
	case ".proto":
		types = []string{"proto"}
	case ".yaml", ".yml":
		trace.Trace("reading %s", file)
		content, err := scenario.fromFileBox(file)
		if err != nil {
			err := fmt.Errorf("error reading %q: %w", file, err)
			trace.Trace(err)
			return types, err
		}
		types, err = yamlDocTypes(content)
		if err != nil {
			return types, err
		}
	}
	types = append(types, "scenario")
	return types, err
}

// yamlDocTypes returns a list of the document types encountered in the (possibly multi-document)
// YAML payload `content`.
func yamlDocTypes(content []byte) (types []string, err error) {
	type docSchema struct {
		Type string
	}

	decoder := yaml.NewDecoder(bytes.NewReader(content))
	for {
		doc := &docSchema{Type: fileTypeUnknown}
		err = decoder.Decode(doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			err = fmt.Errorf("error decoding file: %w", err)
			trace.Trace(err)
			return
		}
		types = append(types, doc.Type)
	}
	return types, nil
}

// createSandbox creates a sandbox directory within which all the generation and checks for this
// scenario will be run.
func (scenario *Scenario) createSandbox() (err error) {
	if scenario.sandbox, err = ioutil.TempDir("", fmt.Sprintf("showcase-qualify.%s.%s.", scenario.timestamp, scenario.name)); err != nil {
		err := fmt.Errorf("could not create sandbox: %w", err)
		trace.Trace(err)
		return err
	}

	err = scenario.copyFiles(scenario.filesByType[fileTypeScenario])
	if err != nil {
		err = fmt.Errorf("could not copy scenario files: %w", err)
		trace.Trace(err)
		return err
	}

	if err = scenario.getIncludes(scenario.filesByType[fileTypeInclude]); err != nil {
		err = fmt.Errorf("could not copy included schema files: %w", err)
		trace.Trace(err)
		return err
	}

	trace.Trace("created sandbox")
	return nil
}

// getIncludes processes include files (whose names start with `includeFilePrefix`) by replacing
// them with the file or directory referenced in the file's contents. The only currently supported
// type of reference is to either an entry in `scenario.schemaBox`, or a directory that resolves to
// one or more entries in `scenario.schemaBox`. If an include file does not resolve to any entries
// in `scenario.schemaBox`, this function returns an error.
//
// If there's need in the future, we could potentially grow to support fetching files or directories
// from public repositories on the web.
func (scenario *Scenario) getIncludes(includeFiles []string) error {
	for _, inclusion := range includeFiles {

		// The following guards against includeFilePrefix being present elsewhere in
		// inclusion. Otherwise, we could simply use
		//
		//   `dstPath := strings.Replace(inclusion, includeFilePrefix, "")`
		prefixStartIdx := strings.LastIndex(inclusion, includeFilePrefix)
		if prefixStartIdx < 0 {
			msg := fmt.Sprintf("logic error: did not find prefix %q in %q", includeFilePrefix, inclusion)
			trace.Trace(msg)
			return fmt.Errorf(msg)
		}
		prefixEndIdx := prefixStartIdx + len(includeFilePrefix)
		dstPath := inclusion[:prefixStartIdx] + inclusion[prefixEndIdx:]

		content, err := scenario.fromFileBox(inclusion)
		if err != nil {
			err = fmt.Errorf("could not read %q: %w", inclusion, err)
			trace.Trace(err)
			return err
		}

		srcPath := strings.TrimSpace(string(content))
		srcPath = strings.ReplaceAll(srcPath, string(inclusion), string(os.PathSeparator))
		files, replacePath, err := GetMatchingFiles(scenario.schemaBox, dstPath, srcPath)
		if err != nil {
			err = fmt.Errorf("could not process %q: %w", inclusion, err)
			trace.Trace(err)
			return err
		}

		scenario.copyFilesTo(files, scenario.schemaBox.Find, srcPath, replacePath)
	}
	return nil
}

// copyFiles copies the files listed in `files` from `scenario.fromFileBox` to `scenario.sandbox`.
func (scenario *Scenario) copyFiles(files []string) error {
	return scenario.copyFilesTo(files, scenario.fromFileBox, "", "")
}

// copyFilesTo copies the files listed in `files` from the source in `fromBox` to
// `scenario.sandbox`. In so doing, it replaces a leading `prefix` in each filename with
// `newPrefix`.
func (scenario *Scenario) copyFilesTo(files []string, fromBox func(string) ([]byte, error), prefix, newPrefix string) (err error) {
	const filePermissions = 0555
	replace := len(prefix) > 0

	trace.Trace("prefix:%q  newPrefix:%q  len(files):%d", prefix, newPrefix, len(files))

	for _, srcFile := range files {
		renamedFile := srcFile
		if replace {
			if !strings.HasPrefix(renamedFile, prefix) {
				err := fmt.Errorf("%q does not begin with %q", srcFile, prefix)
				trace.Trace(err)
				return err
			}
			renamedFile = strings.Replace(renamedFile, prefix, newPrefix, 1)
		}

		renamedDir := filepath.Dir(renamedFile)

		dstDir := scenario.sandboxPath(renamedDir)
		if _, err = os.Stat(dstDir); os.IsNotExist(err) {
			if err = os.MkdirAll(dstDir, os.ModePerm); err != nil {
				err = fmt.Errorf("could not make directory %q: %w", dstDir, err)
				trace.Trace(err)
				return err
			}
		}

		var contents []byte
		if contents, err = fromBox(srcFile); err != nil {
			err = fmt.Errorf("could not find file %q: %w", srcFile, err)
			trace.Trace(err)
			return err
		}

		dstFile := scenario.sandboxPath(renamedFile)
		if err = ioutil.WriteFile(dstFile, contents, filePermissions); err != nil {
			err = fmt.Errorf("could not write file  %q: %w", dstFile, err)
			trace.Trace(err)
			return err
		}
	}
	return nil
}
