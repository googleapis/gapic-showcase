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

type Session interface {
	GetName() string
	GetReport() (*pb.ReportSessionResponse, error)
	AddTest(t *pb.Test) error
	DeleteTest(name string) error
}

var instance Session
var once sync.Once

func GetSessionSingleton() Session {
	once.Do(func() {
    instance = &sessionImpl{
			name: "-",
			mu: &sync.Mutex{},
		}
  })
  return instance
}

type sessionImpl struct {
	name 	string

  mu  	*sync.Mutex
	tests map[string]Test
}

func (s *sessionImpl) GetName() string {
	return s.name
}

func (s *sessionImpl) GetReport() (*pb.ReportSessionResponse, error) {
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
  if (numTests > 0) {
    totalRatio = float32(numTests - numInvalid) / float32(numTests)
  }

  requiredRatio := float32(0.0)
  if (numRequired > 0) {
    requiredRatio = float32(numRequired - numRequiredFailed - numRequiredSkipped) /
                    float32(numRequired)
  }

  recommendedRatio := float32(0.0)
  if (numRecommended > 0) {
    recommendedRatio = float32(numRecommended - numRecommendedFailed - numRecommendedSkipped) /
                       float32(numRecommended)
  }

  optionalRatio := float32(0.0)
  if (numOptional > 0) {
    optionalRatio = float32(numOptional - numOptionalFailed - numOptionalSkipped) /
                    float32(numOptional)
  }

  report := &pb.ReportSessionResponse{
    State: state,
    Errors: errors,
    Warnings: warnings,
    Completion: &pb.ReportSessionResponse_Completion{
      Total: totalRatio,
      Required: requiredRatio,
      Recommended: recommendedRatio,
      Optional: optionalRatio,
    },
  }
	return report, nil
}

func (s *sessionImpl) AddTest(t *pb.Test) error {
  s.mu.Lock()
  defer s.mu.Unlock()

	s.tests[t.Name] = NewTestFromProto(t)
	return nil
}

func (s *sessionImpl) DeleteTest(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tests, name)
	return nil
}
