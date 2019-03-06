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
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Session represents a suite of tests, generally being made in the context
// of testing code generation.
//
// A session defines tests it may expect, based on which version of the
// code generation spec is in use.
type Session interface {
	GetName() string
	GetVersion() pb.Session_Version
	GetReport() *pb.ReportSessionResponse
	RegisterTests(tests []Test)
	ListTests(in *pb.ListTestsRequest) (*pb.ListTestsResponse, error)
	DeleteTest(name string) (*empty.Empty, error)
}

// SessionProto returns a proto representation of the Session.
func SessionProto(s Session) *pb.Session {
	return &pb.Session{
		Name:    s.GetName(),
		Version: s.GetVersion(),
	}
}

// NewSession returns a session for a given name, version and observerRegistry.
func NewSession(name string, version pb.Session_Version, observerRegistry GrpcObserverRegistry) Session {
	session := &sessionImpl{
		name:             name,
		version:          version,
		observerRegistry: observerRegistry,
		token:            NewTokenGenerator(),
		keys:             map[string]int{},
		tests:            []testEntry{},
	}
	return session
}

type testEntry struct {
	test    Test
	deleted bool
}

type sessionImpl struct {
	name             string
	version          pb.Session_Version
	observerRegistry GrpcObserverRegistry
	token            TokenGenerator

	mu    sync.Mutex
	keys  map[string]int
	tests []testEntry
}

func (s *sessionImpl) GetName() string {
	return s.name
}

func (s *sessionImpl) GetVersion() pb.Session_Version {
	return s.version
}

func (s *sessionImpl) GetReport() *pb.ReportSessionResponse {
	result := pb.ReportSessionResponse_PASSED
	for _, entry := range s.tests {
		if entry.deleted {
			continue
		}
		test := entry.test
		issue := test.GetIssue()
		if issue == nil {
			continue
		}
		// TODO: Add severity handling
		if issue.GetType() == pb.Issue_INCORRECT_CONFIRMATION {
			result = pb.ReportSessionResponse_FAILED
			continue
		}
		if result != pb.ReportSessionResponse_FAILED {
			result = pb.ReportSessionResponse_INCOMPLETE
		}
	}

	testRuns := []*pb.TestRun{}
	for _, entry := range s.tests {
		if entry.deleted {
			continue
		}
		testRuns = append(testRuns, TestRunProto(entry.test))
	}
	return &pb.ReportSessionResponse{
		Result:   result,
		TestRuns: testRuns,
	}
}

func (s *sessionImpl) ListTests(in *pb.ListTestsRequest) (*pb.ListTestsResponse, error) {
	start, err := s.token.GetIndex(in.GetPageToken())
	if err != nil {
		return nil, err
	}

	pageSize := in.GetPageSize()
	if pageSize == 0 {
		pageSize = 10
	}

	offset := 0
	tests := []*pb.Test{}
	for _, entry := range s.tests[start:] {
		offset++
		if entry.deleted {
			continue
		}
		tests = append(tests, TestProto(entry.test))
		if len(tests) >= int(pageSize) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(s.tests) {
		nextToken = s.token.ForIndex(start + offset)
	}

	return &pb.ListTestsResponse{Tests: tests, NextPageToken: nextToken}, nil
}

func (s *sessionImpl) DeleteTest(name string) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for n, i := range s.keys {
		if strings.HasSuffix(name, n) && !s.tests[i].deleted {
			s.tests[i] = testEntry{test: s.tests[i].test, deleted: true}

			test := s.tests[i].test
			if _, ok := test.(UnaryObserver); ok {
				s.observerRegistry.DeleteUnaryObserver(test.GetName())
			}
			if _, ok := test.(StreamRequestObserver); ok {
				s.observerRegistry.DeleteStreamRequestObserver(test.GetName())
			}
			if _, ok := test.(StreamResponseObserver); ok {
				s.observerRegistry.DeleteStreamResponseObserver(test.GetName())
			}

			return &empty.Empty{}, nil
		}
	}

	return nil, status.Errorf(
		codes.NotFound,
		"A test with name %s not found.", name)
}

func (s *sessionImpl) RegisterTests(tests []Test) {
	for _, test := range tests {
		i := len(s.tests)
		s.tests = append(s.tests, testEntry{test: test})
		s.keys[test.GetName()] = i
		if obs, ok := test.(UnaryObserver); ok {
			s.observerRegistry.RegisterUnaryObserver(obs)
		}
		if obs, ok := test.(StreamRequestObserver); ok {
			s.observerRegistry.RegisterStreamRequestObserver(obs)
		}
		if obs, ok := test.(StreamResponseObserver); ok {
			s.observerRegistry.RegisterStreamResponseObserver(obs)
		}
	}
}
