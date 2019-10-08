package main

import (
	"fmt"
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

	SchemaSuite = packr.New("schema", "../../schema")
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

// GetMatchingFiles returns the a list of files in `box` matching
// `pattern`. `pattern` should end with `pathSeparator` iff it refers
// to a directory.
//
// This also returns a couple of additional values. These are useful
// for copying files within matched directories, and are provided as a
// convenience:
//
//  * `srcPath`: a version of `pattern` modified so that
//    `pathSeparator`, if any, becomes `os.PathSeparator`
//  * `replacePath`: a version of `dstPath` (which is assumed to have
//     no trailing separator) with `os.PathSeparator` appended.
func GetMatchingFiles(box *packr.Box, dstPath string, pattern string, pathSeparator rune) (files []string, srcPath string, replacePath string, err error) {

	trace.Trace("reading %q with separator %c", pattern, pathSeparator)

	replacePath = dstPath
	srcPath = pattern

	// If `pattern` specifies a single file, match just that.
	if !strings.HasSuffix(pattern, string(pathSeparator)) {
		if !box.Has(pattern) {
			err = fmt.Errorf("file box %q has no file %q", box.Name, pattern)
			return
		}
		files = []string{pattern}
		return
	}

	// Replace the trailing pathSeparator (hard-coded in text)
	// with the os.PathSeparator presumably used by packr, which
	// varies by OS
	srcPath = string(append([]rune(srcPath)[:len(srcPath)-1], os.PathSeparator))

	// Likewise, let the caller replace the directory part of the
	// path with the separator included.
	replacePath = dstPath + string(os.PathSeparator)

	for _, entry := range box.List() {
		if strings.HasPrefix(entry, srcPath) {
			files = append(files, entry)
		}
	}
	if len(files) == 0 {
		err = fmt.Errorf("file box %q has no files matching %q", box.Name, pattern)
		return
	}
	trace.Trace("obtained files %q", files)
	return
}
