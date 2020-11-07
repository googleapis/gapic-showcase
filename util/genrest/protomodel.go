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
	"google.golang.org/protobuf/types/descriptorpb"
)

////////////////////////////////////////
// ProtoModel

type ProtoModel struct {
	ErrorAccumulator
	descInfo pbinfo.Info
	services []*Service
}

func (model *ProtoModel) String() string {
	services := make([]string, len(model.services))
	for idx, svc := range model.services {
		if svc == nil {
			continue
		}
		services[idx] = svc.String()
	}
	return strings.Join(services, "\n\n")
}

func (model *ProtoModel) AddService(service *Service) *Service {
	model.services = append(model.services, service)
	return service
}

func (model *ProtoModel) CheckConsistency() {
	model.CheckBindingsUnique()
}

func (model *ProtoModel) CheckBindingsUnique() {
	allBindings := []*RESTBinding{}
	for _, service := range model.services {
		allBindings = append(allBindings, service.restBindings...)
	}

	for first, firstBinding := range allBindings {
		for _, secondBinding := range allBindings[first:] {
			ambiguousPattern := firstBinding.FindAmbiguityWith(secondBinding)
			if len(ambiguousPattern) == 0 {
				continue
			}
			model.AccumulateError(fmt.Errorf("Pattern %q matches both (%s) and (%s)", ambiguousPattern, firstBinding, secondBinding))
		}
	}
}

////////////////////////////////////////
// Service

type Service struct {
	descriptor   *descriptorpb.ServiceDescriptorProto // maybe not needed
	name         string
	typeName     string
	restBindings []*RESTBinding
}

func (service *Service) String() string {
	handlers := make([]string, len(service.restBindings))
	for idx, h := range service.restBindings {
		handlers[idx] = h.String()
	}
	indent := "  "
	return fmt.Sprintf("%s (%s):\n%s%s", service.name, service.typeName, indent, strings.Join(handlers, "\n"+indent))
}

func (service *Service) AddBinding(binding *RESTBinding) {
	service.restBindings = append(service.restBindings, binding)
}

func fullyQualifiedType(segments ...string) string {
	// In descriptors, putting the dot in front means the name is fully-qualified.
	const dot = "."
	typeName := strings.Join(segments, dot)
	if !strings.HasPrefix(typeName, dot) {
		typeName = dot + typeName
	}
	return typeName
}

////////////////////////////////////////
// RESTBinding

type RESTBinding struct {
	index       int
	protoMethod string
	restPattern *RESTRequestPattern
}

func (binding *RESTBinding) String() string {
	return fmt.Sprintf("%s[%d] : %s", binding.protoMethod, binding.index, binding.restPattern)
}

func (binding *RESTBinding) FindAmbiguityWith(other *RESTBinding) string {
	// TODO: Fill this in. It's a stub for now.
	return ""
}

////////////////////////////////////////
// RESTRequestPattern

type RESTRequestPattern struct {
	httpMethod string // make an enum?
	pattern    string
	// TODO: Add body info
}

func (binding *RESTRequestPattern) String() string {
	return fmt.Sprintf("%s: %q", binding.httpMethod, binding.pattern)
}
