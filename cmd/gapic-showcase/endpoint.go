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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/genrest"
	"github.com/googleapis/gapic-showcase/server/services"
	fallback "github.com/googleapis/grpc-fallback-go/server"
	gmux "github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	locpb "google.golang.org/genproto/googleapis/cloud/location"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	lropb "google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	stdLog.Printf("Showcase listening on port: %s", config.port)

	m := cmux.New(lis)
	httpListener := m.Match(cmux.HTTP1())
	// cmux.Any() is needed below to get mTLS to work for
	// gRPC, and that in turn means the order of the matchers matters. See
	// https://github.com/open-telemetry/opentelemetry-collector/issues/2732
	grpcListener := m.Match(cmux.Any())

	backend := createBackends()
	gRPCServer := newEndpointGRPC(grpcListener, config, backend)
	restServer := newEndpointREST(httpListener, backend)
	cmuxServer := newEndpointMux(m, gRPCServer, restServer)
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
}

func newEndpointMux(cmuxEndpoint cmux.CMux, endpoints ...Endpoint) Endpoint {
	return &endpointMux{
		endpoints: endpoints,
		cmux:      cmuxEndpoint,
	}
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
		grpc.UnaryInterceptor(backend.ObserverRegistry.UnaryInterceptor),
	}

	// load mutual TLS cert/key and root CA cert
	if config.tlsCaCert != "" && config.tlsCert != "" && config.tlsKey != "" {
		keyPair, err := tls.LoadX509KeyPair(config.tlsCert, config.tlsKey)
		if err != nil {
			log.Fatalf("Failed to load server TLS cert/key with error:%v", err)
		}

		cert, err := ioutil.ReadFile(config.tlsCaCert)
		if err != nil {
			log.Fatalf("Failed to load root CA cert file with error:%v", err)
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(cert)

		ta := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{keyPair},
			ClientCAs:    pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		})

		opts = append(opts, grpc.Creds(ta))
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
