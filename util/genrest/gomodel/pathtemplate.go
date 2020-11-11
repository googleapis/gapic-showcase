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

var SlashSegment = &Segment{Kind: Literal, Value: "/"}

////////////////////////////////////////
// PathTemplate

type PathTemplate []*Segment

func NewPathTemplate(pattern string) (PathTemplate, error) {
	return parseTemplate(pattern)
}
