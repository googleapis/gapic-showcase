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
	"os"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
	"github.com/googleapis/gapic-showcase/util/genrest/protomodel"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/genproto/googleapis/cloud/location"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
)

////////////////////////////////////////
// ProtoModel

// NewProtoModel uses the information in `plugin` to create a new protomodel.Model.
func NewProtoModel(plugin *protogen.Plugin) (*protomodel.Model, error) {
	protoFiles := append(plugin.Request.GetProtoFile(), getMixinFiles()...)
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
				addBindingsForMethod(protoModel, serviceModel, method)
			}
		}
	}

	serviceConfig, err := GetServiceConfig()
	if err != nil {
		return nil, err
	}
	mixins := collectMixins(serviceConfig)
	for _, mixinFile := range mixins {
		protoPackage := *mixinFile.file.Package
		for _, mixinService := range mixinFile.services {
			svc := mixinService.service
			serviceModel := protoModel.AddService(NewService(protoPackage, svc))
			for _, method := range mixinService.methods {
				addBindingsForMethod(protoModel, serviceModel, method)
			}
		}
	}

	return protoModel, protoModel.Error()
}

func addBindingsForMethod(protoModel *protomodel.Model, serviceModel *protomodel.Service, method *descriptor.MethodDescriptorProto) {
	options := method.GetOptions()
	if options == nil {
		return
	}

	eHTTP /*, err*/ := proto.GetExtension(method.GetOptions(), annotations.E_Http)
	http := eHTTP.(*annotations.HttpRule)
	rules := []*annotations.HttpRule{http}
	rules = append(rules, http.GetAdditionalBindings()...)
	for idxRule, oneRule := range rules {
		protoModel.AccumulateError(NewServiceBinding(serviceModel, method, oneRule, idxRule))
	}
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

// Note: much of the code below is a modified copy of the mixin code in gapic-generator-go.

// GetServiceConfig reads and returns the Showcase service config file.
func GetServiceConfig() (*serviceconfig.Service, error) {
	// TODO: Consider getting this from the plugin options. On the
	// other hand, there's only one copy of this file, so maybe
	// hard-coding this location isn't terrible.
	serviceConfigPath := "schema/google/showcase/v1beta1/showcase_v1beta1.yaml"

	y, err := os.ReadFile(serviceConfigPath)
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

// Mixins is the collection of files containing methods to be mixed in.
type Mixins []*MixinFile

// MixinFile describes a single file containins methods to be mixed in.
type MixinFile struct {
	file     *descriptor.FileDescriptorProto
	services []*MixinService
}

// MixinService describes a single service containing methods to be filled in
type MixinService struct {
	service *descriptor.ServiceDescriptorProto
	methods []*descriptor.MethodDescriptorProto
}

// indexedRules keys HTTP rules by their selectors
type indexedRules map[string]*annotations.HttpRule

// collectMixins collects the configured mixin APIs from the service config and
// returns the appropriately configured mixin `MethodDescriptorProto`s to generate for each.
func collectMixins(serviceConfig *serviceconfig.Service) Mixins {
	mixinRules := indexedRules{}
	for _, rule := range serviceConfig.GetHttp().GetRules() {
		mixinRules[rule.GetSelector()] = rule
	}
	mixins := Mixins{}
	for _, api := range serviceConfig.GetApis() {
		mixins = append(mixins, getMixinsForAPI(mixinRules, api.GetName())...)
	}
	return mixins
}

// getMixinsForAPI returns the appropriately configured mixin `MethodDescriptorProto`s
// corresponding to the `mixinRules` that refer to `api`. The method descriptors are taken from the
// package-global `mixinDescriptors`.
func getMixinsForAPI(mixinRules indexedRules, api string) Mixins {
	files := Mixins{}
	for _, file := range mixinDescriptors[api] {
		fileToAdd := &MixinFile{
			file: file,
		}
		files = append(files, fileToAdd)
		for _, service := range file.GetService() {
			serviceToAdd := &MixinService{
				service: service,
			}
			fileToAdd.services = append(fileToAdd.services, serviceToAdd)
			for _, method := range service.GetMethod() {
				fqn := fmt.Sprintf("%s.%s.%s", file.GetPackage(), service.GetName(), method.GetName())

				if rule := mixinRules[fqn]; rule != nil {
					proto.SetExtension(method.Options, annotations.E_Http, rule)
					serviceToAdd.methods = append(serviceToAdd.methods, method)
				}
			}
		}
	}
	return files
}

func getMixinFiles() (files []*descriptor.FileDescriptorProto) {
	for _, descriptors := range mixinDescriptors {
		files = append(files, descriptors...)
	}
	return
}

// mixinDescriptors maps fully qualified proto service names of mixins implemented by Showcase to a
// list of FileDescriptors containing the definitions needed for that service.
var mixinDescriptors map[string][]*descriptor.FileDescriptorProto

func init() {
	mixinDescriptors = map[string][]*descriptor.FileDescriptorProto{
		"google.longrunning.Operations": {
			protodesc.ToFileDescriptorProto(longrunningpb.File_google_longrunning_operations_proto),
		},
		"google.cloud.location.Locations": {
			protodesc.ToFileDescriptorProto(location.File_google_cloud_location_locations_proto),
		},
		"google.iam.v1.IAMPolicy": {
			protodesc.ToFileDescriptorProto(iampb.File_google_iam_v1_iam_policy_proto),
			protodesc.ToFileDescriptorProto(iampb.File_google_iam_v1_policy_proto),
			protodesc.ToFileDescriptorProto(iampb.File_google_iam_v1_options_proto),
		},
	}
}
