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

package services

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Sequence_Empty(t *testing.T) {
	s := NewSequenceServer()

	seq, err := s.CreateSequence(context.Background(), &pb.CreateSequenceRequest{})
	if err != nil {
		t.Errorf("CreateSequence(empty): unexpected err %+v", err)
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	e, err := s.AttemptSequence(ctx, &pb.AttemptSequenceRequest{Name: seq.GetName()})
	if err != nil {
		t.Errorf("AttemptSequence(empty): unexpected err %+v", err)
	}

	if e == nil {
		t.Errorf("AttemptSequence(empty): unexpected nil Empty response")
	}

	r := report(seq.GetName())
	report, err := s.GetSequenceReport(context.Background(), &pb.GetSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("GetSequenceReport(empty): unexpected err %+v", err)
	}

	attempts := report.GetAttempts()
	if len(attempts) != 1 {
		t.Errorf("%s: expected number of attempts to be 1 but was %d", t.Name(), len(attempts))
	}

	a := attempts[0]
	ad := a.GetAttemptDeadline().AsTime()
	d, _ := ctx.Deadline()

	if !ad.Equal(d) {
		t.Errorf("%s: server deadline = %v client deadline = %v", t.Name(), ad, d)
	}
}

func Test_Sequence_Retry(t *testing.T) {
	s := NewSequenceServer()
	responses := []*pb.Sequence_Response{
		{
			Status: status.New(codes.Unavailable, "Unavailable").Proto(),
			Delay:  ptypes.DurationProto(1 * time.Second),
		},
		{
			Status: status.New(codes.Unavailable, "Unavailable").Proto(),
			Delay:  ptypes.DurationProto(2 * time.Second),
		},
		{
			Status: status.New(codes.OK, "OK").Proto(),
		},
	}

	seq, err := s.CreateSequence(context.Background(), &pb.CreateSequenceRequest{
		Sequence: &pb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Errorf("CreateSequence(retry): unexpected err %+v", err)
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	delay := 100 * time.Millisecond
	for n, r := range responses {
		res := status.FromProto(r.GetStatus())
		_, err = s.AttemptSequence(ctx, &pb.AttemptSequenceRequest{Name: seq.GetName()})
		if c := status.Code(err); c != res.Code() {
			t.Errorf("%s: status #%d was %v wanted %v", t.Name(), n, c, res.Code())
		}

		if n != len(responses)-1 {
			time.Sleep(delay)
			delay *= 2
		}
	}

	r := report(seq.GetName())
	report, err := s.GetSequenceReport(context.Background(), &pb.GetSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("GetSequenceReport(retry): unexpected err %+v", err)
	}

	attempts := report.GetAttempts()
	if len(attempts) != len(responses) {
		t.Errorf("%s: expected number of attempts to be %d but was %d", t.Name(), len(responses), len(attempts))
	}

	d, _ := ctx.Deadline()
	for n, a := range attempts {
		if a.GetAttemptNumber() != int32(n) {
			t.Errorf("%s: expected attempt #%d but was #%d", t.Name(), n, a.GetAttemptNumber())
		}

		ad := a.GetAttemptDeadline().AsTime()
		if !ad.Equal(d) {
			t.Errorf("%s: server deadline = %v client deadline = %v", t.Name(), ad, d)
		}

		if a.GetStatus().GetCode() != responses[n].GetStatus().GetCode() {
			t.Errorf("%s: expected response %v but was %v", t.Name(), responses[n].GetStatus().GetCode(), a.GetStatus().GetCode())
		}

		if n > 0 && a.GetAttemptDelay().AsDuration() <= attempts[n-1].GetAttemptDelay().AsDuration() {
			t.Errorf("%s: expected attempt delay: %v to be larger than previous: %v", t.Name(), a.GetAttemptDelay().AsDuration(), attempts[n-1].GetAttemptDelay().AsDuration())
		}
	}
}

func Test_GetSequenceReport_NotFound(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.GetSequenceReport(context.Background(), &pb.GetSequenceReportRequest{Name: "foo/bar/baz"})
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.NotFound.String(), c.String())
	}
}

func Test_GetSequenceReport_MissingName(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.GetSequenceReport(context.Background(), &pb.GetSequenceReportRequest{Name: ""})
	if c := status.Code(err); c != codes.InvalidArgument {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.InvalidArgument.String(), c.String())
	}
}

func Test_AttemptSequence_NotFound(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.AttemptSequence(context.Background(), &pb.AttemptSequenceRequest{Name: "foo/bar/baz"})
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.NotFound.String(), c.String())
	}
}

func Test_AttemptSequence_MissingName(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.AttemptSequence(context.Background(), &pb.AttemptSequenceRequest{Name: ""})
	if c := status.Code(err); c != codes.InvalidArgument {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.InvalidArgument.String(), c.String())
	}
}
