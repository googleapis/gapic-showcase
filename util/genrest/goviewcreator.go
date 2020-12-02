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
		file.P(")")
		file.P("")

		for _, handler := range service.Handlers {
			handlerName := namer.Get("Handle" + handler.GoMethod)
			pathMatch := matchingPath(handler.PathTemplate)
			registered = append(registered, &registeredHandler{pathMatch, handlerName})

			file.P("")
			file.P("// %s translates REST requests/responses on the wire to internal proto messages for %s", handlerName, handler.GoMethod)
			file.P("//    Generated for HTTP binding pattern: %s", handler.URIPattern)
			file.P("//         This matches URIs of the form: %s", pathMatch)
			file.P("func (backend *RESTBackend) %s(w http.ResponseWriter, r *http.Request) {", handlerName)
			file.P(`  backend.StdLog.Printf("Received request matching '%s': %%q", r.URL)`, handler.URIPattern)
			if handler.StreamingClient || handler.StreamingServer {
				file.P(`  w.Write([]byte("ERROR: not implementing streaming methods yet"))`)
				file.P("}")
				continue
			}

			file.P("")
			file.P("  var %s *%s.%s", handler.RequestVariable, handler.RequestTypePackage, handler.RequestType)
			file.P("  // TODO: Populate %s with parameters from HTTP request", handler.RequestVariable)
			file.P("")
			// TODO: In the future, we may want to redirect all REST-endpoint requests to the gRPC endpoint so that the gRPC-registered observers get invoked.
			file.P("  %s, err := backend.%sServer.%s(context.Background(), %s)", handler.ResponseVariable, service.ShortName, handler.GoMethod, handler.RequestVariable)
			file.P("  if err != nil {")
			file.P("    // TODO: Properly handle error")
			file.P("    w.Write([]byte(err.Error()))")
			file.P("    return")
			file.P("  }")
			file.P("")
			file.P("  marshaler := &jsonpb.Marshaler{}")
			file.P("  json, err := marshaler.MarshalToString(%s)", handler.ResponseVariable)
			file.P("  if err != nil {")
			file.P("    // TODO: Properly handle error")
			file.P("    w.Write([]byte(err.Error()))")
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
	file.P("")
	file.P(" rest := (*RESTBackend)(backend)")
	for _, handler := range registered {
		file.P(`  router.HandleFunc(%q, rest.%s)`, handler.pattern, handler.function)
	}
	file.P(`}`)
	file.P("")

	return view, nil
}

// registeredHandler pairs a URL path pattern with the name of the associated handler
type registeredHandler struct {
	pattern, function string
}

// matchingPath returns the URL path for a gorilla/mux HTTP handler corresponding to the given `template`.
func matchingPath(template gomodel.PathTemplate) string {
	return extractPath(template, false)
}

// extractPath is a recursive helper function that does the actual work of
// matchingPath(). `insideVariable` denotes whether we're processing segments already inside a
// top-level handler path variable, since nested regexp groups have a different format.
func extractPath(template gomodel.PathTemplate, insideVariable bool) string {
	parts := make([]string, len(template))
	for idx, seg := range template {
		var part string
		switch seg.Kind {
		case gomodel.Literal:
			part = seg.Value
		case gomodel.SingleValue:
			part = `[a-zA-Z_%\-]+`
		case gomodel.MultipleValue:
			part = `[a-zA-Z_%\-/]+`
		case gomodel.Variable:
			if !insideVariable {
				part = fmt.Sprintf("{%s:%s}", seg.Value, extractPath(seg.Subsegments, true))
			} else {
				part = fmt.Sprintf("(?:%s)", extractPath(seg.Subsegments, true))
			}

		}
		parts[idx] = part

	}
	return strings.Join(parts, "")
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
