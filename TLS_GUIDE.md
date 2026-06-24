# GAPIC Showcase: TLS & Post-Quantum Cryptography (PQC) Guide

This guide explains how to configure, run, and connect to the GAPIC Showcase server using TLS.

For most use cases (including verifying PQC), **Auto-TLS** is the recommended mode as it requires zero configuration. Generating certificates manually via OpenSSL is only necessary if you need to test **Mutual TLS (mTLS)**.

## 1. Running the Server with Auto-TLS (Recommended)

With the `--tls` flag, the server automatically generates its own CA and server certificates in-memory at startup. This is the recommended mode for local testing and CI integration for One-Way TLS.

You should use --ca-cert-output-file to write the automatically generated CA certificate to a file to be used by the client test.

Port `:0` automatically assigns a free port, but will need to be read from the logs. 

```sh
./gapic-showcase run \
  --port :0 \
  --tls \
  --ca-cert-output-file showcase.pem
```

*The server log will print the resolved endpoints:*
```
gRPC Endpoint (TLS): https://localhost:45917
HTTP/REST Endpoint (TLS): https://localhost:45917
```

The client needs to load the generated `showcase.pem` file to verify the server connection.
An example of this can be found in cmd/gapic-showcase/tls_test.go.

## 2. Running the Server with Manual Certificates (mTLS Only)

Manual certificate generation is required for Mutual TLS (mTLS) because the client certificate must be signed by the CA trusted by the server. Since Auto-TLS generates the CA key in-memory and does not expose the private key, you cannot sign client certificates with it.

### Step 2.1: Certificate Generation
Create a SAN configuration file named `ext.conf`:
```ini
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
```

Generate the CA and server certificates:
```sh
# 1. Generate CA private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout ca.key -out ca.crt -days 365 -nodes -subj "/CN=ShowcaseCA"

# 2. Generate Server private key and CSR
openssl req -newkey rsa:4096 -keyout server.key -out server.csr -nodes -subj "/CN=localhost"

# 3. Sign the Server CSR with the CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile ext.conf
```
*(Note: For a full mTLS test, you would also generate a client private key and CSR, and sign it using `ca.key` and `ca.crt`)*

### Step 2.2: Start the Server in mTLS Mode
Provide the server cert, key, **and the CA cert** (which the server will use to verify client certificates):

```sh
./gapic-showcase run \
  --tls-cert server.crt \
  --tls-key server.key \
  --tls-ca-cert ca.crt \
  --port 7470
```
*The server log will confirm:* `Configured server with Mutual TLS (mTLS)`

## 3. Verifying Post-Quantum Cryptography (PQC)

When running on **Go 1.24+**, the hybrid post-quantum key exchange **`X25519MLKEM768`** is enabled by default on the Showcase server.

### Disabling PQC (Server-Side)

If you need to verify classical fallback behavior, you can force the Showcase server to disable all Post-Quantum hybrid key exchanges and use only classical cryptography by starting the server with the `--enable-pqc=false` flag:

```sh
./gapic-showcase run \
  --tls-cert certs/server.crt \
  --tls-key certs/server.key \
  --enable-pqc=false
```

## 4. Exposed TLS Response Metadata (Headers)

When a client connects securely, the Showcase server automatically injects the following metadata into the gRPC response headers (and HTTP headers):

*   **`x-showcase-tls-group`**: The negotiated key-exchange group (e.g., `X25519MLKEM768`).
*   **`x-showcase-tls-client-supported-groups`**: A comma-separated list of all key-exchange groups the client offered in its `ClientHello` handshake, ordered by the client's preference.
