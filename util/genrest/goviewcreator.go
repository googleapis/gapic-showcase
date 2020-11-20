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

// NewView creates a a new goview.View (a series of files to be output) from a gomodel.Model.
func NewView(model *gomodel.Model) (*goview.View, error) {
	// TODO: Assert that all services live in the same proto package, which is currently the case and we currently assume.
	namer := NewNamer()

	view := goview.New(len(model.Service))
	registered := []*registeredPath{}

	for idxService, service := range model.Service {
		file := view.Append(goview.NewFile("", strings.ToLower(service.ShortName)+".go"))
		file.P(license)
		file.P("// DO NOT EDIT. This is an auto-generated file containing the REST handlers")
		file.P("// for service #%d: %s.\n", idxService, service.FullName())

		// TODO: add license to the top of the files

		file.P("")
		file.P("package genrest")
		file.P("")

		importStrings := make([]string, 0, len(service.Imports))
		for _, spec := range service.Imports {
			importStrings = append(importStrings, fmt.Sprintf("%s %q", spec.Name, spec.Path))
		}
		sort.Strings(importStrings)

		file.P("import (")
		file.P(`  "fmt"`)
		file.P(`  "net/http"`)
		file.P("")
		file.P("  // %s", strings.Join(importStrings, "\n  // "))
		file.P(")")
		file.P("")

		for _, handler := range service.Handlers {
			handlerName := namer.Get("Handle" + handler.GoMethod)
			pathMatch := matchingPath(handler.PathTemplate)
			registered = append(registered, &registeredPath{pathMatch, handlerName})

			file.P("")
			file.P("// %s translates REST requests/responses on the wire to internal proto messages for %s", handlerName, handler.GoMethod)
			file.P("//    Generated for HTTP binding pattern: %s", handler.URIPattern)
			file.P("//             This matches URIs of form: %s", pathMatch)
			file.P("func %s(w http.ResponseWriter, r *http.Request) {", handlerName)
			file.P("  // serialize from r to: var %s %s.%s", handler.RequestVariable, handler.RequestTypePackage, handler.RequestType)
			file.P("  // %s := %s(%s)", handler.ResponseVariable, handler.GoMethod, handler.RequestVariable)
			file.P("  // serialize back to w")
			file.P(`  fmt.Fprintf(w, "Received request matching '%s': %%q", r.URL)`, handler.URIPattern)
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
	file.P(`  gmux "github.com/gorilla/mux"`)
	file.P(")")
	file.P("")
	file.P("")
	file.P(`func RegisterHandlers(router *gmux.Router) {`)
	for _, handler := range registered {
		file.P(`  router.HandleFunc(%q, %s)`, handler.pattern, handler.function)
	}
	file.P(`}`)

	return view, nil
}

type registeredPath struct {
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
