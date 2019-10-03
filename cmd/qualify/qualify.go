package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	trace "github.com/google/go-trace"
)

/* Development instructions

. ./prepare-to-qualify

=== Sample run ======
go run *.go

Notes for upcoming development:
packing assets into Go files:
https://www.google.com/search?q=golang+ship+data+files+with+binary&oq=golang+ship+data+files+with+binary
https://github.com/go-bindata/go-bindata
https://github.com/gobuffalo/packr
Choosing: https://tech.townsourced.com/post/embedding-static-files-in-go/

*/

const showcaseCmd = "gapic-showcase" // TODO: Consider running in-process

func main() {
	const (
		RetCodeSuccess = iota
		RetCodeInternalError
		RetCodeFailedDependencies
		RetCodeUsageError
		RetCodeCantRunShowcase
		RetSuiteFailure
	)

	debugMe := true // TODO: get from CLI args
	if debugMe {
		trace.On(true)
	}

	if err := checkDependencies(); err != nil {
		os.Exit(RetCodeFailedDependencies)
	}

	showcase, err := startShowcase()
	if err != nil {
		os.Exit(RetCodeCantRunShowcase)
	}
	defer endProcess(showcase)

	generatorData, err := getGeneratorData()
	if err != nil {
		os.Exit(RetCodeUsageError)
	}

	allSuites, err := GetTestSuites(generatorData)
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
	sampleTesterCmd := "sample-tester"
	notFound := []string{}
	trace.Trace("")

	showcasePath, err := exec.LookPath(showcaseCmd)
	if err != nil {
		notFound = append(notFound, showcaseCmd)
	}

	sampleTesterPath, err := exec.LookPath(sampleTesterCmd)
	if err != nil {
		notFound = append(notFound, sampleTesterCmd)
	}

	if len(notFound) > 0 {
		msg := fmt.Sprintf("could not find dependencies in $PATH: %q", notFound)
		log.Printf(msg)
		return fmt.Errorf(msg)
	}
	trace.Trace("found %q: %s", showcaseCmd, showcasePath)
	trace.Trace("found %q: %s", sampleTesterCmd, sampleTesterPath)
	return nil
}

// getGeneratorData obtains the name of the generator from the command
// line, and whether it is a protoc plugin (if not, it is considered
// part of the monolith)
func getGeneratorData() (*GeneratorInfo, error) {
	pluginName := "go"                                                                 // TODO: get from CLI args
	pluginDir := "/tmp/goinstall/bin"                                                  // TODO: get from CLI args
	pluginOpt := "'go-gapic-package=cloud.google.com/go/showcase/apiv1beta1;showcase'" // TODO: get from CLI args
	generator := &GeneratorInfo{
		name:    pluginName,
		dir:     pluginDir,
		options: pluginOpt,
	}
	trace.Trace("hard-coded generator=%#v", generator)
	return generator, nil
}

// GetTestSuites returns a list of Suite as found in the specified
// suite root directory.
func GetTestSuites(generator *GeneratorInfo) ([]*Suite, error) {
	suiteRootPath := "../../acceptance"

	trace.Trace("(generator=%#v)", generator)
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
			generator:    generator,
			files:        suiteFiles,
		}
		trace.Trace("adding suite %#v", newSuite)
		allSuites = append(allSuites, newSuite)
	}
	return allSuites, nil
}

// startShowcase starts the Showcase server and returns its PID
func startShowcase() (*exec.Cmd, error) {
	// TODO: fill in
	trace.Trace("")
	cmd := exec.Command(showcaseCmd)
	err := cmd.Start()
	return cmd, err
}

// endProcess ends the process with the specified PID
func endProcess(cmd *exec.Cmd) error {
	// TODO: fill in
	trace.Trace("cmd=%v)", cmd)
	process := cmd.Process
	if err := process.Signal(os.Interrupt); err != nil {
		msg := fmt.Sprintf("could not end process %v normally", process.Pid)
		if err := process.Kill(); err != nil {
			msg = fmt.Sprintf("%s; killing failed. Please kill manually!", msg)
			log.Printf(msg)
			return fmt.Errorf(msg)
		}
		msg = fmt.Sprintf("%s but killing it succeeded", msg)
		log.Printf(msg)
	}
	trace.Trace("waiting for process to end: %v (%v)", process.Pid, cmd.Path)
	cmd.Wait()
	trace.Trace("process ended:              %v (%v) ended", process.Pid, cmd.Path)
	return nil
}

// getAllFiles returns a list of all the files (excluding directories)
// at any level under `root`.
func getAllFiles(root string) ([]string, error) {
	trace.Trace("root=%q", root)
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
