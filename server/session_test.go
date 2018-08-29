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

func Test_GetSessionSingleton(t *testing.T) {
	s := GetSessionSingleton()
	if s.GetName() != "-" {
		t.Error("Expected session singleton to be named '-'")
	}
}

func Test_GetReport_empty_report(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}
	r := s.GetReport()

	if r.State != pb.ReportSessionResponse_PASSED {
		t.Error("Expected empty report to be marked as passed.")
	}
	if len(r.Errors) != 0 {
		t.Error("Expected empty report to not contain errors.")
	}
	if len(r.Warnings) != 0 {
		t.Error("Expected empty report to not contain warnings.")
	}
	if r.Completion.Total > float32(0) {
		t.Error("Expected completion ratios for empty reports to be 0.")
	}
	if r.Completion.Required > float32(0) {
		t.Error("Expected completion ratios for empty reports to be 0.")
	}
	if r.Completion.Recommended > float32(0) {
		t.Error("Expected completion ratios for empty reports to be 0.")
	}
	if r.Completion.Optional > float32(0) {
		t.Error("Expected completion ratios for empty reports to be 0.")
	}
}

func Test_GetReport_required_tests(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	passing := TestFromProto(&pb.Test{
		Name:             "passing",
		ExpectationLevel: pb.Test_REQUIRED,
	})
	failing := TestFromProto(&pb.Test{
		Name:             "failing",
		ExpectationLevel: pb.Test_REQUIRED,
	})
	skipped := TestFromProto(&pb.Test{
		Name:             "skipped",
		ExpectationLevel: pb.Test_REQUIRED,
	})

	s.AddTest(passing)
	s.AddTest(failing)
	s.AddTest(skipped)

	s.AddAnswers("passing", []string{"eggs", "bacon"})
	s.TestAnswers("passing", []string{"eggs", "bacon"})

	s.AddAnswers("failing", []string{"eggs", "bacon"})
	s.TestAnswers("failing", []string{"spam", "toast"})

	s.AddAnswers("skipped", []string{"eggs", "bacon"})

	r := s.GetReport()
	if r.State != pb.ReportSessionResponse_FAILED {
		t.Error("Expected report to be marked as failed.")
	}
	if len(r.Errors) != 2 {
		t.Error("Expected report to contain two errors.")
	}
	if len(r.Warnings) != 0 {
		t.Error("Expected report to not contain warnings.")
	}
	if r.Completion.Total != (float32(1) / float32(3)) {
		t.Error("Incorrect total completion ratio.")
	}
	if r.Completion.Required != float32(1)/float32(3) {
		t.Error("Incorrect required completionration.")
	}
	if r.Completion.Recommended > float32(0) {
		t.Error("Expected recommended completion ratio to be 0.")
	}
	if r.Completion.Optional > float32(0) {
		t.Error("Expected optional completion ratio to be 0.")
	}
}

func Test_GetReport_recommended_tests(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	passing := TestFromProto(&pb.Test{
		Name:             "passing",
		ExpectationLevel: pb.Test_RECOMMENDED,
	})
	failing := TestFromProto(&pb.Test{
		Name:             "failing",
		ExpectationLevel: pb.Test_RECOMMENDED,
	})
	skipped := TestFromProto(&pb.Test{
		Name:             "skipped",
		ExpectationLevel: pb.Test_RECOMMENDED,
	})

	s.AddTest(passing)
	s.AddTest(failing)
	s.AddTest(skipped)

	s.AddAnswers("passing", []string{"eggs", "bacon"})
	s.TestAnswers("passing", []string{"eggs", "bacon"})

	s.AddAnswers("failing", []string{"eggs", "bacon"})
	s.TestAnswers("failing", []string{"spam", "toast"})

	s.AddAnswers("skipped", []string{"eggs", "bacon"})

	r := s.GetReport()
	if r.State != pb.ReportSessionResponse_PASSED {
		t.Error("Expected report to be marked as passed.")
	}
	if len(r.Errors) != 0 {
		t.Error("Expected report to contain two errors.")
	}
	if len(r.Warnings) != 2 {
		t.Error("Expected report to not contain warnings.")
	}
	if r.Completion.Total != (float32(1) / float32(3)) {
		t.Error("Incorrect total completion ratio.")
	}
	if r.Completion.Required != float32(0) {
		t.Error("Incorrect required completion ratio.")
	}
	if r.Completion.Recommended > float32(1)/float32(3) {
		t.Error("Incorrect recommeded completion ratio")
	}
	if r.Completion.Optional > float32(0) {
		t.Error("Expected optional completion ratio to be 0.")
	}
}

