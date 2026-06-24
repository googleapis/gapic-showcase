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
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// ConnectionState stores the TLS handshake details for a connection.
type ConnectionState struct {
	Version uint16
	Cipher  uint16
}

var (
	tlsRegistryMu sync.RWMutex
	tlsRegistry   = make(map[string]*ConnectionState)
)

// RecordTLSHandshake stores the negotiated TLS state.
func RecordTLSHandshake(remoteAddr string, state tls.ConnectionState) {
	tlsRegistryMu.Lock()
	defer tlsRegistryMu.Unlock()
	tlsRegistry[remoteAddr] = &ConnectionState{
		Version: state.Version,
		Cipher:  state.CipherSuite,
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

// TLSMetadataUnaryInterceptor injects TLS handshake details into unary gRPC response headers.
func TLSMetadataUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if p, ok := peer.FromContext(ctx); ok {
		if state := GetTLSState(p.Addr.String()); state != nil {
			md := metadata.Pairs(
				"x-showcase-tls-version", tls.VersionName(state.Version),
				"x-showcase-tls-cipher", tls.CipherSuiteName(state.Cipher),
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
		}
		next.ServeHTTP(w, r)
	})
}
