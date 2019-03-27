// Copyright 2018 Google LLC
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
	"github.com/googleapis/gapic-showcase/util"
)

// This script regenerates all of the generated source code for the Showcase
// API including the generated messages, gRPC services, go gapic clients,
// and the generated CLI. This script must be ran from the root directory
// of the gapic-showcase repository.
//
// This script should be used whenever any changes are made to any of
// the protos found in schema.
//
// Usage: go run ./util/cmd/compile_protos/main.go
func main() {
	util.CompileProtos("v1alpha3")
}
