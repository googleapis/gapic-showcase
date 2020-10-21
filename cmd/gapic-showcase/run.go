// Copyright 2018 Google LLC
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
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/services"
	fallback "github.com/googleapis/grpc-fallback-go/server"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	lropb "google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func grpcServe(lis net.Listener, config configuration) error {
	// Setup Server.
	logger := &loggerObserver{}
	observerRegistry := server.ShowcaseObserverRegistry()
	observerRegistry.RegisterUnaryObserver(logger)
	observerRegistry.RegisterStreamRequestObserver(logger)
	observerRegistry.RegisterStreamResponseObserver(logger)

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(observerRegistry.StreamInterceptor),
		grpc.UnaryInterceptor(observerRegistry.UnaryInterceptor),
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
	defer s.GracefulStop()

	// Creates services used by both the gRPC and REST servers
	echoServer := services.NewEchoServer()

	// Note: we may be able to share a port
	// for handling gRPC and a regular HTTP
	// port!
	// https://godoc.org/google.golang.org/grpc#Server.ServeHTTP
	// This has many limitations and is
	// experimental, however. A better option:
	// issue a 307 in the specified port if
	// getting gRPC (as per the detection
	// algorithm in the URL above)
	// https://tools.ietf.org/html/rfc2616#page-65
	// https://en.wikipedia.org/wiki/HTTP_302

	// USE THIS: https://medium.com/@drgarcia1986/listen-grpc-and-http-requests-on-the-same-port-263c40cb45ff
	// which employs cmux: https://github.com/soheilhy/cmux

	// Register Services to the server.
	pb.RegisterEchoServer(s, echoServer)
	pb.RegisterSequenceServiceServer(s, services.NewSequenceServer())
	identityServer := services.NewIdentityServer()
	pb.RegisterIdentityServer(s, identityServer)
	messagingServer := services.NewMessagingServer(identityServer)
	pb.RegisterMessagingServer(s, messagingServer)
	operationsServer := services.NewOperationsServer(messagingServer)
	pb.RegisterTestingServer(s, services.NewTestingServer(observerRegistry))
	lropb.RegisterOperationsServer(s, operationsServer)

	fb := fallback.NewServer(config.fallbackPort, "localhost"+config.port)
	fb.StartBackground()
	defer fb.Shutdown()

	// Register reflection service on gRPC server.
	reflection.Register(s)
	stdLog.Printf("  listening for gRPC connections")
	return s.Serve(lis)
}

func httpServe(lis net.Listener, config configuration) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("GAPIC Showcase: HTTP/REST endpoint\n"))
	})

	s := &http.Server{Handler: mux}
	stdLog.Printf("  listening for REST connections")
	return s.Serve(lis)
}

type configuration struct {
	port         string
	fallbackPort string
	tlsCaCert    string
	tlsCert      string
	tlsKey       string
}

func init() {
	config := configuration{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the showcase server",
		Run: func(cmd *cobra.Command, args []string) {
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
			grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
			httpListener := m.Match(cmux.HTTP1Fast())

			g := new(errgroup.Group)
			g.Go(func() error { return grpcServe(grpcListener, config) })
			g.Go(func() error { return httpServe(httpListener, config) })
			g.Go(func() error { return m.Serve() })

			stdLog.Printf("after running server: %x\n", g.Wait())

		},
	}
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(
		&config.port,
		"port",
		"p",
		":7469",
		"The port that showcase will be served on.")
	runCmd.Flags().StringVarP(
		&config.fallbackPort,
		"fallback-port",
		"f",
		":1337",
		"The port that the fallback-proxy will be served on.")
	runCmd.Flags().StringVar(
		&config.tlsCaCert,
		"mtls-ca-cert",
		"",
		"The Root CA certificate path for custom mutual TLS channel.")
	runCmd.Flags().StringVar(
		&config.tlsCert,
		"mtls-cert",
		"",
		"The server certificate path for custom mutual TLS channel.")
	runCmd.Flags().StringVar(
		&config.tlsKey,
		"mtls-key",
		"",
		"The server private key path for custom mutual TLS channel.")
}
