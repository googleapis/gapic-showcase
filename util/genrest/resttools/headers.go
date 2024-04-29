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
	"fmt"
	"net/http"
	"strings"
)

const (
	headerNameContentType      = "Content-Type"
	headerValueContentTypeJSON = "application/json"

	headerNameAPIClient            = "X-Goog-Api-Client"
	headerValueTransportRESTPrefix = "rest/"
	headerValueClientGAPICPrefix   = "gapic/"
)

// PopulateRequestHeaders inspects request and adds the correct headers. This
// is useful for tests where we're not trying to send incorrect
// headers.
func PopulateRequestHeaders(request *http.Request) {
	header := http.Header{}
	header.Set(headerNameAPIClient, fmt.Sprintf("%s0.0.0 %s0.0.0", headerValueTransportRESTPrefix, headerValueClientGAPICPrefix))

	if request.Body != nil {
		header.Set(headerNameContentType, headerValueContentTypeJSON)
	}

	request.Header = header
}

// CheckContentType checks header to ensure the expected JSON content type is specified.
func CheckContentType(header http.Header) error {
	if content, ok := header[headerNameContentType]; !ok || len(content) != 1 || !strings.HasPrefix(strings.ToLower(strings.TrimSpace(content[0])), headerValueContentTypeJSON) {
		return fmt.Errorf("(HeaderContentTypeError) did not find expected HTTP header %q: %q", headerNameContentType, headerValueContentTypeJSON)
	}
	return nil
}

// CheckAPIClientHeader verifies that the "x-goog-api-client" header contains the expected tokens
// ("rest/..." and "gapic/...").
func CheckAPIClientHeader(header http.Header) error {
	content, ok := header[headerNameAPIClient]
	if !ok || len(content) != 1 {
		return fmt.Errorf("(HeaderAPIClientError) did not find expected HTTP header %q: %q\n                found: %q",
			headerNameAPIClient, headerValueTransportRESTPrefix, header)
	}

	var haveREST, haveGAPIC bool
	for _, token := range strings.Split(content[0], " ") {
		trimmed := strings.TrimSpace(token)
		if !haveREST && strings.HasPrefix(trimmed, headerValueTransportRESTPrefix) {
			haveREST = true
		} else if !haveGAPIC && strings.HasPrefix(trimmed, headerValueClientGAPICPrefix) {
			haveGAPIC = true
		} else {
			// nothing changed
			continue
		}
		if haveREST && haveGAPIC {
			return nil
		}
	}
	if !haveREST {
		return fmt.Errorf("(HeaderTransportRESTError) did not find expected HTTP header token %q: %q", headerNameAPIClient, headerValueTransportRESTPrefix)
	}
	if !haveGAPIC {
		return fmt.Errorf("(HeaderClientGAPICError) did not find expected HTTP header token %q: %q", headerNameAPIClient, headerValueClientGAPICPrefix)
	}
	return fmt.Errorf("internal inconsistency")
}

// PrettyPrintHeaders prints all the HTTP headers in `request` in
// lines preceded by `indentation`. Each line has one header key and a
// quoted list of all values for that key.
func PrettyPrintHeaders(request *http.Request, indentation string) string {
	lines := []string{}
	for key, values := range request.Header {
		newLine := fmt.Sprintf("%s%s:", indentation, key)
		for _, oneValue := range values {
			newLine += fmt.Sprintf(" %q", oneValue)
		}
		lines = append(lines, newLine)
	}
	return strings.Join(lines, "\n")
}

// IncludeRequestHeadersInResponse includes all headers in the request `r` and includes them in the response `w`,
// prefixing each of these header keys with a constant to reflect they came from the request.
func IncludeRequestHeadersInResponse(w http.ResponseWriter, r *http.Request) {
	const prefix = "x-showcase-request-"

	responseHeaders := w.Header()
	for requestKey, valueList := range r.Header {
		for _, value := range valueList {
			responseHeaders.Add(prefix+requestKey, value)
		}
	}
}
