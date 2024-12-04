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
	"testing"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func TestSessionProto(t *testing.T) {
	s := NewSession("name", pb.Session_V1_0, ShowcaseObserverRegistry())
	if s.GetName() != "name" {
		t.Errorf("Session.GetName() = %v, want %v", s.GetName(), "name")

	}
	if s.GetVersion() != pb.Session_V1_0 {
		t.Errorf("Session.GetVersion() = %v, want %v", s.GetVersion(), pb.Session_V1_0)
	}
	want := &pb.Session{Name: "name", Version: pb.Session_V1_0}
	got := SessionProto(s)
	if !proto.Equal(got, want) {
		t.Errorf("SessionProto() = %v, want %v", got, want)
	}
}

type mockTest struct {
	name        string
	expectation pb.Test_ExpectationLevel
	blueprints  []*pb.Test_Blueprint

	iType       pb.Issue_Type
	severity    pb.Issue_Severity
	description string

	Test
}

func (t *mockTest) GetName() string {
	return t.name
}

func (t *mockTest) GetIssue() *pb.Issue {
	if t.iType == 0 {
		return nil
	}
	return &pb.Issue{
		Type:        t.iType,
		Severity:    t.severity,
		Description: t.description,
	}
}

func (t *mockTest) GetExpectationLevel() pb.Test_ExpectationLevel {
	return t.expectation
}

func (t *mockTest) GetDescription() string {
	return t.description
}

func (t *mockTest) GetBlueprints() []*pb.Test_Blueprint {
	return t.blueprints
}

// Implements UnaryObserver
func (t *mockTest) ObserveUnary(
	ctx context.Context,
	req interface{},
	resp interface{},
	info *grpc.UnaryServerInfo,
	err error) {
}

// Implements StreamRequestObserver
func (t *mockTest) ObserveStreamRequest(
	ctx context.Context,
	req interface{},
	info *grpc.StreamServerInfo,
	err error) {
}

// Implements StreamResponseObserver
func (t *mockTest) ObserveStreamResponse(
	ctx context.Context,
	resp interface{},
	info *grpc.StreamServerInfo,
	err error) {
}

func Test_sessionImpl_TestLifeCycle(t *testing.T) {
	failed := &mockTest{
		name:        "failedTest",
		iType:       pb.Issue_INCORRECT_CONFIRMATION,
		severity:    pb.Issue_ERROR,
		expectation: pb.Test_REQUIRED,
		description: "This test failed.",
	}
	pending := &mockTest{
		name:        "pendingTest",
		iType:       pb.Issue_PENDING,
		severity:    pb.Issue_ERROR,
		expectation: pb.Test_RECOMMENDED,
		description: "This test failed.",
	}
	skipped := &mockTest{
		name:        "skippedTest",
		expectation: pb.Test_OPTIONAL,
		iType:       pb.Issue_SKIPPED,
		severity:    pb.Issue_ERROR,
		description: "This test failed.",
	}
	passed := &mockTest{
		name:        "passedTest",
		expectation: pb.Test_REQUIRED,
		description: "This test passed.",
	}
	// Register tests and ensure they are all listed.
	session := &sessionImpl{
		observerRegistry: ShowcaseObserverRegistry(),
		token:            &tokenGenerator{salt: ""},
		keys:             map[string]int{},
		tests:            []testEntry{},
	}
	session.RegisterTests([]Test{failed, pending, skipped, passed})
	wantList := &pb.ListTestsResponse{
		Tests: []*pb.Test{
			TestProto(failed),
			TestProto(pending),
			TestProto(skipped),
			TestProto(passed),
		},
	}

	got, _ := session.ListTests(&pb.ListTestsRequest{})
	if !proto.Equal(got, wantList) {
		t.Errorf("sessionImpl.ListTests() = %v, want %v", got, wantList)
	}

	// Delete tests
	if _, err := session.DeleteTest("skippedTest"); err != nil {
		t.Errorf("sessionImpl.DeleteTest() = %v", err)
	}
	if _, err := session.DeleteTest("invalidName"); err == nil {
		t.Error("sessionImpl.DeleteTest() for invalid name wanted err got nil")
	}

	// Ensure list of test does not include the deleted tests.
	wantList = &pb.ListTestsResponse{
		Tests: []*pb.Test{
			TestProto(failed),
			TestProto(pending),
			TestProto(passed),
		},
	}
	got, _ = session.ListTests(&pb.ListTestsRequest{})
	if !proto.Equal(got, wantList) {
		t.Errorf("sessionImpl.ListTests() = %v, want %v", got, wantList)
	}

	// Test pagination
	wantList = &pb.ListTestsResponse{
		Tests: []*pb.Test{
			TestProto(failed),
			TestProto(pending),
		},
		NextPageToken: "Mg==", // Deterministic since we hard coded the page token salt.
	}
	got, _ = session.ListTests(&pb.ListTestsRequest{PageSize: 2})
	if !proto.Equal(got, wantList) {
		t.Errorf("sessionImpl.ListTests() = %v, want %v", got, wantList)
	}

	// Invalid page token
	if _, err := session.ListTests(&pb.ListTestsRequest{PageToken: "invalid"}); err == nil {
		t.Error("sessionImpl.ListTests() with invalid page token got nil.")
	}

}

