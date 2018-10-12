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

package cmd

import (
	"context"

	"google.golang.org/grpc"
)

// This method implements the grpc.UnaryServerInterceptor interface.
func logServerUnary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	stdLog.Printf("Received Unary Request for Method: %s\n", info.FullMethod)
	stdLog.Printf("    Request:  %+v\n", req)
	resp, err := handler(ctx, req)
	if err == nil {
		stdLog.Printf("    Returning Response: %+v\n", resp)
	} else {
		stdLog.Printf("    Returning Error: %+v\n", err)
	}
	stdLog.Println("")
	return resp, err
}

type loggingServerStream struct {
	info *grpc.StreamServerInfo

	grpc.ServerStream
}

func (s *loggingServerStream) SendMsg(m interface{}) error {
	stdLog.Printf("%s Stream for Method: %s\n", s.streamType(), s.info.FullMethod)
	stdLog.Printf("    Sending Message:  %+v\n", m)
	stdLog.Println("")

	return s.ServerStream.SendMsg(m)
}

func (s *loggingServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	stdLog.Printf("%s Stream for Method: %s\n", s.streamType(), s.info.FullMethod)
	stdLog.Printf("    Recieving Message:  %v\n", m)
	stdLog.Println("")

	return err
}

func (s *loggingServerStream) streamType() string {
	if s.info.IsClientStream && s.info.IsServerStream {
		return "Bi-directional"
	} else if s.info.IsClientStream {
		return "Client"
	}
	return "Server"
}

// This method implements the grpc.StreamServerInterceptor interface.
func logServerStreaming(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	loggingStream := &loggingServerStream{info, ss}
	return handler(srv, loggingStream)
}

// This method implements the grpc.UnaryClientInterceptor interface.
func logClientUnary(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	stdLog.Printf("Sending Unary Request for Method: %s\n", method)
	stdLog.Printf("    Request:  %+v\n", req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		stdLog.Printf("    Got Response: %+v\n", reply)
	} else {
		stdLog.Printf("    Got Error: %+v\n", err)
	}
	stdLog.Println("")
	return err
}

type loggingClientStream struct {
	method string
	desc   *grpc.StreamDesc

	grpc.ClientStream
}

func (s *loggingClientStream) SendMsg(m interface{}) error {
	stdLog.Printf("%s Stream for Method: %s\n", s.streamType(), s.method)
	stdLog.Printf("    Sending Message:  %+v\n", m)
	stdLog.Println("")

	return s.ClientStream.SendMsg(m)
}

func (s *loggingClientStream) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	stdLog.Printf("%s Stream for Method: %s\n", s.streamType(), s.method)
	if err != nil {
		stdLog.Printf("    Recieved Error:  %v\n", err)
	} else {
		stdLog.Printf("    Recieving Message:  %v\n", m)
	}
	stdLog.Println("")

	return err
}

func (s *loggingClientStream) streamType() string {
	if s.desc.ClientStreams && s.desc.ServerStreams {
		return "Bi-directional"
	} else if s.desc.ClientStreams {
		return "Client"
	}
	return "Server"
}

// This method implements the grpc.StreamClientInterceptor interface.
func logClientStreaming(ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption) (grpc.ClientStream, error) {
	cs, err := streamer(ctx, desc, cc, method, opts...)
	return &loggingClientStream{method, desc, cs}, err
}
