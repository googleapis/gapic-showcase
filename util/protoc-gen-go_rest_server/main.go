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

package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"

	"github.com/googleapis/gapic-showcase/util/protoc-gen-go_rest_server/server"
)

// Adapted from protoc-gen-go_gapic
func main() {
	reqBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var genReq plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(reqBytes, &genReq); err != nil {
		// log.Fatalf("%s\nError: %s", outBytes, err)
		log.Fatal(err)
	}

	genResp, err := server.Gen(&genReq)
	if err != nil {
		genResp.Error = proto.String(err.Error())
	}

	outBytes, err := proto.Marshal(genResp)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stdout.Write(outBytes); err != nil {
		log.Fatal(err)
	}

	log.Printf("Generated file: %q\n", *genResp.File[0].Name)
}
