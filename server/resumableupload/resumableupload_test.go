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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/googleapis/gapic-showcase/server/resumableupload"
)

func TestHappyPathUpload(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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

func startUploadSession(t *testing.T, handler http.Handler, scenario, config string) string {
	t.Helper()
	req := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	if scenario != "" {
		req.Header.Set("X-Goog-Test-Scenario", scenario)
	}
	if config != "" {
		req.Header.Set("X-Goog-Test-Scenario-Config", config)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on start, got %d", rec.Code)
	}
	uploadURL := rec.Header().Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		t.Fatalf("expected X-Goog-Upload-URL in response, got empty")
	}
	return uploadURL
}

func sendUploadCommand(handler http.Handler, uploadURL, command string, offset int64, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", uploadURL, body)
	req.Header.Set("X-Goog-Upload-Command", command)
	if offset >= 0 {
		req.Header.Set("X-Goog-Upload-Offset", strconv.FormatInt(offset, 10))
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestNonFatalErrorOnChunkUploadScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	uploadURL := startUploadSession(t, handler, "non_fatal_error_on_chunk_upload", `{"error_code":503,"failure_count":1,"after_offset":0}`)

	// First upload attempt fails with injected 503
	recUpload1 := sendUploadCommand(handler, uploadURL, "upload", 0, bytes.NewReader([]byte("chunk1")))
	if recUpload1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on first chunk upload, got %d", recUpload1.Code)
	}

	// Second upload attempt succeeds after exhausting failure_count=1
	recUpload2 := sendUploadCommand(handler, uploadURL, "upload", 0, bytes.NewReader([]byte("chunk1")))
	if recUpload2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on retry chunk upload, got %d", recUpload2.Code)
	}
}

func TestNonFatalErrorOnChunkUploadTerminateScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	uploadURL := startUploadSession(t, handler, "non_fatal_error_on_chunk_upload", `{"error_code":503,"failure_count":1,"action_after_failures":"terminate"}`)

	// First upload attempt fails with injected 503
	recUpload1 := sendUploadCommand(handler, uploadURL, "upload", 0, bytes.NewReader([]byte("chunk1")))
	if recUpload1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on first chunk upload, got %d", recUpload1.Code)
	}

	// Subsequent upload attempt terminates with 500 Internal Server Error
	recUpload2 := sendUploadCommand(handler, uploadURL, "upload", 0, bytes.NewReader([]byte("chunk1")))
	if recUpload2.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error when action_after_failures=terminate, got %d", recUpload2.Code)
	}
}

func TestNonFatalErrorOnQueryScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	uploadURL := startUploadSession(t, handler, "non_fatal_error_on_query", `{"error_code":502,"failure_count":1}`)

	// First query attempt fails with injected 502
	recQuery1 := sendUploadCommand(handler, uploadURL, "query", -1, nil)
	if recQuery1.Code != http.StatusBadGateway {
		t.Fatalf("expected 502 Bad Gateway on first query, got %d", recQuery1.Code)
	}

	// Second query attempt succeeds after exhausting failure_count=1
	recQuery2 := sendUploadCommand(handler, uploadURL, "query", -1, nil)
	if recQuery2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on retry query, got %d", recQuery2.Code)
	}
}

func TestChunkGranularityScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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

// TestFinalizeWithInitialMetadata verifies that finalize returns the metadata name
// from start and the size of the uploaded data payload.
func TestFinalizeWithInitialMetadata(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Start upload session with initial JSON metadata
	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", strings.NewReader(`{"name":"test_file.txt"}`))
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("Content-Type", "application/json")
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	// 2. Finalize upload with 4 bytes of payload data
	reqFinal := httptest.NewRequest("POST", recStart.Header().Get("X-Goog-Upload-URL"), strings.NewReader("data"))
	reqFinal.Header.Set("X-Goog-Upload-Command", "upload, finalize")
	recFinal := httptest.NewRecorder()
	handler.ServeHTTP(recFinal, reqFinal)

	if got, want := recFinal.Body.String(), `{"name":"test_file.txt","size":4}`; got != want {
		t.Fatalf("expected backend response %s, got %s", want, got)
	}
}

