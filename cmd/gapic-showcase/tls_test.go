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

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/googleapis/gapic-showcase/client"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func setupTLSTestServer(t *testing.T, enablePQC *bool) (addr string, certPool *x509.CertPool, cleanup func()) {
	t.Helper()
	caPath := filepath.Join(t.TempDir(), "ca.crt")
	config := RuntimeConfig{
		port:         ":0", // let system choose free port
		fallbackPort: ":0", // avoid conflicts
		autoTLS:      true,
		caCertFile:   caPath,
		enablePQC:    enablePQC,
	}

	serverEndpoint := CreateAllEndpoints(config)

	go func() {
		if err := serverEndpoint.Serve(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("server exited with error: %v", err)
		}
	}()

	addr = fmt.Sprintf("localhost:%d", serverEndpoint.GetAddr().(*net.TCPAddr).Port)

	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		serverEndpoint.Shutdown()
		t.Fatalf("failed to read generated CA cert: %v", err)
	}

	certPool = x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		serverEndpoint.Shutdown()
		t.Fatalf("failed to append CA cert")
	}

	return addr, certPool, func() { serverEndpoint.Shutdown() }
}

func TestConnectWithTLS(t *testing.T) {
	t.Parallel()
	addr, certPool, cleanup := setupTLSTestServer(t, nil)
	defer cleanup()

	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial the server
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		t.Fatalf("failed to connect to server at %s: %v", addr, err)
	}
	defer conn.Close()

	c, err := client.NewEchoClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("failed to create echo client: %v", err)
	}
	defer c.Close()

	content := "TLS Test using Auto-TLS"
	req := &pb.EchoRequest{
		Response: &pb.EchoRequest_Content{
			Content: content,
		},
	}

	var header metadata.MD
	resp, err := c.Echo(ctx, req, gax.WithGRPCOptions(grpc.Header(&header)))
	if err != nil {
		t.Fatalf("Echo call failed: %v", err)
	}

	if resp.GetContent() != content {
		t.Errorf("expected %q, got %q", content, resp.GetContent())
	}

	// Assertions on Negotiated TLS Parameters
	assertHeader(t, header, "x-showcase-tls-group", "X25519MLKEM768") // Default Go 1.25 PQC group

	clientGroups := header.Get("x-showcase-tls-client-supported-groups")
	if len(clientGroups) == 0 {
		t.Fatalf("missing header: x-showcase-tls-client-supported-groups")
	}
	t.Logf("Client supported groups received: %s", clientGroups[0])

	if !strings.Contains(clientGroups[0], "X25519MLKEM768") {
		t.Errorf("expected client supported groups to contain X25519MLKEM768, got %q", clientGroups[0])
	}
}

func TestConnectWithTLS_REST(t *testing.T) {
	t.Parallel()
	addr, certPool, cleanup := setupTLSTestServer(t, nil)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}
	capturer := &headerCapturer{transport: tr}
	hc := &http.Client{Transport: capturer}

	restClient, err := client.NewEchoRESTClient(ctx,
		option.WithHTTPClient(hc),
		option.WithEndpoint("https://"+addr),
	)
	if err != nil {
		t.Fatalf("failed to create echo REST client: %v", err)
	}
	defer restClient.Close()

	content := "TLS REST Test using Auto-TLS"
	req := &pb.EchoRequest{
		Response: &pb.EchoRequest_Content{
			Content: content,
		},
	}

	respREST, err := restClient.Echo(ctx, req)
	if err != nil {
		t.Fatalf("REST Echo call failed: %v", err)
	}
	if respREST.GetContent() != content {
		t.Errorf("REST expected %q, got %q", content, respREST.GetContent())
	}

	assertRESTHeader(t, capturer.headers, "x-showcase-tls-group", "X25519MLKEM768")

	restClientGroups := capturer.headers.Get("x-showcase-tls-client-supported-groups")
	if restClientGroups == "" {
		t.Errorf("REST: missing header: x-showcase-tls-client-supported-groups")
	} else if !strings.Contains(restClientGroups, "X25519MLKEM768") {
		t.Errorf("REST: expected client supported groups to contain X25519MLKEM768, got %q", restClientGroups)
	}
}

func assertHeader(t *testing.T, md metadata.MD, key, expected string) {
	t.Helper()
	val := md.Get(key)
	if len(val) == 0 {
		t.Errorf("missing header: %s", key)
		return
	}
	t.Logf("Header %s: %s", key, val[0])
	if val[0] != expected {
		t.Errorf("header %s: expected %q, got %q", key, expected, val[0])
	}
}

func TestConnectWithTLS_PQCDisabled(t *testing.T) {
	t.Parallel()
	addr, certPool, cleanup := setupTLSTestServer(t, boolPtr(false))
	defer cleanup()

	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial the server
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		t.Fatalf("failed to connect to server at %s: %v", addr, err)
	}
	defer conn.Close()

	c, err := client.NewEchoClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("failed to create echo client: %v", err)
	}
	defer c.Close()

	content := "TLS Test PQC Disabled"
	req := &pb.EchoRequest{
		Response: &pb.EchoRequest_Content{
			Content: content,
		},
	}

	// Call Unary Echo
	var header metadata.MD
	resp, err := c.Echo(ctx, req, gax.WithGRPCOptions(grpc.Header(&header)))
	if err != nil {
		t.Fatalf("Echo call failed: %v", err)
	}

	if resp.GetContent() != content {
		t.Errorf("expected %q, got %q", content, resp.GetContent())
	}

	assertHeader(t, header, "x-showcase-tls-group", "X25519")
}

func boolPtr(b bool) *bool {
	return &b
}

type headerCapturer struct {
	transport http.RoundTripper
	headers   http.Header
}

func (h *headerCapturer) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := h.transport.RoundTrip(req)
	if err == nil {
		h.headers = resp.Header
	}
	return resp, err
}

func assertRESTHeader(t *testing.T, headers http.Header, key, expected string) {
	t.Helper()
	val := headers.Get(key)
	if val == "" {
		t.Errorf("REST: missing header: %s", key)
		return
	}
	t.Logf("REST Header %s: %s", key, val)
	if val != expected {
		t.Errorf("REST header %s: expected %q, got %q", key, expected, val)
	}
}
