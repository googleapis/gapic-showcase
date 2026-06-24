# GAPIC Showcase: TLS Guide

This guide explains how to configure, run, and connect to the GAPIC Showcase server using **One-Way TLS** (server verification only) and **Mutual TLS (mTLS)**.

---

## 1. Certificate Generation (For Manual Testing)

To use TLS, you need certificates. For local testing, you can generate a self-signed Certificate Authority (CA) and a server certificate.

Modern Go versions require **Subject Alternative Names (SANs)** in certificates; otherwise, they will reject them with `x509: certificate relies on legacy Common Name field`.

### Step 1.1: Create a SAN Configuration File
Create a file named `ext.conf` with the following content to define the SANs:
```ini
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
```

### Step 1.2: Generate the Certificates
Run the following commands to generate the CA and the server certificate signed by that CA:

```sh
# 1. Generate CA private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout ca.key -out ca.crt -days 365 -nodes -subj "/CN=ShowcaseCA"

# 2. Generate Server private key and Certificate Signing Request (CSR)
openssl req -newkey rsa:4096 -keyout server.key -out server.csr -nodes -subj "/CN=localhost"

# 3. Sign the Server CSR with the CA, applying the SAN extension
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile ext.conf
```

> [!NOTE]
> The `.gitignore` file is configured to exclude generated `*.crt` and `*.key` files to prevent them from being committed.

---

## 2. Running the Server

First, build the Showcase binary:
```sh
go build ./cmd/gapic-showcase
```

You can then run the Showcase server in one of the following TLS modes.

### Option A: One-Way TLS (Server Verification Only)
In this mode, the server presents its certificate to the client, but does not require the client to present one.

To start the server, provide only the server cert and key:
```sh
./gapic-showcase run \
  --tls-cert certs/server.crt \
  --tls-key certs/server.key \
  --port 7470
```
*The server log will confirm:* `Configured server with One-Way TLS`

### Option B: Mutual TLS (mTLS)
In this mode, the server presents its certificate AND requires the client to present a valid certificate signed by the configured CA.

To start the server, provide the server cert, key, **and the CA cert**:
```sh
./gapic-showcase run \
  --tls-cert certs/server.crt \
  --tls-key certs/server.key \
  --tls-ca-cert certs/ca.crt \
  --port 7470
```
*The server log will confirm:* `Configured server with Mutual TLS (mTLS)`

### Option C: Auto-TLS (Zero-Config / CI Friendly)
In this mode, the server automatically generates its own CA and server certificates in-memory at startup. This is the recommended mode for integration tests, as it eliminates the need to manage certificate files.

You can use a random port (`:0`) to avoid conflicts. The server will print the allocated port to the logs. You can also write the automatically generated CA certificate to a file so clients can load it:

```sh
./gapic-showcase run \
  --port :0 \
  --tls \
  --ca-cert-output-file showcase.pem
```

*The server log will print the resolved endpoints, for example:*
```
Showcase listening securely (TLS) on port: :0
gRPC Endpoint: https://localhost:45917
HTTP/REST Endpoint: https://localhost:45917
```

---

## 3. Connecting with a Go Client

### One-Way TLS Client Example
To connect to a server running One-Way TLS, the client needs to trust the server's CA.

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/googleapis/gapic-showcase/client"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Load the CA certificate that signed the server's certificate
	pemServerCA, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatalf("failed to read CA cert: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		log.Fatalf("failed to append CA cert")
	}

	// Configure TLS to verify the server using the custom CA pool
	config := &tls.Config{
		RootCAs: certPool,
	}
	creds := credentials.NewTLS(config)

	// Dial the server (adjust port as needed)
	conn, err := grpc.Dial("localhost:7470", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Initialize Showcase Client
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	// Make a call
	resp, err := c.Echo(ctx, &pb.EchoRequest{
		Response: &pb.EchoRequest_Content{Content: "Hello One-Way TLS!"},
	})
	if err != nil {
		log.Fatalf("echo failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp.GetContent())
}
```

---

## 4. Exposed TLS Response Metadata (Headers)

When a client connects securely, the Showcase server automatically injects the following metadata into the gRPC response headers (and HTTP headers):

*   **`x-showcase-tls-version`**: The negotiated TLS version (e.g., `TLS 1.3`).
*   **`x-showcase-tls-cipher`**: The negotiated cipher suite (e.g., `TLS_AES_128_GCM_SHA256`).

These headers allow asserting the negotiated TLS parameters.

---

## 5. Hermetic Go Integration Testing

For automated CI environments, a self-contained integration test is available that does not rely on external certificate files or fixed ports.

The test is located at [cmd/gapic-showcase/secure_test.go](file:///usr/local/google/home/hongalex/code/cloud/gapic-showcase/cmd/gapic-showcase/secure_test.go).

### How it works:
1.  **In-Memory Certs**: It dynamically generates a self-signed CA and server certificate on the fly.
2.  **Random Port**: It starts the gRPC server on a random free port (`localhost:0`).
3.  **Address Discovery**: It retrieves the allocated port dynamically using the `GetAddr()` method on the server endpoint, avoiding the need for temporary port files.
4.  **Assertions**: It asserts that the returned headers contain the correct TLS version.

### Running the Test:
```sh
go test -v -run TestSecureConnect ./cmd/gapic-showcase
```
