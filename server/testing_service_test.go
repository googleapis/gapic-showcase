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
	"reflect"
	"testing"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

func Test_NewTestingServer(t *testing.T) {
	// Simply test if this can get called.
	NewTestingServer()
}

type mockReportSession struct {
	reportCalled bool

	sessionImpl
}

func (m *mockReportSession) GetReport() *pb.ReportSessionResponse {
	m.reportCalled = true
	return &pb.ReportSessionResponse{}
}

func Test_ReportSession(t *testing.T) {
	session := &mockReportSession{reportCalled: false}
	server := testingServerImpl{session: session}
	server.ReportSession(context.Background(), &pb.ReportSessionRequest{})
	if !session.reportCalled {
		t.Error("Testing server should have defered reporting to the session.")
	}
}

type mockDeleteTestSession struct {
	deleteName string

	sessionImpl
}

func (m *mockDeleteTestSession) DeleteTest(name string) error {
	m.deleteName = name
	return nil
}

func Test_DeleteTest(t *testing.T) {
	session := &mockDeleteTestSession{deleteName: ""}
	server := testingServerImpl{session: session}
	server.DeleteTest(context.Background(), &pb.DeleteTestRequest{Name: "test"})
	if session.deleteName != "test" {
		t.Error("Testing server should have defered test deletion to the session.")
	}
}

type mockDeleteTestErrorSession struct {
	deleteName string

	sessionImpl
}

func (m *mockDeleteTestErrorSession) DeleteTest(name string) error {
	m.deleteName = name
	return errors.New("Test Error")
}

func Test_DeleteTest_error(t *testing.T) {
	session := &mockDeleteTestErrorSession{deleteName: ""}
	server := testingServerImpl{session: session}
	_, err := server.DeleteTest(context.Background(), &pb.DeleteTestRequest{Name: "test"})
	if session.deleteName != "test" {
		t.Error("Testing server should have defered test deletion to the session.")
	}
	if err == nil {
		t.Error("Testing server should have defered test deletion to the session.")
	}
}

type mockRegisterTestSession struct {
	testName string
	answers  []string

	sessionImpl
}

func (m *mockRegisterTestSession) TestAnswers(name string, answers []string) error {
	m.testName = name
	m.answers = answers
	return nil
}

func Test_RegisterTest(t *testing.T) {
	session := &mockRegisterTestSession{testName: "", answers: []string{}}
	server := testingServerImpl{session: session}
	answers := []string{"test", "answers"}

	server.RegisterTest(
		context.Background(),
		&pb.RegisterTestRequest{
			Name:    "/sessions/-/tests/eggs.bacon.test",
			Answers: answers,
		})

	if session.testName != "eggs.bacon.test" {
		t.Errorf("Testing server expected to search for test "+
			"'eggs.bacon.test' but searched for: %s", session.testName)
	}
	if !reflect.DeepEqual(session.answers, answers) {
		t.Error("Testing server should have deferred testing to the session.")
	}
}

type mockRegisterTestErrorSession struct {
	testName string
	answers  []string

	sessionImpl
}

func (m *mockRegisterTestErrorSession) TestAnswers(name string, answers []string) error {
	m.testName = name
	m.answers = answers
	return errors.New("Test Error")
}

func Test_RegisterTest_error(t *testing.T) {
	session := &mockRegisterTestErrorSession{testName: "", answers: []string{}}
	server := testingServerImpl{session: session}
	answers := []string{"test", "answers"}

	_, err := server.RegisterTest(
		context.Background(),
		&pb.RegisterTestRequest{
			Name:    "/sessions/-/tests/eggs.bacon.test",
			Answers: answers,
		})

	if session.testName != "eggs.bacon.test" {
		t.Errorf("Testing server expected to search for test "+
			"'eggs.bacon.test' but searched for: %s", session.testName)
	}
	if !reflect.DeepEqual(session.answers, answers) {
		t.Error("Testing server should have deferred testing to the session.")
	}
	if err == nil {
		t.Error("Testing server should have deferred testing to the session.")
	}
}
