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
	"sort"
	"strings"
	"time"

	"github.com/googleapis/gapic-showcase/util/genrest/gomodel"
	"github.com/googleapis/gapic-showcase/util/genrest/goview"
	"github.com/googleapis/gapic-showcase/util/genrest/resttools"
)

// NewView creates a a new goview.View (a series of files to be output) from a gomodel.Model. The
// current approach is to generate one file per service, with that file containing all the service's
// RPCs. An additional file `genrest.go` is also created to register all these handlers with a
// gorilla/mux dispatcher.
func NewView(model *gomodel.Model) (*goview.View, error) {
	// TODO: Assert that all services live in the same proto package, which is currently the case and we currently assume.
	namer := NewNamer()

	numServices := len(model.Service)
	view := goview.New(numServices)
	registered := []*registeredHandler{}

	for idxService, service := range model.Service {
		file := view.Append(goview.NewFile("", strings.ToLower(service.ShortName)+".go"))
		file.P(license)
		file.P("// DO NOT EDIT. This is an auto-generated file containing the REST handlers")
		file.P("// for service #%d: %s.\n", idxService, service.FullName())
		file.P("")
		file.P("package genrest")
		file.P("")

		importStrings := make([]string, 0, len(service.Imports))
		for _, spec := range service.Imports {
			importStrings = append(importStrings, fmt.Sprintf("%s %q", spec.Name, spec.Path))
		}
		sort.Strings(importStrings)

		file.P("import (")
		file.P(`  "context"`)
		file.P(`  "net/http"`)
		file.P("")
		// TODO: Properly deal with imports once we actually use them in the code.
		// file.P("  // %s", strings.Join(importStrings, "\n  // "))

		// TODO: Get the following import from the data model, rather than hard-coding, once we deal with all imports.
		file.P(`	  genprotopb "github.com/googleapis/gapic-showcase/server/genproto"`)

		// TODO: Investigate why "google.golang.org/protobuf" doesn't work below
		file.P(`  "github.com/golang/protobuf/jsonpb"`)
		file.P(`  gmux "github.com/gorilla/mux"`)
		file.P("")
		file.P(`  "github.com/googleapis/gapic-showcase/util/genrest/resttools"`)
		file.P(")")
		file.P("")
		for _, handler := range service.Handlers {
			excludedQueryParams := []string{}
			handlerName := namer.Get("Handle" + handler.GoMethod)
			pathMatch, allURLVariables, err := matchingPath(handler.PathTemplate)
			if err != nil {
				return nil, fmt.Errorf("processing %q: %s", handler.PathTemplate, err)
			}
			registered = append(registered, &registeredHandler{pathMatch, handlerName, handler.HTTPMethod})

			file.P("")
			file.P("// %s translates REST requests/responses on the wire to internal proto messages for %s", handlerName, handler.GoMethod)
			file.P("//    Generated for HTTP binding pattern: %q", handler.URIPattern)
			file.P("func (backend *RESTBackend) %s(w http.ResponseWriter, r *http.Request) {", handlerName)
			if handler.StreamingClient || handler.StreamingServer {
				file.P(`  backend.Error(w, http.StatusNotImplemented, "streaming methods not implemented yet (request matched '%s': %%q)", r.URL)`, handler.URIPattern)
				file.P("}")
				continue
			}

			file.P(`  urlPathParams := gmux.Vars(r)`)
			file.P("  numUrlPathParams := len(urlPathParams)")
			file.P("")
			// TODO: Consider factoring out code shared among handlers into a single
			// place, so that handlers only provide the relevant values (eg,, expected
			// number of path variables, etc.)
			file.P(`  backend.StdLog.Printf("Received %%s request matching '%s': %%q", r.Method, r.URL)`, handler.URIPattern)
			file.P(`  backend.StdLog.Printf("  urlPathParams (expect %d, have %%d): %%q", numUrlPathParams, urlPathParams)`, len(allURLVariables))
			file.P("")
			file.P("  if numUrlPathParams!=%d {", len(allURLVariables))
			file.P(`    backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected %d, have %%d: %%#v", numUrlPathParams, urlPathParams)`, len(allURLVariables))
			file.P("    return")
			file.P("  }")

			file.P("")
			file.P("  %s := &%s.%s{}", handler.RequestVariable, handler.RequestTypePackage, handler.RequestType)
			switch handler.RequestBodyFieldSpec {
			case gomodel.BodyFieldAll:
				file.P("  // Intentional: Field values in the URL path override those set in the body.")
				file.P("  if err := jsonpb.Unmarshal(r.Body, %s); err != nil {", handler.RequestVariable)
				file.P(`    backend.Error(w, http.StatusBadRequest, "error reading body params '*': %%s", err)`)
				file.P("    return")
				file.P("  }")
				file.P("")
				file.P("  if queryParams := r.URL.Query(); len(queryParams) > 0 {")
				file.P(`    backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %%v", queryParams)`)
				file.P("    return")
				file.P("  }")

			case gomodel.BodyFieldSingle:
				// TODO: Ensure this works when the specified field is a scalar. We
				// may need to use PopulateFields from the generated code in that
				// case.
				file.P("  // Intentional: Field values in the URL path override those set in the body.")
				file.P("  var %s %s.%s", handler.RequestBodyFieldVariable, handler.RequestBodyFieldPackage, handler.RequestBodyFieldType)
				file.P("  if err := jsonpb.Unmarshal(r.Body, &%s); err != nil {", handler.RequestBodyFieldVariable)
				file.P(`    backend.Error(w, http.StatusBadRequest, "error reading body into request field '%s': %%s", err)`, handler.RequestBodyFieldProtoName)
				file.P("    return")
				file.P("  }")
				file.P("  %s.%s = &%s", handler.RequestVariable, handler.RequestBodyFieldName, handler.RequestBodyFieldVariable)
				file.P("")
				excludedQueryParams = append(excludedQueryParams, handler.RequestBodyFieldProtoName)
			}

			file.P("  if err := resttools.PopulateSingularFields(%s, urlPathParams); err != nil {", handler.RequestVariable)
			file.P(`    backend.Error(w, http.StatusBadRequest, "error reading URL path params: %%s", err)`)
			file.P("    return")
			file.P("  }")
			file.P("")
			excludedQueryParams = append(excludedQueryParams, handler.PathTemplate.ListVariables()...)

			if handler.RequestBodyFieldSpec != gomodel.BodyFieldAll {
				file.P("  // TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both")
				file.P("  queryParams := map[string][]string(r.URL.Query())")
				if len(excludedQueryParams) > 0 {
					file.P("  excludedQueryParams := %#v", excludedQueryParams)
					file.P("  if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {")
					file.P(`    backend.Error(w, http.StatusBadRequest, " found keys that should not appear in query params: %%v", duplicates)`)
					file.P("    return")
					file.P("  }")
				}
				file.P("  if err := resttools.PopulateFields(%s, queryParams); err != nil {", handler.RequestVariable)
				file.P(`    backend.Error(w, http.StatusBadRequest, "error reading query params: %%s", err)`)
				file.P("    return")
				file.P("  }")
				file.P("")
			}
			file.P("")
			file.P("  marshaler := &jsonpb.Marshaler{}")
			file.P("  requestJSON, _ := marshaler.MarshalToString(%s)", handler.RequestVariable)
			file.P(`  backend.StdLog.Printf("  request: %%s", requestJSON)`)
			file.P("")
			// TODO: In the future, we may want to redirect all REST-endpoint requests to the gRPC endpoint so that the gRPC-registered observers get invoked.
			file.P("  %s, err := backend.%sServer.%s(context.Background(), %s)", handler.ResponseVariable, service.ShortName, handler.GoMethod, handler.RequestVariable)
			file.P("  if err != nil {")
			file.P("    // TODO: Properly handle error. Is StatusInternalServerError (500) the right response?")
			file.P(`    backend.Error(w, http.StatusInternalServerError, "server error: %%s", err.Error())`)
			file.P("    return")
			file.P("  }")
			file.P("")
			file.P("  json, err := marshaler.MarshalToString(%s)", handler.ResponseVariable)
			file.P("  if err != nil {")
			file.P(`    backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %%s", err.Error())`)
			file.P("    return")
			file.P("  }")
			file.P("")
			file.P("  w.Write([]byte(json))")
			file.P("}\n")
		}
	}

	file := view.Append(goview.NewFile("", "genrest.go"))
	file.P(license)
	file.P("// DO NOT EDIT. This is an auto-generated file registering the REST handlers.")
	file.P("// for the various Showcase services.")
	file.P("")
	file.P("package genrest")
	file.P("")
	file.P("import (")
	file.P(`  "fmt"`)
	file.P(`   "net/http"`)
	file.P("")
	file.P(`   "github.com/googleapis/gapic-showcase/server/services"`)
	file.P("")
	file.P(`  gmux "github.com/gorilla/mux"`)
	file.P(")")
	file.P("")
	file.P("")
	file.P("type RESTBackend services.Backend")
	file.P("")
	file.P("")
	file.P(`func RegisterHandlers(router *gmux.Router, backend *services.Backend) {`)
	file.P(" rest := (*RESTBackend)(backend)")

	// TODO: Support path-encoded '\n' in strings (%0A), which currently don't work. Probably the way to do this is to add
	//   file.P(" router.UseEncodedPath()")
	// here, and to explicitly path decode in resttools.PopulateSingularFields. We should also
	// add '\n' to the "ExtremeValues" ComplianceGroup in compliance_suite.json.

	// TODO: Fix PATCH requests, like
	//  `curl -X PATCH http://localhost:7469/v1beta1/users/Victor`
	// which don't seem to make it through to the handler. (It doesn't seem to be an issue with
	// them being shadowed by other HTTP verb handlers, since the problem persists even when not
	// registering other handlers with the same URL pattern.) Consider using gorilla/mux
	// subroutes, selecting by HTTP verb before the URL path.
	for _, handler := range registered {
		file.P(`  router.HandleFunc(%q, rest.%s).Methods(%q)`, handler.pattern, handler.function, handler.verb)
	}
	file.P(`}`)
	file.P("")

	file.P("func (backend *RESTBackend) Error(w http.ResponseWriter, status int, format string, args ...interface{}) {")
	file.P("  message := fmt.Sprintf(format, args...)")
	file.P("  backend.StdLog.Print(message)")
	file.P("  w.WriteHeader(status)")
	file.P(`  w.Write([]byte("showcase " + message))`)
	file.P("}")

	return view, nil
}

