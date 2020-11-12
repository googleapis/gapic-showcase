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

type SegmentKind int

const (
	KindUndefined SegmentKind = iota
	Literal
	Variable
	SingleValue
	MultipleValue
	KindEnd
)

func (sk SegmentKind) Valid() bool {
	return sk > KindUndefined && sk < KindEnd
}

func (sk SegmentKind) String() string {
	var names = []string{"(UNDEFINED)", "LITERAL", "VARIABLE", "SINGLEVAL", "MULTIVAL", "(END)"}
	if !sk.Valid() {
		return "INVALID"
	}
	return names[sk]

}

////////////////////////////////////////
// Segment

type Segment struct {
	Kind        SegmentKind
	Value       string // field path if kind==Variable, literal value if kind==Literal, unused otherwise
	Subsegments PathTemplate
}

func (seg *Segment) String() string {
	switch seg.Kind {
	case Literal:
		return fmt.Sprintf("%q", seg.Value)
	case SingleValue, MultipleValue:
		return seg.Value
	case Variable:
		subsegments := "!!ERROR: no subsegments"
		if len(seg.Subsegments) > 0 {
			subsegments = fmt.Sprintf("%s", seg.Subsegments)
		}
		return fmt.Sprintf("{%s = %s}", seg.Value, subsegments)
	}

	// Out of range: print as much info as possible
	return fmt.Sprintf("{%s(%d) %q %s}", seg.Kind, seg.Kind, seg.Value, seg.Subsegments)
}

func (seg *Segment) Flatten() PathTemplate {
	switch seg.Kind {
	case Variable:
		return seg.Subsegments.Flatten()
	default:
		return PathTemplate{seg}
	}
}

var SlashSegment = &Segment{Kind: Literal, Value: "/"}

////////////////////////////////////////
// PathTemplate

type PathTemplate []*Segment

func NewPathTemplate(pattern string) (PathTemplate, error) {
	return parseTemplate(pattern)
}

func (pt PathTemplate) Flatten() PathTemplate {
	flat := PathTemplate{}
	for _, seg := range pt {
		flat = append(flat, seg.Flatten()...)
	}
	return flat
}

type traverser struct {
	pt   PathTemplate
	idx  int
	done bool
}

func (tr *traverser) Segment() *Segment {
	return tr.pt[tr.idx]
}

func (tr *traverser) Len() int {
	return len(tr.pt)
}

func (tr *traverser) Inc() {
	tr.idx++
	if tr.ConsumedAll() {
		tr.SetDone()
	}
}

func (tr *traverser) ConsumedAll() bool {
	return tr.idx >= len(tr.pt)
}

func (tr *traverser) Done() bool {
	return tr.done
}

func (tr *traverser) SetDone() {
	tr.done = true
}

func FindValuesMatching(first, second PathTemplate) (fullMatch bool, longestMatch string, err error) {
	one := &traverser{pt: first.Flatten()}
	two := &traverser{pt: second.Flatten()}

	values := make([]string, max(one.Len(), two.Len()))
	defer func() {
		longestMatch = strings.Join(values, "")
	}()

	for !(one.Done() || two.Done()) {
		match := segmentsMatch(one, two)
		if len(match) == 0 {
			return false, "", nil
		}
		values = append(values, match)
	}
	if one.ConsumedAll() && two.ConsumedAll() {
		return true, "", nil
	}

	if !(one.ConsumedAll() || two.ConsumedAll()) {
		return false, "", fmt.Errorf("did not traverse either full pattern one: %s    two: %s", *one, *two)
	}

	return false, "", nil
}

func segmentsMatch(one, two *traverser) string {
	seg1 := one.Segment()
	seg2 := two.Segment()

	if seg1.Kind != Literal && seg2.Kind == Literal {
		return segmentsMatch(two, one)
	}
	// If at least one of the two is a Literal, the first one is.
	if seg1.Kind == Literal {
		switch seg2.Kind {
		case Literal:
			if seg1.Value != seg2.Value {
				return ""
			}
			one.Inc()
			two.Inc()
			return seg1.Value
		case SingleValue:
			one.Inc()
			two.Inc()
			return seg1.Value
		case MultipleValue:
			one.Inc()
			// no increment for two
			return seg1.Value
		default:
			one.SetDone()
			two.SetDone()
			return "[ERROR]"
		}
	}
	// If we reach here, neither was a Literal

	if seg1.Kind != SingleValue && seg2.Kind == SingleValue {
		return segmentsMatch(two, one)
	}
	// If at least one of the two is a SingleValue, the first one is.
	if seg1.Kind == SingleValue {
		switch seg2.Kind {
		case SingleValue:
			one.Inc()
			two.Inc()
			return "*"
		case MultipleValue:
			one.Inc()
			// no increment for two
			return "*"
		default:
			one.SetDone()
			two.SetDone()
			return "[ERROR]"

		}
	}
	// If we reach here, neither was a Literal or SingleValue

	if seg1.Kind == MultipleValue {
		if seg2.Kind != MultipleValue {
			one.SetDone()
			two.SetDone()
			return fmt.Sprintf("[ERROR: unknown SegmentKind %d]", seg2.Kind)
		}
		// The following should mark both as done, since there shouldn't be any subsequent segments for either
		one.Inc()
		two.Inc()
	}

	// We should never reach here (note that flattened PathTemplates contain no segments with Kind==Variable)
	one.SetDone()
	two.SetDone()
	return fmt.Sprintf("[ERROR: unknown SegmentKind %d]", seg1.Kind)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
