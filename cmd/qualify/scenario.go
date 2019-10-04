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

func init() {
	timestamp = time.Now().Format("20060102.150405")
	trace.Trace("timestamp = %q", timestamp)
}

type Scenario struct {
	// Every f in `files` must begin with `name`:
	// strings.HasPrefix(f, name)
	name        string
	files       []string
	fileBox     *packr.Box
	sandbox     string
	filesByType map[string][]string

	generator    *GeneratorInfo
	showcasePort int
}

func (scenario *Scenario) sandboxFileName(srcPath string) string {
	if len(scenario.name) == 0 {
		trace.Trace("scenario name not set yet when trying to get sandbox name for %q!", srcPath)
		panic("scenario name not set yet when trying to get sandbox name")
	}
	if len(scenario.sandbox) == 0 {
		trace.Trace("sandbox not set yet when trying to get sandbox name for %q!", srcPath)
		panic("sandbox not set yet when trying to get sandbox name")
	}
	if !strings.HasPrefix(srcPath, scenario.name) {
		trace.Trace("file %q is not part of sandbox %q", srcPath, scenario.sandbox)
		panic("file not part of sandbox")
	}
	trace.Trace("srcPath: %s", srcPath)

	relativePath := srcPath[len(scenario.name):]
	if len(relativePath) > 0 {
		// There are other path elements, so remove the leading PathSeparator
		relativePath = relativePath[len(string(os.PathSeparator)):]
	}
	dstPath := filepath.Join(scenario.sandbox, relativePath)

	return dstPath
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

func (scenario *Scenario) getGenerationFiles() (err error) {
	scenario.filesByType = make(map[string][]string)
	// TODO: process `include.*` files first
	for _, thisFile := range scenario.files {
		extension := filepath.Ext(thisFile)
		var fileTypes []string
		switch extension {

		case ".proto":
			fileTypes = []string{"proto"}
		case ".yaml", ".yml":
			trace.Trace("reading %s", thisFile)
			content, err := scenario.fileBox.Find(thisFile)
			if err != nil {
				err := fmt.Errorf("error reading %q: %w", thisFile, err)
				trace.Trace(err)
				return err
			}
			fileTypes, err = yamlDocTypes(content)
			if err != nil {
				return err
			}
		}

		for _, oneType := range fileTypes {
			similarFiles := scenario.filesByType[oneType]
			scenario.filesByType[oneType] = append(similarFiles, thisFile)
			trace.Trace("%s: type %q", thisFile, oneType)
		}
	}
	return nil
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
		return err
	}

	// Copy files needed to generate
	for _, srcFile := range scenario.files {

		// TODO: expand `include.*` files (passed in as params)

		srcDir := filepath.Dir(srcFile)
		dstDir := scenario.sandboxFileName(srcDir)
		trace.Trace("  srcFile: %s", srcFile)
		trace.Trace("   srcDir: %s", srcDir)
		trace.Trace("     name: %s", scenario.name)
		trace.Trace("   dstDir: %s", dstDir)
		if _, err := os.Stat(dstDir); os.IsNotExist(err) {
			if err = os.MkdirAll(dstDir, os.ModePerm); err != nil {
				trace.Trace("could not MkdirAll(%q)", dstDir)
				return err
			}
		}

		var contents []byte
		if contents, err = scenario.fileBox.Find(srcFile); err != nil {
			trace.Trace("could not find file %q", srcFile)
			return err
		}

		dstFile := scenario.sandboxFileName(srcFile)
		if err := ioutil.WriteFile(dstFile, contents, 0555); err != nil {
			trace.Trace("could not WriteFile()")
			return err
		}
	}
	return nil
}
