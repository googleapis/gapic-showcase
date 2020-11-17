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

// ParseTemplate parses according to
// https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#path-template-syntax
//
// Grammar:
//    Template = "/" Segments [ Verb ] ;
//    Segments = Segment { "/" Segment } ;
//    Segment  = "*" | "**" | LITERAL | Variable ;
//    Variable = "{" FieldPath [ "=" Segments ] "}" ;
//    FieldPath = IDENT { "." IDENT } ;
//    Verb     = ":" LITERAL ;
func ParseTemplate(template string) (pt PathTemplate, err error) {
	parser := &Parser{
		source: &Source{
			str: template,
			idx: 0},
	}
	return parser.parse()
}

////////////////////////////////////////
// Parser

type Parser struct {
	source          *Source
	haveLastSegment bool
}

func (parser *Parser) parse() (pt PathTemplate, err error) {
	defer func() {
		if err != nil {
			indent := strings.Repeat(" ", parser.source.idx)
			err = fmt.Errorf("parsing template, position %d: %s\n  %q\n   %s^\n   -> %s", parser.source.idx, err, parser.source.str, indent, pt)
		}
	}()

	if !parser.source.ConsumeIf('/') {
		return nil, fmt.Errorf("template does not start with slash")
	}

	pt, err = parser.parseSegments()
	pt = append(PathTemplate{SlashSegment}, pt...)
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

func (parser *Parser) parseSegments() (PathTemplate, error) {
	pt := PathTemplate{}
	proceed := true

	for proceed {
		if parser.haveLastSegment {
			return pt, fmt.Errorf("already encountered last segment")
		}
		segment, err := parser.parseOneSegment()
		if err != nil {
			return pt, err
		}
		pt = append(pt, segment)
		if proceed = parser.source.ConsumeIf('/'); proceed {
			pt = append(pt, SlashSegment)
		}
	}
	return pt, nil
}

type segmentParser func() (*Segment, error)

func (parser *Parser) parseOneSegment() (*Segment, error) {
	orderedParsers := []segmentParser{parser.parseVariable, parser.parseLiteral, parser.parseMultipleValue, parser.parseSingleValue}
	for _, parser := range orderedParsers {
		seg, err := parser()
		if err != nil || seg != nil {
			return seg, err
		}
	}
	return nil, fmt.Errorf("could not parse path segment")

}

func (parser *Parser) parseSingleValue() (*Segment, error) {
	re := regexp.MustCompile(`\*`)
	return parser.parseToSegment(re, SingleValue)
}

func (parser *Parser) parseMultipleValue() (*Segment, error) {
	re := regexp.MustCompile(`\*\*`)
	seg, err := parser.parseToSegment(re, MultipleValue)
	if seg != nil {
		parser.haveLastSegment = true
	}
	return seg, err
}

func (parser *Parser) parseLiteral() (*Segment, error) {
	re := regexp.MustCompile("([a-zA-Z0-9_%]*)")
	return parser.parseToSegment(re, Literal)
}

func (parser *Parser) parseToSegment(re *regexp.Regexp, kind SegmentKind) (*Segment, error) {
	match := parser.source.ConsumeRegex(re)
	if len(match) == 0 {
		return nil, nil
	}
	return &Segment{Kind: kind, Value: match}, nil
}

func (parser *Parser) parseFieldPath() string {
	re := regexp.MustCompile("([a-zA-Z0-9_.]*)")
	return parser.source.ConsumeRegex(re)
}

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
			return segment, fmt.Errorf("no path segments specified for URI %q", fieldPath)
		}
	} else {
		segment.Subsegments = PathTemplate{&Segment{Kind: SingleValue}}
	}

	if !parser.source.ConsumeIf('}') {
		return segment, fmt.Errorf("expected end-of-variable '}', got %q %q", parser.source.Str(), parser.source.str[parser.source.idx-1:])
	}

	return segment, nil
}

////////////////////////////////////////
// Source

type Source struct {
	// not rune-safe
	str             string
	idx             int
	haveLastSegment bool
}

func (src *Source) Consume(num int) {
	src.idx += num
}

func (src *Source) Str() string {
	if !src.InRange() {
		return ""
	}
	return src.str[src.idx:]
}

func (src *Source) InRange() bool {
	return len(src.str) > src.idx
}

func (src *Source) IsNextByte(query byte) bool {
	return src.InRange() && src.str[src.idx] == query
}

func (src *Source) ConsumeIf(query byte) bool {
	matches := src.IsNextByte(query)
	if matches {
		src.Consume(1)
	}
	return matches
}

func (src *Source) ConsumeRegex(re *regexp.Regexp) string {
	match := re.FindString(src.Str())
	src.Consume(len(match))
	return match
}
