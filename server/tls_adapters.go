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
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// ConnectionState stores the TLS handshake details for a connection.
type ConnectionState struct {
	Version      uint16
	Cipher       uint16
	CurveID      tls.CurveID
	ClientCurves []tls.CurveID
}

var (
	tlsRegistryMu sync.RWMutex
	tlsRegistry   = make(map[string]*ConnectionState)
)

// RecordTLSHandshake stores both the negotiated state and the client's offered curves.
func RecordTLSHandshake(remoteAddr string, state tls.ConnectionState, clientCurves []tls.CurveID) {
	tlsRegistryMu.Lock()
	defer tlsRegistryMu.Unlock()
	tlsRegistry[remoteAddr] = &ConnectionState{
		Version:      state.Version,
		Cipher:       state.CipherSuite,
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
		names = append(names, TLSGroupName(g))
	}
	return strings.Join(names, ",")
}

// TLSGroupName maps a CurveID to a human-readable name, highlighting PQC hybrid groups.
func TLSGroupName(id tls.CurveID) string {
	switch id {
	case 4588:
		return "X25519MLKEM768"
	case 4587:
		return "SecP256r1MLKEM768"
	case tls.X25519:
		return "X25519"
	case tls.CurveP256:
		return "CurveP256"
	case tls.CurveP384:
		return "CurveP384"
	case tls.CurveP521:
		return "CurveP521"
	default:
		return fmt.Sprintf("Unknown-Curve-%d", id)
	}
}

// TLSMetadataUnaryInterceptor injects TLS handshake details into unary gRPC response headers.
func TLSMetadataUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if p, ok := peer.FromContext(ctx); ok {
		if state := GetTLSState(p.Addr.String()); state != nil {
			md := metadata.Pairs(
				"x-showcase-tls-version", tls.VersionName(state.Version),
				"x-showcase-tls-cipher", tls.CipherSuiteName(state.Cipher),
				"x-showcase-tls-group", TLSGroupName(state.CurveID),
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
			w.Header().Set("x-showcase-tls-version", tls.VersionName(state.Version))
			w.Header().Set("x-showcase-tls-cipher", tls.CipherSuiteName(state.Cipher))
			w.Header().Set("x-showcase-tls-group", TLSGroupName(state.CurveID))
			w.Header().Set("x-showcase-tls-client-supported-groups", formatGroups(state.ClientCurves))
		}
		next.ServeHTTP(w, r)
	})
}
