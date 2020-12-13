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
	"regexp"
	"sort"
	"strings"

	"github.com/googleapis/gapic-showcase/util/genrest/errorhandling"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
)

////////////////////////////////////////
// Model

// Model is a data model intended to be capture all the information needed for generating Go source
// files that provide shims between REST messages on the wire and protocol buffer messages in the
// back end, for multiple services.
type Model struct {
	errorhandling.Accumulator
	Service []*ServiceModel
}

// Add adds `service` to this Model.
func (gm *Model) Add(service *ServiceModel) {
	gm.Service = append(gm.Service, service)
}

// String returns a string representation of this Model.
func (gm *Model) String() string {
	shimStrings := make([]string, len(gm.Service))
	for idx, shim := range gm.Service {
		shimStrings[idx] = shim.String()
	}
	sep := "----------------------------------------"
	return fmt.Sprintf("GoModel\n%s\n%s", sep, strings.Join(shimStrings, "\n"+sep+"\n"))
}

// CheckConsistency checks this Model for consistency, accumulating
// any errors found. This means checking that all the HTTP annotations
// across all services resolve to distinct paths.
func (gm *Model) CheckConsistency() { // TODO: CheckCorrectness? Or move the checks to the creation part
	reBodyField := regexp.MustCompile(resttools.RegexField)
	allHandlers := []*RESTHandler{}

	for _, service := range gm.Service {
		allHandlers = append(allHandlers, service.Handlers...)
	}

	for first, firstHandler := range allHandlers {
		if len(firstHandler.BodyField) > 0 && firstHandler.BodyField != "*" {

			// The body field name refers to a top-level field.
			if reBodyField.FindStringIndex(firstHandler.BodyField) == nil {
				gm.AccumulateError(fmt.Errorf("bad syntax in body field spec %q", firstHandler.BodyField))
			}

		}

		// Generators should test that  if `foo` is specified as the `body`, it doesn't appear in the path parameters

		if _, nestedVariables := firstHandler.PathTemplate.HasVariables(); nestedVariables {
			gm.AccumulateError(fmt.Errorf("pattern %q specifies nested variables, which are not allowed as per https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#path-template-syntax", firstHandler.URIPattern))
		}
		for _, secondHandler := range allHandlers[first+1:] {
			if firstHandler.HTTPMethod != secondHandler.HTTPMethod {
				continue
			}

			fullMatch, ambiguousPattern, err := FindValuesMatching(firstHandler.PathTemplate, secondHandler.PathTemplate)
			if err != nil {
				gm.AccumulateError(fmt.Errorf("matching patterns %q and %q (constructed %q): %s", firstHandler, secondHandler, ambiguousPattern, err))
				continue
			}
			if !fullMatch {
				continue
			}
			gm.AccumulateError(fmt.Errorf("pattern %q matches both\n   %s and\n   %s\n\n", ambiguousPattern, firstHandler, secondHandler))
		}

	}
}

////////////////////////////////////////
// ServiceModel

// ServiceModel is a data model for generating a REST/proto shim for
// a single protocol buffer service.
type ServiceModel struct {
	// the fully qualified protocol buffer type name of this service
	ProtoPath string

	// the non-qualified name of this service
	ShortName string

	// a map of import paths to import info for each of the service-related Go imports that will
	// be needed to implement all of the Handlers
	Imports map[string]*pbinfo.ImportSpec

	// a list of all the HTTP handlers that will need to be generated, one for each HTTP
	// annotation for each service RPC
	Handlers []*RESTHandler
}

// FullName pretty-prints the short name and proto path.
func (service *ServiceModel) FullName() string {
	return fmt.Sprintf("%q (%s)", service.ShortName, service.ProtoPath)
}

// String returns a string representation of this ServiceModel.
func (service *ServiceModel) String() string {
	importStrings := make([]string, 0, len(service.Imports))
	for path, spec := range service.Imports {
		importStrings = append(importStrings, fmt.Sprintf("%s: %q %q", spec.Name, spec.Path, path))
	}
	sort.Strings(importStrings)

	handlerStrings := make([]string, len(service.Handlers))
	for idx, handler := range service.Handlers {
		handlerStrings[idx] = handler.String()
	}
	sort.Strings(handlerStrings)

	return fmt.Sprintf("Shim %s\n  Imports:\n    %s\n  Handlers (%d):\n    %s",
		service.FullName(),
		strings.Join(importStrings, "\n    "),
		len(handlerStrings),
		strings.Join(handlerStrings, "\n    "))
}

// AddHandler adds handler to this ServiceModel.
func (service *ServiceModel) AddHandler(handler *RESTHandler) {
	if service.Handlers == nil {
		service.Handlers = []*RESTHandler{}
	}
	service.Handlers = append(service.Handlers, handler)
}

// AddImports adds each element of imports to the imports in this ServiceModel.
func (service *ServiceModel) AddImports(imports ...*pbinfo.ImportSpec) {
	if service.Imports == nil {
		service.Imports = make(map[string]*pbinfo.ImportSpec, len(imports))
	}
	for _, importSpec := range imports {
		service.Imports[importSpec.Path] = importSpec
	}
}

////////////////////////////////////////
// RESTHandler

// RESTHandler contains the information needed to generate a single HTTP handler.
type RESTHandler struct {
	//// Transcoding information

	HTTPMethod      string
	URIPattern      string       // as it appears in the HTTP annotation
	BodyField       string       // single field path, or "*" (all fields), or "" (no field)
	BodyFieldType   string       // single field path, or "*" (all fields), or "" (no field)
	PathTemplate    PathTemplate // parsed version of URIPattern
	StreamingServer bool         // whether this method uses server-side streaming
	StreamingClient bool         // whether this method uses client-side streaming

	// TODO: Fill in with actual information needed to access each field. These are placeholders
	// for now.
	PathFields, QueryFields, BodyFields []*interface{}

	//// Go types

	GoMethod                string
	RequestType             string
	RequestTypePackage      string
	RequestVariable         string
	RequestBodyFieldName    string
	RequestBodyFieldType    string
	RequestBodyFieldPackage string
	ResponseType            string
	ResponseTypePackage     string
	ResponseVariable        string
}

// String returns a string representation of this RESTHandler.
func (rh *RESTHandler) String() string {
	return fmt.Sprintf("%8s %50s func %s(%s %s.%s) (%s %s.%s) {}\n%s\n", rh.HTTPMethod, rh.URIPattern, rh.GoMethod, rh.RequestVariable, rh.RequestTypePackage, rh.RequestType, rh.ResponseVariable, rh.ResponseTypePackage, rh.ResponseType, rh.PathTemplate)
}
