package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	packr "github.com/gobuffalo/packr/v2"
	trace "github.com/google/go-trace"
	yaml "gopkg.in/yaml.v2"
)

var timestamp string

const includeFilePrefix = "include."
const includeFileDirectorySeparator = '/'

func init() {
	timestamp = time.Now().Format("20060102.150405")
	trace.Trace("timestamp = %q", timestamp)
}

type Scenario struct {
	// Every f in `files` corresponds to a file in `fileBox` under
	// a directory called `name`, which is NOT part of the string
	// f itself.
	name        string
	files       []string
	fileBox     *packr.Box
	schemaBox   *packr.Box
	sandbox     string
	filesByType map[string][]string

	generator    *GeneratorInfo
	showcasePort int
}

func (scenario *Scenario) sandboxPath(relativePath string) string {
	if len(scenario.name) == 0 {
		trace.Trace("scenario name not set yet: relativePath==%q!", relativePath)
		panic("scenario name not set yet when trying to get sandbox filename")
	}
	if len(scenario.sandbox) == 0 {
		trace.Trace("sandbox name not set yet: relativePath==%q!", relativePath)
		panic("sandbox name not set yet when trying to get sandbox filename")
	}

	return filepath.Join(scenario.sandbox, relativePath)
}

func (scenario *Scenario) fileBoxPath(relativePath string) string {
	if len(scenario.name) == 0 {
		trace.Trace("scenario name not set yet: relativePath==%q!", relativePath)
		panic("scenario name not set yet when trying to get box filename")
	}
	return filepath.Join(scenario.name, relativePath)
}

func (scenario *Scenario) fromFileBox(relativePath string) ([]byte, error) {
	return scenario.fileBox.Find(scenario.fileBoxPath(relativePath))
}

func (scenario *Scenario) Run() error {
	if err := scenario.Generate(); err != nil {
		return err
	}
	scenario.CheckGeneration()
	scenario.RunTests()
	return nil
}

func (scenario *Scenario) Generate() (err error) {
	if err = scenario.getGenerationFiles(); err != nil {
		return err
	}
	err = scenario.createSandbox()
	if err != nil {
		return fmt.Errorf("could not create sandbox: %w", err)
	}
	trace.Trace("created sandbox: %q", scenario.sandbox)

	var output []byte
	output, err = scenario.generator.Run(scenario.sandbox, scenario.filesByType)
	trace.Trace("run error: %v", err)
	trace.Trace("run output: %s", output)

	/* TODO:
	   - P1. Invoke sample generator with all the protos in the directory
	   - P2. when creating sandbox, expand files of the form `include.FOO` to be file/dir FOO/ with contents cloned from the location specified in the file
	        can be a path rooted at schema/
	        can be an http: which is curled
	*/
	trace.Trace("Generate (TODO)")
	return nil
}

func (scenario *Scenario) CheckGeneration() {
	// TODO: Fill in
	trace.Trace("CheckGeneration (TODO)")
}

func (scenario *Scenario) RunTests() {
	// TODO: Fill in
	trace.Trace("RunTests (TODO)")
}

func (scenario *Scenario) Success() bool {
	// TODO: fill in
	return true
}

// getGenerationFiles classifies all the files in the scenario,
// storing the results in `scenario.filesByType`. Each file is
// classified as either an "include" file or a "scenario"
// file. Scenario files are also classified by the type of data they
// contain. Thus, "scenario" files have at least two labels (ie
// multiple entries in `scenario.filesByType`).
func (scenario *Scenario) getGenerationFiles() (err error) {
	scenario.filesByType = make(map[string][]string)
	// TODO: process `include.*` files first
	for _, thisFile := range scenario.files {
		fileTypes, err := scenario.getFileTypes(thisFile)
		if err != nil {
			return err
		}
		for _, oneType := range fileTypes {
			similarFiles := scenario.filesByType[oneType]
			scenario.filesByType[oneType] = append(similarFiles, thisFile)
		}
		trace.Trace("%s: type %q", thisFile, fileTypes)
	}
	return nil
}

