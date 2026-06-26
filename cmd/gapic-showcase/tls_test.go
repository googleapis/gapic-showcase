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
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	resp, err := c.Echo(ctx, req)
	if err != nil {
		t.Fatalf("Echo call failed: %v", err)
	}

	if resp.GetContent() != content {
		t.Errorf("expected %q, got %q", content, resp.GetContent())
	}
}
