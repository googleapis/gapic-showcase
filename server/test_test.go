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
	"testing"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

func Test_TestGetName(t *testing.T) {
	names := []string{
		"bacon.required.test",
		"",
	}

	for _, name := range names {
		proto := &pb.Test{Name: name}
		test := TestFromProto(proto)
		if name != test.GetName() {
			t.Errorf("Test.GetName: Expected '%s' Got '%s'", name, test.GetName())
		}
	}
}

func Test_GetExpectationLevel(t *testing.T) {
	expectationLevels := []pb.Test_ExpectationLevel{
		pb.Test_REQUIRED,
		pb.Test_OPTIONAL,
		pb.Test_RECOMMENDED,
		pb.Test_EXPECTATION_LEVEL_UNSPECIFIED,
	}

	for _, level := range expectationLevels {
		proto := &pb.Test{Name: "bacon.required.test", ExpectationLevel: level}
		test := TestFromProto(proto)
		if level != test.GetExpectationLevel() {
			t.Errorf("Test.GetExpectationLevel: Expected '%d' Got '%d'", level, test.GetExpectationLevel())
		}
	}
}

func Test_GetState(t *testing.T) {
	states := []pb.ReportSessionResponse_State{
		pb.ReportSessionResponse_INCOMPLETE,
		pb.ReportSessionResponse_PASSED,
		pb.ReportSessionResponse_FAILED,
		pb.ReportSessionResponse_STATE_UNSPECIFIED,
	}

	for _, state := range states {
		test := &testImpl{state: state}
		if state != test.GetState() {
			t.Errorf("Test.GetState: Expected '%d' Got '%d'", state, test.GetState())

		}
	}
}

func Test_GetIssue(t *testing.T) {
	issues := []*pb.ReportSessionResponse_Issue{
		&pb.ReportSessionResponse_Issue{
			Test:        "bacon.required.test",
			Type:        pb.ReportSessionResponse_Issue_INCORRECT_CONFIRMATION,
			Severity:    pb.ReportSessionResponse_Issue_ERROR,
			Description: "Incorrect answer",
		},
		&pb.ReportSessionResponse_Issue{
			Test:        "eggs.required.test",
			Type:        pb.ReportSessionResponse_Issue_SKIPPED,
			Severity:    pb.ReportSessionResponse_Issue_WARNING,
			Description: "Skipped",
		},
		nil,
	}

	for _, issue := range issues {
		test := &testImpl{issue: issue}
		if issue != test.GetIssue() {
			t.Error("Test.GetIssue failed getting the set issue.")
		}
	}
}
