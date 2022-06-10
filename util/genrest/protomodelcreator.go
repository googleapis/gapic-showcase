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
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
	"github.com/googleapis/gapic-showcase/util/genrest/protomodel"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

////////////////////////////////////////
// ProtoModel

// NewProtoModel uses the information in `plugin` to create a new protomodel.Model.
func NewProtoModel(plugin *protogen.Plugin) (*protomodel.Model, error) {
	serviceConfig, err := GetServiceConfig(plugin)
	if err != nil {
		return nil, err
	}
	bindingOverrides := GetBindingOverrides(serviceConfig)

	protoFiles := plugin.Request.GetProtoFile()
	protoModel := &protomodel.Model{
		ProtoInfo: pbinfo.Of(protoFiles),
		Services:  make([]*protomodel.Service, 0, len(protoFiles)),
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

				fqn := fmt.Sprintf("%s.%s.%s", protoFile.GetPackage(), svc.GetName(), method.GetName())
				var http *annotations.HttpRule
				if rule, found := bindingOverrides[fqn]; found && rule != nil {
					http = rule
				} else {
					eHTTP /*, err*/ := proto.GetExtension(method.GetOptions(), annotations.E_Http)
					http = eHTTP.(*annotations.HttpRule)
				}

				rules := []*annotations.HttpRule{http}
				rules = append(rules, http.GetAdditionalBindings()...)
				for idxRule, oneRule := range rules {
					protoModel.AccumulateError(NewServiceBinding(serviceModel, method, oneRule, idxRule))
				}

			}
		}
	}

	return protoModel, protoModel.Error()
}

////////////////////////////////////////
// Service

// NewService creates a protomodel.Service from the given descriptor.
func NewService(protoPackage string, descriptor *descriptorpb.ServiceDescriptorProto) *protomodel.Service {
	// might need *descriptor.ProtoReflect().Type() to instantiate Go type and then reflect to get the name?
	service := &protomodel.Service{
		Descriptor:   descriptor,
		Name:         *descriptor.Name,
		TypeName:     pbinfo.FullyQualifiedType(protoPackage, descriptor.GetName()),
		RESTBindings: make([]*protomodel.RESTBinding, 0),
	}
	return service
}

// NewServiceBinding adds `rule` (as the `index`th binding for `method`) to the specified `service`.
func NewServiceBinding(service *protomodel.Service, method *descriptor.MethodDescriptorProto, rule *annotations.HttpRule, index int) error {
	binding, err := NewRESTBinding(pbinfo.FullyQualifiedType(service.TypeName, method.GetName()), rule, index)
	if err != nil {
		return fmt.Errorf("service %q: %s", service.Name, err)
	}
	service.AddBinding(binding)
	return nil
}

////////////////////////////////////////
// RESTBinding

// NewRESTBinding creates a new protomodel.RESTBinding using the given methodName, rule, and index.
func NewRESTBinding(methodName string, rule *annotations.HttpRule, index int) (*protomodel.RESTBinding, error) {
	restPattern, err := NewRESTRequestPattern(rule)
	if err != nil {
		return nil, fmt.Errorf("method %q, binding %d: %s", methodName, index, err)
	}
	binding := &protomodel.RESTBinding{
		Index:       index,
		ProtoMethod: methodName,
		RESTPattern: restPattern,
		BodyField:   rule.GetBody(),
	}
	return binding, nil
}

////////////////////////////////////////
// RESTRequestPattern

// NewRESTRequestPattern creates a new protomodel.RESTRequestPattern by analyzing the rule provided.
func NewRESTRequestPattern(rule *annotations.HttpRule) (*protomodel.RESTRequestPattern, error) {
	binding := &protomodel.RESTRequestPattern{}
	pattern := rule.GetPattern()
	switch pattern.(type) {
	case *annotations.HttpRule_Get:
		binding.HTTPMethod = "GET"
		binding.Pattern = rule.GetGet()
	case *annotations.HttpRule_Post:
		binding.HTTPMethod = "POST"
		binding.Pattern = rule.GetPost()
	case *annotations.HttpRule_Patch:
		binding.HTTPMethod = "PATCH"
		binding.Pattern = rule.GetPatch()
	case *annotations.HttpRule_Put:
		binding.HTTPMethod = "PUT"
		binding.Pattern = rule.GetPut()
	case *annotations.HttpRule_Delete:
		binding.HTTPMethod = "DELETE"
		binding.Pattern = rule.GetDelete()
	default:
		return nil, fmt.Errorf("unhandled pattern: %#x", pattern)
	}
	return binding, nil
}

////////////////////////////////////////
// Mixins

// GetServiceConfig reads and returns the specified service config file.
func GetServiceConfig(plugin *protogen.Plugin) (*serviceconfig.Service, error) {
	// FIXME: Get this from plugin options
	serviceConfigPath := "schema/google/showcase/v1beta1/showcase_v1beta1.yaml"
	_ = plugin

	y, err := ioutil.ReadFile(serviceConfigPath)
	if err != nil {
		cwd, _ := os.Getwd()
		return nil, fmt.Errorf("error reading service config %q (cwd==%q): %v", serviceConfigPath, cwd, err)
	}

	j, err := yaml.YAMLToJSON(y)
	if err != nil {
		return nil, fmt.Errorf("error converting YAML to JSON: %v", err)
	}

	serviceConfig := &serviceconfig.Service{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(j, serviceConfig); err != nil {
		return nil, fmt.Errorf("error unmarshaling service config: %v", err)
	}

	// An API Service Config will always have a `name` so if it is not populated,
	// it's an invalid config.
	if serviceConfig.GetName() == "" {
		return nil, fmt.Errorf("invalid API service config file %q", serviceConfigPath)
	}
	return serviceConfig, nil
}

// GetBindingOverrides obtains and returns from the service config a map of fully qualified method names to their HTTP bindings.
func GetBindingOverrides(serviceConfig *serviceconfig.Service) (overrides map[string]*annotations.HttpRule) {
	overrides = map[string]*annotations.HttpRule{}
	for _, rule := range serviceConfig.GetHttp().GetRules() {
		overrides[rule.GetSelector()] = rule
	}
	return overrides
}
