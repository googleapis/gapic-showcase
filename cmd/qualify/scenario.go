package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"

	packr "github.com/gobuffalo/packr/v2"
	trace "github.com/google/go-trace"
	yaml "gopkg.in/yaml.v2"
)

type Scenario struct {
	name        string
	files       []string
	fileBox     *packr.Box
	sandbox     string
	filesByType map[string][]string

	generator    *GeneratorInfo
	showcasePort int
}

func (scenario *Scenario) Run() error {
	if err := scenario.Generate(); err != nil {
		return err
	}
	scenario.CheckGeneration()
	scenario.RunTests()
	return nil
}

func (scenario *Scenario) Generate() error {
	if err := scenario.getGenerationFiles(); err != nil {
		return err
	}
	scenario.sandbox = createSandbox()

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
	for _, thisFile := range scenario.files {
		extension := filepath.Ext(thisFile)
		var fileTypes []string
		switch extension {

		case ".proto":
			fileTypes = []string{"proto"}
		case ".yaml", ".yml":
			trace.Trace("reading %s", thisFile)
			content, err := ioutil.ReadFile(thisFile)
			if err != nil {
				trace.Trace("error reading: %s", err)
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
			trace.Trace("error decoding: %s", err)
			return
		}
		docTypes = append(docTypes, doc.Type)
	}
	return docTypes, nil
}

func createSandbox() string {
	// TODO: make a temporary directory that can be deleted
	return "/tmp"
}
