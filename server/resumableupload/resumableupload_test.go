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
	"strconv"
	"strings"
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
	if got, want := rec.Header().Get("X-Goog-Upload-Chunk-Granularity"), strconv.Itoa(256*1024); got != want {
		t.Fatalf("expected chunk granularity %s, got %s", want, got)
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

func TestQueryAndUploadAfterFinalize(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	uploadURL := rec.Header().Get("X-Goog-Upload-URL")
	u, _ := url.Parse(uploadURL)

	reqFinal := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("done")))
	reqFinal.Header.Set("X-Goog-Upload-Command", "upload, finalize")
	reqFinal.Header.Set("X-Goog-Upload-Offset", "0")
	recFinal := httptest.NewRecorder()
	handler.ServeHTTP(recFinal, reqFinal)

	// Query finalized session
	reqQuery := httptest.NewRequest("POST", u.String(), nil)
	reqQuery.Header.Set("X-Goog-Upload-Command", "query")
	recQuery := httptest.NewRecorder()
	handler.ServeHTTP(recQuery, reqQuery)

	if recQuery.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on query finalized session, got %d", recQuery.Code)
	}
	if recQuery.Header().Get("X-Goog-Upload-Status") != "final" {
		t.Fatalf("expected final status on query, got %s", recQuery.Header().Get("X-Goog-Upload-Status"))
	}
	if !strings.Contains(recQuery.Body.String(), `"size":4`) {
		t.Fatalf("expected finalized JSON body on query, got %s", recQuery.Body.String())
	}

	// Attempt upload after finalize
	reqUpload := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("more")))
	reqUpload.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload.Header.Set("X-Goog-Upload-Offset", "4")
	recUpload := httptest.NewRecorder()
	handler.ServeHTTP(recUpload, reqUpload)

	if recUpload.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on upload after finalize, got %d", recUpload.Code)
	}
}

func TestQueryAndUploadAfterCancel(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	uploadURL := rec.Header().Get("X-Goog-Upload-URL")
	u, _ := url.Parse(uploadURL)

	reqCancel := httptest.NewRequest("POST", u.String(), nil)
	reqCancel.Header.Set("X-Goog-Upload-Command", "cancel")
	recCancel := httptest.NewRecorder()
	handler.ServeHTTP(recCancel, reqCancel)

	// Query cancelled session
	reqQuery := httptest.NewRequest("POST", u.String(), nil)
	reqQuery.Header.Set("X-Goog-Upload-Command", "query")
	recQuery := httptest.NewRecorder()
	handler.ServeHTTP(recQuery, reqQuery)

	if recQuery.Code != http.StatusGone {
		t.Fatalf("expected 410 Gone on query cancelled session, got %d", recQuery.Code)
	}
	if recQuery.Header().Get("X-Goog-Upload-Status") != "cancelled" {
		t.Fatalf("expected cancelled status on query, got %s", recQuery.Header().Get("X-Goog-Upload-Status"))
	}

	// Attempt upload after cancel
	reqUpload := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("more")))
	reqUpload.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload.Header.Set("X-Goog-Upload-Offset", "0")
	recUpload := httptest.NewRecorder()
	handler.ServeHTTP(recUpload, reqUpload)

	if recUpload.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on upload after cancel, got %d", recUpload.Code)
	}
}

func TestNonFatalStartErrorWithClientUUID(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	cfgA := `{"client_uuid":"client-a","error_code":503,"failure_count":1}`

	// 1st request for client-a fails with 503
	reqA1 := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	reqA1.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqA1.Header.Set("X-Goog-Upload-Command", "start")
	reqA1.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_start")
	reqA1.Header.Set("X-Goog-Test-Scenario-Config", cfgA)
	recA1 := httptest.NewRecorder()
	handler.ServeHTTP(recA1, reqA1)

	if recA1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on first start for client-a, got %d", recA1.Code)
	}

	// 2nd request for client-a succeeds (failure_count=1 exhausted)
	reqA2 := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	reqA2.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqA2.Header.Set("X-Goog-Upload-Command", "start")
	reqA2.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_start")
	reqA2.Header.Set("X-Goog-Test-Scenario-Config", cfgA)
	recA2 := httptest.NewRecorder()
	handler.ServeHTTP(recA2, reqA2)

	if recA2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on retry start for client-a, got %d", recA2.Code)
	}

	// Request for client-b fails independently with 503 despite same RemoteAddr
	cfgB := `{"client_uuid":"client-b","error_code":503,"failure_count":1}`
	reqB1 := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	reqB1.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqB1.Header.Set("X-Goog-Upload-Command", "start")
	reqB1.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_start")
	reqB1.Header.Set("X-Goog-Test-Scenario-Config", cfgB)
	recB1 := httptest.NewRecorder()
	handler.ServeHTTP(recB1, reqB1)

	if recB1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on first start for client-b, got %d", recB1.Code)
	}
}