func TestUploadURLScheme(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	// Plain HTTP request
	httpReq := httptest.NewRequest("POST", "http://localhost:7469/upload", nil)
	httpReq.Header.Set("X-Goog-Upload-Protocol", "resumable")
	httpReq.Header.Set("X-Goog-Upload-Command", "start")
	httpRec := httptest.NewRecorder()
	handler.ServeHTTP(httpRec, httpReq)

	// Verify returned upload URL uses http
	if got := httpRec.Header().Get("X-Goog-Upload-URL"); !strings.HasPrefix(got, "http://") {
		t.Errorf("expected http:// upload URL, got %q", got)
	}

	// HTTPS request with TLS
	tlsReq := httptest.NewRequest("POST", "https://localhost:7469/upload", nil)
	tlsReq.Header.Set("X-Goog-Upload-Protocol", "resumable")
	tlsReq.Header.Set("X-Goog-Upload-Command", "start")
	tlsRec := httptest.NewRecorder()
	handler.ServeHTTP(tlsRec, tlsReq)

	// Verify returned upload URL uses https
	if got := tlsRec.Header().Get("X-Goog-Upload-URL"); !strings.HasPrefix(got, "https://") {
		t.Errorf("expected https:// upload URL, got %q", got)
	}
}


// TestBinaryPayloadUpload verifies that arbitrary binary payloads (e.g. PNG, octet-stream)
// can be uploaded in the data phase without requiring the binary data to be JSON-formatted.
func TestBinaryPayloadUpload(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Start upload session with initial JSON metadata
	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", strings.NewReader(`{"name":"uploaded_image.png"}`))
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("Content-Type", "application/json")
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	if recStart.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on start, got %d", recStart.Code)
	}

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		t.Fatalf("expected X-Goog-Upload-URL in response, got empty")
	}

	// 2. Upload raw binary PNG bytes (not JSON!)
	binaryPayload := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00")
	reqFinal := httptest.NewRequest("POST", uploadURL, bytes.NewReader(binaryPayload))
	reqFinal.Header.Set("X-Goog-Upload-Command", "upload, finalize")
	reqFinal.Header.Set("X-Goog-Upload-Offset", "0")
	reqFinal.Header.Set("Content-Type", "image/png")
	recFinal := httptest.NewRecorder()
	handler.ServeHTTP(recFinal, reqFinal)

	if recFinal.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on binary upload finalize, got %d with body %s", recFinal.Code, recFinal.Body.String())
	}

	expectedResponse := fmt.Sprintf(`{"name":"uploaded_image.png","size":%d}`, len(binaryPayload))
	if got := strings.TrimSpace(recFinal.Body.String()); got != expectedResponse {
		t.Fatalf("expected final response %s, got %s", expectedResponse, got)
	}
}

// TestDelayMsChunkUploadScenario verifies that delay_ms pauses chunk upload responses.
func TestDelayMsChunkUploadScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Start upload session
	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", strings.NewReader(`{"name":"test.txt"}`))
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("Content-Type", "application/json")
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		t.Fatalf("expected X-Goog-Upload-URL in response, got empty")
	}

	// 2. Upload chunk with delay_ms: 100
	payload := []byte("delayed data chunk")
	reqUpload := httptest.NewRequest("POST", uploadURL, bytes.NewReader(payload))
	reqUpload.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload.Header.Set("X-Goog-Upload-Offset", "0")
	reqUpload.Header.Set("X-Goog-Test-Scenario-Config", `{"delay_ms": 100}`)
	recUpload := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(recUpload, reqUpload)
	elapsed := time.Since(start)

	if recUpload.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", recUpload.Code)
	}
	if elapsed < 100*time.Millisecond {
		t.Fatalf("expected upload to take at least 100ms, elapsed: %v", elapsed)
	}
}

