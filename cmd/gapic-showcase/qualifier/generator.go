// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Language   string // the language in which to generate GAPICs and samples
	Directory  string // the directory in which to find the generator
	Options    string // any options
	isMonolith bool   // TODO: implement this functionality and export this field
}

// Run runs the generator `gen` from `work_dir`, creating it if necessary and obtaining the required
// files from `filesByType`. The generated output is placed in a `generated/` sub-directory.
func (gen *Generator) Run(workDir string, filesByType map[string][]string) ([]byte, *os.ProcessState, error) {
	const generationDir = "generated"

	if gen.isMonolith {
		return nil, nil, fmt.Errorf("monolith not implemented yet")
	}

	if err := os.MkdirAll(path.Join(workDir, generationDir), os.ModePerm); err != nil {
		return nil, nil, err
	}

	// Construct the various arguments to invoke the generator as a protoc plugin.

	pluginOpt := fmt.Sprintf("--%s_gapic_opt", gen.Language)
	opts := []string{}
	if len(gen.Options) > 0 {
		opts = append(opts, pluginOpt, gen.Options)
	}

	sampleConfigs := filesByType[fileTypeSampleConfig]
	if len(sampleConfigs) > 0 {
		opts = append(opts, pluginOpt,
			fmt.Sprintf("samples=%s", strings.Join(sampleConfigs, ",samples=")))
	}

	cmdParts := []string{
		"protoc",
		fmt.Sprintf("--%s_gapic_out", gen.Language), fmt.Sprintf("./%s", generationDir),
	}
	cmdParts = append(cmdParts, opts...)
	if len(gen.Directory) > 0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--plugin=%s/protoc-gen-%s_gapic", gen.Directory, gen.Language))

	}
	cmdParts = append(cmdParts, filesByType[fileTypeProtobuf]...)

	// Execute the command, clear all but internal errors, return.

	trace.Trace("running: %s", strings.Join(cmdParts, " "))

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
