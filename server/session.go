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
	"errors"
	"sync"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

type Session interface {
	GetName() string
	GetReport() *pb.ReportSessionResponse
	AddTest(t Test) error
	DeleteTest(name string) error
	AddAnswers(name string, answers []string) error
	TestAnswers(name string, answers []string) error
}

var Instance Session = &sessionImpl{
	name:  "-",
	mu:    sync.Mutex{},
	tests: map[string]Test{},
}


func GetSessionSingleton() Session {
	return Instance
}

type sessionImpl struct {
	name string

	mu    sync.Mutex
	tests map[string]Test
}

func (s *sessionImpl) GetName() string {
	return s.name
}

type result struct {
	test int
	skipped int
	failed int
	issues []*pb.ReportSessionResponse_Issue
}

func (r *result) ratio() float32 {
	if r.test == 0 {
		return float32(0)
	}
	return float32(r.test - r.skipped - r.failed) / float32(r.test)
}

func (s *sessionImpl) GetReport() *pb.ReportSessionResponse {
	resultTotal := result{0, 0, 0, []*pb.ReportSessionResponse_Issue{}}
	resultRequired := result{0, 0, 0, []*pb.ReportSessionResponse_Issue{}}
	resultRecommended := result{0, 0, 0, []*pb.ReportSessionResponse_Issue{}}
	resultOptional := result{0, 0, 0, []*pb.ReportSessionResponse_Issue{}}

	for _, test := range s.tests {
		expLvl := test.GetExpectationLevel()
		issue := test.GetIssue()

		var r *result
		switch expLvl {
			case pb.Test_REQUIRED:
			  r = &resultRequired
			case pb.Test_RECOMMENDED:
				r = &resultRecommended
			default:
				r = &resultOptional
		}

		r.test++
		resultTotal.test++
		if issue != nil {
		  r.issues = append(r.issues, issue)
		  if issue.Type == pb.ReportSessionResponse_Issue_SKIPPED {
		    r.skipped++
				resultTotal.skipped++
		  } else {
		    r.failed++
				resultTotal.skipped++
		  }
		}
	}

	var state pb.ReportSessionResponse_State
	if resultRequired.failed == 0 && resultRequired.skipped == 0 {
		state = pb.ReportSessionResponse_PASSED
	} else if resultRequired.failed == 0 {
		state = pb.ReportSessionResponse_INCOMPLETE
	} else {
		state = pb.ReportSessionResponse_FAILED
	}

	report := &pb.ReportSessionResponse{
		State:    state,
		Errors:   resultRequired.issues,
		Warnings: append(resultRecommended.issues, resultOptional.issues...),
		Completion: &pb.ReportSessionResponse_Completion{
			Total:       resultTotal.ratio(),
			Required:    resultRequired.ratio(),
			Recommended: resultRecommended.ratio(),
			Optional:    resultOptional.ratio(),
		},
	}
	return report
}

func (s *sessionImpl) AddTest(t Test) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tests[t.GetName()] = t
	return nil
}

func (s *sessionImpl) DeleteTest(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tests, name)
	return nil
}

func (s *sessionImpl) AddAnswers(name string, answers []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tests[name]
	if !ok {
		return errors.New("Could not find test: " + name)
	}

	answers = append(t.GetAnswers(), answers...)
	state := t.GetState()
	if !t.HasFailed() {
		state = pb.ReportSessionResponse_INCOMPLETE
	}

	s.tests[name] = NewTest(
		t.GetProto(),
		t.GetIssue(),
		state,
		answers)
	return nil
}

func (s *sessionImpl) TestAnswers(name string, answers []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tests[name]
	if !ok {
		return errors.New("Could not find test: " + name)
	}

	if t.HasFailed() {
		return nil
	}

	testAnswers := t.GetAnswers()
	issue := t.GetIssue()
	state := t.GetState()
	for _, a := range answers {
		if len(testAnswers) == 0 {
			issue = createIssue(
				t.GetProto(),
				pb.ReportSessionResponse_Issue_INCORRECT_CONFIRMATION,
				"More answers registered than Showcase expected.")
			state = pb.ReportSessionResponse_FAILED
			break
		}

		expected := testAnswers[0]
		testAnswers = testAnswers[1:]
		if a != expected {
			issue = createIssue(
				t.GetProto(),
				pb.ReportSessionResponse_Issue_INCORRECT_CONFIRMATION,
				"Incorrect answer registered. Expected '%s' but got '%s'.")
			state = pb.ReportSessionResponse_FAILED
			break
		}
		if len(testAnswers) == 0 {
			issue = nil
			state = pb.ReportSessionResponse_PASSED
		}
	}
	s.tests[name] = NewTest(
		t.GetProto(),
		issue,
		state,
		testAnswers)

	return nil
}

func createIssue(t *pb.Test, issueType pb.ReportSessionResponse_Issue_Type, desc string) *pb.ReportSessionResponse_Issue {
	var severity pb.ReportSessionResponse_Issue_Severity
	if t.ExpectationLevel <= pb.Test_REQUIRED {
		severity = pb.ReportSessionResponse_Issue_ERROR
	} else {
		severity = pb.ReportSessionResponse_Issue_WARNING
	}
	return &pb.ReportSessionResponse_Issue{
		Test:        t.Name,
		Type:        issueType,
		Severity:    severity,
		Description: desc,
	}
}
