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
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/googleapis/gapic-showcase/server"
	showcasepb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/takama/daemon"
	lropb "google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	name        = "gapic-showcase"
	description = "Gapic Showcase V1Alpha1 Service"
	version     = "0.0.4"
	port        = ":7469"
)

var stdlog, errlog *log.Logger
var dependencies = []string{}

// Service has embedded daemon
type process struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (p *process) manage() (string, error) {

	usage := fmt.Sprintf(
		"Usage: %s version | install | remove | start | stop | status", os.Args[0])

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "version":
			return version, nil
		case "install":
			return p.Install()
		case "remove":
			return p.Remove()
		case "start":
			return p.Start()
		case "stop":
			return p.Stop()
		case "status":
			return p.Status()
		default:
			return usage, nil
		}
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set start listening.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	stdlog.Printf("Gapic Showcase V1Alpha1 listening on port: %s", port)

	// Setup Server.
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logRequests),
	}
	s := grpc.NewServer(opts...)
	defer s.GracefulStop()

	opStore := server.NewOperationStore()
	showcasepb.RegisterShowcaseServer(s, server.NewShowcaseServer(opStore))
	lropb.RegisterOperationsServer(s, server.NewOperationsServer(opStore))

	// Register reflection service on gRPC server.
	reflection.Register(s)
	go s.Serve(lis)

	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
	return usage, nil
}

func logRequests(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	stdlog.Printf("Received Request for Method: %s\n", info.FullMethod)
	stdlog.Printf("    Request:  %+v\n", req)
	resp, err := handler(ctx, req)
	if err == nil {
		stdlog.Printf("    Returning Response: %+v\n", resp)
	} else {
		stdlog.Printf("    Returning Error: %+v\n", err)
	}
	stdlog.Println("")
	return resp, err
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	p := &process{srv}
	status, err := p.manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
