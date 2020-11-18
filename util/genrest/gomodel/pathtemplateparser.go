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
	"regexp"
	"strings"
)

// ParseTemplate parses a path template string according to
// https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#path-template-syntax
//
// Grammar:
//    Template = "/" Segments [ Verb ] ;
//    Segments = Segment { "/" Segment } ;
//    Segment  = "*" | "**" | LITERAL | Variable ;
//    Variable = "{" FieldPath [ "=" Segments ] "}" ;
//    FieldPath = IDENT { "." IDENT } ;
//    Verb     = ":" LITERAL ;
// with "**" matching the last part of the path template string except for the Verb.
func ParseTemplate(template string) (parsed PathTemplate, err error) {
	return NewParser(template).parse()
}

////////////////////////////////////////
// Parser

// Parser contains the context for parsing a path template string.
type Parser struct {
	source *Source

	// whether we've encountered the last "**" segment
	haveLastSegment bool

	// regexes
	regex map[string]*regexp.Regexp
}

// NewParser returns a Parser ready to process template.
func NewParser(template string) *Parser {
	return &Parser{
		source: &Source{
			str: template,
		},
		regex: make(map[string]*regexp.Regexp),
	}
}

// parse returns the parsed PathTemplate.
func (parser *Parser) parse() (pt PathTemplate, err error) {
	defer func() {
		if err != nil {
			indent := strings.Repeat(" ", parser.source.idx)
			err = fmt.Errorf("parsing template, position %d: %s   haveLastSegment: %v\n  %q\n   %s^\n   -> %s", parser.source.idx, err, parser.haveLastSegment, parser.source.str, indent, pt)
		}
	}()

	if !parser.source.ConsumeIf('/') {
		return nil, fmt.Errorf("template does not start with slash")
	}

	pt, err = parser.parseSegments()
	pt = append(PathTemplate{slashSegment}, pt...)
	if err != nil {
		return pt, err
	}

	if parser.source.ConsumeIf(':') {
		verb, err := parser.parseLiteral()
		if err != nil {
			return pt, fmt.Errorf("could not parse verb")
		}

		pt = append(pt, &Segment{Kind: Literal, Value: ":"}, verb)
	}

	if parser.source.InRange() {
		return pt, fmt.Errorf("unexpected character")
	}

	return pt, nil

}

// parseSegments parses a sequence of slash-delimited segments.
func (parser *Parser) parseSegments() (PathTemplate, error) {
	pt := PathTemplate{}
	proceed := true

	for proceed {
		if parser.haveLastSegment {
			return pt, fmt.Errorf("already encountered last segment")
		}
		segment, err := parser.parseOneSegment()
		pt = append(pt, segment)
		if err != nil {
			return pt, err
		}

		if proceed = parser.source.ConsumeIf('/'); proceed {
			pt = append(pt, slashSegment)
		}
	}
	return pt, nil
}

// segmentParser is a function type that parses a specific type of segment. It returns both nil
// error and nil Segment if the parser does not apply to the next stream of characters in the
// source. It returns a non-nil Segment if the characters from the point at which called matched the
// segment type.
type segmentParser func() (*Segment, error)

// parseOneSegment parses exactly one segment of the recognized types.
func (parser *Parser) parseOneSegment() (*Segment, error) {
	orderedParsers := []segmentParser{parser.parseVariable, parser.parseLiteral, parser.parseMultipleValue, parser.parseSingleValue}
	for _, parse := range orderedParsers {
		seg, err := parse()
		if err != nil || seg != nil {
			return seg, err
		}
	}
	return nil, fmt.Errorf("could not parse path segment")
}

// parseSingleValue parses a segment with `Kind==SingleValue` (i.e. a single-segment placeholder), returning nil if not possible.
func (parser *Parser) parseSingleValue() (*Segment, error) {
	re := parser.GetRegexp(`^\*`)
	return parser.parseToSegment(re, SingleValue), nil
}

