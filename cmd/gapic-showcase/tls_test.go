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

func TestConnectWithTLS(t *testing.T) {
	// Setup Server Config using Auto-TLS
	caPath := filepath.Join(t.TempDir(), "ca.crt")
	config := RuntimeConfig{
		port:         ":0", // let system choose free port
		fallbackPort: ":0", // avoid conflicts
		autoTLS:      true,
		caCertFile:   caPath,
	}

	// CreateAllEndpoints synchronously generates certs and writes caCertFile when autoTLS is true
	serverEndpoint := CreateAllEndpoints(config)

	// Start server in background
	go func() {
		if err := serverEndpoint.Serve(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("server exited with error: %v", err)
		}
	}()
	defer serverEndpoint.Shutdown()

	// Get resolved port
	addr := fmt.Sprintf("localhost:%d", serverEndpoint.GetAddr().(*net.TCPAddr).Port)
	t.Logf("Server resolved to address: %s", addr)

	// Load the synchronously generated CA cert
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		t.Fatalf("failed to read generated CA cert: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		t.Fatalf("failed to append CA cert")
	}

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

	// Call Unary Echo
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
	// Setup Server Config using Auto-TLS and PQC Disabled
	caPath := filepath.Join(t.TempDir(), "ca.crt")
	config := RuntimeConfig{
		port:         ":0", // let system choose free port
		fallbackPort: ":0", // avoid conflicts
		autoTLS:      true,
		caCertFile:   caPath,
		enablePQC:    boolPtr(false), // Disable PQC
	}

	// CreateAllEndpoints synchronously generates certs and writes caCertFile when autoTLS is true
	serverEndpoint := CreateAllEndpoints(config)

	// Start server in background
	go func() {
		if err := serverEndpoint.Serve(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("server exited with error: %v", err)
		}
	}()
	defer serverEndpoint.Shutdown()

	addr := fmt.Sprintf("localhost:%d", serverEndpoint.GetAddr().(*net.TCPAddr).Port)
	t.Logf("Server resolved to address: %s", addr)

	// Load the synchronously generated CA cert
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		t.Fatalf("failed to read generated CA cert: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		t.Fatalf("failed to append CA cert")
	}

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
