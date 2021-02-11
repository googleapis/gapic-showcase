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

import (
	"testing"
)

func TestKeysMatchPath(t *testing.T) {
	for idx, testCase := range []struct {
		examine     map[string][]string
		lookFor     []string
		wantMatches []string
	}{
		{
			examine: map[string][]string{
				"loc":           nil,
				"location":      nil,
				"loc.lat":       nil,
				"extra.loc.lat": nil,
				"location.lat":  nil,
			},
			lookFor:     []string{"loc"},
			wantMatches: []string{"loc", "loc.lat"},
		},
		{
			examine: map[string][]string{
				"loc":           nil,
				"location":      nil,
				"loc.lat":       nil,
				"extra.loc.lat": nil,
				"location.lat":  nil,
			},
			lookFor:     []string{"location", "loc"},
			wantMatches: []string{"loc", "location", "loc.lat", "location.lat"},
		},
	} {
		matches := KeysMatchPath(testCase.examine, testCase.lookFor)
		if got, want := len(matches), len(testCase.wantMatches); got != want {
			t.Errorf("testCase = %d: unexpected number of variables returned: got %v, want %v: returned elements: %v",
				idx, got, want, matches)
			continue
		}
		for matchIdx, got := range matches {
			if want := testCase.wantMatches[matchIdx]; got != want {
				t.Errorf("testCase = %d: variable %d unexpected: got %v, want %v", idx, matchIdx, got, want)
			}
		}

	}
}
