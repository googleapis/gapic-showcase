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
	"log"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// TODO(vchudnov-g): Continue filling this in. It's a an initial empty
// stub at the moment.
func Generate(plugin *protogen.Plugin) error {
	log.Printf("Generating REST!")

	file := plugin.NewGeneratedFile("showcase-rest-sample-response.txt", "github.com/googleapis/gapic-showcase/server/genrest")

	// The typecasting below appears to be idiomatic as per
	// https://github.com/protocolbuffers/protobuf-go/blob/master/cmd/protoc-gen-go/internal_gengo/main.go#L31
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	file.P("Generated via \"google.golang.org/protobuf/compiler/protogen\" via ProtoModel!")
	file.P("Files:\n", strings.Join(plugin.Request.GetFileToGenerate(), "\n"))

	protoModel, err := NewProtoModel(plugin)
	if err != nil {
		return err
	}

	file.P("\nProto Model:")
	file.P(protoModel.String())

	file.P("\n\n")
	goModel, err := NewGoModel(protoModel)
	if err != nil {
		return err
	}
	file.P(goModel.String())

	return nil
}
