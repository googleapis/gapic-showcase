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
	"net/http"
)

const (
	headerNameContentType      = "Content-Type"
	headerValueContentTypeJSON = "application/json"
)

// PopulateRequestHeaders inspects request and adds the correct headers. This
// is useful for tests where we're not trying to send incorrect
// headers.
func PopulateRequestHeaders(request *http.Request) {
	if request.Body != nil {
		request.Header = http.Header{headerNameContentType: []string{headerValueContentTypeJSON}}
	}
}
