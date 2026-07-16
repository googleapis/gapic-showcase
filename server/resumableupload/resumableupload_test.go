// Copyright 2026 Google LLC
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

package resumableupload_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/googleapis/gapic-showcase/server/resumableupload"
)

func TestHappyPathUpload(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	// 1. Start command
	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on start, got %d", rec.Code)
	}
	if rec.Header().Get("X-Goog-Upload-Status") != "active" {
		t.Fatalf("expected active status, got %s", rec.Header().Get("X-Goog-Upload-Status"))
	}
	uploadURL := rec.Header().Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		t.Fatal("expected X-Goog-Upload-URL header")
	}

	u, err := url.Parse(uploadURL)
	if err != nil {
		t.Fatalf("failed to parse upload URL: %v", err)
	}

	// 2. Upload chunk
	req2 := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("hello")))
	req2.Header.Set("X-Goog-Upload-Command", "upload")
	req2.Header.Set("X-Goog-Upload-Offset", "0")
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on chunk upload, got %d: %s", rec2.Code, rec2.Body.String())
	}

	// 3. Finalize
	req3 := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte(" world")))
	req3.Header.Set("X-Goog-Upload-Command", "upload, finalize")
	req3.Header.Set("X-Goog-Upload-Offset", "5")
	rec3 := httptest.NewRecorder()

	handler.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on finalize, got %d", rec3.Code)
	}
	if rec3.Header().Get("X-Goog-Upload-Status") != "final" {
		t.Fatalf("expected final status, got %s", rec3.Header().Get("X-Goog-Upload-Status"))
	}
}

func TestNonFatalErrorOnStart(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	// 1. First start should fail with 503
	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_start")
	req.Header.Set("X-Goog-Test-Scenario-Config", `{"error_code": 503, "failure_count": 1}`)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 on first start, got %d", rec.Code)
	}

	// 2. Second start should succeed
	req2 := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req2.RemoteAddr = "1.2.3.4:1234"
	req2.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req2.Header.Set("X-Goog-Upload-Command", "start")
	req2.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_start")
	req2.Header.Set("X-Goog-Test-Scenario-Config", `{"error_code": 503, "failure_count": 1}`)
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 on retry start, got %d", rec2.Code)
	}
}

func TestNonFatalErrorOnQuery(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_query")
	req.Header.Set("X-Goog-Test-Scenario-Config", `{"error_code": 503, "failure_count": 1}`)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on start, got %d", rec.Code)
	}
	uploadURL := rec.Header().Get("X-Goog-Upload-URL")

	reqQuery1 := httptest.NewRequest("POST", uploadURL, nil)
	reqQuery1.Header.Set("X-Goog-Upload-Command", "query")
	recQuery1 := httptest.NewRecorder()
	handler.ServeHTTP(recQuery1, reqQuery1)
	if recQuery1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 on first query, got %d", recQuery1.Code)
	}

	reqQuery2 := httptest.NewRequest("POST", uploadURL, nil)
	reqQuery2.Header.Set("X-Goog-Upload-Command", "query")
	recQuery2 := httptest.NewRecorder()
	handler.ServeHTTP(recQuery2, reqQuery2)
	if recQuery2.Code != http.StatusOK {
		t.Fatalf("expected 200 on retry query, got %d", recQuery2.Code)
	}
}

func TestChunkGranularityScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Test-Scenario", "chunk_granularity")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on start, got %d", rec.Code)
	}
	if rec.Header().Get("X-Goog-Upload-Chunk-Granularity") != "256" {
		t.Fatalf("expected granularity 256, got %s", rec.Header().Get("X-Goog-Upload-Chunk-Granularity"))
	}
}

