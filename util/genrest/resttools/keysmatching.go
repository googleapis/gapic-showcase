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

func matchesSelfOrParent(examine, lookFor string) bool {
	if !strings.HasPrefix(examine, lookFor) {
		return false
	}
	return len(examine) == len(lookFor) || examine[len(lookFor)] == '.'
}
