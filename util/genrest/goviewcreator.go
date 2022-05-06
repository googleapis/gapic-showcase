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

		fileImports := map[string]string{
			"context":  "",
			"net/http": "",
			"github.com/googleapis/gapic-showcase/util/genrest/resttools": "",
			"github.com/gorilla/mux":                               "gmux",
			"github.com/googleapis/gapic-showcase/server/genproto": "genprotopb",
		}

		// TODO: Properly deal with import strings. They may need to be taken out of the gomodel

		// Accumulate source code for each method corresponding to an RPC, as well as the imports the code requires.
		methodSources := []*goview.Source{}

		// Accumulate helper methods required by RPC methods, taking care to not duplicate them.
		helperSources := sourceMap{}

		for _, handler := range service.Handlers {
			source := goview.NewSource()
			methodSources = append(methodSources, source)

			excludedQueryParams := []string{}
			handlerName := namer.Get("Handle" + handler.GoMethod)
			pathMatch, allURLVariables, err := matchingPath(handler.PathTemplate)
			if err != nil {
				return nil, fmt.Errorf("processing %q: %s", handler.PathTemplate, err)
			}
			registered = append(registered, &registeredHandler{pathMatch, handlerName, handler.HTTPMethod})

			source.P("")
			source.P("// %s translates REST requests/responses on the wire to internal proto messages for %s", handlerName, handler.GoMethod)
			source.P("//    Generated for HTTP binding pattern: %q", handler.URIPattern)
			source.P("func (backend *RESTBackend) %s(w http.ResponseWriter, r *http.Request) {", handlerName)
			if handler.StreamingClient {
				source.P(`  backend.Error(w, http.StatusNotImplemented, "client-streaming methods not implemented yet (request matched '%s': %%q)", r.URL)`, handler.URIPattern)
				source.P("}")
				continue
			}

			source.P(`  urlPathParams := gmux.Vars(r)`)
			source.P("  numUrlPathParams := len(urlPathParams)")
			source.P("")
			// TODO: Consider factoring out code shared among handlers into a single
			// place, so that handlers only provide the relevant values (e.g. expected
			// number of path variables, etc.)
			source.P(`  backend.StdLog.Printf("Received %%s request matching '%s': %%q", r.Method, r.URL)`, handler.URIPattern)
			source.P(`  backend.StdLog.Printf("  urlPathParams (expect %d, have %%d): %%q", numUrlPathParams, urlPathParams)`, len(allURLVariables))
			source.P("")
			source.P("  if numUrlPathParams!=%d {", len(allURLVariables))
			source.P(`    backend.Error(w, http.StatusBadRequest, "found unexpected number of URL variables: expected %d, have %%d: %%#v", numUrlPathParams, urlPathParams)`, len(allURLVariables))
			source.P("    return")
			source.P("  }")

			source.P("")
			source.P("  %s := &%s.%s{}", handler.RequestVariable, handler.RequestTypePackage, handler.RequestType)
			switch handler.RequestBodyFieldSpec {
			case gomodel.BodyFieldAll:
				fileImports["bytes"] = ""
				fileImports["io"] = ""

				source.P("  // Intentional: Field values in the URL path override those set in the body.")
				source.P("  var jsonReader bytes.Buffer")
				source.P("  bodyReader := io.TeeReader(r.Body, &jsonReader)")
				source.P("  rBytes := make([]byte, r.ContentLength)")
				source.P("  if _, err := bodyReader.Read(rBytes); err != nil && err != io.EOF {")
				source.P(`    backend.Error(w, http.StatusBadRequest, "error reading body content: %%s", err)`)
				source.P("    return")
				source.P("  }")
				source.P("")
				source.P("  if err := resttools.FromJSON().Unmarshal(rBytes, %s); err != nil {", handler.RequestVariable)
				source.P(`    backend.Error(w, http.StatusBadRequest, "error reading body params '*': %%s", err)`)
				source.P("    return")
				source.P("  }")
				source.P("")
				source.P("  if err := resttools.CheckRequestFormat(&jsonReader, r, %s.ProtoReflect()); err != nil {", handler.RequestVariable)
				source.P(`    backend.Error(w, http.StatusBadRequest, "REST request failed format check: %%s", err)`)
				source.P("    return")
				source.P("  }")
				source.P("")
				source.P("  if queryParams := r.URL.Query(); len(queryParams) > 0 {")
				source.P(`    backend.Error(w, http.StatusBadRequest, "encountered unexpected query params: %%v", queryParams)`)
				source.P("    return")
				source.P("  }")

			case gomodel.BodyFieldSingle:
				fileImports["bytes"] = ""
				fileImports["io"] = ""

				// TODO: Ensure this works when the specified field is a scalar. We
				// may need to use PopulateFields from the generated code in that
				// case.
				source.P("  // Intentional: Field values in the URL path override those set in the body.")
				source.P("  var %s %s.%s", handler.RequestBodyFieldVariable, handler.RequestBodyFieldPackage, handler.RequestBodyFieldType)
				source.P("  var jsonReader bytes.Buffer")
				source.P("  bodyReader := io.TeeReader(r.Body, &jsonReader)")
				source.P("  rBytes := make([]byte, r.ContentLength)")
				source.P("  if _, err := bodyReader.Read(rBytes); err != nil && err != io.EOF {")
				source.P(`    backend.Error(w, http.StatusBadRequest, "error reading body content: %%s", err)`)
				source.P("    return")
				source.P("  }")
				source.P("")
				source.P("  if err := resttools.FromJSON().Unmarshal(rBytes, &%s); err != nil {", handler.RequestBodyFieldVariable)
				source.P(`    backend.Error(w, http.StatusBadRequest, "error reading body into request field '%s': %%s", err)`, handler.RequestBodyFieldProtoName)
				source.P("    return")
				source.P("  }")
				source.P("")
				source.P("  if err := resttools.CheckRequestFormat(&jsonReader, r, %s.ProtoReflect()); err != nil {", handler.RequestVariable)
				source.P(`    backend.Error(w, http.StatusBadRequest, "REST request failed format check: %%s", err)`)
				source.P("    return")
				source.P("  }")
				source.P("  %s.%s = &%s", handler.RequestVariable, handler.RequestBodyFieldName, handler.RequestBodyFieldVariable)
				source.P("")
				excludedQueryParams = append(excludedQueryParams, handler.RequestBodyFieldProtoName)

			default:
				source.P("  if err := resttools.CheckRequestFormat(nil, r, %s.ProtoReflect()); err != nil {", handler.RequestVariable)
				source.P(`    backend.Error(w, http.StatusBadRequest, "REST request failed format check: %%s", err)`)
				source.P("    return")
				source.P("  }")
			}

			source.P("  if err := resttools.PopulateSingularFields(%s, urlPathParams); err != nil {", handler.RequestVariable)
			source.P(`    backend.Error(w, http.StatusBadRequest, "error reading URL path params: %%s", err)`)
			source.P("    return")
			source.P("  }")
			source.P("")
			excludedQueryParams = append(excludedQueryParams, handler.PathTemplate.ListVariables()...)

			if handler.RequestBodyFieldSpec != gomodel.BodyFieldAll {
				source.P("  // TODO: Decide whether query-param value or URL-path value takes precedence when a field appears in both")
				source.P("  queryParams := map[string][]string(r.URL.Query())")
				if len(excludedQueryParams) > 0 {
					source.P("  excludedQueryParams := %#v", excludedQueryParams)
					source.P("  if duplicates := resttools.KeysMatchPath(queryParams, excludedQueryParams); len(duplicates) > 0 {")
					source.P(`    backend.Error(w, http.StatusBadRequest, "(QueryParamsInvalidFieldError) found keys that should not appear in query params: %%v", duplicates)`)
					source.P("    return")
					source.P("  }")
				}
				source.P("  if err := resttools.PopulateFields(%s, queryParams); err != nil {", handler.RequestVariable)
				source.P(`    backend.Error(w, http.StatusBadRequest, "error reading query params: %%s", err)`)
				source.P("    return")
				source.P("  }")
				source.P("")
			}
			source.P("")
			source.P("  marshaler := resttools.ToJSON()")
			source.P("  requestJSON, _ := marshaler.Marshal(%s)", handler.RequestVariable)
			source.P(`  backend.StdLog.Printf("  request: %%s", requestJSON)`)
			source.P("")

			if handler.StreamingServer {
				streamerType := constructServerStreamer(service, handler, fileImports, helperSources)

				source.P(`  serverStreamer, err := resttools.NewServerStreamer(w, resttools.ServerStreamingChunkSize)`)
				source.P(`  if err != nil {`)
				source.P(`    backend.Error(w,http.StatusInternalServerError, "server error: could not construct server streamer: %%s", err.Error())`)
				source.P(`    return`)
				source.P(`  }`)
				source.P(`  defer serverStreamer.End()`)

				source.P(` streamer := &%s{serverStreamer}`, streamerType)

				source.P(" if err := backend.%sServer.%s(%s, streamer); err != nil {", service.ShortName, handler.GoMethod, handler.RequestVariable)
				source.P("   backend.ReportGRPCError(w, err)")
				source.P(" }")

			} else { // regular unary call
				// TODO: In the future, we may want to redirect all REST-endpoint requests to the gRPC endpoint so that the gRPC-registered observers get invoked.
				source.P("  %s, err := backend.%sServer.%s(context.Background(), %s)", handler.ResponseVariable, service.ShortName, handler.GoMethod, handler.RequestVariable)
				source.P("  if err != nil {")
				source.P("    backend.ReportGRPCError(w, err)")
				source.P("    return")
				source.P("  }")
				source.P("")
				source.P("  json, err := marshaler.Marshal(%s)", handler.ResponseVariable)
				source.P("  if err != nil {")
				source.P(`    backend.Error(w, http.StatusInternalServerError, "error json-encoding response: %%s", err.Error())`)
				source.P("    return")
				source.P("  }")
				source.P("")
				source.P("  w.Write(json)")
			}
			source.P("}\n")
		}

		// Use the accumulated method source code and import data to write out the actual file.

		file.P("import (")
		for path, name := range fileImports {
			file.P("  %s %q", name, path)
		}
		file.P(")")
		file.P("")
		for _, method := range methodSources {
			file.Append(method)
		}

		for _, k := range helperSources.sortedKeys() {
			file.Append(helperSources[k])
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
	file.P(`   "github.com/googleapis/gapic-showcase/util/genrest/resttools"`)
	file.P("")
	file.P(`  gmux "github.com/gorilla/mux"`)
	file.P(`  "google.golang.org/grpc/status"`)
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

	for _, handler := range registered {
		file.P(`  router.HandleFunc(%q, rest.%s).Methods(%q)`, handler.pattern, handler.function, handler.verb)
	}
	file.P(`  router.PathPrefix("/").HandlerFunc(rest.catchAllHandler)`)
	file.P(`}`)
	file.P("")

	file.P("func (backend *RESTBackend) catchAllHandler(w http.ResponseWriter, r *http.Request) {")
	file.P(`  backend.Error(w, http.StatusBadRequest, "unrecognized request: %%s %%q", r.Method, r.URL)`)
	file.P("}")
	file.P("")

	file.P("func (backend *RESTBackend) Error(w http.ResponseWriter, status int, format string, args ...interface{}) {")
	file.P("  message := fmt.Sprintf(format, args...)")
	file.P("  backend.ErrLog.Print(message)")
	file.P("  resttools.ErrorResponse(w, status, message)")
	file.P("}")

	file.P("func (backend *RESTBackend) ReportGRPCError(w http.ResponseWriter, err error) {")
	file.P("  st, ok := status.FromError(err)")
	file.P("  if !ok {")
	file.P(`  	backend.Error(w, http.StatusInternalServerError, "server error: %%s", err.Error())`)
	file.P("    return")
	file.P("  }")
	file.P("")
	file.P("  backend.ErrLog.Print(st.Message())")
	file.P("  code := resttools.GRPCToHTTP[st.Code()]")
	file.P("  resttools.ErrorResponse(w, code, st.Message(), st.Details()...)")
	file.P("}")

	return view, nil
}

