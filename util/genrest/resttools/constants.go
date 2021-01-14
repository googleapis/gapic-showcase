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

package resttools

import "fmt"

const (
	// CharsField contains the characters allowed in a field name (URL path or body)
	CharsField = `-_.0-9a-zA-Z`

	// CharsLiteral contains the the characters allowed in a URL path literal segment.
	CharsLiteral = `-_.~0-9a-zA-Z%`

	// CharsURLPathValue contains the characters allowed in an actual URL path segment encoding
	// user-supplied values.
	//
	// TODO: Gorilla/mux seems to already un-encode these characters, so maybe the right rexep
	// here is [[:ascii:]] or [[:graph:]] or [[:print:]], minus slashes. See
	// https://golang.org/pkg/regexp/syntax/
	CharsURLPathValue = `0-9a-zA-Z_% "\-`
)

var (
	// CharsFieldPath contains the characters allowed in a dotted field path.
	CharsFieldPath string

	// RegexField contains the regex expression for matching a single (not nested) field name.
	RegexField string

	// RegexField contains the regex expression for matching a dotted field path.
	RegexFieldPath string

	// RegexLiteral contains the regex expression for matching a URL path literal segment.
	RegexLiteral string

	// RegexURLPathSingleSegmentValue contains the regex expression for matching a single URL
	// path segment (i.e. `/` is not allowed)
	RegexURLPathSingleSegmentValue string

	// RegexURLPathSingleSegmentValue contains the regex expression for matching multiple URL
	// path segments (i.e. `/` is allowed)
	RegexURLPathMultipleSegmentValue string
)

func init() {
	RegexField = fmt.Sprintf(`[%s]+`, CharsField)

	CharsFieldPath = CharsField + `.`
	RegexFieldPath = fmt.Sprintf(`^[%s]+`, CharsFieldPath)

	RegexLiteral = fmt.Sprintf(`^[%s]+`, CharsLiteral)

	RegexURLPathSingleSegmentValue = fmt.Sprintf("[%s]+", CharsURLPathValue)
	RegexURLPathMultipleSegmentValue = fmt.Sprintf("[%s/]+", CharsURLPathValue)
}
