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

package services

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gapic-showcase/server/spec"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewTestingServer returns a new TestingServer for the Showcase API.
func NewTestingServer(observerRegistry server.GrpcObserverRegistry) pb.TestingServer {
	name := fmt.Sprintf("sessions/-")
	defaultSession := server.NewSession(name, pb.Session_V1_LATEST, observerRegistry)
	defaultSession.RegisterTests(spec.ShowcaseTests(name, pb.Session_V1_LATEST))
	sessions := []sessionEntry{sessionEntry{session: defaultSession}}
	keys := map[string]int{name: len(sessions) - 1}

	s := &testingServerImpl{
		token:            server.NewTokenGenerator(),
		observerRegistry: observerRegistry,
		keys:             keys,
		sessions:         sessions,
	}

	return s
}

type sessionEntry struct {
	session server.Session
	deleted bool
}

type testingServerImpl struct {
	uid              server.UniqID
	token            server.TokenGenerator
	observerRegistry server.GrpcObserverRegistry

	mu       sync.Mutex
	keys     map[string]int
	sessions []sessionEntry
}

func (s *testingServerImpl) CreateSession(_ context.Context, req *pb.CreateSessionRequest) (*pb.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seshProto := req.GetSession()
	id := s.uid.Next()
	name := fmt.Sprintf("sessions/%d", id)
	sesh := server.NewSession(name, seshProto.GetVersion(), s.observerRegistry)
	sesh.RegisterTests(spec.ShowcaseTests(name, seshProto.GetVersion()))

	index := len(s.sessions)
	s.sessions = append(s.sessions, sessionEntry{session: sesh})
	s.keys[name] = index

	return server.SessionProto(sesh), nil
}

func (s *testingServerImpl) GetSession(_ context.Context, req *pb.GetSessionRequest) (*pb.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := req.GetName()
	if i, ok := s.keys[name]; ok && !s.sessions[i].deleted {
		return server.SessionProto(s.sessions[i].session), nil
	}

	return nil, status.Errorf(
		codes.NotFound, "A session with name %s not found.",
		name)
}

func (s *testingServerImpl) ListSessions(_ context.Context, in *pb.ListSessionsRequest) (*pb.ListSessionsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	start, err := s.token.GetIndex(in.GetPageToken())
	if err != nil {
		return nil, err
	}
	if start >= len(s.sessions) {
		return nil, server.InvalidTokenErr
	}

	offset := 0
	sessions := []*pb.Session{}
	for _, entry := range s.sessions[start:] {
		offset++
		if entry.deleted {
			continue
		}
		sessions = append(sessions, server.SessionProto(entry.session))
		if len(sessions) >= int(in.GetPageSize()) {
			break
		}
	}

	nextToken := ""
	if next := start + offset; next < len(s.sessions) {
		nextToken = s.token.ForIndex(next)
	}

	return &pb.ListSessionsResponse{Sessions: sessions, NextPageToken: nextToken}, nil
}

func (s *testingServerImpl) DeleteSession(_ context.Context, req *pb.DeleteSessionRequest) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.keys[req.GetName()]

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"A session with name %s not found.", req.GetName())
	}

	entry := s.sessions[i]
	s.sessions[i] = sessionEntry{session: entry.session, deleted: true}

	return &empty.Empty{}, nil
}

func (s *testingServerImpl) ReportSession(_ context.Context, req *pb.ReportSessionRequest) (*pb.ReportSessionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if i, ok := s.keys[req.GetName()]; ok && !s.sessions[i].deleted {
		return s.sessions[i].session.GetReport(), nil
	}

	return nil, status.Errorf(
		codes.NotFound, "A session with name %s not found.",
		req.GetName())
}

func (s *testingServerImpl) ListTests(_ context.Context, in *pb.ListTestsRequest) (*pb.ListTestsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := in.GetParent()

	if i, ok := s.keys[name]; ok && !s.sessions[i].deleted {
		return s.sessions[i].session.ListTests(in)
	}

	return nil, status.Errorf(
		codes.NotFound, "A session with name %s not found.",
		name)
}

func (s *testingServerImpl) DeleteTest(_ context.Context, req *pb.DeleteTestRequest) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, i := range s.keys {
		if strings.HasPrefix(req.GetName(), name) && !s.sessions[i].deleted {
			return s.sessions[i].session.DeleteTest(req.GetName())
		}
	}

	return nil, status.Errorf(
		codes.NotFound,
		"A test with name %s not found.", req.GetName())
}

func (s *testingServerImpl) VerifyTest(context.Context, *pb.VerifyTestRequest) (*pb.VerifyTestResponse, error) {
	// This should be handled by the test observers.
	return &pb.VerifyTestResponse{}, nil
}
