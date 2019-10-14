package qualifier

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	trace "github.com/google/go-trace"
)

/* Development instructions

. ./prepare-to-qualify

=== Sample run ======
go run *.go


*/

type Settings struct {
	Generator
	ShowcasePort int
	Timestamp    string
	Verbose      bool
}

const showcaseCmd = "gapic-showcase"

func Run(settings *Settings) {
	// TODO: Return an error rather than exiting. We can return an error type ErrorCode that
	// wraps the current errors and includes the appropriate return code, and change main() so
	// that if the error it gets is of type ErrorCode, it exits with that code.
	const (
		RetCodeSuccess = iota
		RetCodeInternalError
		RetCodeFailedDependencies
		RetCodeUsageError
		RetScenarioFailure
	)

	if err := GetAssets(); err != nil {
		os.Exit(RetCodeInternalError)
	}

	if err := checkDependencies(); err != nil {
		os.Exit(RetCodeFailedDependencies)
	}

	allScenarios, err := GetTestScenarios(settings)
	if err != nil {
		os.Exit(RetCodeInternalError)
	}

	success := true
	for _, scenario := range allScenarios {
		if err := scenario.Run(); err != nil {
			os.Exit(RetCodeInternalError)
		}
		status := "SUCCESS"
		if !scenario.Success() {
			success = false
			status = "FAILURE"
		}
		fmt.Printf("scenario %q: %s    %s\n", scenario.name, status, scenario.sandbox)
	}
	if !success {
		os.Exit(RetScenarioFailure)
	}
}

func checkDependencies() error {
	sampleTesterCmd := "sample-tester"
	notFound := []string{}
	trace.Trace("")

	sampleTesterPath, err := exec.LookPath(sampleTesterCmd)
	if err != nil {
		notFound = append(notFound, sampleTesterCmd)
	}

	if len(notFound) > 0 {
		msg := fmt.Sprintf("could not find dependencies in $PATH: %q", notFound)
		log.Printf(msg)
		return fmt.Errorf(msg)
	}
	trace.Trace("found %q: %s", sampleTesterCmd, sampleTesterPath)
	return nil
}

// GetTestScenarios returns a list of Scenario as found in the specified
// scenario root directory.
func GetTestScenarios(settings *Settings) ([]*Scenario, error) {
	allScenarios := []*Scenario{}
	allScenarioConfigs := GetFilesByDir(AcceptanceSuite)
	for _, config := range allScenarioConfigs {
		newScenario := &Scenario{
			name:         config.Directory,
			timestamp:    settings.Timestamp,
			showcasePort: settings.ShowcasePort,
			generator:    &settings.Generator,
			files:        config.Files,
			fileBox:      AcceptanceSuite,
			schemaBox:    SchemaSuite,
		}
		trace.Trace("adding scenario %#v", newScenario)
		allScenarios = append(allScenarios, newScenario)
	}
	return allScenarios, nil
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
