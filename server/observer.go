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

package server

import (
	"context"
	"errors"
	"fmt"
	"sync"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc"
)

// UnaryObserver provides an interface for observing unary requests and responses.
type UnaryObserver interface {
	GetName() string
	ObserveUnary(
		ctx context.Context,
		req interface{},
		resp interface{},
		info *grpc.UnaryServerInfo,
		err error)
}

// StreamRequestObserver provides an interface for observing streaming requests.
type StreamRequestObserver interface {
	GetName() string
	ObserveStreamRequest(
		ctx context.Context,
		req interface{},
		info *grpc.StreamServerInfo,
		err error)
}

// StreamResponseObserver provides an interface for observing streaming responses.
type StreamResponseObserver interface {
	GetName() string
	ObserveStreamResponse(
		ctx context.Context,
		resp interface{},
		info *grpc.StreamServerInfo,
		err error)
}

// GrpcObserverRegistry is a registry of observers. These observers are hooked into the
// grpc interceptors that are provided by this interface.
type GrpcObserverRegistry interface {
	// UnaryInterceptor implements the grpc.UnaryServerInterceptor type to allow the
	// registry to hook into unary grpc methods.
	UnaryInterceptor(
		context.Context,
		interface{},
		*grpc.UnaryServerInfo,
		grpc.UnaryHandler) (interface{}, error)
	// StreamInterceptor implements the grpc.StreamServerInterceptor type to allow the
	// registry to hook into streaming grpc methods.
	StreamInterceptor(
		interface{},
		grpc.ServerStream,
		*grpc.StreamServerInfo,
		grpc.StreamHandler) error
	RegisterUnaryObserver(UnaryObserver)
	RegisterStreamRequestObserver(StreamRequestObserver)
	RegisterStreamResponseObserver(StreamResponseObserver)
}

// ShowcaseObserverRegistry returns the showcase specific observer registry.
func ShowcaseObserverRegistry() GrpcObserverRegistry {
	return &showcaseObserverRegistry{
		uObservers:     map[string]UnaryObserver{},
		sReqObservers:  map[string]StreamRequestObserver{},
		sRespObservers: map[string]StreamResponseObserver{},
	}
}

// showcaseObserverRegistry is an implementation of the ObserverRegistry. This registry
// automatically handles DeleteTest requests and deletes the appropriate observers
// for that request.
type showcaseObserverRegistry struct {
	mu             sync.Mutex
	uObservers     map[string]UnaryObserver
	sReqObservers  map[string]StreamRequestObserver
	sRespObservers map[string]StreamResponseObserver
}

func (r *showcaseObserverRegistry) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)

	if info.FullMethod == "/google.showcase.v1alpha3.Testing/DeleteTest" {
		err = r.deleteTestHandler(req)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, obs := range r.uObservers {
		obs.ObserveUnary(ctx, req, resp, info, err)
	}

	return resp, err
}

func (r *showcaseObserverRegistry) deleteTestHandler(req interface{}) error {
	deleteTestRequest, ok := req.(*pb.DeleteTestRequest)
	if !ok {
		return errors.New("Failed to delete the test")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, uFound := r.uObservers[deleteTestRequest.GetName()]
	_, sReqFound := r.sReqObservers[deleteTestRequest.GetName()]
	_, sRespFound := r.sRespObservers[deleteTestRequest.GetName()]
	if !(uFound || sReqFound || sRespFound) {
		return fmt.Errorf("Could not find test: %s", deleteTestRequest.GetName())
	}

	delete(r.uObservers, deleteTestRequest.GetName())
	delete(r.sReqObservers, deleteTestRequest.GetName())
	delete(r.sRespObservers, deleteTestRequest.GetName())
	return nil
}

type showcaseStream struct {
	info     *grpc.StreamServerInfo
	registry *showcaseObserverRegistry

	grpc.ServerStream
}

func (s *showcaseStream) SendMsg(m interface{}) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	err := s.ServerStream.SendMsg(m)
	for _, obs := range s.registry.sRespObservers {
		obs.ObserveStreamResponse(s.ServerStream.Context(), m, s.info, err)
	}
	return err
}

func (s *showcaseStream) RecvMsg(m interface{}) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	err := s.ServerStream.RecvMsg(m)
	for _, obs := range s.registry.sReqObservers {
		obs.ObserveStreamRequest(s.ServerStream.Context(), m, s.info, err)
	}
	return err
}

func (r *showcaseObserverRegistry) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return handler(srv, &showcaseStream{info, r, ss})
}

// RegisterUnaryObserver registers a unary observer. If an observer of the same name
// has already been registered, the new observer will override it.
func (r *showcaseObserverRegistry) RegisterUnaryObserver(obs UnaryObserver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.uObservers[obs.GetName()] = obs
}

// RegisterStreamRequestObserver registers a stream observer. If an observer of the same name
// has already been registered, the new observer will override it.
func (r *showcaseObserverRegistry) RegisterStreamRequestObserver(obs StreamRequestObserver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sReqObservers[obs.GetName()] = obs
}

// RegisterStreamResponseObserver registers a stream observer. If an observer of the same name
// has already been registered, the new observer will override it.
func (r *showcaseObserverRegistry) RegisterStreamResponseObserver(obs StreamResponseObserver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sRespObservers[obs.GetName()] = obs
}
