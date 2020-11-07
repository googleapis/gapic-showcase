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

	"google.golang.org/protobuf/types/descriptorpb"
)

func NewGoModel(protoModel *ProtoModel) (*GoModel, error) {
	goModel := &GoModel{
		shim: make([]*GoServiceShim, 0, len(protoModel.services)),
	}

	descInfo := protoModel.descInfo

	for _, service := range protoModel.services {
		shim := &GoServiceShim{path: fmt.Sprintf("%q (%s)", service.name, service.typeName)}
		goModel.Add(shim)
		for _, binding := range service.restBindings {
			protoMethodType := binding.protoMethod
			protoMethodDesc, ok := descInfo.Type[protoMethodType].(*descriptorpb.MethodDescriptorProto)
			if !ok {
				goModel.AccumulateError(fmt.Errorf("could not get descriptor for %q: %#x", protoMethodType, descInfo.Type[protoMethodType]))
				continue
			}
			inProtoType := descInfo.Type[*protoMethodDesc.InputType]
			inGoType, inImports, err := descInfo.NameSpec(inProtoType)
			goModel.AccumulateError(err)

			outProtoType := descInfo.Type[*protoMethodDesc.OutputType]
			goModel.AccumulateError(err)
			outGoType, outImports, err := descInfo.NameSpec(outProtoType)

			restHandler := &RESTHandler{
				httpMethod: binding.restPattern.httpMethod,
				urlMatcher: binding.restPattern.pattern,

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
