package genrest

import (
	"fmt"
	"regexp"
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
	kind        SegmentKind // LITERAL, VARIABLE
	value       string      // field path if variable, literal value if literal, unused if value
	subsegments PathTemplate
}

func (seg *Segment) String() string {
	subsegments := ""
	if len(seg.subsegments) > 0 {
		subsegments = fmt.Sprintf(" %s", seg.subsegments)
	}
	return fmt.Sprintf("{%s %q%s}", seg.kind, seg.value, subsegments)
}

var SlashSegment = &Segment{kind: Literal, value: "/"}

////////////////////////////////////////
// PathTemplate

type PathTemplate []*Segment

func NewPathTemplate(pattern string) (PathTemplate, error) {
	return parseTemplate(pattern)

}

////////////////////////////////////////
// PathTemplate

type scanner struct {
	// not rune-safe
	str             string
	idx             int
	haveLastSegment bool
}

func (src *scanner) Consume(num int) {
	src.idx += num
}

func (src *scanner) Str() string {
	if !src.InRange() {
		return ""
	}
	return src.str[src.idx:]
}

func (src *scanner) InRange() bool {
	return len(src.str) > src.idx
}

func (src *scanner) IsNextByte(query byte) bool {
	return src.InRange() && src.str[src.idx] == query
}

// Make parser an object; move haveLastSegment here

func parseTemplate(template string) (pt PathTemplate, err error) {
	src := &scanner{
		str: template,
		idx: 0}
	defer func() {
		if err != nil {
			indent := strings.Repeat(" ", src.idx)
			err = fmt.Errorf("parsing template, position %d: %s\n  %q\n   %s^\n   -> %s", src.idx, err, template, indent, pt)
		}
	}()

	if !src.IsNextByte('/') {
		return nil, fmt.Errorf("template does not start with slash")
	}
	src.Consume(1)

	pt, err = parseSegments(src)
	pt = append(PathTemplate{SlashSegment}, pt...)
	if err != nil {
		return pt, err
	}

	if src.IsNextByte(':') {
		src.Consume(1)
		verb, err := parseLiteral(src)
		if err != nil {
			return pt, fmt.Errorf("could not parse verb")
		}

		pt = append(pt, &Segment{kind: Literal, value: ":"}, verb)

	}

	if src.InRange() {
		return pt, fmt.Errorf("unexpected character")
	}

	return pt, nil

}

func parseSegments(src *scanner) (PathTemplate, error) {
	pt := PathTemplate{}
	proceed := true

	for proceed {
		if src.haveLastSegment {
			return pt, fmt.Errorf("already encountered last segment")
		}
		segment, err := parseOneSegment(src)
		if err != nil {
			return pt, err
		}
		pt = append(pt, segment)
		if proceed = src.IsNextByte('/'); proceed {
			src.Consume(1)
			pt = append(pt, SlashSegment)
		}
	}
	return pt, nil
}

type segmentParser func(*scanner) (*Segment, error)

func parseOneSegment(src *scanner) (*Segment, error) {
	orderedParsers := []segmentParser{parseVariable, parseLiteral, parseMultipleValue, parseSingleValue}
	for _, parser := range orderedParsers {
		seg, err := parser(src)
		if err != nil || seg != nil {
			return seg, err
		}
	}
	return nil, fmt.Errorf("could not parse path segment")

}

func parseSingleValue(src *scanner) (*Segment, error) {
	re := regexp.MustCompile(`\*`)
	return parseToSegment(src, re, SingleValue)
}

func parseMultipleValue(src *scanner) (*Segment, error) {
	re := regexp.MustCompile(`\*\*`)
	seg, err := parseToSegment(src, re, MultipleValue)
	if seg != nil {
		src.haveLastSegment = true
	}
	return seg, err
}

func parseLiteral(src *scanner) (*Segment, error) {
	re := regexp.MustCompile("([a-zA-Z0-9_%]*)")
	return parseToSegment(src, re, Literal)
}

func parseToSegment(src *scanner, re *regexp.Regexp, kind SegmentKind) (*Segment, error) {
	match := consumeRegex(src, re)
	if len(match) == 0 {
		return nil, nil
	}
	return &Segment{kind: kind, value: match}, nil
}

func parseFieldPath(src *scanner) string {
	re := regexp.MustCompile("([a-zA-Z0-9_.]*)")
	return consumeRegex(src, re)
}

func consumeRegex(src *scanner, re *regexp.Regexp) string {
	match := re.FindString(src.Str())
	src.Consume(len(match))
	return match
}

func parseVariable(src *scanner) (*Segment, error) {
	if !src.IsNextByte('{') {
		return nil, nil
	}
	src.Consume(1)

	fieldPath := parseFieldPath(src)
	if len(fieldPath) == 0 {
		return nil, fmt.Errorf("no field path specified")
	}

	segment := &Segment{
		kind:  Variable,
		value: fieldPath,
	}

	//var subsegments PathTemplate
	if src.IsNextByte('=') {
		src.Consume(1)

		var err error
		segment.subsegments, err = parseSegments(src)
		if err != nil {
			return segment, err
		}
		if len(segment.subsegments) == 0 {
			return segment, fmt.Errorf("no path segments specified for URI %q", fieldPath)
		}
	} else {
		segment.subsegments = PathTemplate{&Segment{kind: SingleValue}}
	}

	//	segment.subsegments = subsegments

	if !src.IsNextByte('}') {
		return segment, fmt.Errorf("expected end-of-variable '}', got %q %q", src.Str(), src.str[src.idx-1:])
	}
	src.Consume(1)

	return segment, nil
}

/*	idxEquals := strings.Index(src.str[src.idx+1:], "=")
	idxOpen := strings.Index(src.str[src.idx+1:], "{")
	idxClose := strings.Index(src.str[src.idx+1:], "}")
	if idxEquals < 0 ||
		(idxOpen >= 0 && idxEquals > idxOpen); idxClose >= 0 && idxEquals > idxClose {

	}
*/
