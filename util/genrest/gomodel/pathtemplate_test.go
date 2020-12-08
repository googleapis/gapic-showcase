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

package gomodel

import (
	"testing"
)

func TestHasVariables(t *testing.T) {
	for idx, testCase := range []struct {
		stringTemplate   string
		expectVars       bool
		expectNestedVars bool
	}{
		{
			stringTemplate:   "/aa/cc/ee/*/gg/ii/jj/*/kk/**:ll",
			expectVars:       false,
			expectNestedVars: false,
		},
		{
			stringTemplate:   "/aa/{bb}/cc/{dd=ee/*/gg}/{hh=ii/jj/*/kk/**}:ll",
			expectVars:       true,
			expectNestedVars: false,
		},
		{
			stringTemplate:   "/aa/{bb}/cc/{dd=ee/*/gg/{hh=ii/jj/*/kk}/**}:ll",
			expectVars:       true,
			expectNestedVars: true,
		},
	} {
		parsed, err := ParseTemplate(testCase.stringTemplate)
		if err != nil {
			t.Errorf("testCase = %d: ParseTemplate failed: %s \n   Test case input: %v", idx, err, testCase)
		}

		hasVars, hasNestedVars := parsed.HasVariables()
		if got, want := hasVars, testCase.expectVars; got != want {
			t.Errorf("testCase = %d: HasVars() failed checking variables: got %v, want %v", idx, got, want)
		}
		if got, want := hasNestedVars, testCase.expectNestedVars; got != want {
			t.Errorf("testCase = %d: HasVars() failed checking nested variables: got %v, want %v", idx, got, want)
		}
	}
}