// registeredHandler pairs a URL path pattern with the name of the associated handler
type registeredHandler struct {
	pattern  string // URL pattern
	function string // handler function
	verb     string // HTTP verb
}

// matchingPath returns the URL path for a gorilla/mux HTTP handler corresponding to the given
// `template`, as well as a list of all the template variables (identified by their proto field
// path, which are also their key values in the variable map that will be returned by gorilla/mux.Vars()).
func matchingPath(template gomodel.PathTemplate) (string, []string, error) {
	return extractPath(template, false)
}

// extractPath is a one-level recursive helper function that does the actual work of
// matchingPath(). `insideVariable` denotes whether we're processing segments already inside a
// top-level handler path variable, since nested variables are not allowed. The list of variable
// keys (which will be provided by gorilla/mux, and which also match their proto field path in the
// request object) is returned together with the actual gorilla/mux matching path.
func extractPath(template gomodel.PathTemplate, insideVariable bool) (string, []string, error) {
	var allVariables []string
	parts := make([]string, len(template))
	for idx, seg := range template {
		var part string
		switch seg.Kind {
		case gomodel.Literal:
			part = seg.Value
		case gomodel.SingleValue:
			part = resttools.RegexURLPathSingleSegmentValue
		case gomodel.MultipleValue:
			part = resttools.RegexURLPathMultipleSegmentValue
		case gomodel.Variable:
			if insideVariable {
				return "", nil, fmt.Errorf("nested variables are disallowed: https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#path-template-syntax")
			}
			subParts, _, err := extractPath(seg.Subsegments, true)
			if err != nil {
				return "", nil, err

			}
			part = fmt.Sprintf("{%s:%s}", seg.Value, subParts)
			allVariables = append(allVariables, seg.Value)
		}
		parts[idx] = part

	}
	return strings.Join(parts, ""), allVariables, nil
}

////////////////////////////////////////
// Namer

// Namer keeps track of a series of symbol names being used in order to disambiguate new names.
type Namer struct {
	registered map[string]int
}

// NewNamer returns a new Namer.
func NewNamer() *Namer {
	return &Namer{registered: make(map[string]int)}
}

// Get registers and returns a non-previously registered name that is as close to newName as
// possible. Disambiguation, if needed, is accomplished by adding a numeric suffix.
func (namer *Namer) Get(newName string) string {
	for {
		numSeen := namer.registered[newName]
		namer.registered[newName] = numSeen + 1
		if numSeen == 0 {
			return newName
		}

		newName = fmt.Sprintf("%s_%d", newName, numSeen)
		// run through the loop again to ensure the new name hasn't been previously registered either.
	}
}

var license string

func init() {
	license = fmt.Sprintf(`// Copyright %d Google LLC
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
`, time.Now().Year())
}