func constructServerStreamer(service *gomodel.ServiceModel, handler *gomodel.RESTHandler, fileImports map[string]string, helperSources sourceMap) (streamerType string) {
	streamerType = fmt.Sprintf("%s_%sServer", service.ShortName, handler.GoMethod)
	streamerInterface := fmt.Sprintf("%s.%s_%sServer", handler.RequestTypePackage, service.ShortName, handler.GoMethod)
	streamerElement := fmt.Sprintf("*%s.%s", handler.ResponseTypePackage, handler.ResponseType)

	helper := goview.NewSource()
	helperSources[streamerType] = helper

	helper.P(`// %s implements %s to provide server-side streaming over REST, returning all the`, streamerType, streamerInterface)
	helper.P(`// individual responses as part of a long JSON list.`)
	helper.P(`type %s struct{`, streamerType)
	helper.P(`   *resttools.ServerStreamer`)
	helper.P(`}`)
	helper.P(``)
	helper.P(` // Send accumulates a response to be fetched later as part of response list returned over REST.`)
	helper.P(`func (streamer *%s) Send(response %s) error {`, streamerType, streamerElement)
	helper.P(`  return streamer.ServerStreamer.Send(response)`)
	helper.P(`}`)

	return streamerType
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

			// Here we convert the proto-cased (snake-cased) field path to be JSON-cased
			// (lower-camel-cased) so that we can keep the resttools.Populate*Field*()
			// functions simple, only dealing with JSON-cased field names,
			part = fmt.Sprintf("{%s:%s}", resttools.ToDottedLowerCamel(seg.Value), subParts)
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

// sourceMap implements a method to return the keys in sorted order.
type sourceMap map[string]*goview.Source

// sortedKeys returns a slice containing all the keys in sm, sorted.
func (sm sourceMap) sortedKeys() []string {
	keys := []string{}
	for k := range sm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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
