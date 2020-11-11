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

package gomodel

import (
	"fmt"
	"sort"
	"strings"

	"github.com/googleapis/gapic-showcase/util/genrest/errorhandling"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
)

////////////////////////////////////////
// GoModel

type Model struct {
	errorhandling.Accumulator
	Shim []*GoServiceShim
}

func (gm *Model) Add(shim *GoServiceShim) {
	gm.Shim = append(gm.Shim, shim)
}

func (gm *Model) String() string {
	shimStrings := make([]string, len(gm.Shim))
	for idx, shim := range gm.Shim {
		shimStrings[idx] = shim.String()
	}
	sep := "----------------------------------------"
	return fmt.Sprintf("GoModel\n%s\n%s", sep, strings.Join(shimStrings, "\n"+sep+"\n"))
}

////////////////////////////////////////
// GoServiceShim

type GoServiceShim struct {
	ProtoPath string
	ShortName string
	Imports   map[string]*pbinfo.ImportSpec
	Handlers  []*RESTHandler
}

func (shim *GoServiceShim) FullName() string {
	return fmt.Sprintf("%q (%s)", shim.ShortName, shim.ProtoPath)
}

func (shim *GoServiceShim) String() string {
	importStrings := make([]string, 0, len(shim.Imports))
	for path, spec := range shim.Imports {
		importStrings = append(importStrings, fmt.Sprintf("%s: %q %q", spec.Name, spec.Path, path))
	}
	sort.Strings(importStrings)

	handlerStrings := make([]string, len(shim.Handlers))
	for idx, handler := range shim.Handlers {
		handlerStrings[idx] = handler.String()
	}
	sort.Strings(handlerStrings)

	return fmt.Sprintf("Shim %s\n  Imports:\n    %s\n  Handlers (%d):\n    %s",
		shim.FullName(),
		strings.Join(importStrings, "\n    "),
		len(handlerStrings),
		strings.Join(handlerStrings, "\n    "))
}

func (shim *GoServiceShim) AddHandler(handler *RESTHandler) {
	if shim.Handlers == nil {
		shim.Handlers = []*RESTHandler{}
	}
	shim.Handlers = append(shim.Handlers, handler)
}

func (shim *GoServiceShim) AddImports(imports ...*pbinfo.ImportSpec) {
	if shim.Imports == nil {
		shim.Imports = make(map[string]*pbinfo.ImportSpec, len(imports))
	}
	for _, importSpec := range imports {
		shim.Imports[importSpec.Path] = importSpec
	}
}

////////////////////////////////////////
// RESTHandler

type RESTHandler struct {
	HTTPMethod   string
	URIPattern   string
	PathTemplate PathTemplate

	GoMethod                            string
	RequestType                         string
	RequestTypePackage                  string
	RequestVariable                     string
	PathFields, QueryFields, BodyFields []*RESTFieldMatcher

	ResponseType        string
	ResponseTypePackage string
	ResponseVariable    string
}

func (rh *RESTHandler) String() string {
	return fmt.Sprintf("%8s %50s func %s(%s %s.%s) (%s %s.%s) {}\n%s\n", rh.HTTPMethod, rh.URIPattern, rh.GoMethod, rh.RequestVariable, rh.RequestTypePackage, rh.RequestType, rh.ResponseVariable, rh.ResponseTypePackage, rh.ResponseType, rh.PathTemplate)
}

////////////////////////////////////////
// RESTFieldMatcher

type RESTFieldMatcher struct {
	RESTName   string
	GoAccessor string
}
