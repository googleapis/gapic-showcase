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

package genrest

import (
	"fmt"
	"strings"

	"github.com/googleapis/gapic-showcase/util/genrest/errorhandling"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
)

////////////////////////////////////////
// GoModel

type GoModel struct {
	errorhandling.Accumulator
	shim []*GoServiceShim
}

func (gm *GoModel) Add(shim *GoServiceShim) {
	gm.shim = append(gm.shim, shim)
}

func (gm *GoModel) String() string {
	shimStrings := make([]string, len(gm.shim))
	for idx, shim := range gm.shim {
		shimStrings[idx] = shim.String()
	}
	sep := "----------------------------------------"
	return fmt.Sprintf("GoModel\n%s\n%s", sep, strings.Join(shimStrings, "\n"+sep+"\n"))
}

////////////////////////////////////////
// GoServiceShim

type GoServiceShim struct {
	path     string
	imports  map[string]*pbinfo.ImportSpec
	handlers []*RESTHandler
}

func (shim *GoServiceShim) String() string {
	importStrings := make([]string, 0, len(shim.imports))
	for path, spec := range shim.imports {
		importStrings = append(importStrings, fmt.Sprintf("%s: %q %q", spec.Name, spec.Path, path))
	}

	handlerStrings := make([]string, len(shim.handlers))
	for idx, handler := range shim.handlers {
		handlerStrings[idx] = handler.String()
	}

	return fmt.Sprintf("Shim %s\n  Imports:\n    %s\n  Handlers:\n    %s",
		shim.path,
		strings.Join(importStrings, "\n    "),
		strings.Join(handlerStrings, "\n    "))
}

func (shim *GoServiceShim) AddHandler(handler *RESTHandler) {
	if shim.handlers == nil {
		shim.handlers = []*RESTHandler{}
	}
	shim.handlers = append(shim.handlers, handler)
}

func (shim *GoServiceShim) AddImports(imports ...*pbinfo.ImportSpec) {
	if shim.imports == nil {
		shim.imports = make(map[string]*pbinfo.ImportSpec, len(imports))
	}
	for _, importSpec := range imports {
		shim.imports[importSpec.Path] = importSpec
	}
}

////////////////////////////////////////
// RESTHandler

type RESTHandler struct {
	httpMethod   string
	urlMatcher   string
	pathTemplate PathTemplate

	goMethod                            string
	requestType                         string
	requestTypePackage                  string
	requestVariable                     string
	pathFields, queryFields, bodyFields []*RESTFieldMatcher

	responseType        string
	responseTypePackage string
	responseVariable    string
}

func (rh *RESTHandler) String() string {
	return fmt.Sprintf("%8s %50s func %s(%s %s.%s) (%s %s.%s) {}\n%s\n", rh.httpMethod, rh.urlMatcher, rh.goMethod, rh.requestVariable, rh.requestTypePackage, rh.requestType, rh.responseVariable, rh.responseTypePackage, rh.responseType, rh.pathTemplate)
}

////////////////////////////////////////
// RESTFieldMatcher

type RESTFieldMatcher struct {
	restName   string
	goAccessor string
}
