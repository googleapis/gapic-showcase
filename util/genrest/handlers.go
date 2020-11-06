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

	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

////////////////////////////////////////
// HandlerGenerator

type HandlerGenerator struct {
	descInfo pbinfo.Info
	services []*Service
}

func (hg *HandlerGenerator) String() string {
	services := make([]string, len(hg.services))
	for idx, svc := range hg.services {
		if svc == nil {
			continue
		}
		services[idx] = svc.String()
	}
	return strings.Join(services, "\n\n")
}

func (hg *HandlerGenerator) AddService(service *Service) {
	hg.services = append(hg.services, service)
}

////////////////////////////////////////
// Service

type Service struct {
	descriptor      *descriptorpb.ServiceDescriptorProto
	name            string
	typeName        string
	patternHandlers []*PatternHandler
}

func (service *Service) String() string {
	handlers := make([]string, len(service.patternHandlers))
	for idx, h := range service.patternHandlers {
		handlers[idx] = h.String()
	}
	indent := "  "
	return fmt.Sprintf("%s (%s):\n%s%s", service.name, service.typeName, indent, strings.Join(handlers, "\n"+indent))
}

func NewService(descriptor *descriptorpb.ServiceDescriptorProto) *Service {
	// might need *descriptor.ProtoReflect().Type() to instantiate Go type and then reflect to get the name?
	service := &Service{
		descriptor: descriptor,
		name:       *descriptor.Name,
		typeName:   string(descriptor.ProtoReflect().Descriptor().FullName()),
	}
	return service
}

////////////////////////////////////////
// PatternHandler

type PatternHandler struct {
	prefix   string
	handlers []*Handler
}

func (ph *PatternHandler) String() string {
	handlers := make([]string, len(ph.handlers))
	for idx, h := range ph.handlers {
		handlers[idx] = h.String()
	}
	return fmt.Sprintf("%s {%s}", ph.prefix, strings.Join(handlers, "} {"))
}

////////////////////////////////////////
// Handler

type Handler struct {
	methodType  string
	httpBinding HTTPBinding
}

func (handler *Handler) String() string {
	return fmt.Sprintf("%s : %s", handler.methodType, handler.httpBinding)
}

////////////////////////////////////////
// HTTPBinding

type HTTPBinding struct {
	method  string // make an enum?
	pattern string
	// maybe RPC, although maybe that should be an enveloping struct
	// maybe separate patterns for URI and for body
}

func (binding *HTTPBinding) String() string {
	return fmt.Sprintf("%s: %q", binding.method, binding.pattern)
}

func NewGenerator(plugin *protogen.Plugin) (*HandlerGenerator, error) {
	protoFiles := plugin.Request.GetProtoFile()
	hg := &HandlerGenerator{
		descInfo: pbinfo.Of(protoFiles),
		services: make([]*Service, 0, len(protoFiles)),
	}

	file := plugin.NewGeneratedFile("showcase-rest-sample-response.txt", "github.com/googleapis/gapic-showcase/server/genrest")

	// The typecasting below appears to be idiomatic as per
	// https://github.com/protocolbuffers/protobuf-go/blob/master/cmd/protoc-gen-go/internal_gengo/main.go#L31
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	file.P("Generated via \"google.golang.org/protobuf/compiler/protogen\" via HandlerGenerator!")
	file.P("Files:\n", strings.Join(plugin.Request.FileToGenerate, "\n"))

	for idxProto, protoFile := range protoFiles {
		file.P(fmt.Sprintf("ProtoFile[%02d]: %q (%s)\n  Services:", idxProto, *protoFile.Name, *protoFile.Package))
		for idxSvc, svc := range protoFile.Service {
			hg.AddService(NewService(svc))
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

				out := processMethod(method)
				file.P(fmt.Sprintf("\n        ----------------------------------------\n        %02d %s", idxMethod, out))

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

	file.P("\n\nGenerator structs:")
	file.P(hg.String())

	return hg, nil
}