func Test_GetReport_optional_tests(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	passing := TestFromProto(&pb.Test{
		Name:             "passing",
		ExpectationLevel: pb.Test_OPTIONAL,
	})
	failing := TestFromProto(&pb.Test{
		Name:             "failing",
		ExpectationLevel: pb.Test_OPTIONAL,
	})
	skipped := TestFromProto(&pb.Test{
		Name:             "skipped",
		ExpectationLevel: pb.Test_OPTIONAL,
	})

	s.AddTest(passing)
	s.AddTest(failing)
	s.AddTest(skipped)

	s.AddAnswers("passing", []string{"eggs", "bacon"})
	s.TestAnswers("passing", []string{"eggs", "bacon"})

	s.AddAnswers("failing", []string{"eggs", "bacon"})
	s.TestAnswers("failing", []string{"spam", "toast"})

	s.AddAnswers("skipped", []string{"eggs", "bacon"})

	r := s.GetReport()
	if r.State != pb.ReportSessionResponse_PASSED {
		t.Error("Expected report to be marked as passed.")
	}
	if len(r.Errors) != 0 {
		t.Error("Expected report to contain two errors.")
	}
	if len(r.Warnings) != 2 {
		t.Error("Expected report to not contain warnings.")
	}
	if r.Completion.Total != (float32(1) / float32(3)) {
		t.Error("Incorrect total completion ratio.")
	}
	if r.Completion.Required != float32(0) {
		t.Error("Incorrect required completion ratio.")
	}
	if r.Completion.Recommended > float32(0) {
		t.Error("Expected recommended completion ratio to be 0.")
	}
	if r.Completion.Optional > float32(1)/float32(3) {
		t.Error("Incorrect optional completion ratio")
	}
}

func Test_GetReport_incomplete(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	skipped := TestFromProto(&pb.Test{
		Name:             "skipped",
		ExpectationLevel: pb.Test_REQUIRED,
	})

	s.AddTest(skipped)

	s.AddAnswers("skipped", []string{"eggs", "bacon"})

	r := s.GetReport()
	if r.State != pb.ReportSessionResponse_INCOMPLETE {
		t.Error("Expected report to be marked as incomplete")
	}
	if len(r.Errors) != 1 {
		t.Error("Expected report to contain one error.")
	}
	if len(r.Warnings) != 0 {
		t.Error("Expected report to not contain warnings.")
	}
	if r.Completion.Total > float32(0) {
		t.Error("Incorrect total completion ratio.")
	}
	if r.Completion.Required > float32(0) {
		t.Error("Incorrect required completionration.")
	}
	if r.Completion.Recommended > float32(0) {
		t.Error("Expected recommended completion ratio to be 0.")
	}
	if r.Completion.Optional > float32(0) {
		t.Error("Expected optional completion ratio to be 0.")
	}
}

func Test_TestHandling(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}
	name := "bacon.required.test"

	_, ok := s.tests[name]
	if ok {
		t.Error("Expected no tests to be present in the session.")
	}

	test := NewTest(
		&pb.Test{Name: name, ExpectationLevel: pb.Test_REQUIRED},
		nil,
		pb.ReportSessionResponse_PASSED,
		[]string{})
	s.AddTest(test)

	got, ok := s.tests[name]
	if !ok {
		t.Error("Error getting the test from the session")
	}
	if test != got {
		t.Error("Expected to get the test that was added")
	}

	err := s.DeleteTest(name)
	if err != nil {
		t.Error("Expected test deletion to succeed")
	}

	_, ok = s.tests[name]
	if ok {
		t.Error("Expected no tests to be present in the session.")
	}
}

func Test_AddAnswers_not_found(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}
	if err := s.AddAnswers("test", []string{"Shouldn't be found"}); err == nil {
		t.Error("Adding answers to a not present test should have failed.")
	}
}

func Test_AddAnswers_incomplete_tests(t *testing.T) {
	testCases := []struct {
		initial    []string
		additional []string
	}{
		{[]string{}, []string{"bacon"}},
		{[]string{"eggs"}, []string{"bacon"}},
		{[]string{"grits"}, []string{}},
	}

	for _, testCase := range testCases {
		s := &sessionImpl{
			name:  "-",
			mu:    sync.Mutex{},
			tests: map[string]Test{},
		}
		s.AddTest(&testImpl{
			t:       &pb.Test{Name: "test"},
			answers: testCase.initial,
			state:   pb.ReportSessionResponse_PASSED,
		})
		s.AddAnswers("test", testCase.additional)
		if !reflect.DeepEqual(append(testCase.initial, testCase.additional...), s.tests["test"].GetAnswers()) {
			t.Error("Test.AddAnswers failed to append answers.")
		}
		if s.tests["test"].GetState() != pb.ReportSessionResponse_INCOMPLETE {
			t.Errorf("Test.AddAnswers expected to mark the test as incomplete.")
		}
	}
}

