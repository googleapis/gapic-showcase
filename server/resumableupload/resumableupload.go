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

package resumableupload

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const defaultHost = "localhost:7469"

type status string

const (
	statusActive    status = "active"
	statusFinal     status = "final"
	statusCancelled status = "cancelled"
)

type uploadSession struct {
	ID            string
	CurrentOffset int64
	Buffer        bytes.Buffer
	Status        status
	mu            sync.Mutex
}

type Manager struct {
	sessions sync.Map
	idGen    int64
}

// NewManager creates a new Resumable Upload session manager.
func NewManager() *Manager {
	return &Manager{}
}

// Middleware returns an HTTP middleware that handles Resumable Upload protocol requests.
func (m *Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protocol := r.Header.Get("X-Goog-Upload-Protocol")
		command := r.Header.Get("X-Goog-Upload-Command")

		// If neither X-Goog-Upload-Protocol nor X-Goog-Upload-Command is present, pass through.
		if !strings.EqualFold(protocol, "resumable") && command == "" {
			next.ServeHTTP(w, r)
			return
		}

		m.handleRequest(w, r)
	})
}

func (sess *uploadSession) query(w http.ResponseWriter) {
	w.Header().Set("X-Goog-Upload-Status", string(sess.Status))
	if sess.Status == statusCancelled {
		w.WriteHeader(http.StatusGone)
		return
	}
	if sess.Status == statusFinal {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"name":"%s","size":%d}`, sess.ID, sess.CurrentOffset)
		w.Write([]byte(resp))
		return
	}
	w.Header().Set("X-Goog-Upload-Size-Received", strconv.FormatInt(sess.CurrentOffset, 10))
	w.WriteHeader(http.StatusOK)
}

func (sess *uploadSession) cancel(w http.ResponseWriter) {
	if sess.Status == statusFinal {
		sendError(w, http.StatusBadRequest, "Cannot cancel a finalized session", sess.Status)
		return
	}
	sess.Status = statusCancelled
	w.Header().Set("X-Goog-Upload-Status", string(statusCancelled))
	w.WriteHeader(http.StatusOK)
}

func (sess *uploadSession) upload(w http.ResponseWriter, r *http.Request, offset int64) bool {
	if sess.Status != statusActive {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Cannot upload to session with status %s", sess.Status), sess.Status)
		return false
	}

	if offset != sess.CurrentOffset {
		sendError(w, http.StatusConflict, fmt.Sprintf("Invalid offset: expected %d, got %d", sess.CurrentOffset, offset), sess.Status)
		return false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Error reading request body", sess.Status)
		return false
	}
	sess.Buffer.Write(body)
	sess.CurrentOffset += int64(len(body))
	return true
}

func (sess *uploadSession) finalize(w http.ResponseWriter) {
	if sess.Status == statusCancelled {
		sendError(w, http.StatusBadRequest, "Cannot finalize a cancelled session", sess.Status)
		return
	}
	sess.Status = statusFinal
	w.Header().Set("X-Goog-Upload-Status", string(statusFinal))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := fmt.Sprintf(`{"name":"%s","size":%d}`, sess.ID, sess.CurrentOffset)
	w.Write([]byte(resp))
}

func (m *Manager) handleRequest(w http.ResponseWriter, r *http.Request) {
	commands := parseCommands(r.Header.Get("X-Goog-Upload-Command"))

	if slices.Contains(commands, "start") {
		m.handleStart(w, r)
		return
	}

	sid := r.URL.Query().Get("sid")
	val, ok := m.sessions.Load(sid)
	if !ok || sid == "" {
		sendError(w, http.StatusNotFound, "Invalid or missing session ID", "")
		return
	}
	sess := val.(*uploadSession)

	sess.mu.Lock()
	defer sess.mu.Unlock()

	if slices.Contains(commands, "query") {
		sess.query(w)
		return
	}

	if slices.Contains(commands, "cancel") {
		sess.cancel(w)
		return
	}

	offsetStr := r.Header.Get("X-Goog-Upload-Offset")
	var offset int64 = 0
	if offsetStr != "" {
		var err error
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Invalid X-Goog-Upload-Offset header", sess.Status)
			return
		}
	}

	if slices.Contains(commands, "upload") {
		if !sess.upload(w, r, offset) {
			return
		}
	}

	if slices.Contains(commands, "finalize") {
		sess.finalize(w)
		return
	}

	if slices.Contains(commands, "upload") {
		w.Header().Set("X-Goog-Upload-Status", string(statusActive))
		w.WriteHeader(http.StatusOK)
		return
	}

	sendError(w, http.StatusBadRequest, "Unsupported command sequence", "")
}

func (m *Manager) handleStart(w http.ResponseWriter, r *http.Request) {
	id := atomic.AddInt64(&m.idGen, 1)
	sid := fmt.Sprintf("scotty-sid-%d", id)

	sess := &uploadSession{
		ID:     sid,
		Status: statusActive,
	}
	m.sessions.Store(sid, sess)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = defaultHost
	}
	path := r.URL.Path
	if path == "" {
		path = "/upload"
	}
	query := url.Values{}
	query.Add("sid", sid)
	newUrl := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: query.Encode(),
	}
	uploadURL := newUrl.String()

	w.Header().Set("X-Goog-Upload-Status", string(statusActive))
	w.Header().Set("X-Goog-Upload-URL", uploadURL)
	w.Header().Set("X-Goog-Upload-Control-URL", uploadURL)
	w.Header().Set("X-Goog-Upload-Chunk-Granularity", strconv.Itoa(256*1024))
	w.WriteHeader(http.StatusOK)
}

func parseCommands(header string) []string {
	parts := strings.Split(header, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.ToLower(strings.TrimSpace(p))
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

func sendError(w http.ResponseWriter, code int, message string, st status) {
	if st != "" {
		w.Header().Set("X-Goog-Upload-Status", string(st))
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(message))
}
