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

func TestFindValuesMatching(t *testing.T) {
	for idx, testCase := range []struct {
		pattern1, pattern2, longestMatch string
		fullMatch                        bool
	}{
		// Literal, Literal

		{"/zz/yy", "/zz/yy", "/zz/yy", true},
		{"/zz", "/zz/yy", "/zz", false},

		{"/zz/yy:xx", "/zz/yy:xx", "/zz/yy:xx", true},
		{"/zz:xx", "/zz/yy:xx", "/zz", false},

		{"/zz/yy:xx", "/zz/yy", "/zz/yy", false},
		{"/zz", "/zz/yy:xx", "/zz", false},

		// Literal, SingleValue

		{"/zz/yy/{ww}/vv", "/zz/yy/xx/vv", "/zz/yy/xx/vv", true},
		{"/zz/yy/xx/vv", "/zz/yy/{ww}/vv", "/zz/yy/xx/vv", true},

		{"/zz/yy/{ww}/vv:uu", "/zz/yy/xx/vv:uu", "/zz/yy/xx/vv:uu", true},
		{"/zz/yy/{ww}:uu", "/zz/yy/xx:uu", "/zz/yy/xx:uu", true},

		{"/zz/yy/{xx=ww/*}/vv:uu", "/zz/yy/ww/tt/vv:uu", "/zz/yy/ww/tt/vv:uu", true},
		{"/zz/yy/{xx=ww/*}:uu", "/zz/yy/ww/tt:uu", "/zz/yy/ww/tt:uu", true},

		// Literal, MultipleValue
		{"/zz/yy/xx/ww/vv", "/zz/{pp=yy/xx/**}", "/zz/yy/xx/ww/vv", true},
		{"/zz/{pp=yy/xx/**}", "/zz/yy/xx/ww/vv", "/zz/yy/xx/ww/vv", true},

		{"/zz/{pp=yy/xx/**}:rr", "/zz/yy/xx/ww/vv:rr", "/zz/yy/xx/ww/vv:rr", true},
		{"/zz/{pp=yy/xx/**}:rr", "/zz/yy/xx/ww/vv", "/zz/yy/xx/ww/vv", false},
		{"/zz/{pp=yy/xx/**}", "/zz/yy/xx/ww/vv:rr", "/zz/yy/xx/ww/vv", false},
		{"/zz/{pp=yy/xx/**}:rr", "/zz/yy/xx/ww/vv:ss", "/zz/yy/xx/ww/vv:", false},

		// SingleValue, SingleValue

		{"/zz/yy/{xx=ww/*}/vv", "/zz/yy/{tt=ww/*}/vv", "/zz/yy/ww/*/vv", true},
		{"/zz/yy/{xx=ww/*}:uu", "/zz/yy/{tt=ww/*}:uu", "/zz/yy/ww/*:uu", true},

		{"/zz/yy/{xx=ww/*}/vv/{ss=rr/*/pp}", "/zz/yy/ww/{xx}/vv/{oo=rr/*/pp}", "/zz/yy/ww/*/vv/rr/*/pp", true},

		// SingleValue, MultipleValue

		{"/zz/yy/{xx=ww/**}", "/zz/yy/ww/{vv}/uu", "/zz/yy/ww/*/uu", true},
		{"/zz/yy/{xx=ww/**}:tt", "/zz/yy/ww/{vv}/uu:tt", "/zz/yy/ww/*/uu:tt", true},
		{"/zz/yy/{xx=ww/**}", "/zz/yy/{vv=ww/*}/uu", "/zz/yy/ww/*/uu", true},
		{"/zz/yy/{xx=ww/**}:tt", "/zz/yy/{vv=ww/*}/uu:tt", "/zz/yy/ww/*/uu:tt", true},
		{"/zz/yy/{xx=ww/**}:tt", "/zz/yy/ww/{vv}/uu", "/zz/yy/ww/*/uu", false},
		{"/zz/yy/{xx=ww/**}", "/zz/yy/ww/{vv}/uu:tt", "/zz/yy/ww/*/uu", false},

		// MultipleValue, MultipleValue
		{"/zz/yy/{xx=ww/**}", "/zz/yy/{vv=ww/**}", "/zz/yy/ww/**", true},
		{"/zz/yy/{xx=ww/**}:ss", "/zz/yy/{vv=ww/**}:ss", "/zz/yy/ww/**:ss", true},
		{"/zz/yy/{xx=ww/**}", "/zz/yy/{vv=ww/**}:ss", "/zz/yy/ww/**", false},
		{"/zz/yy/{xx=ww/**}:ss", "/zz/yy/{vv=ww/**}", "/zz/yy/ww/**", false},

		// Mix
	} {
		template1, err := NewPathTemplate(testCase.pattern1)
		if err != nil {
			t.Errorf("testCase %2d: unexpected error constructing template1: %s:\n   Test case input: %v", idx, err, testCase)
		}

		template2, err := NewPathTemplate(testCase.pattern2)
		if err != nil {
			t.Errorf("testCase %2d: unexpected error constructing template2: %s:\n   Test case input: %v", idx, err, testCase)
		}

		fullMatch, longestMatch, err := FindValuesMatching(template1, template2)
		if got, want := longestMatch, testCase.longestMatch; got != want {
			t.Errorf("testCase %2d: longestMatch failed: got %q   want %q:\n   Test case input: %v\n   err = %v", idx, got, want, testCase, err)
		}
		if got, want := fullMatch, testCase.fullMatch; got != want {
			t.Errorf("testCase %2d: fullMatch failed: got %v   want %v:\n   Test case input: %v\n   err = %v", idx, got, want, testCase, err)
		}
		_ = err

	}
}
