package qualifier

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	trace "github.com/google/go-trace"
)

// Generator has the information needed to run a generator in a given language. The generator
// will be run as a protoc plugin unless `isMonolith` is set, in which case gapic-generator will be
// invoked directly.
type Generator struct {
	Language        string
	PluginDirectory string
	PluginOptions   string
	isMonolith      bool
}

// Run will run the generator from `work_dir`, obtaining the required files from filesByType.
func (gen *Generator) Run(workDir string, filesByType map[string][]string) ([]byte, *os.ProcessState, error) {
	generationDir := "generated"
	if gen.isMonolith {
		return nil, nil, fmt.Errorf("monolith not implemented yet")
	}

	if err := os.MkdirAll(path.Join(workDir, generationDir), os.ModePerm); err != nil {
		return nil, nil, err
	}

	optionFlag := fmt.Sprintf("--%s_gapic_opt", gen.Language)
	allOptions := []string{}
	if len(gen.PluginOptions) > 0 {
		allOptions = append(allOptions, optionFlag, gen.PluginOptions)
	}

	sampleConfigFiles, _ := filesByType[fileTypeSampleConfig]
	if len(sampleConfigFiles) > 0 {
		allOptions = append(allOptions, optionFlag,
			fmt.Sprintf("samples=%s", strings.Join(sampleConfigFiles, ",samples=")))
	}

	cmdParts := []string{
		"protoc",
		fmt.Sprintf("--%s_gapic_out", gen.Language), fmt.Sprintf("./%s", generationDir),
	}
	cmdParts = append(cmdParts, allOptions...)
	if len(gen.PluginDirectory) > 0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--plugin=%s/protoc-gen-%s_gapic", gen.PluginDirectory, gen.Language))

	}
	cmdParts = append(cmdParts, filesByType[fileTypeProtobuf]...)

	cmdString := strings.Join(cmdParts, " ")
	trace.Trace("running: %s", cmdString)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		err = nil
		trace.Trace("clearing exit error: %v", exitError)
	}
	return output, cmd.ProcessState, err
}