// getFileTypes gets the various type labels for `thisFile`. The type
// labels are either the singleton list {"include"}, or a list
// {"scenario", ...} that includes the various types of data included
// in that file.
func (scenario *Scenario) getFileTypes(thisFile string) (fileTypes []string, err error) {

	for idx := strings.Index(thisFile, includeFilePrefix); idx >= 0; idx = strings.Index(thisFile[idx+1:], includeFilePrefix) {
		if !(idx == 0 || thisFile[idx-1] == os.PathSeparator) {
			// The directory entry does not start with
			// this pattern, so it's not an include file.
			continue
		}
		if strings.Index(thisFile[idx+1:], string(os.PathListSeparator)) != -1 {
			// There is a sub-directory component, so this is not a file.
			continue
		}

		// This looks like a legitimate include file. Cease classifying.
		return []string{"include"}, nil
	}

	extension := filepath.Ext(thisFile)
	switch extension {
	case ".proto":
		fileTypes = []string{"proto"}
	case ".yaml", ".yml":
		trace.Trace("reading %s", thisFile)
		content, err := scenario.fromFileBox(thisFile)
		if err != nil {
			err := fmt.Errorf("error reading %q: %w", thisFile, err)
			trace.Trace(err)
			return fileTypes, err
		}
		fileTypes, err = yamlDocTypes(content)
		if err != nil {
			return fileTypes, err
		}
	}
	fileTypes = append(fileTypes, "scenario")
	return fileTypes, err
}

func yamlDocTypes(fileContent []byte) (docTypes []string, err error) {
	type docSchema struct {
		Type string
	}

	decoder := yaml.NewDecoder(bytes.NewReader(fileContent))
	for {
		doc := &docSchema{Type: "(UNKNOWN)"}
		err = decoder.Decode(doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			err = fmt.Errorf("error decoding file: %w", err)
			trace.Trace(err)
			return
		}
		docTypes = append(docTypes, doc.Type)
	}
	return docTypes, nil
}

func (scenario *Scenario) createSandbox() (err error) {

	if scenario.sandbox, err = ioutil.TempDir("", fmt.Sprintf("showcase-qualify.%s.%s.", timestamp, scenario.name)); err != nil {
		err := fmt.Errorf("could not create sandbox: %w", err)
		trace.Trace(err)
		return err
	}

	err = scenario.copyFiles(scenario.filesByType["scenario"])
	if err != nil {
		err = fmt.Errorf("could not copy scenario files: %w", err)
		trace.Trace(err)
		return err
	}

	if err = scenario.getIncludes(scenario.filesByType["include"]); err != nil {
		err = fmt.Errorf("could not copy included schema files: %w", err)
		trace.Trace(err)
		return err
	}

	trace.Trace("created sandbox")
	return nil
}

func (scenario *Scenario) getIncludes(includes []string) error {
	for _, includeFile := range includes {

		// The following guards against includeFilePrefix being present elsewhere in includeFile. Otherwise, we could use `dstPath := strings.Replace(includeFile, includeFilePrefix, "")`
		prefixStartIdx := strings.LastIndex(includeFile, includeFilePrefix)
		prefixEndIdx := prefixStartIdx + len(includeFilePrefix)
		if prefixStartIdx < 0 || prefixEndIdx >= len(includeFile) {
			msg := fmt.Sprintf("logic error: start %d, end %d outside of range [0, %d] for %q", prefixStartIdx, prefixEndIdx, len(includeFile), includeFile)
			trace.Trace(msg)
			panic(msg)
		}
		dstPath := includeFile[0:prefixStartIdx] + includeFile[prefixEndIdx:]

		content, err := scenario.fromFileBox(includeFile)
		if err != nil {
			err = fmt.Errorf("could not read %q: %w", includeFile, err)
			trace.Trace(err)
			return err
		}

		srcPath := strings.TrimSpace(string(content))
		srcPath = strings.ReplaceAll(srcPath, string(includeFileDirectorySeparator), string(os.PathSeparator))
		files, replacePath, err := GetMatchingFiles(scenario.schemaBox, dstPath, srcPath)
		if err != nil {
			err = fmt.Errorf("could not process %q: %w", includeFile, err)
			trace.Trace(err)
			return err
		}

		scenario.copyFilesReplacePath(files, scenario.schemaBox.Find, srcPath, replacePath)
	}
	return nil
}

func (scenario *Scenario) copyFiles(files []string) error {
	return scenario.copyFilesReplacePath(files, scenario.fromFileBox, "", "")
}

func (scenario *Scenario) copyFilesReplacePath(files []string, fromBox func(string) ([]byte, error), replaceSrc, replaceDst string) (err error) {
	const filePermissions = 0555
	replace := len(replaceSrc) > 0

	trace.Trace("replaceSrc:%q  replaceDst:%q  len(files):%d", replaceSrc, replaceDst, len(files))

	for _, srcFile := range files {
		renamedFile := srcFile
		if replace {
			if !strings.HasPrefix(renamedFile, replaceSrc) {
				panic(fmt.Sprintf("%q does not begin with %q", srcFile, replaceSrc))
			}
			renamedFile = strings.Replace(renamedFile, replaceSrc, replaceDst, 1)
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
