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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

////////////////////////////////////////
// ProtoModel

func NewProtoModel(plugin *protogen.Plugin) (*ProtoModel, error) {
	protoFiles := plugin.Request.GetProtoFile()
	protoModel := &ProtoModel{
		descInfo: pbinfo.Of(protoFiles),
		services: make([]*Service, 0, len(protoFiles)),
	}

	generateFile := map[string]bool{}
	for _, fileName := range plugin.Request.GetFileToGenerate() {
		generateFile[fileName] = true
	}

	for _, protoFile := range protoFiles {
		if !generateFile[*protoFile.Name] {
			continue
		}
		protoPackage := *protoFile.Package
		for _, svc := range protoFile.GetService() {
			serviceModel := protoModel.AddService(NewService(protoPackage, svc))
			for _, method := range svc.GetMethod() {
				options := method.GetOptions()
				if options == nil {
					continue
				}

				eHTTP /*, err*/ := proto.GetExtension(method.GetOptions(), annotations.E_Http)
				http := eHTTP.(*annotations.HttpRule)
				rules := []*annotations.HttpRule{http}
				rules = append(rules, http.GetAdditionalBindings()...)
				for idxRule, oneRule := range rules {
					protoModel.AccumulateError(NewServiceBinding(serviceModel, method, oneRule, idxRule))
				}

			}
		}
	}

	protoModel.CheckConsistency()
	return protoModel, protoModel.Error()
}

////////////////////////////////////////
// Service

func NewServiceBinding(service *Service, method *descriptor.MethodDescriptorProto, rule *annotations.HttpRule, index int) error {
	binding, err := NewRESTBinding(fullyQualifiedType(service.typeName, method.GetName()), rule, index)
	if err != nil {
		return fmt.Errorf("service %q: %s", service.name, err)
	}
	service.AddBinding(binding)
	return nil
}

func NewService(protoPackage string, descriptor *descriptorpb.ServiceDescriptorProto) *Service {
	// might need *descriptor.ProtoReflect().Type() to instantiate Go type and then reflect to get the name?
	service := &Service{
		descriptor:   descriptor,
		name:         *descriptor.Name,
		typeName:     fullyQualifiedType(protoPackage, descriptor.GetName()),
		restBindings: make([]*RESTBinding, 0),
	}
	return service
}

////////////////////////////////////////
// RESTBinding

func NewRESTBinding(methodName string, rule *annotations.HttpRule, index int) (*RESTBinding, error) {
	restPattern, err := NewRESTRequestPattern(rule)
	if err != nil {
		return nil, fmt.Errorf("method %q, binding %d: %s", methodName, index, err)
	}
	binding := &RESTBinding{
		index:       index,
		protoMethod: methodName,
		restPattern: restPattern,
	}
	return binding, nil
}

////////////////////////////////////////
// RESTRequestPattern

func NewRESTRequestPattern(rule *annotations.HttpRule) (*RESTRequestPattern, error) {
	binding := &RESTRequestPattern{}
	pattern := rule.GetPattern()
	switch pattern.(type) {
	case *annotations.HttpRule_Get:
		binding.httpMethod = "GET"
		binding.pattern = rule.GetGet()
	case *annotations.HttpRule_Post:
		binding.httpMethod = "POST"
		binding.pattern = rule.GetPost()
	case *annotations.HttpRule_Patch:
		binding.httpMethod = "PATCH"
		binding.pattern = rule.GetPatch()
	case *annotations.HttpRule_Put:
		binding.httpMethod = "PUT"
		binding.pattern = rule.GetPut()
	case *annotations.HttpRule_Delete:
		binding.httpMethod = "DELETE"
		binding.pattern = rule.GetDelete()
	default:
		return nil, fmt.Errorf("unhandled pattern: %#x", pattern)
	}
	return binding, nil
}
