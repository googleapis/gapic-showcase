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
	"reflect"
	"testing"
)

func TestParseTemplate(t *testing.T) {
	for idx, testCase := range []struct {
		stringTemplate       string
		expectSuccess        bool
		expectParsedTemplate PathTemplate
	}{
		// general segment parsing errors
		{"/aa/", false, nil},
		{"//bb", false, nil},
		{"/@aa/bb", false, nil},
		{"/a@/", false, nil},
		{"/aa/:bb", false, nil},
		{"/aa/:bb:cc", false, nil},
		{"/aa/:bb/cc/dd:ee", false, nil},

		// parse() errors
		{"aa/bb", false, nil},
		{"/aa/bb:@dd", false, nil},
		{"/aa/bb:cc@", false, nil},
		{"/aa/bb:d/d", false, nil},

		// parseVariable() errors
		{"/aa/{}", false, nil},
		{"/aa/{bb=}", false, nil},
		{"/aa/{bb=/cc}", false, nil},
		{"/aa/{bb=cc/}", false, nil},
		{"/aa/{bb=cc/@}", false, nil},
		{"/aa/{@bb=cc}", false, nil},

		// successful cases
		{
			stringTemplate: "/aa/{bb}/cc/{dd=ee/*/gg/{hh=ii/jj/*/kk}/**}:ll",
			expectSuccess:  true,
			expectParsedTemplate: PathTemplate{
				&Segment{Literal, "/", nil},
				&Segment{Literal, "aa", nil},
				&Segment{Literal, "/", nil},
				&Segment{Variable, "bb", PathTemplate{
					&Segment{SingleValue, "", nil},
				}},
				&Segment{Literal, "/", nil},
				&Segment{Literal, "cc", nil},
				&Segment{Literal, "/", nil},
				&Segment{Variable, "dd", PathTemplate{
					&Segment{Literal, "ee", nil},
					&Segment{Literal, "/", nil},
					&Segment{SingleValue, "*", nil},
					&Segment{Literal, "/", nil},
					&Segment{Literal, "gg", nil},
					&Segment{Literal, "/", nil},
					&Segment{Variable, "hh", PathTemplate{
						&Segment{Literal, "ii", nil},
						&Segment{Literal, "/", nil},
						&Segment{Literal, "jj", nil},
						&Segment{Literal, "/", nil},
						&Segment{SingleValue, "*", nil},
						&Segment{Literal, "/", nil},
						&Segment{Literal, "kk", nil}},
					},
					&Segment{Literal, "/", nil},
					&Segment{MultipleValue, "**", nil}}},
				&Segment{Literal, ":", nil},
				&Segment{Literal, "ll", nil}},
		},
		{
			stringTemplate: "/aa/{bb}/cc/{dd=ee/*/gg/{hh=ii/jj/*/kk}/**}",
			expectSuccess:  true,
			expectParsedTemplate: PathTemplate{
				&Segment{Literal, "/", nil},
				&Segment{Literal, "aa", nil},
				&Segment{Literal, "/", nil},
				&Segment{Variable, "bb", PathTemplate{
					&Segment{SingleValue, "", nil}},
				},
				&Segment{Literal, "/", nil},
				&Segment{Literal, "cc", nil},
				&Segment{Literal, "/", nil},
				&Segment{Variable, "dd", PathTemplate{
					&Segment{Literal, "ee", nil},
					&Segment{Literal, "/", nil},
					&Segment{SingleValue, "*", nil},
					&Segment{Literal, "/", nil},
					&Segment{Literal, "gg", nil},
					&Segment{Literal, "/", nil},
					&Segment{Variable, "hh", PathTemplate{
						&Segment{Literal, "ii", nil},
						&Segment{Literal, "/", nil},
						&Segment{Literal, "jj", nil},
						&Segment{Literal, "/", nil},
						&Segment{SingleValue, "*", nil},
						&Segment{Literal, "/", nil},
						&Segment{Literal, "kk", nil}},
					},
					&Segment{Literal, "/", nil},
					&Segment{MultipleValue, "**", nil},
				}},
			},
		},
	} {
		parsed, err := ParseTemplate(testCase.stringTemplate)
		if got, want := err == nil, testCase.expectSuccess; got != want {
			t.Errorf("testCase = %d: ParseTemplate failed: want success: %v;   got error: %s   \n   Test case input: %v", idx, want, err, testCase)
		}
		if !testCase.expectSuccess {
			continue

		}
		if got, want := parsed, testCase.expectParsedTemplate; !reflect.DeepEqual(got, want) {
			t.Errorf("parsed template incorrect:\n    got: %s\n   want: %s", got.asGoLiteral(), want.asGoLiteral())
		}
		_ = parsed
	}
}
