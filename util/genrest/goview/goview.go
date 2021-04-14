// Copyright 2020 Google LLC
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

package goview

import (
	"fmt"
	"strings"
)

// View contains a list of files to be output.
type View struct {
	Files []*SourceFile
}

// SourceFile contains a single file to be output, including both its content and location.
type SourceFile struct {
	Directory string
	Name      string // without any directory components
	source    *Source
}

// New returns a new, empty View.
func New(capacity int) *View {
	return &View{Files: make([]*SourceFile, 0, capacity)}
}

// Append appends file to this View.
func (view *View) Append(file *SourceFile) *SourceFile {
	view.Files = append(view.Files, file)
	return file
}

// NewFile creates an empty SourceFile with the specified Directory and Name.
func NewFile(directory, name string) *SourceFile {
	return &SourceFile{
		Directory: directory,
		Name:      name,
		source:    NewSource(),
	}
}

type Source struct {
	lines []string // a list of lines
}

func NewSource() *Source {
	return &Source{lines: []string{}}
}

// Contents returns the stringified contents this SourceFile.
func (sf *SourceFile) Contents() string {
	return sf.source.Contents()

}

// Append appends the lines in Source to the lines in SourceFile.
func (sf *SourceFile) Append(source *Source) {
	sf.source.lines = append(sf.source.lines, source.lines...)
}

// P appends a printf-formatted line to SourcerFile.
func (sf *SourceFile) P(format string, args ...interface{}) {
	sf.source.P(format, args...)
}

// Contents returns the stringified contents this SourceFile.
func (source *Source) Contents() string {
	return strings.Join(source.lines, "\n") + "\n"
}

// P writes a new line of content to this SourceFile. The arguments are treated exactly as in
// fmt.Printf. Note that there is an implicit in the SourceFile contents "\n" after each call to
// P().
func (source *Source) P(format string, args ...interface{}) {
	source.lines = append(source.lines, fmt.Sprintf(format, args...))
}
