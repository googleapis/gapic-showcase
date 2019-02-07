// Copyright 2019 Google LLC
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

package server

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc"
)

type testUnaryObserver struct {
	name string
	req  interface{}
	resp interface{}
	info *grpc.UnaryServerInfo
	err  error
}

func (o *testUnaryObserver) GetName() string { return o.name }

func (o *testUnaryObserver) ObserveUnary(
	ctx context.Context,
	req interface{},
	resp interface{},
	info *grpc.UnaryServerInfo,
	err error) {
	o.req = req
	o.resp = resp
	o.info = info
	o.err = err
}

func Test_showcaseObserverRegistry_UnaryInterceptor(t *testing.T) {
	observerName := "observerName"
	tests := []struct {
		name            string
		req             interface{}
		resp            interface{}
		err             error
		info            *grpc.UnaryServerInfo
		observerDeleted bool
	}{
		{
			"Passes through request and response",
			"test req",
			"test resp",
			nil,
			&grpc.UnaryServerInfo{},
			false,
		},
		{
			"Passes through request and error",
			"test req",
			nil,
			errors.New("test error"),
			&grpc.UnaryServerInfo{},
			false,
		},
		{
			"Deletes observers by name",
			&pb.DeleteTestRequest{Name: observerName},
			&empty.Empty{},
			nil,
			&grpc.UnaryServerInfo{FullMethod: "/google.showcase.v1alpha3.Testing/DeleteTest"},
			true,
		},
	}
	for _, tt := range tests {
		obs := &testUnaryObserver{name: observerName}
		t.Run(tt.name, func(t *testing.T) {
			r := &showcaseObserverRegistry{
				uObservers: map[string]UnaryObserver{obs.GetName(): obs},
			}
			handler := func(_ context.Context, req interface{}) (interface{}, error) {
				if req != tt.req {
					t.Errorf("showcaseObserverRegistry.UnaryInterceptor() want to invoke handler with %v, got %v", tt.req, req)
				}
				return tt.resp, tt.err
			}
			got, err := r.UnaryInterceptor(context.Background(), tt.req, tt.info, handler)
			if err != tt.err {
				t.Errorf("showcaseObserverRegistry.UnaryInterceptor() error = %v, want %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.resp) {
				t.Errorf("showcaseObserverRegistry.UnaryInterceptor() = %v, want %v", got, tt.resp)
			}
			if tt.observerDeleted && r.uObservers[observerName] != nil {
				t.Error("showcaseObserverRegistry.UnaryInterceptor() want delete observers but did not")
			}
			if !tt.observerDeleted && obs.req != tt.req {
				t.Errorf("showcaseObserverRegistry.UnaryInterceptor() want to invoke observers with %v, got %v", tt.req, obs.req)
			}
			if !tt.observerDeleted && obs.resp != tt.resp {
				t.Errorf("showcaseObserverRegistry.UnaryInterceptor() want to invoke observers with %v, got %v", tt.resp, obs.resp)

			}
			if !tt.observerDeleted && obs.err != tt.err {
				t.Errorf("showcaseObserverRegistry.UnaryInterceptor() want to invoke observers with %v, got %v", tt.err, obs.err)

			}
			if !tt.observerDeleted && obs.info != tt.info {
				t.Errorf("showcaseObserverRegistry.UnaryInterceptor() want to invoke observers with %v, got %v", tt.info, obs.info)
			}
		})
	}
}

func Test_showcaseObserverRegistry_deleteTestHandler(t *testing.T) {
	type fields struct {
		uObservers     []UnaryObserver
		sReqObservers  []StreamRequestObserver
		sRespObservers []StreamResponseObserver
	}
	tests := []struct {
		name    string
		req     interface{}
		fields  fields
		wantErr bool
	}{
		{
			"Deletes Unary Observers",
			&pb.DeleteTestRequest{Name: "unaryObserver"},
			fields{uObservers: []UnaryObserver{&testUnaryObserver{name: "unaryObserver"}}},
			false,
		},
		{
			"Deletes Stream Request Observers",
			&pb.DeleteTestRequest{Name: "streamRequestObserver"},
			fields{uObservers: []UnaryObserver{&testUnaryObserver{name: "streamRequestObserver"}}},
			false,
		},
		{
			"Deletes Stream Response Observers",
			&pb.DeleteTestRequest{Name: "streamResponseObserver"},
			fields{uObservers: []UnaryObserver{&testUnaryObserver{name: "streamResponseObserver"}}},
			false,
		},
		{
			"Returns Error when an observer was not found",
			&pb.DeleteTestRequest{Name: "notFound"},
			fields{},
			true,
		},
		{
			"Returns Error request is not a DeleteTestRequest",
			"not DeleteTestRequest",
			fields{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &showcaseObserverRegistry{
				uObservers:     map[string]UnaryObserver{},
				sReqObservers:  map[string]StreamRequestObserver{},
				sRespObservers: map[string]StreamResponseObserver{},
			}
			for _, obs := range tt.fields.uObservers {
				r.RegisterUnaryObserver(obs)
			}
			for _, obs := range tt.fields.sReqObservers {
				r.RegisterStreamRequestObserver(obs)
			}
			for _, obs := range tt.fields.sRespObservers {
				r.RegisterStreamResponseObserver(obs)
			}
			if err := r.deleteTestHandler(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("showcaseObserverRegistry.deleteTestHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type testServerStream struct {
	err      error
	sent     interface{}
	received interface{}

	grpc.ServerStream
}

func (ss *testServerStream) Context() context.Context { return context.Background() }

func (ss *testServerStream) SendMsg(m interface{}) error {
	ss.sent = m
	return ss.err
}

func (ss *testServerStream) RecvMsg(m interface{}) error {
	ss.received = m
	return ss.err
}

type testStreamResponseObserver struct {
	name string
	resp interface{}
	info *grpc.StreamServerInfo
	err  error
}

func (o *testStreamResponseObserver) GetName() string { return o.name }

func (o *testStreamResponseObserver) ObserveStreamResponse(
	ctx context.Context,
	resp interface{},
	info *grpc.StreamServerInfo,
	err error) {
	o.resp = resp
	o.info = info
	o.err = err
}

func Test_showcaseStream_SendMsg(t *testing.T) {
	tests := []struct {
		name string
		info *grpc.StreamServerInfo
		msg  interface{}
		err  error
	}{
		{
			"Passes msg, info, and error to observer",
			&grpc.StreamServerInfo{},
			"sent msg",
			errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := &showcaseObserverRegistry{sRespObservers: map[string]StreamResponseObserver{}}
			obs := &testStreamResponseObserver{name: "streamObserver"}
			registry.RegisterStreamResponseObserver(obs)
			ss := &testServerStream{err: tt.err}
			s := &showcaseStream{tt.info, registry, ss}
			if err := s.SendMsg(tt.msg); err != tt.err {
				t.Errorf("showcaseStream.SendMsg() error = %v, want %v", err, tt.err)
			}
			if obs.resp != tt.msg {
				t.Errorf("showcaseStream.SendMsg() want to invoke observers with %v, got %v", tt.msg, obs.resp)
			}
			if ss.sent != tt.msg {
				t.Errorf("showcaseStream.SendMsg() want to invoke server stream with %v, got %v", tt.msg, ss.sent)
			}
			if obs.info != tt.info {
				t.Errorf("showcaseStream.SendMsg() want to invoke observers with %v, got %v", tt.info, obs.info)
			}
			if obs.err != tt.err {
				t.Errorf("showcaseStream.SendMsg() want to invoke observers with %v, got %v", tt.err, obs.err)
			}
		})
	}
}

type testStreamRequestObserver struct {
	name string
	req  interface{}
	info *grpc.StreamServerInfo
	err  error
}

func (o *testStreamRequestObserver) GetName() string { return o.name }

func (o *testStreamRequestObserver) ObserveStreamRequest(
	ctx context.Context,
	req interface{},
	info *grpc.StreamServerInfo,
	err error) {
	o.req = req
	o.info = info
	o.err = err
}

func Test_showcaseStream_RecvMsg(t *testing.T) {
	tests := []struct {
		name string
		info *grpc.StreamServerInfo
		msg  interface{}
		err  error
	}{
		{
			"Passes msg, info, and error to observer",
			&grpc.StreamServerInfo{},
			"received msg",
			errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := &showcaseObserverRegistry{sReqObservers: map[string]StreamRequestObserver{}}
			obs := &testStreamRequestObserver{name: "streamObserver"}
			registry.RegisterStreamRequestObserver(obs)
			ss := &testServerStream{err: tt.err}
			s := &showcaseStream{tt.info, registry, ss}
			if err := s.RecvMsg(tt.msg); err != tt.err {
				t.Errorf("showcaseStream.RecvMsg() error = %v, want %v", err, tt.err)
			}
			if obs.req != tt.msg {
				t.Errorf("showcaseStream.SendMsg() want to invoke observers with %v, got %v", tt.msg, obs.req)
			}
			if ss.received != tt.msg {
				t.Errorf("showcaseStream.SendMsg() want to invoke server stream with %v, got %v", tt.msg, ss.received)
			}
			if obs.info != tt.info {
				t.Errorf("showcaseStream.SendMsg() want to invoke observers with %v, got %v", tt.info, obs.info)
			}
			if obs.err != tt.err {
				t.Errorf("showcaseStream.SendMsg() want to invoke observers with %v, got %v", tt.err, obs.err)
			}
		})
	}
}

func Test_showcaseObserverRegistry_StreamInterceptor(t *testing.T) {
	r := ShowcaseObserverRegistry()
	srv := "server"
	ss := &testServerStream{}
	info := &grpc.StreamServerInfo{}
	tErr := errors.New("test error")
	handler := func(gotSrv interface{}, gotSs grpc.ServerStream) error {
		if srv != gotSrv {
			t.Errorf("showcaseObserverRegistry.StreamInterceptor() want to invoke handler with %v got %v", srv, gotSrv)
		}
		showcaseStream, ok := gotSs.(*showcaseStream)
		if !ok {
			t.Error("showcaseObserverRegistry.StreamInterceptor() expected to wrap server stream with showcase stream")
		}
		if info != showcaseStream.info {
			t.Errorf("showcaseObserverRegistry.StreamInterceptor() want to instantiate showcase stream with %v got %v", info, showcaseStream.info)
		}
		if r != showcaseStream.registry {
			t.Errorf("showcaseObserverRegistry.StreamInterceptor() want to instantiate showcase stream with %v got %v", r, showcaseStream.registry)
		}
		return tErr
	}
	if err := r.StreamInterceptor(srv, ss, info, handler); err != tErr {
		t.Errorf("showcaseObserverRegistry.StreamInterceptor() error = %v, wantErr %v", err, tErr)
	}
}
