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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-showcase/util/genrest/gomodel"
	"github.com/googleapis/gapic-showcase/util/genrest/internal/pbinfo"
	"github.com/googleapis/gapic-showcase/util/genrest/protomodel"
	"google.golang.org/protobuf/types/descriptorpb"
)

// NewGoModel creates a new goModel.Model from the given protomodel.Model. It essentially extracts
// and organizes the data needed to later generate Go source files.
func NewGoModel(protoModel *protomodel.Model) (*gomodel.Model, error) {
	goModel := &gomodel.Model{
		Service: make([]*gomodel.ServiceModel, 0, len(protoModel.Services)),
	}

	protoInfo := protoModel.ProtoInfo

	for _, service := range protoModel.Services {
		serviceModel := &gomodel.ServiceModel{ProtoPath: service.TypeName, ShortName: service.Name}
		goModel.Add(serviceModel)
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

			var (
				requestBodyFieldType string
				requestBodyFieldName string
				bodyFieldImports     pbinfo.ImportSpec
				bodyFieldSpec        gomodel.BodyFieldSpec
			)

			if binding.BodyField == "*" {
				bodyFieldSpec = gomodel.BodyFieldAll
			} else if len(binding.BodyField) > 0 {
				bodyFieldSpec = gomodel.BodyFieldSingle
				var bodyFieldDesc *descriptorpb.FieldDescriptorProto
				inProtoTypeDescriptor, ok := inProtoType.(*descriptor.DescriptorProto)
				if !ok {
					goModel.AccumulateError(fmt.Errorf("could not type assert inProtoType %v to *descriptor.DescriptorProto", inProtoType))
				}

				// Intentional: the following indirectly enforces that the field
				// specified is top-level (is not a dotted path), as periods are not
				// allowed in field names.
				for _, fd := range inProtoTypeDescriptor.GetField() {
					if fd.GetName() == binding.BodyField {
						bodyFieldDesc = fd
						break
					}
				}
				if bodyFieldDesc == nil {
					goModel.AccumulateError(fmt.Errorf("could not find body field %q in %q", binding.BodyField, inProtoType.GetName()))
				}
				bodyFieldTypeDesc, ok := protoInfo.Type[*bodyFieldDesc.TypeName]
				if !ok {
					goModel.AccumulateError(fmt.Errorf("could not read protoInfo[%q]", inProtoType))
				}
				requestBodyFieldType, bodyFieldImports, err = protoInfo.NameSpec(bodyFieldTypeDesc)
				// TODO: test for HTTP body encoding a single field whose names is different than its type
				// TODO: Test for HTTP body encoding a single field that is a scalar, not a message
				requestBodyFieldName = strings.Title(bodyFieldDesc.GetName())
				goModel.AccumulateError(err)
			}

			outProtoType := protoInfo.Type[*protoMethodDesc.OutputType]
			outGoType, outImports, err := protoInfo.NameSpec(outProtoType)
			goModel.AccumulateError(err)

			pathTemplate, err := gomodel.NewPathTemplate(binding.RESTPattern.Pattern)
			goModel.AccumulateError(err)

			// TODO: Check that each field path in the handler path template refers to
			// an actual field in the request. We can use the functionality in
			// resttools.PopulateOneField() (after some refactoring) to do this,
			// starting from FieldDescriptor.ProtoReflect(). This will allow us to error
			// at build time rather than at server run time.

			restHandler := &gomodel.RESTHandler{
				HTTPMethod: binding.RESTPattern.HTTPMethod,
				URIPattern: binding.RESTPattern.Pattern,

				PathTemplate:    pathTemplate,
				StreamingServer: protoMethodDesc.GetServerStreaming(),
				StreamingClient: protoMethodDesc.GetClientStreaming(),

				GoMethod:                  protoMethodDesc.GetName(),
				RequestType:               inGoType,
				RequestTypePackage:        inImports.Name,
				RequestVariable:           "request",
				RequestBodyFieldSpec:      bodyFieldSpec,
				RequestBodyFieldProtoName: binding.BodyField,
				RequestBodyFieldName:      requestBodyFieldName,
				RequestBodyFieldType:      requestBodyFieldType,
				RequestBodyFieldVariable:  "bodyField",
				RequestBodyFieldPackage:   bodyFieldImports.Name,

				ResponseType:        outGoType,
				ResponseTypePackage: outImports.Name,
				ResponseVariable:    "response",
			}

			serviceModel.AddImports(&inImports, &outImports)
			serviceModel.AddHandler(restHandler)
		}
	}

	goModel.CheckConsistency()
	return goModel, goModel.Error()
}
