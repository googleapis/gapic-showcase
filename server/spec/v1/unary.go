// Copyright 2019 Google LLC
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

package v1

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type unaryTest struct {
	sessionName string

	mu                    sync.Mutex
	responses             []interface{}
	attemptedVerification bool
	verified              bool
}

func NewUnaryTest(sessionName string) server.Test {
	return &unaryTest{
		sessionName: sessionName,
		responses:   []interface{}{},
	}
}

func (t *unaryTest) GetName() string {
	return fmt.Sprintf("%s/gapic.v1p0.unary_unary.ok", t.sessionName)
}

func (t *unaryTest) GetExpectationLevel() pb.Test_ExpectationLevel {
	return pb.Test_REQUIRED
}

func (t *unaryTest) GetDescription() string {
	return "The generator generates a unary-unary RPC, is able to call it," +
		", and is able to handle an OK server response."
}

// TODO: Figure out a good way to represent all unary methods as a bluerprint.
func (t *unaryTest) GetBlueprints() []*pb.Test_Blueprint {
	return []*pb.Test_Blueprint{}
}

func (t *unaryTest) GetIssue() *pb.Issue {
	if t.verified {
		return nil
	}
	if t.attemptedVerification {
		return &pb.Issue{
			Type:        pb.Issue_INCORRECT_CONFIRMATION,
			Severity:    pb.Issue_ERROR,
			Description: "An incorrect answer was supplied to verify this test.",
		}
	}

	if len(t.responses) > 0 {
		return &pb.Issue{
			Type:        pb.Issue_PENDING,
			Severity:    pb.Issue_ERROR,
			Description: "This test has not been verified.",
		}
	}
	return &pb.Issue{
		Type:        pb.Issue_SKIPPED,
		Severity:    pb.Issue_ERROR,
		Description: "This test has not been started. Make a unary request to start this test.",
	}
}

func (t *unaryTest) ObserveUnary(
	ctx context.Context,
	req interface{},
	resp interface{},
	info *grpc.UnaryServerInfo,
	err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if resp != nil {
		t.responses = append(t.responses, resp)
	}

	if info.FullMethod == "/google.showcase.v1beta1.Testing/VerifyTest" {
		// Only validate for this test.
		vtReq := req.(*pb.VerifyTestRequest)
		if vtReq.GetName() != t.GetName() {
			return
		}

		t.attemptedVerification = true
		for _, r := range t.responses {
			bs, _ := proto.Marshal(r.(proto.Message))

			if bytes.Equal(bs, vtReq.GetAnswer()) {
				t.verified = true
			}
		}
	}

}
