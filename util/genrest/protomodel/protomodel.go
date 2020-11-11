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

package protomodel

import (
	"fmt"
	"strings"

	"github.com/googleapis/gapic-showcase/util/genrest/errorhandling"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
	"google.golang.org/protobuf/types/descriptorpb"
)

////////////////////////////////////////
// ProtoModel

type Model struct {
	errorhandling.Accumulator
	ProtoInfo pbinfo.Info
	Services  []*Service
}

func (model *Model) String() string {
	services := make([]string, len(model.Services))
	for idx, svc := range model.Services {
		if svc == nil {
			continue
		}
		services[idx] = svc.String()
	}
	return strings.Join(services, "\n\n")
}

func (model *Model) AddService(service *Service) *Service {
	model.Services = append(model.Services, service)
	return service
}

func (model *Model) CheckConsistency() { // move to gomodel
	model.CheckBindingsUnique()
}

func (model *Model) CheckBindingsUnique() { // move to gomodel
	allBindings := []*RESTBinding{}
	for _, service := range model.Services {
		allBindings = append(allBindings, service.RESTBindings...)
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
	Descriptor   *descriptorpb.ServiceDescriptorProto // maybe not needed
	Name         string
	TypeName     string
	RESTBindings []*RESTBinding
}

func (service *Service) String() string {
	handlers := make([]string, len(service.RESTBindings))
	for idx, h := range service.RESTBindings {
		handlers[idx] = h.String()
	}
	indent := "  "
	return fmt.Sprintf("%s (%s):\n%s%s", service.Name, service.TypeName, indent, strings.Join(handlers, "\n"+indent))
}

func (service *Service) AddBinding(binding *RESTBinding) {
	service.RESTBindings = append(service.RESTBindings, binding)
}

////////////////////////////////////////
// RESTBinding

type RESTBinding struct {
	Index       int
	ProtoMethod string
	RESTPattern *RESTRequestPattern
}

func (binding *RESTBinding) String() string {
	return fmt.Sprintf("%s[%d] : %s", binding.ProtoMethod, binding.Index, binding.RESTPattern)
}

func (binding *RESTBinding) FindAmbiguityWith(other *RESTBinding) string {
	// TODO: Fill this in. It's a stub for now.
	return ""
}

////////////////////////////////////////
// RESTRequestPattern

type RESTRequestPattern struct {
	HTTPMethod string // make an enum?
	Pattern    string
	// TODO: Add body info
}

func (binding *RESTRequestPattern) String() string {
	return fmt.Sprintf("%s: %q", binding.HTTPMethod, binding.Pattern)
}
