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
	"context"
	"testing"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func Test_unaryTest_GetName(t *testing.T) {
	ut := NewUnaryTest("sessions/-")
	got := ut.GetName()
	want := "sessions/-/gapic.v1p0.unary_unary.ok"
	if got != want {
		t.Errorf("GetName: got %s, want %s", got, want)
	}
}

func Test_unaryTest_GetExpectationLevel(t *testing.T) {
	ut := NewUnaryTest("sessions/-")
	got := ut.GetExpectationLevel()
	want := pb.Test_REQUIRED
	if got != want {
		t.Errorf("GetExpectationLevel: got %+v, want %+v", got, want)
	}
}

func Test_unaryTest_GetDescription(t *testing.T) {
	ut := NewUnaryTest("sessions/-")
	got := ut.GetDescription()
	if got == "" {
		t.Errorf("GetDescription: expected non-empty string")
	}
}

func Test_unaryTest_GetBlueprints(t *testing.T) {
	ut := NewUnaryTest("sessions/-")
	got := ut.GetBlueprints()
	if len(got) != 0 {
		t.Errorf("GetBluprints: expected empty list")
	}
}

func Test_unaryTest_GetIssue_doesntVerifyForOtherSessions(t *testing.T) {
	ut := &unaryTest{
		sessionName: "sessions/-",
		responses:   []interface{}{},
	}
	other := NewUnaryTest("sessions/1")
	resp := &pb.EchoResponse{Content: "hello world"}
	ut.ObserveUnary(
		context.Background(),
		nil,
		resp,
		serverInfo("/google.showcase.v1beta1.Echo/Echo"),
		nil)
	data, _ := proto.Marshal(resp)
	ut.ObserveUnary(
		context.Background(),
		&pb.VerifyTestRequest{
			Name:   other.GetName(),
			Answer: data,
		},
		&pb.VerifyTestResponse{},
		serverInfo("/google.showcase.v1beta1.Testing/VerifyTest"),
		nil)

	got := ut.GetIssue()
	want := &pb.Issue{
		Type:        pb.Issue_PENDING,
		Severity:    pb.Issue_ERROR,
		Description: "This test has not been verified.",
	}
	if !proto.Equal(got, want) {
		t.Errorf("GetIssue: got %+v, want %+v", got, want)
	}
}

func Test_unaryTest_GetIssue_verified(t *testing.T) {
	ut := &unaryTest{
		sessionName: "sessions/-",
		responses:   []interface{}{},
	}
	resp := &pb.EchoResponse{Content: "hello world"}
	ut.ObserveUnary(
		context.Background(),
		nil,
		resp,
		serverInfo("/google.showcase.v1beta1.Echo/Echo"),
		nil)
	data, _ := proto.Marshal(resp)
	ut.ObserveUnary(
		context.Background(),
		&pb.VerifyTestRequest{
			Name:   ut.GetName(),
			Answer: data,
		},
		&pb.VerifyTestResponse{},
		serverInfo("/google.showcase.v1beta1.Testing/VerifyTest"),
		nil)
	got := ut.GetIssue()

	if got != nil {
		t.Errorf("GetIssue: got %+v, want nil", got)
	}
}

func Test_unaryTest_GetIssue_failedVerification(t *testing.T) {
	ut := &unaryTest{
		sessionName: "sessions/-",
		responses:   []interface{}{},
	}
	ut.ObserveUnary(
		context.Background(),
		nil,
		&pb.EchoResponse{Content: "hello world"},
		serverInfo("/google.showcase.v1beta1.Echo/Echo"),
		nil)
	data, _ := proto.Marshal(&pb.EchoResponse{Content: "hi world"})
	ut.ObserveUnary(
		context.Background(),
		&pb.VerifyTestRequest{
			Name:   ut.GetName(),
			Answer: data,
		},
		&pb.VerifyTestResponse{},
		serverInfo("/google.showcase.v1beta1.Testing/VerifyTest"),
		nil)
	got := ut.GetIssue()
	want := &pb.Issue{
		Type:        pb.Issue_INCORRECT_CONFIRMATION,
		Severity:    pb.Issue_ERROR,
		Description: "An incorrect answer was supplied to verify this test.",
	}

	if !proto.Equal(got, want) {
		t.Errorf("GetIssue: got %+v, want %+v", got, want)
	}
}

func Test_unaryTest_GetIssue_needsVerification(t *testing.T) {
	ut := &unaryTest{
		sessionName: "sessions/-",
		responses:   []interface{}{},
	}
	ut.ObserveUnary(
		context.Background(),
		nil,
		&pb.EchoResponse{Content: "hello world"},
		serverInfo("/google.showcase.v1beta1.Echo/Echo"),
		nil)
	got := ut.GetIssue()
	want := &pb.Issue{
		Type:        pb.Issue_PENDING,
		Severity:    pb.Issue_ERROR,
		Description: "This test has not been verified.",
	}
	if !proto.Equal(got, want) {
		t.Errorf("GetIssue: got %+v, want %+v", got, want)
	}
}

func Test_unaryTest_GetIssue_skipped(t *testing.T) {
	ut := &unaryTest{
		sessionName: "sessions/-",
		responses:   []interface{}{},
	}
	got := ut.GetIssue()
	want := &pb.Issue{
		Type:        pb.Issue_SKIPPED,
		Severity:    pb.Issue_ERROR,
		Description: "This test has not been started. Make a unary request to start this test.",
	}
	if !proto.Equal(got, want) {
		t.Errorf("GetIssue: got %+v, want %+v", got, want)
	}
}

func serverInfo(methodName string) *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{FullMethod: methodName}
}
