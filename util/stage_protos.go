// Copyright 2019 Google LLC
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

package util

import (
	"log"
	"os"
	"path/filepath"
)

func StageProtos(version, showcaseProtoDir, outDir string) {
	if err := os.RemoveAll(outDir); err != nil {
		log.Fatalf("Failed to remove the directory %s: %v", outDir, err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Failed to make the directory %s: %v", outDir, err)
	}

	// Get proto dependencies
	Execute(
		"git",
		"clone",
		"https://github.com/googleapis/api-common-protos.git",
		outDir,
	)

	// Move showcase protos alongside its dependencies.
	protoDest := filepath.Join(
		outDir,
		"google",
		"showcase",
		version)
	if err := os.MkdirAll(protoDest, 0755); err != nil {
		log.Fatalf("Failed to make the dir %s: %v", protoDest, err)
	}

	files, err := filepath.Glob(filepath.Join(showcaseProtoDir, "*.proto"))
	if err != nil {
		log.Fatal("Error: failed to find protos in: " + showcaseProtoDir)
	}

	for _, f := range files {
		Execute("cp", f, protoDest)
	}
}
