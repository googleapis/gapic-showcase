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

type View struct {
	Files []*SourceFile
}

type SourceFile struct {
	Directory string
	Name      string
	contents  []string
}

func New(capacity int) *View {
	return &View{Files: make([]*SourceFile, 0, capacity)}
}

func (view *View) Append(file *SourceFile) *SourceFile {
	view.Files = append(view.Files, file)
	return file
}

func NewFile(directory, name string) *SourceFile {
	return &SourceFile{
		Directory: directory,
		Name:      name,
		contents:  []string{},
	}
}

func (sf *SourceFile) Contents() string {
	return strings.Join(sf.contents, "\n")
}

func (sf *SourceFile) P(format string, args ...interface{}) {
	sf.contents = append(sf.contents, fmt.Sprintf(format, args...))
}
