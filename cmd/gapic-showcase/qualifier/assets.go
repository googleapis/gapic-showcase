package qualifier

import (
	"fmt"
	"os"
	"strings"

	packr "github.com/gobuffalo/packr/v2"
	trace "github.com/google/go-trace"
)

var AcceptanceSuite *packr.Box
var SchemaSuite *packr.Box

const fileTypeSampleConfig = "com.google.api.codegen.SampleConfigProto"
const fileTypeProtobuf = "proto"

func GetAssets() {
	// I believe we can't pass the arguments into a function
	// because otherwise packr won't be able to recognize these
	// paths should be packed.
	AcceptanceSuite = packr.New("acceptance suite", "acceptance_suite")
	SchemaSuite = packr.New("schema", "../../../schema")

	traceBox(AcceptanceSuite)
	traceBox(SchemaSuite)
}

func traceBox(box *packr.Box) {
	trace.Trace("suite %q has %d entries", box.Name, len(box.List()))
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
// `srcPath`. `srcPath` should end with `os.PathSeparator` iff it refers
// to a directory.
//
// As a convenience, this also returns an additional value, useful
// for copying files within matched directories:
//
//  * `replacePath`: a copy of `dstPath` (which is assumed to have no
//     trailing separator), with `os.PathSeparator` appended if
//     `srcPath` refers to a directory
func GetMatchingFiles(box *packr.Box, dstPath string, srcPath string) (files []string, replacePath string, err error) {

	trace.Trace("reading %q", srcPath)

	replacePath = dstPath

	// If `srcPath` specifies a single file, match just that.
	if !strings.HasSuffix(srcPath, string(os.PathSeparator)) {
		if !box.Has(srcPath) {
			err = fmt.Errorf("file box %q has no file %q", box.Name, srcPath)
			return
		}
		files = []string{srcPath}
		return
	}

	// Let the caller replace the directory part of the path with
	// the separator included.
	replacePath = dstPath + string(os.PathSeparator)

	for _, entry := range box.List() {
		if strings.HasPrefix(entry, srcPath) {
			files = append(files, entry)
		}
	}
	if len(files) == 0 {
		err = fmt.Errorf("file box %q has no files matching %q", box.Name, srcPath)
		return
	}

	trace.Trace("obtained %d files", len(files))
	return
}
