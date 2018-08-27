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
	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

type Test interface {
	GetName() string
	GetExpectationLevel() pb.Test_ExpectationLevel
	GetProto() *pb.Test
	GetState() pb.ReportSessionResponse_State
	GetIssue() *pb.ReportSessionResponse_Issue
	GetAnswers() []string

	HasFailed() bool
}

type testImpl struct {
	t       *pb.Test
	issue   *pb.ReportSessionResponse_Issue
	state   pb.ReportSessionResponse_State
	answers []string
}

func NewTest(t *pb.Test, issue *pb.ReportSessionResponse_Issue, state pb.ReportSessionResponse_State, answers []string) Test {
	return &testImpl{
		t:       t,
		issue:   issue,
		state:   state,
		answers: answers,
	}
}

func TestFromProto(t *pb.Test) Test {
	return NewTest(
		t,
		createIssue(t, pb.ReportSessionResponse_Issue_SKIPPED, "This has not been tested."),
		pb.ReportSessionResponse_INCOMPLETE,
		[]string{})
}

func (t *testImpl) GetName() string {
	return t.t.Name
}

func (t *testImpl) GetExpectationLevel() pb.Test_ExpectationLevel {
	return t.t.ExpectationLevel
}

func (t *testImpl) GetProto() *pb.Test {
	return t.t
}

func (t *testImpl) GetIssue() *pb.ReportSessionResponse_Issue {
	return t.issue
}

func (t *testImpl) GetState() pb.ReportSessionResponse_State {
	return t.state
}

func (t *testImpl) GetAnswers() []string {
	return t.answers
}

func (t *testImpl) HasFailed() bool {
	return t.state == pb.ReportSessionResponse_FAILED
}
