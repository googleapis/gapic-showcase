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
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type wireBuffer struct {
	bytes.Buffer
	chunks []string
}

func (wb *wireBuffer) Flush() {
	wb.chunks = append(wb.chunks, wb.String())
	wb.Reset()
}

func TestServerStreamer(t *testing.T) {
	for idx, testCase := range []struct {
		name           string
		messages       []string
		expectedChunks []string
	}{
		{
			name:     "Single chunk",
			messages: []string{"greetings"},
			expectedChunks: []string{
				"[greetings",
				"]",
			},
		},
		{
			name:     "Two chunks",
			messages: []string{"greetings", "  earthling"},
			expectedChunks: []string{
				"[greetings",
				",  earthling",
				"]",
			},
		},
		{
			name:     "Many chunks",
			messages: []string{"greetings", "  people", "of ", "Earth"},
			expectedChunks: []string{
				"[greetings",
				",  people",
				",of ",
				",Earth",
				"]",
			},
		},
		{
			name:           "No chunks",
			messages:       []string{},
			expectedChunks: nil,
		},
		{
			name:           "Single empty chunk",
			messages:       []string{""},
			expectedChunks: nil,
		},
		{
			name:     "Intermediate empty chunk",
			messages: []string{"greetings", "", "earthlings"},
			expectedChunks: []string{
				"[greetings",
				",earthlings",
				"]",
			},
		},
	} {
		label := fmt.Sprintf("[%d:%s]", idx, testCase.name)

		wire := &wireBuffer{}
		streamer, err := NewServerStreamer(wire)
		if err != nil {
			t.Fatalf("%s: could not construct ServerStreamer: %s", label, err)
		}

		for msgIdx, msg := range testCase.messages {
			if err := streamer.sendJSONArrayChunk([]byte(msg)); err != nil {
				t.Errorf("%s: error sending message #%d (%q): %s", label, msgIdx, msg, err)
				break
			}
		}
		if err := streamer.End(); err != nil {
			t.Errorf("%s: error ending stream: %s", label, err)
		}

		if got, want := wire.chunks, testCase.expectedChunks; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: did not received expected chunks\n== got ===\n%#v\n== want ==\n%#v\n",
				label, got, want)
		}
	}
}
