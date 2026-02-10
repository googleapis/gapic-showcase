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

package resttools

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessQueryString(t *testing.T) {
	for idx, testCase := range []struct {
		queryString string
		wantSystem  *SystemParameters
		wantParams  map[string][]string
		wantError   bool
	}{
		{queryString: ""},
		{
			queryString: "foo=bar",
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "$foo=bar",
			wantParams: map[string][]string{
				"$foo": {"bar"},
			},
		},
		{
			queryString: "%24foo=bar",
			wantParams: map[string][]string{
				"$foo": {"bar"},
			},
		},
		{
			queryString: "foo%3Dbar",
			wantParams: map[string][]string{
				"foo=bar": {""},
			},
		},
		{
			queryString: "foo%3Dbar=xyz",
			wantParams: map[string][]string{
				"foo=bar": {"xyz"},
			},
		},
		{
			queryString: "%24foo%3Dbar",
			wantParams: map[string][]string{
				"$foo=bar": {""},
			},
		},
		{
			queryString: "%24foo%3Dbar=xyz",
			wantParams: map[string][]string{
				"$foo=bar": {"xyz"},
			},
		},

		// system param incorrect
		{
			queryString: "%24alt%3Djson",
			wantParams: map[string][]string{
				"$alt=json": {""},
			},
		},
		{
			queryString: "$ALT=JSON",
			wantParams: map[string][]string{
				"$ALT": {"JSON"},
			},
		},
		{
			queryString: "%24ALT=JSON",
			wantParams: map[string][]string{
				"$ALT": {"JSON"},
			},
		},

		// system param by itself
		{queryString: "alt=json"},
		{queryString: "$alt=json"},
		{queryString: "%24alt=json"},
		{
			queryString: "alt=json%3Benum-encoding=int",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
		},
		{
			queryString: "$alt=json%3Benum-encoding=int",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
		},
		{
			queryString: "%24alt=json%3Benum-encoding=int",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
		},
		{
			queryString: "alt=json%3Benum-encoding%3Dint",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
		},
		{
			queryString: "$alt=json%3Benum-encoding%3Dint",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
		},
		{
			queryString: "%24alt=json%3Benum-encoding%3Dint",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
		},

		// system param+query params in front
		{
			queryString: "foo=bar&alt=json",
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "foo=bar&$alt=json",
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "foo=bar&%24alt=json",
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "foo=bar&alt=json%3Benum-encoding=int",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "foo=bar&$alt=json%3Benum-encoding=int",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "foo=bar&%24alt=json%3Benum-encoding=int",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},

		// system param+query params in rear
		{
			queryString: "alt=json&foo=bar",
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "$alt=json&foo=bar",
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "%24alt=json&foo=bar",
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "alt=json%3Benum-encoding=int&foo=bar",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "$alt=json%3Benum-encoding=int&foo=bar",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "%24alt=json%3Benum-encoding=int&foo=bar",
			wantSystem:  &SystemParameters{EnumEncodingAsInt: true},
			wantParams: map[string][]string{
				"foo": {"bar"},
			},
		},
		{
			queryString: "$apiVersion=v7_20260210",
			wantSystem:  &SystemParameters{APIVersion: "v7_20260210"},
			wantParams:  map[string][]string{},
		},
		{
			queryString: "alt=json%3Benum-encoding=int&apiVersion=v7_20260210",
			wantSystem: &SystemParameters{
				APIVersion:        "v7_20260210",
				EnumEncodingAsInt: true,
			},
			wantParams: map[string][]string{},
		},
		{
			queryString: "alt=json%3Benum-encoding=int&apiVersion=v7_20260210&foo=bar",
			wantSystem: &SystemParameters{
				APIVersion:        "v7_20260210",
				EnumEncodingAsInt: true,
			},
			wantParams: map[string][]string{
				"foo": []string{"bar"},
			},
		},

		// system param errors
		{
			queryString: "$alt=foo",
			wantError:   true,
		},
		{
			queryString: "$alt",
			wantError:   true,
		},
		{
			queryString: "$alt=JSON",
			wantError:   true,
		},
		{
			queryString: "$alt=json%3Benum-encoding=string",
			wantError:   true,
		},
		{
			queryString: "$alt=json;enum-encoding=int", // unencoded semicolon
			wantError:   true,
		},
		{
			queryString: "$alt=json%3Benum-encoding=INT",
			wantError:   true,
		},
		{
			queryString: "foo&$alt=json&bar&alt=json", // repeated
			wantError:   true,
		},
		{
			queryString: "$apiVersion=v8&apiVersion=v9", // repeated
			wantError:   true,
		},
	} {
		t.Run(fmt.Sprintf("[%2d %q]", idx, testCase.queryString), func(t *testing.T) {
			systemParams, queryParams, err := processQueryString(testCase.queryString)

			if got, want := (err != nil), testCase.wantError; got != want {
				t.Errorf("error condition not met: want error: %v, got error:%v", testCase.wantError, err)
			}
			if err != nil {
				return
			}
			wantParams := testCase.wantParams
			if wantParams == nil {
				wantParams = map[string][]string{}
			}
			if diff := cmp.Diff(queryParams, wantParams); diff != "" {
				t.Errorf("query params mismatch (+want, -got):\n%s", diff)
			}
			wantSystem := testCase.wantSystem
			if wantSystem == nil {
				wantSystem = &SystemParameters{}
			}
			if diff := cmp.Diff(systemParams, wantSystem); diff != "" {
				t.Errorf("system params mismatch (+want, -got):\n%s", diff)
			}
		})
	}
}
