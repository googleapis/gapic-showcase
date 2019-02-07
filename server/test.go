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

// Test represents a test case that is run. This interfaces exposes the
// properties of the test as well as the observers used to run this test.
//
// A Test will also implement at least one of UnaryObserver, StreamRequestObserver,
// StreamResponseOvserver in order to track requests made to the showcase server.
type Test interface {
	GetName() string
	GetExpectationLevel() pb.Test_ExpectationLevel
	GetDescription() string
	GetBlueprints() []*pb.Test_Blueprint
	GetIssue() *pb.Issue
}

// TestProto returns a proto representation of the Test.
func TestProto(t Test) *pb.Test {
	return &pb.Test{
		Name:             t.GetName(),
		ExpectationLevel: t.GetExpectationLevel(),
		Description:      t.GetDescription(),
		Blueprints:       t.GetBlueprints(),
	}
}

// TestRunProto returns a proto representation of a test run.
func TestRunProto(t Test) *pb.TestRun {
	return &pb.TestRun{
		Test:  t.GetName(),
		Issue: t.GetIssue(),
	}
}
