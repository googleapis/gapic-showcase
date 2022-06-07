// Copyright 2022 Google LLC
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
	"net/url"
)

type SystemParameters struct {
	EnumEncodingAsInt bool
}

func GetSystemParams(request *http.Request) (queryParams map[string][]string, systemParams *SystemParameters, err error) {
	return processQueryString(request.URL.RawQuery)

}
func processQueryString(queryString string) (queryParams map[string][]string, systemParams *SystemParameters, err error) {
	// If and when we support additional system parameters that could be specified in request
	// headers, we can make this a package-private function called from a m re general entry
	// point `GetSystemParameters(http.Request) (queryParams, systemParams)`.

	// TODO: Run https://pkg.go.dev/net/url#ParseQuery to check for un-encoded ampersands
	// Justification: Since semicolons can be valid query string delimiters (cf
	// https://github.com/golang/go/issues/25192,
	// https://en.wikipedia.org/wiki/Query_string#Web_forms), we insist on URL-encoded system
	// parameters

	// We parse the raw query string manually rather than relying on request.URL.Query() so that
	// we can error out in the case of malformed strings (e.g. unencoded ampersands), rather
	// than having them ignored with potentially incorrect results.
	queryPairs, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, nil, err
	}
	queryParams = map[string][]string(queryPairs)
	systemParams = &SystemParameters{}
	for param, values := range queryPairs {
		switch param {
		case "alt", "$alt":
			for _, val := range values {
				switch val {
				case "json":
					// no op
				case "json;enum-encoding=int": // already URL-decoded
					systemParams.EnumEncodingAsInt = true
				default:
					return queryParams, systemParams, fmt.Errorf("unhandled value %q for system parameter %q", val, param)
				}
			}
			delete(queryParams, param)
		}
	}
	return queryParams, systemParams, nil
}
