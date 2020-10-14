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

package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func Gen(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	fileName := "showcase-rest-sample-response.txt"

	resp := &plugin.CodeGeneratorResponse{}
	resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
		Name: &fileName,
		Content: proto.String(
			fmt.Sprintf("Generated at %s\nFiles:\n%s",
				time.Now().Format(time.UnixDate),
				strings.Join(genReq.FileToGenerate, "\n"))),
	})
	resp.SupportedFeatures = proto.Uint64(uint64(plugin.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))

	return resp, nil
}
