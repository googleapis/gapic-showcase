package main

import (
	"os"
	"strings"

	packr "github.com/gobuffalo/packr/v2"
	trace "github.com/google/go-trace"
)

var AcceptanceSuite *packr.Box
var SchemaSuite *packr.Box

func GetAssets() {
	AcceptanceSuite = packr.New("acceptance suite", "../../acceptance")
	trace.Trace("packr suite: %v", AcceptanceSuite)
	contents := AcceptanceSuite.List()
	trace.Trace("in the suite: \n  %v", strings.Join(contents, "\n  "))

	SchemaSuite = packr.New("schema suite", "../../schema")
	trace.Trace("packr suite: %v", SchemaSuite)
	contents = SchemaSuite.List()
	trace.Trace("in the suite: \n  %v", strings.Join(contents, "\n  "))
}

type FilesByDir struct {
	Directory string
	Files     []string // filenames are relative paths from `Directory`
}

func GetFilesByDir(box *packr.Box) []*FilesByDir {
	filesInAllDirs := []*FilesByDir{}
	var filesInThisDir *FilesByDir
	commitThisDir := func() {
		if filesInThisDir != nil {
			filesInAllDirs = append(filesInAllDirs, filesInThisDir)
			filesInThisDir = nil
		}
	}

	// We assume allFiles is returned in a top-down order, so that
	// all files under each first level directory appear
	// contiguously.
	allFiles := box.List()
	previousDir := ""
	for _, oneFile := range allFiles {
		dir := strings.Split(oneFile, string(os.PathSeparator))[0]
		if dir != previousDir {
			commitThisDir()
			filesInThisDir = &FilesByDir{Directory: dir}
		}
		previousDir = dir
		filesInThisDir.Files = append(filesInThisDir.Files, oneFile[len(dir)+len(string(os.PathSeparator)):])
	}
	commitThisDir()

	return filesInAllDirs
}