// TestPartialCommitChunkUploadScenario verifies that partial_commit_on_chunk_upload
// commits partial_bytes, returns HTTP 503 Service Unavailable, and allows resuming after query.
func TestPartialCommitChunkUploadScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Start session with partial_commit_on_chunk_upload (commit 4 bytes, fail with 503)
	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", strings.NewReader(`{"name":"test.txt"}`))
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("Content-Type", "application/json")
	reqStart.Header.Set("X-Goog-Test-Scenario", "partial_commit_on_chunk_upload")
	reqStart.Header.Set("X-Goog-Test-Scenario-Config", `{"partial_bytes": 4, "failure_count": 1}`)
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		t.Fatalf("expected X-Goog-Upload-URL in response, got empty")
	}

	// 2. Upload an 8-byte chunk at offset 0
	chunk1 := []byte("12345678")
	reqChunk1 := httptest.NewRequest("POST", uploadURL, bytes.NewReader(chunk1))
	reqChunk1.Header.Set("X-Goog-Upload-Command", "upload")
	reqChunk1.Header.Set("X-Goog-Upload-Offset", "0")
	recChunk1 := httptest.NewRecorder()
	handler.ServeHTTP(recChunk1, reqChunk1)

	// Expect HTTP 503 Service Unavailable with X-Goog-Upload-Status: active
	if recChunk1.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on partial chunk upload, got %d: %s", recChunk1.Code, recChunk1.Body.String())
	}
	if got := recChunk1.Header().Get("X-Goog-Upload-Status"); got != "active" {
		t.Fatalf("expected X-Goog-Upload-Status active, got %q", got)
	}

	// 3. Query current offset to verify server committed 4 bytes
	reqQuery := httptest.NewRequest("POST", uploadURL, nil)
	reqQuery.Header.Set("X-Goog-Upload-Command", "query")
	recQuery := httptest.NewRecorder()
	handler.ServeHTTP(recQuery, reqQuery)

	if recQuery.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on query, got %d: %s", recQuery.Code, recQuery.Body.String())
	}
	if got := recQuery.Header().Get("X-Goog-Upload-Size-Received"); got != "4" {
		t.Fatalf("expected query to report offset 4, got %q", got)
	}

	// 4. Resume upload with remaining 4 bytes at offset 4 and finalize
	chunk2 := []byte("5678")
	reqChunk2 := httptest.NewRequest("POST", uploadURL, bytes.NewReader(chunk2))
	reqChunk2.Header.Set("X-Goog-Upload-Command", "upload, finalize")
	reqChunk2.Header.Set("X-Goog-Upload-Offset", "4")
	recChunk2 := httptest.NewRecorder()
	handler.ServeHTTP(recChunk2, reqChunk2)

	if recChunk2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on resume finalize, got %d: %s", recChunk2.Code, recChunk2.Body.String())
	}
	expectedBody := `{"name":"test.txt","size":8}`
	if got := strings.TrimSpace(recChunk2.Body.String()); got != expectedBody {
		t.Fatalf("expected body %s, got %s", expectedBody, got)
	}
}

// TestDelayMsAfterOffsetChunkUploadScenario verifies that delay_ms with after_offset
// delays chunks at or after after_offset without requiring a dedicated scenario.
func TestDelayMsAfterOffsetChunkUploadScenario(t *testing.T) {
	mgr := resumableupload.NewManager()
	handler := mgr.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Start upload session with delay_ms and after_offset configured
	reqStart := httptest.NewRequest("POST", "http://localhost:7469/upload", strings.NewReader(`{"name":"delay_offset_test.txt"}`))
	reqStart.Header.Set("X-Goog-Upload-Protocol", "resumable")
	reqStart.Header.Set("X-Goog-Upload-Command", "start")
	reqStart.Header.Set("Content-Type", "application/json")
	reqStart.Header.Set("X-Goog-Test-Scenario-Config", `{"delay_ms": 100, "after_offset": 5}`)
	recStart := httptest.NewRecorder()
	handler.ServeHTTP(recStart, reqStart)

	uploadURL := recStart.Header().Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		t.Fatalf("expected X-Goog-Upload-URL in response, got empty")
	}

	// 2. Upload first chunk at offset 0 (length 5) -> should NOT delay since offset < 5
	reqUpload1 := httptest.NewRequest("POST", uploadURL, bytes.NewReader([]byte("hello")))
	reqUpload1.Header.Set("X-Goog-Upload-Command", "upload")
	reqUpload1.Header.Set("X-Goog-Upload-Offset", "0")
	recUpload1 := httptest.NewRecorder()

	start1 := time.Now()
	handler.ServeHTTP(recUpload1, reqUpload1)
	elapsed1 := time.Since(start1)

	if recUpload1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", recUpload1.Code)
	}
	if elapsed1 >= 100*time.Millisecond {
		t.Fatalf("expected first chunk to not delay, but took %v", elapsed1)
	}

	// 3. Upload second chunk at offset 5 (length 5) -> should delay at least 100ms
	reqUpload2 := httptest.NewRequest("POST", uploadURL, bytes.NewReader([]byte("world")))
	reqUpload2.Header.Set("X-Goog-Upload-Command", "upload, finalize")
	reqUpload2.Header.Set("X-Goog-Upload-Offset", "5")
	recUpload2 := httptest.NewRecorder()

	start2 := time.Now()
	handler.ServeHTTP(recUpload2, reqUpload2)
	elapsed2 := time.Since(start2)

	if recUpload2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", recUpload2.Code)
	}
	if elapsed2 < 100*time.Millisecond {
		t.Fatalf("expected second chunk to delay at least 100ms, elapsed: %v", elapsed2)
	}
}
