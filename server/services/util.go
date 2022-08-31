// Copyright 2022 Google LLC
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

package services

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

func strContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// applyFieldMask applies the values from the src message to the values of the
// dst message according to the contents of the geiven field mask paths.
// If paths is empty/nil, or contains *, it is considered a full update.
//
// TODO: Does not support nested message paths. Currently only used with flat
// resource messages.
func applyFieldMask(src, dst protoreflect.Message, paths []string) {
	fullUpdate := len(paths) == 0 || strContains(paths, "*")

	dst.Range(func(f protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fullUpdate || strContains(paths, string(f.Name())) {
			dst.Set(f, src.Get(f))
		}
		return true
	})
}
