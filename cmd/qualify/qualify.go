package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
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

	allSuites, err := GetTestSuites(generatorName, viaProtoc)
	if err != nil {
		os.Exit(RetCodeInternalError)
	}

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

// GetTestSuites returns a list of Suite as found in the specified
// suite root directory.
func GetTestSuites(generatorName string, viaProtoc bool) ([]*Suite, error) {
	suiteRootPath := "../../acceptance"

	debugLog.Printf("GetTestSuites: generator=%q   protoc-%v", generatorName, viaProtoc)
	allSuites := []*Suite{}
	suiteEntries, err := ioutil.ReadDir(suiteRootPath)
	if err != nil {
		return nil, err
	}

	for _, suiteDir := range suiteEntries {
		if !suiteDir.IsDir() {
			continue
		}
		name := suiteDir.Name()
		location, err := filepath.Abs(path.Join(suiteRootPath, name))
		if err != nil {
			return nil, err
		}

		suiteFiles, err := getAllFiles(location)
		if err != nil {
			return nil, err
		}

		defaultShowcasePort := 123 // TODO fix
		newSuite := &Suite{
			name:         name,
			location:     location,
			showcasePort: defaultShowcasePort,
			viaProtoc:    viaProtoc,
			generator:    generatorName,
			files:        suiteFiles,
			debugLog:     debugLog,
		}
		debugLog.Printf("adding suite %#v", newSuite)
		allSuites = append(allSuites, newSuite)
	}
	return allSuites, nil
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

// getAllFiles returns a list of all the files (excluding directories)
// at any level under `root`.
func getAllFiles(root string) ([]string, error) {
	debugLog.Printf("getAllFiles: root=%q", root)
	allFiles := []string{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			allFiles = append(allFiles, path)
			return nil
		})
	return allFiles, err
}
