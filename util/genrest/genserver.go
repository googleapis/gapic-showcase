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
	"log"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type HTTPBinding struct {
	method  string // make an enum?
	pattern string
	// maybe RPC, although maybe that should be an enveloping struct
	// maybe separate patterns for URI and for body
}

func extractBinding(rule *annotations.HttpRule) (binding *HTTPBinding, err error) {
	binding = &HTTPBinding{}
	pattern := rule.GetPattern()
	switch pattern.(type) {
	case *annotations.HttpRule_Get:
		binding.method = "GET"
		binding.pattern = rule.GetGet()
	case *annotations.HttpRule_Post:
		binding.method = "POST"
		binding.pattern = rule.GetPost()
	case *annotations.HttpRule_Patch:
		binding.method = "PATCH"
		binding.pattern = rule.GetPatch()
	case *annotations.HttpRule_Put:
		binding.method = "PUT"
		binding.pattern = rule.GetPut()
	case *annotations.HttpRule_Delete:
		binding.method = "DELETE"
		binding.pattern = rule.GetDelete()
	default:
		return nil, fmt.Errorf("unhandled pattern: %#x", pattern)
	}
	return binding, nil
}

func (binding *HTTPBinding) String() string {
	return fmt.Sprintf("%s: %q", binding.method, binding.pattern)
}

// TODO(vchudnov-g): Continue filling this in. It's a an initial empty
// stub at the moment.
func Generate(plugin *protogen.Plugin) error {
	log.Printf("Generating REST!")
	file := plugin.NewGeneratedFile("showcase-rest-sample-response.txt", "github.com/googleapis/gapic-showcase/server/genrest")

	// https://godoc.org/google.golang.org/protobuf/types/pluginpb
	// The HTTP annotation proto is defined in https://github.com/googleapis/googleapis/blob/master/google/api/http.proto

	// The typecasting below appears to be idiomatic as per
	// https://github.com/protocolbuffers/protobuf-go/blob/master/cmd/protoc-gen-go/internal_gengo/main.go#L31
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	file.P("Generated via \"google.golang.org/protobuf/compiler/protogen\" now")
	file.P("Files:\n", strings.Join(plugin.Request.FileToGenerate, "\n"))

	for idxProto, protoFile := range plugin.Request.GetProtoFile() {
		file.P(fmt.Sprintf("ProtoFile[%02d]: %q (%s)\n  Services:", idxProto, *protoFile.Name, *protoFile.Package))
		for idxSvc, svc := range protoFile.Service {
			file.P(fmt.Sprintf("    %2d: %q", idxSvc, *svc.Name))
			eDefaultHost := proto.GetExtension(svc.GetOptions(), annotations.E_DefaultHost)
			file.P(fmt.Sprintf("       Default host: %q", eDefaultHost))

			for idxSvcOptions, svcOption := range svc.GetOptions().GetUninterpretedOption() {
				file.P(fmt.Sprintf("       Option[%d]: %s", idxSvcOptions, svcOption))
			}
			file.P("        Methods")
			for idxMethod, method := range svc.GetMethod() {
				file.P(fmt.Sprintf("        %02d %s", idxMethod, method.GetName()))
				options := method.GetOptions()
				if options == nil {
					continue
				}

				file.P(fmt.Sprintf("           %s",
					method.Options.GetIdempotencyLevel()))

				uninterpreted := options.GetUninterpretedOption()
				file.P(fmt.Sprintf("           uninterpreted length: %d", len(uninterpreted)))
				for idxOpt, opt := range uninterpreted {
					file.P(fmt.Sprintf("          %d: %x : %s %q", idxOpt, opt.GetName(), opt.GetIdentifierValue(), string(opt.StringValue)))
				}

				// TODO: explore how this works; look up references
				// This was adapted from dev/googleapis/gapic-generator-go/internal/gengapic/gengapic.go:parseRequestHeaders()
				// Need to see what this means for the new library we're supposed to use
				eHTTP /*, err*/ := proto.GetExtension(method.GetOptions(), annotations.E_Http)
				http := eHTTP.(*annotations.HttpRule)
				rules := []*annotations.HttpRule{http}
				rules = append(rules, http.GetAdditionalBindings()...)
				file.P("           Rules:")
				for idxRule, oneRule := range rules {
					binding, err := extractBinding(oneRule)
					if err != nil {
						file.P(fmt.Sprintf("           %d: ERROR extracting binding from: %s", idxRule, oneRule))
						continue

					}
					file.P(fmt.Sprintf("           %d: %s", idxRule, binding))
				}

			}
		}
	}
	return nil
}
