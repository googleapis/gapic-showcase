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

	"github.com/googleapis/gapic-showcase/util/genrest/protomodel"
	"google.golang.org/protobuf/types/descriptorpb"
)

func NewGoModel(protoModel *protomodel.Model) (*GoModel, error) {
	goModel := &GoModel{
		shim: make([]*GoServiceShim, 0, len(protoModel.Services)),
	}

	protoInfo := protoModel.ProtoInfo

	for _, service := range protoModel.Services {
		shim := &GoServiceShim{path: fmt.Sprintf("%q (%s)", service.Name, service.TypeName)}
		goModel.Add(shim)
		for _, binding := range service.RESTBindings {
			protoMethodType := binding.ProtoMethod
			protoMethodDesc, ok := protoInfo.Type[protoMethodType].(*descriptorpb.MethodDescriptorProto)
			if !ok {
				goModel.AccumulateError(fmt.Errorf("could not get descriptor for %q: %#x", protoMethodType, protoInfo.Type[protoMethodType]))
				continue
			}
			inProtoType := protoInfo.Type[*protoMethodDesc.InputType]
			inGoType, inImports, err := protoInfo.NameSpec(inProtoType)
			goModel.AccumulateError(err)

			outProtoType := protoInfo.Type[*protoMethodDesc.OutputType]
			goModel.AccumulateError(err)
			outGoType, outImports, err := protoInfo.NameSpec(outProtoType)

			pathTemplate, err := NewPathTemplate(binding.RESTPattern.Pattern)
			goModel.AccumulateError(err)

			restHandler := &RESTHandler{
				httpMethod:   binding.RESTPattern.HTTPMethod,
				urlMatcher:   binding.RESTPattern.Pattern,
				pathTemplate: pathTemplate,

				goMethod:           protoMethodDesc.GetName(),
				requestType:        inGoType,
				requestTypePackage: inImports.Name,
				requestVariable:    "request",

				responseType:        outGoType,
				responseTypePackage: outImports.Name,
				responseVariable:    "response",
			}

			shim.AddImports(&inImports, &outImports)
			shim.AddHandler(restHandler)

			_ = outImports

		}
	}
	return goModel, goModel.Error()
}
