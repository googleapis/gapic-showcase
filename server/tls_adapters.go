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

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// ConnectionState stores the TLS handshake details for a connection.
type ConnectionState struct {
	CurveID      tls.CurveID
	ClientCurves []tls.CurveID
}

var (
	tlsRegistryMu sync.RWMutex
	tlsRegistry   = make(map[string]*ConnectionState)
)

// RecordTLSHandshake stores the negotiated state and the client's offered curves.
func RecordTLSHandshake(remoteAddr string, state tls.ConnectionState, clientCurves []tls.CurveID) {
	tlsRegistryMu.Lock()
	defer tlsRegistryMu.Unlock()
	tlsRegistry[remoteAddr] = &ConnectionState{
		CurveID:      state.CurveID,
		ClientCurves: clientCurves,
	}
}

// GetTLSState retrieves the recorded TLS state for a connection.
func GetTLSState(remoteAddr string) *ConnectionState {
	tlsRegistryMu.RLock()
	defer tlsRegistryMu.RUnlock()
	return tlsRegistry[remoteAddr]
}

// RemoveTLSState deletes the connection record to prevent memory leaks.
func RemoveTLSState(remoteAddr string) {
	tlsRegistryMu.Lock()
	defer tlsRegistryMu.Unlock()
	delete(tlsRegistry, remoteAddr)
}

// Helper to format CurveIDs into a comma-separated string
func formatGroups(groups []tls.CurveID) string {
	var names []string
	for _, g := range groups {
		names = append(names, g.String())
	}
	return strings.Join(names, ",")
}

// ParseTLSGroup parses a single TLS key exchange group ID represented as a hex (e.g. "0x11ec")
// or decimal (e.g. "4588") integer into a tls.CurveID.
func ParseTLSGroup(raw string) (tls.CurveID, error) {
	raw = strings.TrimSpace(raw)
	val, err := strconv.ParseUint(raw, 0, 16)
	if err != nil || val == 0 {
		return 0, fmt.Errorf("invalid TLS group ID %q: expected non-zero hex (e.g. 0x11ec) or decimal (e.g. 4588) IANA group ID", raw)
	}
	return tls.CurveID(val), nil
}

// ParseTLSGroups parses a comma-separated list of TLS key exchange group IDs.
func ParseTLSGroups(csv string) ([]tls.CurveID, error) {
	if strings.TrimSpace(csv) == "" {
		return nil, nil
	}
	parts := strings.Split(csv, ",")
	groups := make([]tls.CurveID, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		g, err := ParseTLSGroup(p)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// TLSMetadataUnaryInterceptor injects TLS handshake details into unary gRPC response headers.
func TLSMetadataUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if p, ok := peer.FromContext(ctx); ok {
		if state := GetTLSState(p.Addr.String()); state != nil {
			md := metadata.Pairs(
				"x-showcase-tls-group", state.CurveID.String(),
				"x-showcase-tls-client-supported-groups", formatGroups(state.ClientCurves),
			)
			_ = grpc.SetHeader(ctx, md)
		}
	}
	return handler(ctx, req)
}

// TLSHTTPMiddleware injects TLS handshake details into HTTP response headers.
func TLSHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if state := GetTLSState(r.RemoteAddr); state != nil {
			w.Header().Set("x-showcase-tls-group", state.CurveID.String())
			w.Header().Set("x-showcase-tls-client-supported-groups", formatGroups(state.ClientCurves))
		}
		next.ServeHTTP(w, r)
	})
}
