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
	"fmt"
	"strings"
)

// FindValuesMatching returns the the longest string prefix that matches the templates `first` and
// `second`, and a bool indicating whether this prefix if a full match for both expressions.
func FindValuesMatching(first, second PathTemplate) (fullMatch bool, longestMatch string, err error) {
	one := &traverser{template: first.Flatten()}
	two := &traverser{template: second.Flatten()}

	values := make([]string, max(one.Len(), two.Len()))
	defer func() {
		longestMatch = strings.Join(values, "")
	}()

	for !(one.Done() || two.Done()) {
		matches, matching := segmentsMatch(one, two)
		if !matches {
			return false, "", nil
		}
		values = append(values, matching)
	}

	// A Segment of kind MultipleValue can match many other segments, so it doesn't auto-advance
	// in all scenarios.
	one.ConsumeIfKind(MultipleValue)
	two.ConsumeIfKind(MultipleValue)

	if one.ConsumedAll() && two.ConsumedAll() {
		return true, "", nil
	}

	if !(one.ConsumedAll() || two.ConsumedAll()) {
		return false, "", fmt.Errorf("did not traverse either full pattern one: %v    two: %v", *one, *two)
	}

	return false, "", nil
}

// segmentsMatch is the main work function for FindValuesMatching. It returns true iff one and two
// can be matched, as well as the matching string between `one` and `two`, which can be a literal, a
// wildcard "*", or be empty (for some matches in progress). This typically advanced the indices of
// `one` and/or `two`.
func segmentsMatch(one, two *traverser) (bool, string) {
	seg1 := one.Segment()
	seg2 := two.Segment()

	if seg1.Kind != Literal && seg2.Kind == Literal {
		return segmentsMatch(two, one)
	}
	// If at least one of the arguments is a Literal, the first one is.
	if seg1.Kind == Literal {
		switch seg2.Kind {
		case Literal:
			if seg1.Value != seg2.Value {
				return false, ""
			}
			one.Inc()
			two.Inc()
			return true, seg1.Value
		case SingleValue:
			one.Inc()
			two.Inc()
			return true, seg1.Value
		case MultipleValue:
			// special case for the verb, which is the only thing that can come after a
			// multiple value. Advance past the multiple value and reprocess.
			if seg1.Value == ":" {
				two.Inc()
				return true, ""
			}

			one.Inc()
			// no increment for two
			return true, seg1.Value
		default:
			one.SetDone()
			two.SetDone()
			return false, "[ERROR]"
		}
	}
	// If we reach here, neither argument was a Literal

	if seg1.Kind != SingleValue && seg2.Kind == SingleValue {
		return segmentsMatch(two, one)
	}
	// If at least one of the arguments is a SingleValue, the first one is.
	if seg1.Kind == SingleValue {
		switch seg2.Kind {
		case SingleValue:
			one.Inc()
			two.Inc()
			return true, "*"
		case MultipleValue:
			one.Inc()
			// no increment for two
			return true, "*"
		default:
			one.SetDone()
			two.SetDone()
			return false, "[ERROR]"

		}
	}
	// If we reach here, neither argument was a Literal or SingleValue

	if seg1.Kind == MultipleValue {
		if seg2.Kind != MultipleValue {
			one.SetDone()
			two.SetDone()
			return false, fmt.Sprintf("[ERROR: unknown SegmentKind %d]", seg2.Kind)
		}
		// The following should mark both as done, since there shouldn't be any subsequent segments for either
		one.Inc()
		two.Inc()
		return true, "**"
	}

	// We should never reach here (note that flattened PathTemplates contain no segments with Kind==Variable)
	one.SetDone()
	two.SetDone()
	return false, fmt.Sprintf("[ERROR: unknown SegmentKind %d]", seg1.Kind)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

////////////////////////////////////////
// traverser

// traverser is a helper struct for FindValuesMatching. It keeps track of which segment of a
// PathTemplate we're current processing, and whether the traversal has been determined to be
// finished.
type traverser struct {
	// template being traversed
	template PathTemplate

	// index into the segment of template being examined
	idx int

	// whether we're done traversing this template
	done bool
}

// Segment returns the current Segment in the traversal.
func (tr *traverser) Segment() *Segment {
	return tr.template[tr.idx]
}

// Len returns the length of the template.
func (tr *traverser) Len() int {
	return len(tr.template)
}

// Inc advances to the next Segment in the template.
func (tr *traverser) Inc() {
	tr.idx++
	if tr.ConsumedAll() {
		tr.SetDone()
	}
}

// ConsumeIfKind calls Inc() iff the current Segment's `Kind` equals `kind`.
func (tr *traverser) ConsumeIfKind(kind SegmentKind) {
	if !tr.Done() && tr.template[tr.idx].Kind == kind {
		tr.Inc()
	}
}

// ConsumedAll returns true iff we have traversed all the segments of the template.
func (tr *traverser) ConsumedAll() bool {
	return tr.idx >= len(tr.template)
}

// Done returns true iff this traversal has been marked finished.
func (tr *traverser) Done() bool {
	return tr.done
}

// SetDone returns marks this traversal finished.
func (tr *traverser) SetDone() {
	tr.done = true
}
