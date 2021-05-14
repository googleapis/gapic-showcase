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

// CheckRESTHeader checks header to ensure that "x-goog-api-client" contains the "rest/" token.
func CheckRESTHeader(header http.Header) error {
	content, ok := header[headerNameAPIClient]
	if !ok || len(content) != 1 {
		return fmt.Errorf("(HeaderAPIClientError) did not find expected HTTP header %q: %q", headerNameAPIClient, headerValueTransportRESTPrefix)
	}

	for _, token := range strings.Split(content[0], " ") {
		if strings.HasPrefix(strings.TrimSpace(token), headerValueTransportRESTPrefix) {
			return nil
		}
	}
	return fmt.Errorf("(HeaderTransportRESTError) did not find expected HTTP header token %q: %q", headerNameAPIClient, headerValueTransportRESTPrefix)
}

// CheckGAPICHeader checks header to ensure that "x-goog-api-client" contains the "gapic/" token.
func CheckGAPICHeader(header http.Header) error {
	content, ok := header[headerNameAPIClient]
	if !ok || len(content) != 1 {
		return fmt.Errorf("(HeaderAPIClientError) did not find expected HTTP header %q: %q", headerNameAPIClient, headerValueTransportRESTPrefix)
	}

	for _, token := range strings.Split(content[0], " ") {
		if strings.HasPrefix(strings.TrimSpace(token), headerValueClientGAPICPrefix) {
			return nil
		}
	}
	return fmt.Errorf("(HeaderClientGAPICError) did not find expected HTTP header token %q: %q", headerNameAPIClient, headerValueClientGAPICPrefix)
}
