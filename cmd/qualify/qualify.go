package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var debugLog *log.Logger

func main() {
	const (
		RetCodeSuccess = iota
		RetCodeFailedDependencies
		RetCodeUsageError
		RetCodeInternalError
		RetSuiteFailure
	)

	debugMe := true // TODO: get from CLI args
	debugStream := io.Writer(os.Stderr)
	if !debugMe {
		debugStream = ioutil.Discard
	}
	debugLog = log.New(debugStream, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	if err := checkDependencies(); err != nil {
		os.Exit(RetCodeFailedDependencies)
	}

	showcasePID := startShowcase()
	defer endProcess(showcasePID)

	generatorName, viaProtoc, err := getGeneratorData()
	if err != nil {
		os.Exit(RetCodeUsageError)
	}

	allSuites := GetTestSuites(generatorName, viaProtoc)
	success := true
	for _, suite := range allSuites {
		if err := suite.Run(); err != nil {
			os.Exit(RetCodeInternalError)
		}
		if !suite.Success() {
			success = false
		}
	}
	if !success {
		os.Exit(RetSuiteFailure)
	}
}

func checkDependencies() error {
	// TODO: add check for sample-tester
	debugLog.Printf("checkDependencies (TODO)")
	return nil
}

// getGeneratorData obtains the name of the generaor from the command
// line, and whether it is a protoc plugin (if not, it is considered
// part of the monolith)
func getGeneratorData() (string, bool, error) {
	// TODO: Get from the command line
	return "python", true, nil
}

func GetTestSuites(generatorName string, viaProtoc bool) []*Suite {
	allSuites := []*Suite{}
	defaultShowcasePort := 123 // TODO fix
	// TODO: iterate over test suites
	// TODO: get files in each suite
	suiteFiles := []string{""}
	newSuite := &Suite{
		showcasePort: defaultShowcasePort,
		viaProtoc:    viaProtoc,
		generator:    generatorName,
		files:        suiteFiles,
		debugLog:     debugLog,
	}
	debugLog.Printf("adding suite %#v", newSuite)
	allSuites = append(allSuites, newSuite)
	return allSuites
}

// startShowcase starts the Showcase server and returns its PID
func startShowcase() int {
	// TODO: fill in
	debugLog.Printf("startShowcase (TODO)")
	return 0
}

// endProcess ends the process with the specified PID
func endProcess(pid int) {
	// TODO: fill in
	debugLog.Printf("endProcess (TODO): %v", pid)
}
