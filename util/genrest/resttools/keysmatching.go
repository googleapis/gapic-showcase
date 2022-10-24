// Copyright 2021 Google LLC
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

package resttools

import "strings"

// KeysMatchPath returns the keys in `examine` that match any of the elements in `lookFor`,
// interpreting the latter as dotted-field paths. This means a match occurs when (a) the two are
// identical, or (b) when any element of `lookFor` is a prefix of the `examine` key and is followed
// by a period. For example:
// KeysMatchPath(map[string][]string{"loc": nil, "loc.lat": nil, "location":nil},
//
//	         []string{"loc"})
//	== []string{"loc","loc.lat"}
func KeysMatchPath(examine map[string][]string, lookFor []string) []string {
	matchingKeys := []string{}
	for key, _ := range examine {
		for _, target := range lookFor {
			if matchesSelfOrParent(key, target) {
				matchingKeys = append(matchingKeys, key)
				break
			}
		}
	}
	return matchingKeys
}

// matchesSelfOrParent returns whether any element of `lookFor` is or contains a full or partial
// path (in the dotted-field sense) to `examine`. In other words, this returns true when (a) the two
// are identical, or (b) when `examine` starts with `lookFor` and is followed by a period.
func matchesSelfOrParent(examine, lookFor string) bool {
	if !strings.HasPrefix(examine, lookFor) {
		return false
	}
	return len(examine) == len(lookFor) || examine[len(lookFor)] == '.'
}
