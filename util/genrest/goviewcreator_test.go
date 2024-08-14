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

package genrest

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-showcase/util/genrest/gomodel"
)

func TestMatchingPath(t *testing.T) {
	for idx, testCase := range []struct {
		template    string
		expectError bool
		expectMatch string
		expectVars  []string
	}{
		{
			template:    "/aa/{bb}/cc/{dd=ee/*/gg/{hh=ii/jj/*/kk}/**}:ll",
			expectError: true,
		},
		{
			template:    "/aa/{bb}/cc/{dd=ee/*/gg}/{hh=ii/jj/*/kk/**}",
			expectMatch: "/aa/{bb:[^:]+}/cc/{dd:ee/[^:]+/gg}/{hh:ii/jj/[^:]+/kk/[^:]+}",
			expectVars:  []string{"bb", "dd", "hh"},
		},
	} {
		pathTemplate, err := gomodel.ParseTemplate(testCase.template)
		if err != nil {
			t.Errorf("testCase %2d: unexpected error constructing template: %s:\n   Test case input: %v", idx, err, testCase)
			continue
		}

		path, allVars, err := matchingPath(pathTemplate)

		if got, want := (err != nil), testCase.expectError; got != want {
			t.Errorf("testCase %2d: matchingPath error:\n    got: %q\n   want: %v", idx, err, want)
		}
		if err != nil {
			continue
		}
		if got, want := path, testCase.expectMatch; got != want {
			t.Errorf("testCase %2d: matchingPath path:\n    got: %q\n   want: %q", idx, got, want)
		}
		if got, want := allVars, testCase.expectVars; !reflect.DeepEqual(got, want) {
			t.Errorf("testCase %2d: matchingPath path:\n    got: %#v\n   want: %#v", idx, got, want)
		}

	}
}

func TestNamer(t *testing.T) {
	namer := NewNamer()
	for idx, testCase := range []struct {
		requested string
		expected  string
	}{
		// Order matters, since we're testing disambiguation with previously seen items.
		{"rainbow", "rainbow"},
		{"rainbow", "rainbow_1"},
		{"rainbow", "rainbow_2"},
		{"rainbow_1", "rainbow_1_1"},
		{"rainbow_1", "rainbow_1_2"},
		{"rainbow_1", "rainbow_1_3"},
		{"rainbow_2", "rainbow_2_1"},
		{"sun_1", "sun_1"},
		{"sun_1", "sun_1_1"},
	} {
		if got, want := namer.Get(testCase.requested), testCase.expected; got != want {
			t.Errorf("testCase %2d: got %q, want %q", idx, got, want)
		}
	}
}

func TestConstructStreamingServer(t *testing.T) {
	fileImports := map[string]string{}
	helperSources := sourceMap{}

	constructServerStreamer(&gomodel.ServiceModel{ShortName: "Catalog"},
		&gomodel.RESTHandler{
			RequestTypePackage:  "catalogpb",
			ResponseTypePackage: "responsepb",
			GoMethod:            "StreamAuthors",
			ResponseType:        "AuthorEntry",
		},
		fileImports, helperSources)
	if got, want := len(helperSources), 1; got != want {
		t.Errorf("unexpected length of helperSources: got %d, want %d", got, want)
	}

	constructServerStreamer(&gomodel.ServiceModel{ShortName: "Catalog"},
		&gomodel.RESTHandler{
			RequestTypePackage:  "catalogpb",
			ResponseTypePackage: "responsepb",
			GoMethod:            "StreamTitles",
			ResponseType:        "TitleEntry",
		},
		fileImports, helperSources)
	if got, want := len(helperSources), 2; got != want {
		t.Errorf("unexpected length of helperSources: got %d, want %d", got, want)
	}

	constructServerStreamer(&gomodel.ServiceModel{ShortName: "Media"},
		&gomodel.RESTHandler{
			RequestTypePackage:  "mediapb",
			ResponseTypePackage: "responsepb",
			GoMethod:            "StreamAudio",
			ResponseType:        "AudioEntry",
		},
		fileImports, helperSources)
	if got, want := len(helperSources), 3; got != want {
		t.Errorf("unexpected length of helperSources: got %d, want %d", got, want)
	}

	constructServerStreamer(&gomodel.ServiceModel{ShortName: "Media"},
		&gomodel.RESTHandler{
			RequestTypePackage:  "mediapb",
			ResponseTypePackage: "responsepb",
			GoMethod:            "StreamVideo",
			ResponseType:        "VideoEntry",
		},
		fileImports, helperSources)
	if got, want := len(helperSources), 4; got != want {
		t.Errorf("unexpected length of helperSources: got %d, want %d", got, want)
	}

	actualSources := ""
	for _, key := range helperSources.sortedKeys() {
		actualSources += helperSources[key].Contents() + "\n"
	}

	expectedSources, err := os.ReadFile("testdata/TestConstructServerStreamer.go.baseline")
	if err != nil {
		t.Fatalf("could not load file: %s", err)
	}

	if got, want := actualSources, string(expectedSources); got != want {
		t.Errorf("unexpected helper sources:\n = got: ===\n%s\n= want: ===\n%s\n= diff ===\n%s", got, want, cmp.Diff(got, want))
	}
}
