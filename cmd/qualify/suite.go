package main

import "log"

type Suite struct {
	name     string
	location string
	files    []string
	sandbox  string
	debugLog *log.Logger

	generator    string
	showcasePort int
	viaProtoc    bool
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
	protos, gapic_yamls, sample_yamls := suite.getGenerationFiles(suite.files)
	suite.sandbox = createSandbox()

	suite.debugLog.Printf("Generate (TODO): protos: %v    gapic_yamls: %v    sample_yamls: %v", protos, gapic_yamls, sample_yamls)
	return nil
}

func (suite *Suite) CheckGeneration() {
	// TODO: Fill in
	suite.debugLog.Printf("CheckGeneration (TODO)")
}

func (suite *Suite) RunTests() {
	// TODO: Fill in
	suite.debugLog.Printf("RunTests (TODO)")
}

func (suite *Suite) Success() bool {
	// TODO: fill in
	return true
}

func (suite *Suite) getGenerationFiles(paths []string) (protos, gapic_yamls, sample_yamls []string) {
	// TODO: extract all protos
	// TODO: extract all gapic configs
	// TODO: extract all sample_configs

	suite.debugLog.Printf("getGenerationFiles (TODO)")
	return nil, nil, nil
}

func createSandbox() string {
	// TODO: make a temporary directory that can be deleted
	return "/tmp"
}
