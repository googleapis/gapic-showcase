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
  "reflect"
  "sync"
  "testing"

  pb "github.com/googleapis/gapic-showcase/server/genproto"
)

func Test_GetName(t *testing.T) {
  names := []string{
    "bacon.required.test",
    "",
  }

  for _, name := range names {
    proto := &pb.Test{ Name: name }
    test := NewTestFromProto(proto)
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
    proto := &pb.Test{ Name: "bacon.required.test", ExpectationLevel: level }
    test := NewTestFromProto(proto)
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
      Test: "bacon.required.test",
      Type: pb.ReportSessionResponse_Issue_INCORRECT_CONFIRMATION,
      Severity: pb.ReportSessionResponse_Issue_ERROR,
      Description: "Incorrect answer",
    },
    &pb.ReportSessionResponse_Issue{
      Test: "eggs.required.test",
      Type: pb.ReportSessionResponse_Issue_SKIPPED,
      Severity: pb.ReportSessionResponse_Issue_WARNING,
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

func Test_AddAnswers(t *testing.T) {
  testCases := []struct {
    initial     []string
    additional  []string
  }{
    {[]string{}, []string{"bacon"}},
    {[]string{"eggs"}, []string{"bacon"}},
    {[]string{"grits"}, []string{}},
  }

  for _, testCase := range testCases {
    test := testImpl{
      mu: &sync.Mutex{},
      answers: testCase.initial,
      state: pb.ReportSessionResponse_PASSED,
    }
    test.AddAnswers(testCase.additional)
    if !reflect.DeepEqual(append(testCase.initial, testCase.additional...), test.answers) {
      t.Error("Test.AddAnswers failed to append answers.")
    }
    if test.GetState() != pb.ReportSessionResponse_INCOMPLETE {
      t.Errorf("Test.AddAnswers expected to mark the test as incomplete.")
    }
  }

  for _, testCase := range testCases {
    test := testImpl{
      mu: &sync.Mutex{},
      answers: testCase.initial,
      state: pb.ReportSessionResponse_FAILED,
    }
    test.AddAnswers(testCase.additional)
    if !reflect.DeepEqual(append(testCase.initial, testCase.additional...), test.answers) {
      t.Error("Test.AddAnswers failed to append answers.")
    }
    if test.GetState() != pb.ReportSessionResponse_FAILED {
      t.Error("Test.AddAnswers expected to keep the test as failed.")
    }
  }
}

func Test_TestAnswers_already_failed(t *testing.T) {
  initialAnswers := []string{"Should", "not", "change"}

  tmp := make([]string, len(initialAnswers))
  copy(tmp, initialAnswers)
  test := &testImpl{
    mu: &sync.Mutex{},
    answers: tmp,
    state: pb.ReportSessionResponse_FAILED,
  }

  answers := [][]string{
    []string{"bacon", "eggs"},
    []string{},
  }
  for _, a := range answers {
    test.TestAnswers(a)
    if !reflect.DeepEqual(initialAnswers, test.answers) {
      t.Error("Test.TestAnswers expected to not alter the test after already failing.")
    }
  }
}

func Test_TestAnswers_more_answers_than_expected(t *testing.T) {
  proto := &pb.Test{ Name: "eggos.test" }
  test := NewTestFromProto(proto)

  test.TestAnswers([]string{"durian", "for", "breakfast"})
  if test.GetState() != pb.ReportSessionResponse_FAILED {
    t.Error("Test.TestAnswers should fail when testing more answers than expecrted.")
  }
  if test.GetIssue() == nil {
    t.Error("Test.TestAnswers should have made an issue for this failure.")
  }
}

func Test_TestAnswers_failed_tests(t *testing.T) {
  proto := &pb.Test{ Name: "breakfast.test" }
  test := NewTestFromProto(proto)

  test.AddAnswers([]string{"mangos", "for", "lunch!"})
  test.TestAnswers([]string{"mangos", "for", "dinner?"})

  if test.GetState() != pb.ReportSessionResponse_FAILED {
    t.Error("Test.TestAnswers should fail when testing an unexpected answer")
  }
  if test.GetIssue() == nil {
    t.Error("Test.TestAnswers should have made an issue for this failure.")
  }
}

func Test_TestAnswers_correct_tests(t *testing.T) {
  proto := &pb.Test{ Name: "breakfast.test" }
  test := NewTestFromProto(proto)

  test.AddAnswers([]string{"mangos", "for", "lunch!"})
  test.TestAnswers([]string{"mangos"})
  if test.GetState() != pb.ReportSessionResponse_INCOMPLETE {
    t.Error("Test.TestAnswers should only pass when all expected values have been tested.")
  }
  if test.GetIssue() == nil {
    t.Error("Test.TestAnswers should only set issue to nil when the test passes.")
  }

  test.TestAnswers([]string{"for", "lunch!"})
  if test.GetState() != pb.ReportSessionResponse_PASSED {
    t.Error("Test.TestAnswers should pass when all expected values have been tested.")
  }
  if test.GetIssue() != nil {
    t.Error("Test.TestAnswers should set issue to nil when the test passes.")
  }
}
