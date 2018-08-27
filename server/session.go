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

var instance Session
var once sync.Once

func GetSessionSingleton() Session {
	once.Do(func() {
		instance = &sessionImpl{
			name:  "-",
			mu:    &sync.Mutex{},
			tests: map[string]Test{},
		}
	})
	return instance
}

type sessionImpl struct {
	name string

	mu    *sync.Mutex
	tests map[string]Test
}

func (s *sessionImpl) GetName() string {
	return s.name
}

func (s *sessionImpl) GetReport() *pb.ReportSessionResponse {
	numRequired := 0
	numRequiredSkipped := 0
	numRequiredFailed := 0

	numRecommended := 0
	numRecommendedSkipped := 0
	numRecommendedFailed := 0

	numOptional := 0
	numOptionalSkipped := 0
	numOptionalFailed := 0

	errors := []*pb.ReportSessionResponse_Issue{}
	warnings := []*pb.ReportSessionResponse_Issue{}

	for _, test := range s.tests {
		expectationLevel := test.GetExpectationLevel()
		issue := test.GetIssue()

		if expectationLevel == pb.Test_REQUIRED {
			numRequired = numRequired + 1

			if issue != nil {
				errors = append(errors, issue)

				if issue.Type == pb.ReportSessionResponse_Issue_SKIPPED {
					numRequiredSkipped = numRequiredSkipped + 1
				} else {
					numRequiredFailed = numRequiredFailed + 1
				}
			}
		}

		if expectationLevel == pb.Test_RECOMMENDED {
			numRecommended = numRecommended + 1

			if issue != nil {
				warnings = append(warnings, issue)

				if issue.Type == pb.ReportSessionResponse_Issue_SKIPPED {
					numRecommendedSkipped = numRecommendedSkipped + 1
				} else {
					numRecommendedFailed = numRecommendedFailed + 1
				}
			}
		}

		if expectationLevel == pb.Test_OPTIONAL {
			numOptional = numOptional + 1

			if issue != nil {
				warnings = append(warnings, issue)

				if issue.Type == pb.ReportSessionResponse_Issue_SKIPPED {
					numOptionalSkipped = numOptionalSkipped + 1
				} else {
					numOptionalFailed = numOptionalFailed + 1
				}
			}
		}
	}

	var state pb.ReportSessionResponse_State
	if numRequired == 0 || (numRequiredFailed == 0 && numRequiredSkipped == 0) {
		state = pb.ReportSessionResponse_PASSED
	} else if numRequiredFailed == 0 {
		state = pb.ReportSessionResponse_INCOMPLETE
	} else {
		state = pb.ReportSessionResponse_FAILED
	}

	numTests := numRequired + numRecommended + numOptional
	numInvalid := numRequiredFailed + numRequiredSkipped +
		numRecommendedFailed + numRecommendedSkipped +
		numOptionalFailed + numOptionalSkipped

	totalRatio := float32(0.0)
	if numTests > 0 {
		totalRatio = float32(numTests-numInvalid) / float32(numTests)
	}

	requiredRatio := float32(0.0)
	if numRequired > 0 {
		requiredRatio = float32(numRequired-numRequiredFailed-numRequiredSkipped) /
			float32(numRequired)
	}

	recommendedRatio := float32(0.0)
	if numRecommended > 0 {
		recommendedRatio = float32(numRecommended-numRecommendedFailed-numRecommendedSkipped) /
			float32(numRecommended)
	}

	optionalRatio := float32(0.0)
	if numOptional > 0 {
		optionalRatio = float32(numOptional-numOptionalFailed-numOptionalSkipped) /
			float32(numOptional)
	}

	report := &pb.ReportSessionResponse{
		State:    state,
		Errors:   errors,
		Warnings: warnings,
		Completion: &pb.ReportSessionResponse_Completion{
			Total:       totalRatio,
			Required:    requiredRatio,
			Recommended: recommendedRatio,
			Optional:    optionalRatio,
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
