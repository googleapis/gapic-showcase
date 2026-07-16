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
