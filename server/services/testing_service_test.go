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
	"encoding/base64"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func Test_Session_lifecycle(t *testing.T) {
	s := NewTestingServer(server.ShowcaseObserverRegistry())

	first, err := s.CreateSession(
		context.Background(),
		&pb.CreateSessionRequest{
			Session: &pb.Session{Version: pb.Session_V1_0},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	delete, err := s.CreateSession(
		context.Background(),
		&pb.CreateSessionRequest{
			Session: &pb.Session{Version: pb.Session_V1_0},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteSession(
		context.Background(),
		&pb.DeleteSessionRequest{Name: delete.Name})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	created, err := s.CreateSession(
		context.Background(),
		&pb.CreateSessionRequest{
			Session: &pb.Session{Version: pb.Session_V1_0},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	got, err := s.GetSession(
		context.Background(),
		&pb.GetSessionRequest{Name: created.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	if !proto.Equal(created, got) {
		t.Error("Expected to get created session.")
	}

	r, err := s.ListSessions(
		context.Background(),
		&pb.ListSessionsRequest{PageSize: 1, PageToken: ""})
	if len(r.GetSessions()) != 1 {
		t.Errorf("List want: page size %d, got %d", 1, len(r.GetSessions()))
	}
	if r.GetSessions()[0].GetName() != "sessions/-" {
		t.Errorf("List want: first session 'sessions/-', got %+v", r.GetSessions()[0])
	}
	if r.GetNextPageToken() == "" {
		t.Error("List want: non empty next page token")
	}

	r, err = s.ListSessions(
		context.Background(),
		&pb.ListSessionsRequest{PageSize: 10, PageToken: r.GetNextPageToken()})
	if len(r.GetSessions()) != 2 {
		t.Errorf("List want: page size %d, got %d", 2, len(r.GetSessions()))
	}
	if !proto.Equal(first, r.GetSessions()[0]) {
		t.Errorf("List want: first session %+v, got %+v", got, r.GetSessions()[0])
	}
	if !proto.Equal(got, r.GetSessions()[1]) {
		t.Errorf("List want: second session %+v, got %+v", got, r.GetSessions()[1])
	}
	if r.GetNextPageToken() != "" {
		t.Error("List want: empty next page token")
	}
}

func Test_GetSession_deleted(t *testing.T) {
	s := NewTestingServer(server.ShowcaseObserverRegistry())
	created, err := s.CreateSession(
		context.Background(),
		&pb.CreateSessionRequest{
			Session: &pb.Session{Version: pb.Session_V1_0},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteSession(
		context.Background(),
		&pb.DeleteSessionRequest{Name: created.GetName()})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	_, err = s.GetSession(
		context.Background(),
		&pb.GetSessionRequest{Name: created.GetName()})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Get deleted: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_ListSessions_invalidToken(t *testing.T) {
	sessions := []sessionEntry{{session: &sessionMock{}}}
	keys := map[string]int{"name": len(sessions) - 1}

	s := &testingServerImpl{
		token:    server.TokenGeneratorWithSalt("salt"),
		keys:     keys,
		sessions: sessions,
	}

	tests := []string{
		"0", // Not base64 encoded
		base64.StdEncoding.EncodeToString([]byte("0")),        // No salt
		base64.StdEncoding.EncodeToString([]byte("saltblah")), // Invalid index
		base64.StdEncoding.EncodeToString([]byte("salt1000")), // index out of range.
	}

	for _, token := range tests {
		_, err := s.ListSessions(
			context.Background(),
			&pb.ListSessionsRequest{
				PageSize:  1,
				PageToken: token})
		status, _ := status.FromError(err)
		if status.Code() != codes.InvalidArgument {
			t.Errorf(
				"List: Want error code %d got %d",
				codes.InvalidArgument,
				status.Code())
		}
	}
}

func Test_DeleteSession_notFound(t *testing.T) {
	s := NewTestingServer(server.ShowcaseObserverRegistry())
	_, err := s.DeleteSession(
		context.Background(),
		&pb.DeleteSessionRequest{Name: "invalid"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Delete: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_ReportSession_notFound(t *testing.T) {
	s := NewTestingServer(server.ShowcaseObserverRegistry())
	_, err := s.ReportSession(
		context.Background(),
		&pb.ReportSessionRequest{Name: "not found"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"ReportSession: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

type sessionMock struct {
	wantReport   *pb.ReportSessionResponse
	wantList     *pb.ListTestsResponse
	deleteCalled bool

	server.Session
}

func (s *sessionMock) GetReport() *pb.ReportSessionResponse { return s.wantReport }

func (s *sessionMock) ListTests(in *pb.ListTestsRequest) (*pb.ListTestsResponse, error) {
	return s.wantList, nil
}

func (s *sessionMock) DeleteTest(name string) (*empty.Empty, error) {
	s.deleteCalled = true
	return &empty.Empty{}, nil
}

func Test_ReportSession(t *testing.T) {
	want := &pb.ReportSessionResponse{
		Result: pb.ReportSessionResponse_PASSED,
	}
	sessions := []sessionEntry{{session: &sessionMock{wantReport: want}}}
	keys := map[string]int{"name": 0}

	s := &testingServerImpl{
		keys:     keys,
		sessions: sessions,
	}

	got, err := s.ReportSession(
		context.Background(),
		&pb.ReportSessionRequest{Name: "name"})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	if !proto.Equal(got, want) {
		t.Errorf(
			"ReportSession: Want %+v got %+v",
			want,
			got)
	}
}

func Test_ListTests_notFound(t *testing.T) {
	s := NewTestingServer(server.ShowcaseObserverRegistry())
	_, err := s.ListTests(
		context.Background(),
		&pb.ListTestsRequest{Parent: "not found"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"ListTests: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_ListTests(t *testing.T) {
	want := &pb.ListTestsResponse{
		Tests: []*pb.Test{},
	}
	sessions := []sessionEntry{{session: &sessionMock{wantList: want}}}
	keys := map[string]int{"name": 0}

	s := &testingServerImpl{
		token:    server.NewTokenGenerator(),
		keys:     keys,
		sessions: sessions,
	}

	got, err := s.ListTests(
		context.Background(),
		&pb.ListTestsRequest{Parent: "name"})
	if err != nil {
		t.Errorf("ListTests: unexpected err %+v", err)
	}

	if !proto.Equal(got, want) {
		t.Errorf(
			"ListTests: Want %+v got %+v",
			want,
			got)
	}
}

func Test_DeleteTest_notFound(t *testing.T) {
	s := NewTestingServer(server.ShowcaseObserverRegistry())
	_, err := s.DeleteTest(
		context.Background(),
		&pb.DeleteTestRequest{Name: "not found"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"DeleteTest: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_DeleteTests(t *testing.T) {
	sesh := &sessionMock{}
	sessions := []sessionEntry{{session: sesh}}
	keys := map[string]int{"name": 0}

	s := &testingServerImpl{
		token:    server.NewTokenGenerator(),
		keys:     keys,
		sessions: sessions,
	}

	_, err := s.DeleteTest(
		context.Background(),
		&pb.DeleteTestRequest{Name: "name"})
	if err != nil {
		t.Errorf("DeleteTest: unexpected err %+v", err)
	}

	if !sesh.deleteCalled {
		t.Error("DeleteTest: expected to delegate call to session.")
	}
}

func Test_VerifyTest(t *testing.T) {
	s := NewTestingServer(server.ShowcaseObserverRegistry())
	got, err := s.VerifyTest(context.Background(), &pb.VerifyTestRequest{})
	if err != nil {
		t.Errorf("VerifyTest: unexpected err %+v", err)
	}
	if !proto.Equal(got, &pb.VerifyTestResponse{}) {
		t.Errorf("VerifyTest want %+v got %+v", &pb.VerifyTestResponse{}, got)
	}
}
