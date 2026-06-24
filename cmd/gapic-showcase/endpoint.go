// Copyright 2020 Google LLC
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
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	iampb "cloud.google.com/go/iam/apiv1/iampb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/genrest"
	"github.com/googleapis/gapic-showcase/server/services"
	fallback "github.com/googleapis/grpc-fallback-go/server"
	gmux "github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	locpb "google.golang.org/genproto/googleapis/cloud/location"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RuntimeConfig has the run-time settings necessary to run the
// Showcase servers.
type RuntimeConfig struct {
	port         string
	fallbackPort string
	tlsCaCert    string
	tlsCert      string
	tlsKey       string
	autoTls      bool
	caCertFile   string
}

// Endpoint defines common operations for any of the various types of
// transport-specific network endpoints Showcase supports
type Endpoint interface {
	// Serve beings the listen-and-serve loop for this
	// Endpoint. It typically blocks until the server is shut
	// down. The error it returns depends on the underlying
	// implementation.
	Serve() error

	// Shutdown causes the currently running Endpoint to
	// terminate. The error it returns depends on the underlying
	// implementation.
	Shutdown() error

	// GetAddr returns the network address that this Endpoint is listening on.
	GetAddr() net.Addr
}

// generateInMemCerts generates a self-signed CA and a server certificate signed by it.
// It returns the PEM encoded CA cert, server cert, and server private key.
func generateInMemCerts() (caPEM, certPEM, keyPEM []byte, err error) {
	// 1. Generate CA Key and Cert
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate CA key: %w", err)
	}

	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Showcase Auto TLS CA"},
			CommonName:   "Showcase Auto TLS CA",
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
		return nil, nil, nil, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	caPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caBytes})

	// 2. Generate Server Key and Cert signed by CA
	serverKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate server key: %w", err)
	}

	serverTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Showcase Auto TLS Server"},
			CommonName:   "localhost",
		},
		NotBefore:   time.Now().Add(-1 * time.Hour),
		NotAfter:    time.Now().Add(24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	serverBytes, err := x509.CreateCertificate(rand.Reader, &serverTemplate, &caTemplate, &serverKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create server certificate: %w", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverBytes})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})

	return caPEM, certPEM, keyPEM, nil
}

func createTLSConfig(config RuntimeConfig) *tls.Config {
	var keyPair tls.Certificate
	var caPEM []byte
	var err error

	if config.autoTls {
		var certPEM, keyPEM []byte
		caPEM, certPEM, keyPEM, err = generateInMemCerts()
		if err != nil {
			log.Fatalf("Failed to generate in-memory TLS certs: %v", err)
		}
		keyPair, err = tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			log.Fatalf("Failed to load generated TLS keypair: %v", err)
		}
		stdLog.Printf("Automatically generated in-memory TLS certificates.")

		if config.caCertFile != "" {
			if err := os.WriteFile(config.caCertFile, caPEM, 0600); err != nil {
				log.Fatalf("Failed to write CA certificate to %s: %v", config.caCertFile, err)
			}
			stdLog.Printf("Wrote automatically generated CA certificate to %s", config.caCertFile)
		}
	} else {
		keyPair, err = tls.LoadX509KeyPair(config.tlsCert, config.tlsKey)
		if err != nil {
			log.Fatalf("Failed to load server TLS cert/key with error:%v", err)
		}
	}

	baseConfig := &tls.Config{
		Certificates: []tls.Certificate{keyPair},
		NextProtos:   []string{"h2", "http/1.1"},
	}



	// Handle Client CA for mTLS
	if config.autoTls {
		// For autoTls, we don't enforce client certs unless tlsCaCert is explicitly passed (which would be rare for autoTls)
		if config.tlsCaCert != "" {
			cert, err := os.ReadFile(config.tlsCaCert)
			if err != nil {
				log.Fatalf("Failed to load root CA cert file with error:%v", err)
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(cert)
			baseConfig.ClientCAs = pool
			baseConfig.ClientAuth = tls.RequireAndVerifyClientCert
			stdLog.Printf("Configured server with Mutual TLS (mTLS)")
		} else {
			baseConfig.ClientAuth = tls.NoClientCert
			stdLog.Printf("Configured server with One-Way TLS")
		}
	} else {
		if config.tlsCaCert != "" {
			cert, err := os.ReadFile(config.tlsCaCert)
			if err != nil {
				log.Fatalf("Failed to load root CA cert file with error:%v", err)
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(cert)
			baseConfig.ClientCAs = pool
			baseConfig.ClientAuth = tls.RequireAndVerifyClientCert
			stdLog.Printf("Configured server with Mutual TLS (mTLS)")
		} else {
			baseConfig.ClientAuth = tls.NoClientCert
			stdLog.Printf("Configured server with One-Way TLS")
		}
	}

	// Clone config per connection to bind RemoteAddr to VerifyConnection
	baseConfig.GetConfigForClient = func(info *tls.ClientHelloInfo) (*tls.Config, error) {
		remoteAddr := info.Conn.RemoteAddr().String()
		conf := baseConfig.Clone()
		conf.VerifyConnection = func(state tls.ConnectionState) error {
			server.RecordTLSHandshake(remoteAddr, state)
			stdLog.Printf("TLS Handshake Complete on Server for %s:", remoteAddr)
			stdLog.Printf("  Protocol: %s", tls.VersionName(state.Version))
			stdLog.Printf("  Cipher Suite: %s", tls.CipherSuiteName(state.CipherSuite))
			return nil
		}
		return conf, nil
	}

	return baseConfig
}

// CreateAllEndpoints returns an Endpoint that can serve gRPC and
// HTTP/REST connections (on config.port) and gRPC-fallback
// connections (on config.fallbackPort)
func CreateAllEndpoints(config RuntimeConfig) Endpoint {
	// Ensure port is of the right form.
	if !strings.HasPrefix(config.port, ":") {
		config.port = ":" + config.port
	}

	// Start listening.
	lis, err := net.Listen("tcp", config.port)
	if err != nil {
		log.Fatalf("Showcase failed to listen on port '%s': %v", config.port, err)
	}

	// 1. Setup TLS if enabled (either via certs or autoTls)
	isTLS := (config.tlsCert != "" && config.tlsKey != "") || config.autoTls
	if isTLS {
		tlsConfig := createTLSConfig(config)
		lis = tls.NewListener(lis, tlsConfig)
		stdLog.Printf("Showcase listening securely (TLS) on port: %s", config.port)
	} else {
		stdLog.Printf("Showcase listening insecurely (Plaintext) on port: %s", config.port)
	}

	// 2. Wrap in cleanup listener for connection tracking and leak prevention
	lis = &cleanupListener{Listener: lis}

	// 3. Get allocated port for logging
	addr := lis.Addr().(*net.TCPAddr)
	portStr := fmt.Sprintf("%d", addr.Port)

	scheme := "http"
	if isTLS {
		scheme = "https"
	}
	stdLog.Printf("gRPC Endpoint: %s://localhost:%s", scheme, portStr)
	stdLog.Printf("HTTP/REST Endpoint: %s://localhost:%s", scheme, portStr)

	m := cmux.New(lis)

	// Since TLS is decrypted at the listener level, we match on decrypted protocols.
	// HTTP1 for REST (HTTP/1.1), HTTP2 for gRPC.
	httpListener := m.Match(cmux.HTTP1())
	// Match HTTP2 for gRPC, fallback to Any if needed
	grpcListener := m.Match(cmux.HTTP2(), cmux.Any())

	backend := createBackends()

	// Pass a config copy with TLS disabled to gRPC because the listener already decrypted it.
	grpcConfig := config
	grpcConfig.tlsCert = ""
	grpcConfig.tlsKey = ""
	grpcConfig.autoTls = false
	grpcConfig.tlsCaCert = "" // also disable client CA at gRPC level

	gRPCServer := newEndpointGRPC(grpcListener, grpcConfig, backend)
	restServer := newEndpointREST(httpListener, backend)
	cmuxServer := newEndpointMux(lis, m, gRPCServer, restServer)
	return cmuxServer
}

// endpointMux is an Endpoint for cmux, the connection multiplexer
// allowing different types of connections on the same port.
//
// We choose not to use grpc.Server.ServeHTTP because it is
// experimental and does not support some gRPC features available
// through grpc.Server.Serve. (cf
// https://godoc.org/google.golang.org/grpc#Server.ServeHTTP)
type endpointMux struct {
	endpoints []Endpoint
	cmux      cmux.CMux
	mux       sync.Mutex
	listener  net.Listener
}

func newEndpointMux(lis net.Listener, cmuxEndpoint cmux.CMux, endpoints ...Endpoint) Endpoint {
	return &endpointMux{
		endpoints: endpoints,
		cmux:      cmuxEndpoint,
		listener:  lis,
	}
}

func (em *endpointMux) GetAddr() net.Addr {
	return em.listener.Addr()
}

func (em *endpointMux) String() string {
	return "endpoint multiplexer"
}

func (em *endpointMux) Serve() error {
	g := new(errgroup.Group)
	for idx, endpt := range em.endpoints {
		if endpt != nil {
			stdLog.Printf("Starting endpoint %d: %s", idx, endpt)
			endpoint := endpt
			g.Go(func() error {
				err := endpoint.Serve()
				err2 := em.Shutdown()
				if err != nil {
					return err
				}
				return err2
			})
		}
	}
	if em.cmux != nil {
		stdLog.Printf("Starting %s", em)

		g.Go(func() error {
			err := em.cmux.Serve()
			err2 := em.Shutdown()
			if err != nil {
				return err
			}
			return err2

		})
	}
	return g.Wait()
}

func (em *endpointMux) Shutdown() error {
	em.mux.Lock()
	defer em.mux.Unlock()

	var err error
	if em.cmux != nil {
		// TODO: Wait for https://github.com/soheilhy/cmux/pull/69 (due to
		// https://github.com/soheilhy/cmux/pull/69#issuecomment-712928041.)
		//
		// err = em.mux.Close()
		em.cmux = nil
	}

	for idx, endpt := range em.endpoints {
		if endpt != nil {
			// TODO: Wait for https://github.com/soheilhy/cmux/pull/69
			// newErr := endpt.Shutdown()
			// if err==nil {
			// 	err = newErr
			// }
			em.endpoints[idx] = nil
		}
	}
	return err
}

// endpointGRPC is an Endpoint for gRPC connections to the Showcase
// server.
type endpointGRPC struct {
	server         *grpc.Server
	fallbackServer *fallback.FallbackServer
	listener       net.Listener
	mux            sync.Mutex
}

// createBackends creates services used by both the gRPC and REST servers.
func createBackends() *services.Backend {
	logger := &loggerObserver{}
	observerRegistry := server.ShowcaseObserverRegistry()
	observerRegistry.RegisterUnaryObserver(logger)
	observerRegistry.RegisterStreamRequestObserver(logger)
	observerRegistry.RegisterStreamResponseObserver(logger)

	identityServer := services.NewIdentityServer()
	messagingServer := services.NewMessagingServer(identityServer)
	return &services.Backend{
		EchoServer:            services.NewEchoServer(),
		SequenceServiceServer: services.NewSequenceServer(),
		IdentityServer:        identityServer,
		MessagingServer:       messagingServer,
		ComplianceServer:      services.NewComplianceServer(),
		TestingServer:         services.NewTestingServer(observerRegistry),
		OperationsServer:      services.NewOperationsServer(messagingServer),
		LocationsServer:       services.NewLocationsServer(),
		IAMPolicyServer:       services.NewIAMPolicyServer(),
		StdLog:                stdLog,
		ErrLog:                errLog,
		ObserverRegistry:      observerRegistry,
	}
}

func newEndpointGRPC(lis net.Listener, config RuntimeConfig, backend *services.Backend) Endpoint {
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(backend.ObserverRegistry.StreamInterceptor),
		grpc.ChainUnaryInterceptor(
			backend.ObserverRegistry.UnaryInterceptor,
			server.TLSMetadataUnaryInterceptor,
		),
	}

	s := grpc.NewServer(opts...)

	// Register Services to the server.
	pb.RegisterEchoServer(s, backend.EchoServer)
	pb.RegisterSequenceServiceServer(s, backend.SequenceServiceServer)
	pb.RegisterIdentityServer(s, backend.IdentityServer)
	pb.RegisterMessagingServer(s, backend.MessagingServer)
	pb.RegisterComplianceServer(s, backend.ComplianceServer)
	pb.RegisterTestingServer(s, backend.TestingServer)
	lropb.RegisterOperationsServer(s, backend.OperationsServer)
	locpb.RegisterLocationsServer(s, backend.LocationsServer)
	iampb.RegisterIAMPolicyServer(s, backend.IAMPolicyServer)

	fb := fallback.NewServer(config.fallbackPort, "localhost"+config.port)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	return &endpointGRPC{
		server:         s,
		fallbackServer: fb,
		listener:       lis,
	}
}

func (eg *endpointGRPC) String() string {
	return "gRPC endpoint"
}

func (eg *endpointGRPC) Serve() error {
	defer eg.Shutdown()
	if eg.fallbackServer != nil {
		stdLog.Printf("Listening for gRPC-fallback connections")
		eg.fallbackServer.StartBackground()
	}
	if eg.server != nil {
		stdLog.Printf("Listening for gRPC connections")
		return eg.server.Serve(eg.listener)
	}
	return fmt.Errorf("gRPC server not set up")
}

func (eg *endpointGRPC) Shutdown() error {
	eg.mux.Lock()
	defer eg.mux.Unlock()

	if eg.fallbackServer != nil {
		stdLog.Printf("Stopping gRPC-fallback connections")
		eg.fallbackServer.Shutdown()
		eg.fallbackServer = nil
	}

	if eg.server != nil {
		stdLog.Printf("Stopping gRPC connections")
		eg.server.GracefulStop()
		eg.server = nil
	}
	stdLog.Printf("Stopped gRPC")
	return nil
}

func (eg *endpointGRPC) GetAddr() net.Addr {
	return eg.listener.Addr()
}

// endpointREST is an Endpoint for HTTP/REST connections to the Showcase
// server.
type endpointREST struct {
	server   *http.Server
	listener net.Listener
	mux      sync.Mutex
}

func newEndpointREST(lis net.Listener, backend *services.Backend) *endpointREST {
	router := gmux.NewRouter()
	router.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("GAPIC Showcase: HTTP/REST endpoint using gorilla/mux\n"))
	})
	genrest.RegisterHandlers(router, backend)

	// Register TLS HTTP Middleware
	router.Use(server.TLSHTTPMiddleware)

	return &endpointREST{
		server:   &http.Server{Handler: router},
		listener: lis,
	}
}

func (er *endpointREST) String() string {
	return "HTTP/REST endpoint"
}

func (er *endpointREST) Serve() error {
	defer er.Shutdown()
	if er.server != nil {
		stdLog.Printf("Listening for REST connections")
		return er.server.Serve(er.listener)
	}
	return fmt.Errorf("REST server not set up")
}

func (er *endpointREST) Shutdown() error {
	er.mux.Lock()
	defer er.mux.Unlock()
	var err error
	if er.server != nil {
		stdLog.Printf("Stopping REST connections")
		err = er.server.Shutdown(context.Background())
		er.server = nil
	}
	stdLog.Printf("Stopped REST")
	return err
}

func (er *endpointREST) GetAddr() net.Addr {
	return er.listener.Addr()
}

type cleanupConn struct {
	net.Conn
}

func (c *cleanupConn) Close() error {
	server.RemoveTLSState(c.RemoteAddr().String())
	return c.Conn.Close()
}

type cleanupListener struct {
	net.Listener
}

func (l *cleanupListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &cleanupConn{Conn: c}, nil
}