func Test_sessionImpl_GetReport(t *testing.T) {
	failed := &mockTest{
		name:        "failedTest",
		iType:       pb.Issue_INCORRECT_CONFIRMATION,
		severity:    pb.Issue_ERROR,
		description: "This test failed.",
	}
	pending := &mockTest{
		name:        "pendingTest",
		iType:       pb.Issue_PENDING,
		severity:    pb.Issue_ERROR,
		description: "This test failed.",
	}
	skipped := &mockTest{
		name:        "skippedTest",
		iType:       pb.Issue_SKIPPED,
		severity:    pb.Issue_ERROR,
		description: "This test failed.",
	}
	passed := &mockTest{name: "passedTest"}
	tests := []struct {
		name    string
		entries []testEntry
		want    *pb.ReportSessionResponse
	}{
		{
			"Passes if all tests passed",
			[]testEntry{
				testEntry{test: &mockTest{name: "passed 1"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 2"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 3"}, deleted: false},
			},
			&pb.ReportSessionResponse{
				Result: pb.ReportSessionResponse_PASSED,
				TestRuns: []*pb.TestRun{
					&pb.TestRun{Test: "passed 1"},
					&pb.TestRun{Test: "passed 2"},
					&pb.TestRun{Test: "passed 3"},
				},
			},
		},
		{
			"Marks as incomplete if there are pending and skipped tests.",
			[]testEntry{
				testEntry{test: &mockTest{name: "passed 1"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 2"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 3"}, deleted: false},
				testEntry{test: pending, deleted: false},
				testEntry{test: skipped, deleted: false},
			},
			&pb.ReportSessionResponse{
				Result: pb.ReportSessionResponse_INCOMPLETE,
				TestRuns: []*pb.TestRun{
					&pb.TestRun{Test: "passed 1"},
					&pb.TestRun{Test: "passed 2"},
					&pb.TestRun{Test: "passed 3"},
					&pb.TestRun{Test: pending.GetName(), Issue: pending.GetIssue()},
					&pb.TestRun{Test: skipped.GetName(), Issue: skipped.GetIssue()},
				},
			},
		},
		{
			"Fails if one tests failed",
			[]testEntry{
				testEntry{test: &mockTest{name: "passed 1"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 2"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 3"}, deleted: false},
				testEntry{test: pending, deleted: false},
				testEntry{test: skipped, deleted: false},
				testEntry{test: failed, deleted: false},
			},
			&pb.ReportSessionResponse{
				Result: pb.ReportSessionResponse_FAILED,
				TestRuns: []*pb.TestRun{
					&pb.TestRun{Test: "passed 1"},
					&pb.TestRun{Test: "passed 2"},
					&pb.TestRun{Test: "passed 3"},
					&pb.TestRun{Test: pending.GetName(), Issue: pending.GetIssue()},
					&pb.TestRun{Test: skipped.GetName(), Issue: skipped.GetIssue()},
					&pb.TestRun{Test: failed.GetName(), Issue: failed.GetIssue()},
				},
			},
		},
		{
			"Marks as incomplete if there are pending and skipped tests.",
			[]testEntry{
				testEntry{test: &mockTest{name: "passed 1"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 2"}, deleted: false},
				testEntry{test: &mockTest{name: "passed 3"}, deleted: false},
				testEntry{test: pending, deleted: false},
				testEntry{test: skipped, deleted: false},
			},
			&pb.ReportSessionResponse{
				Result: pb.ReportSessionResponse_INCOMPLETE,
				TestRuns: []*pb.TestRun{
					&pb.TestRun{Test: "passed 1"},
					&pb.TestRun{Test: "passed 2"},
					&pb.TestRun{Test: "passed 3"},
					&pb.TestRun{Test: pending.GetName(), Issue: pending.GetIssue()},
					&pb.TestRun{Test: skipped.GetName(), Issue: skipped.GetIssue()},
				},
			},
		},
		{
			"Skips Deleted Tests",
			[]testEntry{
				testEntry{test: passed, deleted: false},
				testEntry{test: failed, deleted: true},
			},
			&pb.ReportSessionResponse{
				Result:   pb.ReportSessionResponse_PASSED,
				TestRuns: []*pb.TestRun{TestRunProto(passed)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sessionImpl{
				tests: tt.entries,
			}
			if got := s.GetReport(); !proto.Equal(got, tt.want) {
				t.Errorf("sessionImpl.GetReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
