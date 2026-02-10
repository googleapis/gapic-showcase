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
	"strings"
)

// SystemParameters encapsulates the system parameters recognized by Showcase. These are a subset of
// Google's accepted system parameters described in
// https://cloud.google.com/apis/docs/system-parameters.
type SystemParameters struct {
	EnumEncodingAsInt bool
	APIVersion        string
}

// GetSystemParameters returns the SystemParameters encoded in request, and the query parameters in
// the request's query string that are not themselves system parameters.
func GetSystemParameters(request *http.Request) (systemParams *SystemParameters, queryParams map[string][]string, err error) {
	return processQueryString(request.URL.RawQuery)
}

// processQueryString returns the SystemParameters encoded in queryString, and the query parameters in
// queryString that are not themselves system parameters.
//
// Since we want GAPICs to be strict in what they emit, and Showcase helps guarantee that, this
// function is strict in what it accepts:
// - no more than one instance of the $alt system parameter
// - $alt may be appear as "$alt" or "alt", always lower case
// - the only possible values for $alt are "json" or "json;enum-encoding=int"; again, lower case
// - the semicolon in "json;enum-encoding=int" must be URL escaped as "%3B" or "%3b"
// - the equal sign in "json;enum-encoding=int" may or may not be URL-escaped
func processQueryString(queryString string) (systemParams *SystemParameters, queryParams map[string][]string, err error) {

	// We parse the raw query string manually rather than relying on request.URL.Query() so that
	// we can error out in the case of malformed strings (e.g. unencoded ampersands), rather
	// than having them ignored with potentially incorrect results.
	queryPairs, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Try removing this workaround when we update the Go version in the Showcase CI to
	// 1.17+. As of this writing, CI uses 1.16.3, and the tests fail
	// (https://github.com/googleapis/gapic-showcase/runs/6798834903?check_suite_focus=true#step:6:60)
	// without this explicit check. The tests pass without this workaround on local machines or in the Go
	// Playground (https://go.dev/play/p/ewyv5qj55an) using Go >=1.17
	if strings.Contains(queryString, ";") {
		return nil, nil, fmt.Errorf("found unescaped semicolon in query string %q", queryString)
	}

	queryParams = map[string][]string(queryPairs)
	systemParams = &SystemParameters{}
	sawAltParam := false
	for param, values := range queryPairs {
		switch param {
		case "alt", "$alt":
			for _, val := range values {
				if sawAltParam {
					return systemParams, queryParams, fmt.Errorf("multiple instances of $alt system parameter")
				}

				switch val {
				case "json":
					// no op
				case "json;enum-encoding=int": // already URL-decoded
					systemParams.EnumEncodingAsInt = true
				default:
					return systemParams, queryParams, fmt.Errorf("unhandled value %q for system parameter %q", val, param)
				}
			}
			delete(queryParams, param)
			sawAltParam = true
		case "apiVersion", "$apiVersion":
			if systemParams.APIVersion != "" {
				return systemParams, queryParams, fmt.Errorf("multiple instances of $apiVersion system parameter")
			}
			for _, val := range values {
				systemParams.APIVersion = val
			}
			delete(queryParams, param)
		}
	}
	return systemParams, queryParams, nil
}
