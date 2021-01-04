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

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRESTCalls(t *testing.T) {
	server := httptest.NewUnstartedServer(nil)
	backend := createBackends()
	restServer := newEndpointREST(nil, backend)

	server.Config = restServer.server
	server.Start()
	defer server.Close()

	for idx, testCase := range []struct {
		verb       string
		path       string
		body       string
		want       string
		statusCode int // 0 taken to mean 200 for simplicty
	}{
		{
			verb: "GET",
			path: "/hello",
			want: "GAPIC Showcase: HTTP/REST endpoint using gorilla/mux\n",
		},
		{
			verb: "GET",
			path: "/v1beta1/repeat:query?info.f_string=jonas^mila",

			// TODO: This should return an error because `^` is not URL-escaped
			want:       `{"info":{"fString":"jonas^mila"}}`,
			statusCode: 200,
		},
	} {

		request, err := http.NewRequest(testCase.verb, server.URL+testCase.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}

		want := testCase.statusCode
		if want == 0 {
			want = 200
		}
		if got := response.StatusCode; got != want {
			t.Errorf("testcase %2d: status code: got %d, want %d", idx, got, want)
			t.Errorf("  request: %v", request)
		}

		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if got, want := string(body), testCase.want; got != want {
			t.Errorf("testcase %2d: body: got `%s`, want %q", idx, got, want)
			t.Errorf("  request: %v", request)
		}
	}

}
