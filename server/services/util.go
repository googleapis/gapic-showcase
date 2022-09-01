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
// dst message according to the contents of the given field mask paths.
// If paths is empty/nil, or contains *, it is considered a full update.
//
// TODO: Does not support nested message paths. Currently only used with flat
// resource messages.
func applyFieldMask(src, dst protoreflect.Message, paths []string) {
	fullUpdate := len(paths) == 0 || strContains(paths, "*")

	fields := dst.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		isOneof := field.ContainingOneof() != nil && !field.ContainingOneof().IsSynthetic()

		// Set field in dst with value from src, skipping true oneofs, while
		// handling proto3_optional fields represented as synthetic oneofs.
		if (fullUpdate || strContains(paths, string(field.Name()))) && !isOneof {
			dst.Set(field, src.Get(field))
		}
	}

	oneofs := dst.Descriptor().Oneofs()
	for i := 0; i < oneofs.Len(); i++ {
		oneof := oneofs.Get(i)
		// Skip proto3_optional synthetic oneofs.
		if oneof.IsSynthetic() {
			continue
		}

		setOneof := src.WhichOneof(oneof)
		if setOneof == nil && fullUpdate {
			// Full update with no field set in this oneof of
			// src means clear all fields for this oneof in dst.
			fields := oneof.Fields()
			for j := 0; j < fields.Len(); j++ {
				dst.Clear(fields.Get(j))
			}
		} else if setOneof != nil && (fullUpdate || strContains(paths, string(setOneof.Name()))) {
			// Full update or targeted updated with a field set in this oneof of
			// src means set that field for the same oneof in dst, which implicitly
			// clears any previously set field for this oneof.
			dst.Set(setOneof, src.Get(setOneof))
		}
	}
}