// parseMultipleValue parses a segment with `Kind==MultipleValue` (i.e. a multiple-segment placeholder), returning nil if not possible.
func (parser *Parser) parseMultipleValue() (*Segment, error) {
	re := parser.GetRegexp(`^\*\*`)
	seg := parser.parseToSegment(re, MultipleValue)
	if seg != nil {
		parser.haveLastSegment = true
	}
	return seg, nil
}

// parseLiteral parses a segment with `Kind==Literal`, returning nil if not possible.
func (parser *Parser) parseLiteral() (*Segment, error) {
	re := parser.GetRegexp(`^([-_.~0-9a-zA-Z%]+)`)
	return parser.parseToSegment(re, Literal), nil
}

// parseToSegment is a helper function that creates a Segment of the specified kind if the next
// characters in the parse stream match the expression re.
func (parser *Parser) parseToSegment(re *regexp.Regexp, kind SegmentKind) *Segment {
	match := parser.source.ConsumeRegex(re)
	if len(match) == 0 {
		return nil
	}
	return &Segment{Kind: kind, Value: match}
}

// parseFieldPath parses a field path, which is the "field" in a "{field=segments}" declaration.
func (parser *Parser) parseFieldPath() string {
	re := parser.GetRegexp(`^([-_.9a-zA-Z]+)`)
	return parser.source.ConsumeRegex(re)
}

// parseVariable parses a variable spec, which is a sequence of field names, segments, and/or
// placeholders, all enclosed by matching braces "{}".
func (parser *Parser) parseVariable() (*Segment, error) {
	if !parser.source.ConsumeIf('{') {
		return nil, nil
	}

	fieldPath := parser.parseFieldPath()
	if len(fieldPath) == 0 {
		return nil, fmt.Errorf("no field path specified")
	}

	segment := &Segment{
		Kind:  Variable,
		Value: fieldPath,
	}

	if parser.source.ConsumeIf('=') {
		var err error
		segment.Subsegments, err = parser.parseSegments()
		if err != nil {
			return segment, err
		}
		if len(segment.Subsegments) == 0 {
			return segment, fmt.Errorf("no path segments specified for field path %q", fieldPath)
		}
	} else {
		segment.Subsegments = PathTemplate{&Segment{Kind: SingleValue}}
	}

	if !parser.source.ConsumeIf('}') {
		return segment, fmt.Errorf("expected end-of-variable '}', got %q", parser.source.Str())
	}

	return segment, nil
}

// GetRegexp returns the memoized compiled regex corresponding to expr. This allows defining regexes
// where they're used while still compiling them only once.
func (parser *Parser) GetRegexp(expr string) *regexp.Regexp {
	compiled := parser.regex[expr]
	if compiled == nil {
		compiled = regexp.MustCompile(expr)
		parser.regex[expr] = compiled
	}
	return compiled
}

var slashSegment = &Segment{Kind: Literal, Value: "/"}

////////////////////////////////////////
// Source

// Source contains the context for the source string being parsed. Note that Source and its methods
// are NOT rune-safe and operate on each individual character, not each rune.
type Source struct {
	str string
	idx int
}

// Consume advances source by num characters.
func (src *Source) Consume(num int) {
	src.idx += num
}

// Str returns the unparsed part of the original source string.
func (src *Source) Str() string {
	if !src.InRange() {
		return ""
	}
	return src.str[src.idx:]
}

// InRange returns true iff there are more characters that can be read from the source string.
func (src *Source) InRange() bool {
	return len(src.str) > src.idx
}

// IsNextByte returns true iff the next character to be read is `query`.
func (src *Source) IsNextByte(query byte) bool {
	return src.InRange() && src.str[src.idx] == query
}

// ConsumeIf advances the source and returns true iff the next character matches `query`.
func (src *Source) ConsumeIf(query byte) bool {
	matches := src.IsNextByte(query)
	if matches {
		src.Consume(1)
	}
	return matches
}

// ConsumeMatch consumes and returns characters matching re, and returns them.
func (src *Source) ConsumeRegex(re *regexp.Regexp) string {
	match := re.FindString(src.Str())
	src.Consume(len(match))
	return match
}
