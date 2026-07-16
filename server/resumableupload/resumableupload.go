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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type ScenarioConfig struct {
	ErrorCode           int    `json:"error_code"`
	FailureCount        int    `json:"failure_count"`
	ActionAfterFailures string `json:"action_after_failures"`
	AfterOffset         int64  `json:"after_offset"`
}

type uploadSession struct {
	ID             string
	CurrentOffset  int64
	Buffer         bytes.Buffer
	Status         string
	Scenario       string
	ScenarioConfig ScenarioConfig
	StartFailures  int
	UploadFailures int
	mu             sync.Mutex
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

func (m *Manager) handleRequest(w http.ResponseWriter, r *http.Request) {
	commands := parseCommands(r.Header.Get("X-Goog-Upload-Command"))

	if containsCommand(commands, "start") {
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

	if containsCommand(commands, "query") {
		if sess.Scenario == "non_fatal_error_on_query" {
			if sess.UploadFailures < sess.ScenarioConfig.FailureCount {
				sess.UploadFailures++
				sendError(w, sess.ScenarioConfig.ErrorCode, "Injected non-fatal query error", "active")
				return
			}
		}
		w.Header().Set("X-Goog-Upload-Status", sess.Status)
		w.Header().Set("X-Goog-Upload-Size-Received", strconv.FormatInt(sess.CurrentOffset, 10))
		w.WriteHeader(http.StatusOK)
		return
	}

	if containsCommand(commands, "cancel") {
		sess.Status = "cancelled"
		w.Header().Set("X-Goog-Upload-Status", "cancelled")
		w.WriteHeader(http.StatusOK)
		return
	}

	offsetStr := r.Header.Get("X-Goog-Upload-Offset")
	var offset int64 = 0
	if offsetStr != "" {
		var err error
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Invalid X-Goog-Upload-Offset header", "active")
			return
		}
	}

	if containsCommand(commands, "upload") {
		if sess.Scenario == "non_fatal_error_on_chunk_upload" {
			if offset >= sess.ScenarioConfig.AfterOffset {
				if sess.UploadFailures < sess.ScenarioConfig.FailureCount {
					sess.UploadFailures++
					sendError(w, sess.ScenarioConfig.ErrorCode, "Injected non-fatal chunk upload error", "active")
					return
				}
				if sess.ScenarioConfig.ActionAfterFailures == "terminate" {
					sendError(w, http.StatusInternalServerError, "Scenario requested termination", "")
					return
				}
			}
		}

		if sess.Scenario == "chunk_granularity" && !containsCommand(commands, "finalize") {
			bodyLen := r.ContentLength
			if bodyLen > 0 && bodyLen%256 != 0 {
				sendError(w, http.StatusBadRequest, fmt.Sprintf("Chunk size %d is not a multiple of granularity 256", bodyLen), "active")
				return
			}
		}

		if offset != sess.CurrentOffset {
			sendError(w, http.StatusConflict, fmt.Sprintf("Invalid offset: expected %d, got %d", sess.CurrentOffset, offset), "active")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Error reading request body", "active")
			return
		}
		sess.Buffer.Write(body)
		sess.CurrentOffset += int64(len(body))
	}

	if containsCommand(commands, "finalize") {
		sess.Status = "final"
		w.Header().Set("X-Goog-Upload-Status", "final")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"name":"%s","size":%d}`, sess.ID, sess.CurrentOffset)
		w.Write([]byte(resp))
		return
	}

	if containsCommand(commands, "upload") {
		w.Header().Set("X-Goog-Upload-Status", "active")
		w.WriteHeader(http.StatusOK)
		return
	}

	sendError(w, http.StatusBadRequest, "Unsupported command sequence", "")
}

func (m *Manager) handleStart(w http.ResponseWriter, r *http.Request) {
	scenario := r.Header.Get("X-Goog-Test-Scenario")
	if scenario == "" {
		scenario = "happy_path"
	}

	config := ScenarioConfig{
		ErrorCode:           http.StatusServiceUnavailable,
		FailureCount:        1,
		ActionAfterFailures: "succeed",
		AfterOffset:         0,
	}
	if cfgStr := r.Header.Get("X-Goog-Test-Scenario-Config"); cfgStr != "" {
		_ = json.Unmarshal([]byte(cfgStr), &config)
	}

	// Use scenario+client specific failure counter key or session state
	// For start requests before session creation, track start failures on manager if needed,
	// or use scenario header failure tracking per client IP/scenario.
	key := scenario + "-" + r.RemoteAddr
	failures := m.incrementStartFailure(key, scenario, config.FailureCount)
	if scenario == "non_fatal_error_on_start" && failures <= config.FailureCount {
		sendError(w, config.ErrorCode, "Injected non-fatal start error", "")
		return
	}
	if scenario == "fatal_error_on_start" {
		code := config.ErrorCode
		if code == 0 {
			code = http.StatusForbidden
		}
		sendError(w, code, "Injected fatal start error", "")
		return
	}

	id := atomic.AddInt64(&m.idGen, 1)
	sid := fmt.Sprintf("scotty-sid-%d", id)

	sess := &uploadSession{
		ID:             sid,
		Status:         "active",
		Scenario:       scenario,
		ScenarioConfig: config,
	}
	m.sessions.Store(sid, sess)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:7469"
	}
	path := r.URL.Path
	if path == "" {
		path = "/upload"
	}

	uploadURL := fmt.Sprintf("%s://%s%s?sid=%s", scheme, host, path, sid)
	granularity := "1"
	if scenario == "chunk_granularity" {
		granularity = "256"
	}
	w.Header().Set("X-Goog-Upload-Status", "active")
	w.Header().Set("X-Goog-Upload-URL", uploadURL)
	w.Header().Set("X-Goog-Upload-Control-URL", uploadURL)
	w.Header().Set("X-Goog-Upload-Chunk-Granularity", granularity)
	w.WriteHeader(http.StatusOK)
}

var startFailureMap sync.Map

func (m *Manager) incrementStartFailure(key, scenario string, limit int) int {
	if scenario != "non_fatal_error_on_start" {
		return limit + 1
	}
	val, _ := startFailureMap.LoadOrStore(key, new(int32))
	ptr := val.(*int32)
	current := int(atomic.AddInt32(ptr, 1))
	if current > limit {
		// Reset after success so subsequent tests can run independently
		atomic.StoreInt32(ptr, 0)
	}
	return current
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

func containsCommand(commands []string, cmd string) bool {
	for _, c := range commands {
		if c == cmd {
			return true
		}
	}
	return false
}

func sendError(w http.ResponseWriter, code int, msg string, status string) {
	if status != "" {
		w.Header().Set("X-Goog-Upload-Status", status)
	}
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
