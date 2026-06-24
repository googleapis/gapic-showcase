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
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
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

func TestSecureConnect(t *testing.T) {
	// 1. Generate certs hermetically
	caPEM, certPath, keyPath := generateCerts(t)

	// 2. Setup Server Config using CreateAllEndpoints
	config := RuntimeConfig{
		port:         ":0", // let system choose free port
		fallbackPort: ":0", // avoid conflicts
		tlsCert:      certPath,
		tlsKey:       keyPath,
	}

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

	// 3. Setup Client trusting our in-memory CA
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

	// 4. Call Unary Echo and capture response headers (metadata)
	content := "TLS & PQC Test"
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

	// 5. Assertions on Negotiated TLS Parameters
	assertHeader(t, header, "x-showcase-tls-version", "TLS 1.3")
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

// generateCerts generates a self-signed CA and a server certificate signed by it.
// It writes them to temp files in t.TempDir() and returns the CA PEM bytes and paths to the server cert and key.
func generateCerts(t *testing.T) (caPEM []byte, certPath, keyPath string) {
	t.Helper()

	// 1. Generate CA Key and Cert
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("failed to generate CA key: %v", err)
	}

	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Showcase Test CA"},
			CommonName:   "Showcase Test CA",
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("failed to create CA certificate: %v", err)
	}

	caPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	caPEM = pem.EncodeToMemory(caPEMBlock)

	// 2. Generate Server Key and Cert signed by CA
	serverKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("failed to generate server key: %v", err)
	}

	serverTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Showcase Test Server"},
			CommonName:   "localhost",
		},
		NotBefore:   time.Now().Add(-1 * time.Hour),
		NotAfter:    time.Now().Add(24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}

	serverBytes, err := x509.CreateCertificate(rand.Reader, &serverTemplate, &caTemplate, &serverKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("failed to create server certificate: %v", err)
	}

	serverCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverBytes})
	serverKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})

	// 3. Write Server Cert and Key to Temp Files
	tempDir := t.TempDir()
	certPath = filepath.Join(tempDir, "server.crt")
	keyPath = filepath.Join(tempDir, "server.key")

	if err := os.WriteFile(certPath, serverCertPEM, 0600); err != nil {
		t.Fatalf("failed to write server cert: %v", err)
	}
	if err := os.WriteFile(keyPath, serverKeyPEM, 0600); err != nil {
		t.Fatalf("failed to write server key: %v", err)
	}

	return caPEM, certPath, keyPath
}

func TestConnectExistingServer(t *testing.T) {
	if os.Getenv("SHOWCASE_CONNECT_EXTERNAL") != "true" {
		t.Skip("Skipping TestConnectExistingServer: SHOWCASE_CONNECT_EXTERNAL=true not set")
	}

	addr := "localhost:7470"
	caPath := filepath.Join("..", "..", "certs", "ca.crt")

	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		t.Fatalf("failed to read CA cert from %s: %v. Make sure you generated certs using TLS_GUIDE.md", caPath, err)
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

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		t.Fatalf("failed to connect to server at %s: %v. Make sure server is running.", addr, err)
	}
	defer conn.Close()

	c, err := client.NewEchoClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("failed to create echo client: %v", err)
	}
	defer c.Close()

	content := "External TLS Test"
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

	assertHeader(t, header, "x-showcase-tls-version", "TLS 1.3")
}

