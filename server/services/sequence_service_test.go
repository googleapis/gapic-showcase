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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestSequenceEmpty(t *testing.T) {
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

func TestSequenceRetry(t *testing.T) {
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
		if got, want := a.GetAttemptNumber(), int32(n); got != want {
			t.Errorf("%s: expected attempt #%d but was #%d", t.Name(), want, got)
		}

		if got, want := a.GetAttemptDeadline().AsTime(), d; !got.Equal(want) {
			t.Errorf("%s: server deadline = %v client deadline = %v", t.Name(), got, want)
		}

		if got, want := a.GetStatus().GetCode(), responses[n].GetStatus().GetCode(); got != want {
			t.Errorf("%s: expected response %v but was %v", t.Name(), want, got)
		}

		// Check that perceived delay between attempts was changing as expected.
		if n > 0 {
			if cur, prev := a.GetAttemptDelay().AsDuration(), attempts[n-1].GetAttemptDelay().AsDuration(); cur <= prev {
				t.Errorf("%s: expected attempt delay: %v to be larger than previous: %v", t.Name(), cur, prev)
			}
		}
	}
}

func TestSequenceOutOfRange(t *testing.T) {
	s := NewSequenceServer()

	seq, err := s.CreateSequence(context.Background(), &pb.CreateSequenceRequest{})
	if err != nil {
		t.Errorf("CreateSequence(out_of_range): unexpected err %+v", err)
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	e, err := s.AttemptSequence(ctx, &pb.AttemptSequenceRequest{Name: seq.GetName()})
	if err != nil {
		t.Errorf("AttemptSequence(out_of_range): unexpected err %+v", err)
	}

	if e == nil {
		t.Errorf("AttemptSequence(out_of_range): unexpected nil Empty response")
	}

	_, err = s.AttemptSequence(ctx, &pb.AttemptSequenceRequest{Name: seq.GetName()})
	if c := status.Code(err); c != codes.OutOfRange {
		t.Errorf("%s: status was %v wanted %v", t.Name(), c, codes.OutOfRange)
	}

	_, err = s.AttemptSequence(ctx, &pb.AttemptSequenceRequest{Name: seq.GetName()})
	if c := status.Code(err); c != codes.OutOfRange {
		t.Errorf("%s: status was %v wanted %v", t.Name(), c, codes.OutOfRange)
	}

	r := report(seq.GetName())
	report, err := s.GetSequenceReport(context.Background(), &pb.GetSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("GetSequenceReport(out_of_range): unexpected err %+v", err)
	}

	attempts := report.GetAttempts()
	if len(attempts) != 3 {
		t.Errorf("%s: expected number of attempts to be 3 but was %d", t.Name(), len(attempts))
	}

	a := attempts[0]
	ad := a.GetAttemptDeadline().AsTime()
	d, _ := ctx.Deadline()

	if !ad.Equal(d) {
		t.Errorf("%s: server deadline = %v client deadline = %v", t.Name(), ad, d)
	}
}

func TestGetSequenceReportNotFound(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.GetSequenceReport(context.Background(), &pb.GetSequenceReportRequest{Name: "foo/bar/baz"})
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.NotFound, c)
	}
}

func TestGetSequenceReportMissingName(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.GetSequenceReport(context.Background(), &pb.GetSequenceReportRequest{Name: ""})
	if c := status.Code(err); c != codes.InvalidArgument {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.InvalidArgument, c)
	}
}

func TestAttemptSequenceNotFound(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.AttemptSequence(context.Background(), &pb.AttemptSequenceRequest{Name: "foo/bar/baz"})
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.NotFound, c)
	}
}

func TestAttemptSequenceMissingName(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.AttemptSequence(context.Background(), &pb.AttemptSequenceRequest{Name: ""})
	if c := status.Code(err); c != codes.InvalidArgument {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.InvalidArgument, c)
	}
}

type mockStreamSequence struct {
	exp   []string
	head  []string
	trail []string
	t     *testing.T
	pb.SequenceService_AttemptStreamingSequenceServer
}

func (m *mockStreamSequence) Context() context.Context {
	return appendTestOutgoingMetadata(context.Background(), &mockSTS{stream: m, t: m.t})
}

func (m *mockStreamSequence) SetTrailer(md metadata.MD) {
	m.trail = append(m.trail, md.Get("showcase-trailer")...)
}

func (m *mockStreamSequence) SetHeader(md metadata.MD) error {
	m.head = append(m.head, md.Get("x-goog-request-params")...)
	return nil
}

func (m *mockStreamSequence) verify(expectHeadersAndTrailers bool) {
	if len(m.exp) > 0 {
		m.t.Errorf("Expand did not stream all expected values. %d expected values remaining.", len(m.exp))
	}
	if expectHeadersAndTrailers && !reflect.DeepEqual([]string{"show", "case"}, m.trail) && !reflect.DeepEqual([]string{"showcaseHeader", "anotherHeader"}, m.head) {
		m.t.Errorf("Expand did not get all expected headers and trailers.\nGot these headers: %+v\nGot these trailers: %+v", m.head, m.trail)
	}
}

