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

func TestInvalidScenarioConfigJSON(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Test-Scenario-Config", "{invalid-json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on invalid scenario config JSON, got %d", rec.Code)
	}
}

func TestFatalErrorOnStartScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	// Default error code (503 Service Unavailable)
	req1 := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req1.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req1.Header.Set("X-Goog-Upload-Command", "start")
	req1.Header.Set("X-Goog-Test-Scenario", "fatal_error_on_start")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on fatal_error_on_start default, got %d", rec1.Code)
	}

	// Configured custom error code (403 Forbidden)
	req2 := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req2.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req2.Header.Set("X-Goog-Upload-Command", "start")
	req2.Header.Set("X-Goog-Test-Scenario", "fatal_error_on_start")
	req2.Header.Set("X-Goog-Test-Scenario-Config", `{"error_code":403}`)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden on fatal_error_on_start custom config, got %d", rec2.Code)
	}
}

func TestNonFatalErrorOnChunkUploadScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_chunk_upload")
	reqStart.Header.Set("X-Goog-Test-Scenario-Config", `{"error_code":503,"failure_count":1,"after_offset":0}`)
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	u, _ := url.Parse(uploadURL)

	// First upload attempt fails with injected 503
	reqUpload1 := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("chunk1")))
	reqUpload1.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload1.Header.Set("X-Goog-Upload-Offset", "0")
	recUpload1 := httptest.NewRecorder()
	handler.ServeHTTP(recUpload1, reqUpload1)

	if recUpload1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on first chunk upload, got %d", recUpload1.Code)
	}

	// Second upload attempt succeeds after exhausting failure_count=1
	reqUpload2 := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("chunk1")))
	reqUpload2.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload2.Header.Set("X-Goog-Upload-Offset", "0")
	recUpload2 := httptest.NewRecorder()
	handler.ServeHTTP(recUpload2, reqUpload2)

	if recUpload2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on retry chunk upload, got %d", recUpload2.Code)
	}
}

func TestNonFatalErrorOnChunkUploadTerminateScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_chunk_upload")
	reqStart.Header.Set("X-Goog-Test-Scenario-Config", `{"error_code":503,"failure_count":1,"action_after_failures":"terminate"}`)
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	u, _ := url.Parse(uploadURL)

	// First upload attempt fails with injected 503
	reqUpload1 := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("chunk1")))
	reqUpload1.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload1.Header.Set("X-Goog-Upload-Offset", "0")
	recUpload1 := httptest.NewRecorder()
	handler.ServeHTTP(recUpload1, reqUpload1)

	if recUpload1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on first chunk upload, got %d", recUpload1.Code)
	}

	// Subsequent upload attempt terminates with 500 Internal Server Error
	reqUpload2 := httptest.NewRequest("POST", u.String(), bytes.NewReader([]byte("chunk1")))
	reqUpload2.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload2.Header.Set("X-Goog-Upload-Offset", "0")
	recUpload2 := httptest.NewRecorder()
	handler.ServeHTTP(recUpload2, reqUpload2)

	if recUpload2.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error when action_after_failures=terminate, got %d", recUpload2.Code)
	}
}

func TestNonFatalErrorOnQueryScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("X-Goog-Test-Scenario", "non_fatal_error_on_query")
	reqStart.Header.Set("X-Goog-Test-Scenario-Config", `{"error_code":502,"failure_count":1}`)
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	u, _ := url.Parse(uploadURL)

	// First query attempt fails with injected 502
	reqQuery1 := httptest.NewRequest("POST", u.String(), nil)
	reqQuery1.Header.Set("X-Goog-Upload-Command", "query")
	recQuery1 := httptest.NewRecorder()
	handler.ServeHTTP(recQuery1, reqQuery1)

	if recQuery1.Code != http.StatusBadGateway {
		t.Fatalf("expected 502 Bad Gateway on first query, got %d", recQuery1.Code)
	}

	// Second query attempt succeeds after exhausting failure_count=1
	reqQuery2 := httptest.NewRequest("POST", u.String(), nil)
	reqQuery2.Header.Set("X-Goog-Upload-Command", "query")
	recQuery2 := httptest.NewRecorder()
	handler.ServeHTTP(recQuery2, reqQuery2)

	if recQuery2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on retry query, got %d", recQuery2.Code)
	}
}

func TestChunkGranularityScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("X-Goog-Test-Scenario", "chunk_granularity")
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	if got, want := recStart.Header().Get("X-Goog-Upload-Chunk-Granularity"), "256"; got != want {
		t.Fatalf("expected chunk granularity %s on chunk_granularity scenario, got %s", want, got)
	}

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	u, _ := url.Parse(uploadURL)

	// Non-final chunk not a multiple of 256 bytes should fail with 400 Bad Request
	reqInvalid := httptest.NewRequest("POST", u.String(), bytes.NewReader(make([]byte, 100)))
	reqInvalid.Header.Set("X-Goog-Upload-Command", "upload")
	reqInvalid.Header.Set("X-Goog-Upload-Offset", "0")
	recInvalid := httptest.NewRecorder()
	handler.ServeHTTP(recInvalid, reqInvalid)

	if recInvalid.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on non-aligned chunk upload, got %d", recInvalid.Code)
	}

	// Non-final chunk that is a multiple of 256 bytes should succeed
	reqValid := httptest.NewRequest("POST", u.String(), bytes.NewReader(make([]byte, 256)))
	reqValid.Header.Set("X-Goog-Upload-Command", "upload")
	reqValid.Header.Set("X-Goog-Upload-Offset", "0")
	recValid := httptest.NewRecorder()
	handler.ServeHTTP(recValid, reqValid)

	if recValid.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on 256-byte aligned chunk upload, got %d", recValid.Code)
	}

	// Finalized chunk of arbitrary size should succeed
	reqFinal := httptest.NewRequest("POST", u.String(), bytes.NewReader(make([]byte, 100)))
	reqFinal.Header.Set("X-Goog-Upload-Command", "upload, finalize")
	reqFinal.Header.Set("X-Goog-Upload-Offset", "256")
	recFinal := httptest.NewRecorder()
	handler.ServeHTTP(recFinal, reqFinal)

	if recFinal.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on unaligned final chunk upload, got %d", recFinal.Code)
	}
}


