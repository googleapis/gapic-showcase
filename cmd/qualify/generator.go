package main

import (
	"fmt"
	"os/exec"
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
	if gi.isMonolith {
		return nil, fmt.Errorf("monolith not implemented yet")
	}

	//	allProtos := strings.Join(filesByType["proto"], " ")
	cmdParts := []string{
		"protoc",
		fmt.Sprintf("--%s_gapic_out", gi.name), "./generated",
		fmt.Sprintf("--%s_gapic_opt", gi.name), gi.options,
		fmt.Sprintf("--plugin=%s/protoc-gen-%s_gapic", gi.dir, gi.name),
	}
	cmdParts = append(cmdParts, filesByType["proto"]...)
	cmdString := strings.Join(cmdParts, " ")

	// MAYBE we need to store the list of files as relative dirs, and keep the source dir separate

	// cmdString := fmt.Sprintf("protoc --%s_gapic_out ./gen --%s_gapic_opt %s --plugin=%s/protoc-gen-%s_gapic %s",
	// 	gi.name, gi.name, gi.options,
	// 	gi.dir, gi.name,
	// 	allProtos)
	trace.Trace("running: %s", cmdString)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = wdir
	return cmd.CombinedOutput()
}
