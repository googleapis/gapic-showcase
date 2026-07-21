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

type ScenarioConfig struct {
	ClientUUID          string `json:"client_uuid"`
	ErrorCode           int    `json:"error_code"`
	FailureCount        int    `json:"failure_count"`
	ActionAfterFailures string `json:"action_after_failures"`
	AfterOffset         int64  `json:"after_offset"`
}

type uploadSession struct {
	ID             string
	CurrentOffset  int64
	Buffer         bytes.Buffer
	Status         status
	Scenario       string
	ScenarioConfig ScenarioConfig
	UploadFailures int
	mu             sync.Mutex
}

func (sess *uploadSession) handleScenario(cmd string, w http.ResponseWriter, r *http.Request, offset int64) bool {
	switch sess.Scenario {
	case "non_fatal_error_on_query":
		return sess.nonFatalErrorOnQuery(cmd, w)
	case "non_fatal_error_on_chunk_upload":
		return sess.nonFatalErrorOnChunkUpload(cmd, w, offset)
	case "chunk_granularity":
		return sess.chunkGranularity(cmd, w, r)
	default:
		return false
	}
}

func (sess *uploadSession) nonFatalErrorOnQuery(cmd string, w http.ResponseWriter) bool {
	if cmd == "query" && sess.UploadFailures < sess.ScenarioConfig.FailureCount {
		sess.UploadFailures++
		sendError(w, sess.ScenarioConfig.ErrorCode, "Injected non-fatal query error", statusActive)
		return true
	}
	return false
}

func (sess *uploadSession) nonFatalErrorOnChunkUpload(cmd string, w http.ResponseWriter, offset int64) bool {
	if cmd == "upload" && offset >= sess.ScenarioConfig.AfterOffset {
		if sess.UploadFailures < sess.ScenarioConfig.FailureCount {
			sess.UploadFailures++
			sendError(w, sess.ScenarioConfig.ErrorCode, "Injected non-fatal chunk upload error", statusActive)
			return true
		}
		if sess.ScenarioConfig.ActionAfterFailures == "terminate" {
			sendError(w, http.StatusInternalServerError, "Scenario requested termination", "")
			return true
		}
	}
	return false
}

func (sess *uploadSession) chunkGranularity(cmd string, w http.ResponseWriter, r *http.Request) bool {
	if cmd == "upload" {
		commands := parseCommands(r.Header.Get("X-Goog-Upload-Command"))
		if !slices.Contains(commands, "finalize") {
			bodyLen := r.ContentLength
			if bodyLen > 0 && bodyLen%256 != 0 {
				sendError(w, http.StatusBadRequest, fmt.Sprintf("Chunk size %d is not a multiple of granularity 256", bodyLen), statusActive)
				return true
			}
		}
	}
	return false
}

type Manager struct {
	sessions        sync.Map
	startFailureMap sync.Map
	idGen           int64
}

// NewManager creates a new Resumable Upload session manager.
func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) incrementStartFailure(key, scenario string, limit int) int {
	if scenario != "non_fatal_error_on_start" {
		return limit + 1
	}
	val, _ := m.startFailureMap.LoadOrStore(key, new(int32))
	ptr := val.(*int32)
	current := int(atomic.AddInt32(ptr, 1))
	if current > limit {
		atomic.StoreInt32(ptr, 0)
	}
	return current
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
	if sess.handleScenario("query", w, nil, 0) {
		return
	}
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

	if sess.handleScenario("upload", w, r, offset) {
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
		if err := json.Unmarshal([]byte(cfgStr), &config); err != nil {
			sendError(w, http.StatusBadRequest, "Invalid JSON in X-Goog-Test-Scenario-Config header", "")
			return
		}
	}

	lookupKey := scenario + "-" + r.RemoteAddr
	if config.ClientUUID != "" {
		lookupKey = scenario + "-" + config.ClientUUID
	}

	failures := m.incrementStartFailure(lookupKey, scenario, config.FailureCount)
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
		Status:         statusActive,
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

	granularity := strconv.Itoa(256 * 1024)
	if scenario == "chunk_granularity" {
		granularity = "256"
	}
	w.Header().Set("X-Goog-Upload-Status", string(statusActive))
	w.Header().Set("X-Goog-Upload-URL", uploadURL)
	w.Header().Set("X-Goog-Upload-Control-URL", uploadURL)
	w.Header().Set("X-Goog-Upload-Chunk-Granularity", granularity)
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