func Test_AddAnswers_failed_tests(t *testing.T) {
	testCases := []struct {
		initial    []string
		additional []string
	}{
		{[]string{}, []string{"bacon"}},
		{[]string{"eggs"}, []string{"bacon"}},
		{[]string{"grits"}, []string{}},
	}

	for _, testCase := range testCases {
		s := &sessionImpl{
			name:  "-",
			mu:    sync.Mutex{},
			tests: map[string]Test{},
		}
		s.AddTest(&testImpl{
			t:       &pb.Test{Name: "test"},
			answers: testCase.initial,
			state:   pb.ReportSessionResponse_FAILED,
		})
		s.AddAnswers("test", testCase.additional)
		if !reflect.DeepEqual(append(testCase.initial, testCase.additional...), s.tests["test"].GetAnswers()) {
			t.Error("Test.AddAnswers failed to append answers.")
		}
		if s.tests["test"].GetState() != pb.ReportSessionResponse_FAILED {
			t.Error("Test.AddAnswers expected to keep the test as failed.")
		}
	}
}

func Test_TestAnswers_not_found(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}
	if err := s.TestAnswers("test", []string{"Shouldn't be found"}); err == nil {
		t.Error("Adding answers to a not present test should have failed.")
	}
}

func Test_TestAnswers_already_failed(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	initialAnswers := []string{"Should", "not", "change"}

	tmp := make([]string, len(initialAnswers))
	copy(tmp, initialAnswers)
	s.AddTest(&testImpl{
		t:       &pb.Test{Name: "test"},
		answers: tmp,
		state:   pb.ReportSessionResponse_FAILED,
	})

	answers := [][]string{
		[]string{"bacon", "eggs"},
		[]string{},
	}
	for _, a := range answers {
		s.TestAnswers("test", a)
		if !reflect.DeepEqual(initialAnswers, s.tests["test"].GetAnswers()) {
			t.Error("Test.TestAnswers expected to not alter the test after already failing.")
		}
	}
}

func Test_TestAnswers_more_answers_than_expected(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	s.AddTest(TestFromProto(&pb.Test{Name: "test"}))
	s.TestAnswers("test", []string{"durian", "for", "breakfast"})
	if s.tests["test"].GetState() != pb.ReportSessionResponse_FAILED {
		t.Error("Session.TestAnswers should fail when testing more answers than expecrted.")
	}
	if s.tests["test"].GetIssue() == nil {
		t.Error("Session.TestAnswers should have made an issue for this failure.")
	}
}

func Test_TestAnswers_failed_tests(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	s.AddTest(TestFromProto(&pb.Test{Name: "test"}))

	s.AddAnswers("test", []string{"mangos", "for", "lunch!"})
	s.TestAnswers("test", []string{"mangos", "for", "dinner?"})

	if s.tests["test"].GetState() != pb.ReportSessionResponse_FAILED {
		t.Error("Test.TestAnswers should fail when testing an unexpected answer")
	}
	if s.tests["test"].GetIssue() == nil {
		t.Error("Test.TestAnswers should have made an issue for this failure.")
	}
}

func Test_TestAnswers_correct_tests(t *testing.T) {
	s := &sessionImpl{
		name:  "-",
		mu:    sync.Mutex{},
		tests: map[string]Test{},
	}

	s.AddTest(TestFromProto(&pb.Test{Name: "test"}))
	s.AddAnswers("test", []string{"mangos", "for", "lunch!"})
	s.TestAnswers("test", []string{"mangos"})
	if s.tests["test"].GetState() != pb.ReportSessionResponse_INCOMPLETE {
		t.Error("Test.TestAnswers should only pass when all expected values have been tested.")
	}
	if s.tests["test"].GetIssue() == nil {
		t.Error("Test.TestAnswers should only set issue to nil when the test passes.")
	}

	s.TestAnswers("test", []string{"for", "lunch!"})
	if s.tests["test"].GetState() != pb.ReportSessionResponse_PASSED {
		t.Error("Test.TestAnswers should pass when all expected values have been tested.")
	}
	if s.tests["test"].GetIssue() != nil {
		t.Error("Test.TestAnswers should set issue to nil when the test passes.")
	}
}
