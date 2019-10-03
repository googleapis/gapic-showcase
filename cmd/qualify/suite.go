package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"

	trace "github.com/google/go-trace"
	yaml "gopkg.in/yaml.v2"
)

type Suite struct {
	name        string
	location    string
	files       []string
	sandbox     string
	filesByType map[string][]string

	generator    *GeneratorInfo
	showcasePort int
}

func (suite *Suite) Run() error {
	if err := suite.Generate(); err != nil {
		return err
	}
	suite.CheckGeneration()
	suite.RunTests()
	return nil
}

func (suite *Suite) Generate() error {
	if err := suite.getGenerationFiles(); err != nil {
		return err
	}
	suite.sandbox = createSandbox()

	trace.Trace("Generate (TODO)")
	return nil
}

func (suite *Suite) CheckGeneration() {
	// TODO: Fill in
	trace.Trace("CheckGeneration (TODO)")
}

func (suite *Suite) RunTests() {
	// TODO: Fill in
	trace.Trace("RunTests (TODO)")
}

func (suite *Suite) Success() bool {
	// TODO: fill in
	return true
}

func (suite *Suite) getGenerationFiles() (err error) {
	suite.filesByType = make(map[string][]string)
	for _, thisFile := range suite.files {
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
			similarFiles := suite.filesByType[oneType]
			suite.filesByType[oneType] = append(similarFiles, thisFile)
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
