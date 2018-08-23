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
	"sync"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

type Test interface {
	GetName() string
  GetExpectationLevel() pb.Test_ExpectationLevel
	GetState() pb.ReportSessionResponse_State
	GetIssue() *pb.ReportSessionResponse_Issue
	AddAnswers(answers []string)
	TestAnswers(answers []string)
}

type testImpl struct {
	t 					*pb.Test
	issue				*pb.ReportSessionResponse_Issue
	state				pb.ReportSessionResponse_State

  mu          *sync.Mutex
  answers 		[]string
}

func NewTestFromProto(t *pb.Test) Test {
  return &testImpl{
		t: t,
		answers: []string{},
		state: pb.ReportSessionResponse_INCOMPLETE,
		issue: createIssue(t, pb.ReportSessionResponse_Issue_SKIPPED, "This has not been tested."),
	}
}

func (t *testImpl) GetName() string {
	return t.t.Name
}

func (t *testImpl) GetExpectationLevel() pb.Test_ExpectationLevel {
  return t.t.ExpectationLevel
}

func (t *testImpl) GetState() pb.ReportSessionResponse_State {
	return t.state
}

func (t *testImpl) GetIssue() *pb.ReportSessionResponse_Issue {
	return t.issue
}

func (t *testImpl) AddAnswers(answers []string) {
  t.mu.Lock()
  defer t.mu.Unlock()

	t.answers = append(t.answers, answers...)

	if !t.hasFailed() {
		t.state = pb.ReportSessionResponse_INCOMPLETE
	}
}

func (t *testImpl) TestAnswers(answers []string) {
  t.mu.Lock()
  defer t.mu.Unlock()

	if t.hasFailed() {
		return
	}

	if len(answers) > len(t.answers) {
		t.state = pb.ReportSessionResponse_FAILED
    t.issue = createIssue(t.t, pb.ReportSessionResponse_Issue_INCORRECT_CONFIRMATION, "More answers registered than Showcase expected.")
		return
	}
	for _, a := range answers {
		if a != t.answers[0] {
			t.state = pb.ReportSessionResponse_FAILED
      t.issue = createIssue(t.t, pb.ReportSessionResponse_Issue_INCORRECT_CONFIRMATION, "Incorrect answer registered. Expected '%s' but got '%s'.")
			return
		}
		t.answers = t.answers[1:]
    if len(t.answers) == 0 {
      t.state = pb.ReportSessionResponse_PASSED
      t.issue = nil
    }
	}
}

func (t *testImpl) hasFailed() bool {
	return t.state == pb.ReportSessionResponse_FAILED
}

func createIssue(t *pb.Test, issueType pb.ReportSessionResponse_Issue_Type, desc string) *pb.ReportSessionResponse_Issue {
  var severity pb.ReportSessionResponse_Issue_Severity
  if (t.ExpectationLevel <= pb.Test_REQUIRED) {
    severity = pb.ReportSessionResponse_Issue_ERROR
  } else {
    severity = pb.ReportSessionResponse_Issue_WARNING
  }
  return &pb.ReportSessionResponse_Issue{
    Test: t.Name,
    Type: issueType,
    Severity: severity,
    Description: desc,
  }
}
