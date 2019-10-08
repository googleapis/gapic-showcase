package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	trace "github.com/google/go-trace"
)

type GeneratorInfo struct {
	name       string
	dir        string
	options    string
	isMonolith bool
}

func (gi *GeneratorInfo) Run(wdir string, filesByType map[string][]string) ([]byte, error) {
	generationDir := "generated"
	if gi.isMonolith {
		return nil, fmt.Errorf("monolith not implemented yet")
	}

	if err := os.MkdirAll(path.Join(wdir, generationDir), os.ModePerm); err != nil {
		return nil, err
	}

	cmdParts := []string{
		"protoc",
		fmt.Sprintf("--%s_gapic_out", gi.name), fmt.Sprintf("./%s", generationDir),
		fmt.Sprintf("--%s_gapic_opt", gi.name), gi.options,
		fmt.Sprintf("--plugin=%s/protoc-gen-%s_gapic", gi.dir, gi.name),
	}
	cmdParts = append(cmdParts, filesByType["proto"]...)
	cmdString := strings.Join(cmdParts, " ")

	trace.Trace("running: %s", cmdString)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = wdir
	return cmd.CombinedOutput()
}