func TestStreamingSequenceEmpty(t *testing.T) {
	s := NewSequenceServer()

	seq, err := s.CreateStreamingSequence(context.Background(), &pb.CreateStreamingSequenceRequest{})
	if err != nil {
		t.Errorf("CreateSequence(empty): unexpected err %+v", err)
	}

	stream := &mockStreamSequence{}

	timeout := 5 * time.Second
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = s.AttemptStreamingSequence(&pb.AttemptStreamingSequenceRequest{Name: seq.GetName()}, stream)
	if err != nil {
		t.Errorf("AttemptSequence(empty): unexpected err %+v", err)
	}

	r := streamingReport(seq.GetName())
	report, err := s.GetStreamingSequenceReport(context.Background(), &pb.GetStreamingSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("GetSequenceReport(empty): unexpected err %+v", err)
	}

	attempts := report.GetAttempts()

	if len(attempts) != 1 {
		t.Errorf("%s: expected number of attempts to be 1 but was %d", t.Name(), len(attempts))
	}
}

func TestStreamingSequenceRetry(t *testing.T) {
	s := NewSequenceServer()
	responses := []*pb.StreamingSequence_Response{
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

	seq, err := s.CreateStreamingSequence(context.Background(), &pb.CreateStreamingSequenceRequest{
		Streamingsequence: &pb.StreamingSequence{Responses: responses, Content: "Hello World"},
	})
	if err != nil {
		t.Errorf("CreateSequence(retry): unexpected err %+v", err)
	}

	timeout := 5 * time.Second
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	delay := 100 * time.Millisecond
	stream := &mockStreamSequence{}

	for n, r := range responses {
		res := status.FromProto(r.GetStatus())
		err = s.AttemptStreamingSequence(&pb.AttemptStreamingSequenceRequest{Name: seq.GetName()}, stream)
		if c := status.Code(err); c != res.Code() {
			t.Errorf("%s: status #%d was %v wanted %v", t.Name(), n, c, res.Code())
		}

		if n != len(responses)-1 {
			time.Sleep(delay)
			delay *= 2
		}
	}

	r := streamingReport(seq.GetName())
	report, err := s.GetStreamingSequenceReport(context.Background(), &pb.GetStreamingSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("GetSequenceReport(retry): unexpected err %+v", err)
	}

	attempts := report.GetAttempts()
	if len(attempts) != len(responses) {
		t.Errorf("%s: expected number of attempts to be %d but was %d", t.Name(), len(responses), len(attempts))
	}

	for n, a := range attempts {
		if got, want := a.GetAttemptNumber(), int32(n); got != want {
			t.Errorf("%s: expected attempt #%d but was #%d", t.Name(), want, got)
		}

		if got, want := a.GetStatus().GetCode(), responses[n].GetStatus().GetCode(); got != want {
			t.Errorf("%s: expected response %v but was %v", t.Name(), want, got)
		}
	}
}

func TestStreamingSequenceOutOfRange(t *testing.T) {
	s := NewSequenceServer()

	seq, err := s.CreateStreamingSequence(context.Background(), &pb.CreateStreamingSequenceRequest{})
	if err != nil {
		t.Errorf("CreateSequence(out_of_range): unexpected err %+v", err)
	}

	timeout := 5 * time.Second
	stream := &mockStreamSequence{}

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = s.AttemptStreamingSequence(&pb.AttemptStreamingSequenceRequest{Name: seq.GetName()}, stream)
	if err != nil {
		t.Errorf("AttemptSequence(out_of_range): unexpected err %+v", err)
	}

	err = s.AttemptStreamingSequence(&pb.AttemptStreamingSequenceRequest{Name: seq.GetName()}, stream)
	if c := status.Code(err); c != codes.OutOfRange {
		t.Errorf("%s: status was %v wanted %v", t.Name(), c, codes.OutOfRange)
	}

	err = s.AttemptStreamingSequence(&pb.AttemptStreamingSequenceRequest{Name: seq.GetName()}, stream)
	if c := status.Code(err); c != codes.OutOfRange {
		t.Errorf("%s: status was %v wanted %v", t.Name(), c, codes.OutOfRange)
	}

	r := streamingReport(seq.GetName())
	report, err := s.GetStreamingSequenceReport(context.Background(), &pb.GetStreamingSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("GetSequenceReport(out_of_range): unexpected err %+v", err)
	}

	attempts := report.GetAttempts()
	if len(attempts) != 3 {
		t.Errorf("%s: expected number of attempts to be 3 but was %d", t.Name(), len(attempts))
	}
}

func TestAttemptStreamingSequenceNotFound(t *testing.T) {
	s := NewSequenceServer()
	stream := &mockStreamSequence{exp: strings.Fields(""), t: t}
	attemptRequest := &pb.AttemptStreamingSequenceRequest{Name: "sequences/0"}
	err := s.AttemptStreamingSequence(attemptRequest, stream)
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.NotFound, c)
	}
}

func TestAttemptStreamingSequenceMissingName(t *testing.T) {
	s := NewSequenceServer()
	stream := &mockStreamSequence{exp: strings.Fields(""), t: t}
	err := s.AttemptStreamingSequence(&pb.AttemptStreamingSequenceRequest{Name: ""}, stream)
	if c := status.Code(err); c != codes.InvalidArgument {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.InvalidArgument, c)
	}
}

func TestGetStreamingSequenceReportMissingName(t *testing.T) {
	s := NewSequenceServer()
	_, err := s.GetStreamingSequenceReport(context.Background(), &pb.GetStreamingSequenceReportRequest{Name: ""})
	if c := status.Code(err); c != codes.InvalidArgument {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.InvalidArgument, c)
	}
}
