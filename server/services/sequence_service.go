// Copyright 2020 Google LLC
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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewSequenceServer returns a new SequenceServer for the Showcase API.
func NewSequenceServer() pb.SequenceServiceServer {
	return &sequenceServerImpl{
		token:              server.NewTokenGenerator(),
		sequences:          sync.Map{},
		reports:            sync.Map{},
		streamingsequences: sync.Map{},
		streamingreports:   sync.Map{},
	}
}

type sequenceServerImpl struct {
	uid   server.UniqID
	token server.TokenGenerator

	sequences sync.Map
	reports   sync.Map

	streamingsequences sync.Map
	streamingreports   sync.Map
	sentContent        string
}

func (s *sequenceServerImpl) CreateSequence(ctx context.Context, in *pb.CreateSequenceRequest) (*pb.Sequence, error) {
	seq := clone(in.GetSequence())

	// Assign Name.
	id := s.uid.Next()
	seq.Name = fmt.Sprintf("sequences/%d", id)
	report := &pb.SequenceReport{
		Name: report(seq.GetName()),
	}

	s.sequences.Store(seq.GetName(), seq)
	s.reports.Store(report.GetName(), report)

	return seq, nil
}

func (s *sequenceServerImpl) AttemptSequence(ctx context.Context, in *pb.AttemptSequenceRequest) (*empty.Empty, error) {
	received := time.Now()
	name := in.GetName()
	if name == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"The field `name` is required.")
	}

	// Retrieve Sequence and associated SequenceReport.
	i, ok := s.sequences.Load(name)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"The Sequence with %q does not exist.",
			name,
		)
	}
	seq := i.(*pb.Sequence)

	i, _ = s.reports.Load(report(name))
	rep, _ := i.(*pb.SequenceReport)

	// Retrieve the attempt deadline.
	deadline, _ := ctx.Deadline()
	dpb, err := ptypes.TimestampProto(deadline)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	// Get the number of attempts, which coincides with this attempt's number.
	n := len(rep.Attempts)

	// Prepare the attempt response defined by the Sequence.
	st := status.New(codes.OK, "Successful attempt")
	var delay time.Duration
	responses := seq.GetResponses()
	if l := len(responses); l > 0 && n < l {
		resp := responses[n]
		delay = resp.GetDelay().AsDuration()
		st = status.FromProto(resp.GetStatus())
	} else if n > l {
		st = status.New(codes.OutOfRange, "Attempt exceeded predefined responses")
	}

	// A delay of 0 returns immediately.
	time.Sleep(delay)

	// Calculate the perceived delay since the last RPC attempt.
	attDelay := &duration.Duration{}
	if n > 0 {
		prev := rep.GetAttempts()[n-1]
		respTime := prev.GetResponseTime()
		d := received.Sub(respTime.AsTime())
		attDelay = ptypes.DurationProto(d)
	}

	// Clock the time that the server is sending the response
	responseTime := time.Now()
	rpb, err := ptypes.TimestampProto(responseTime)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	rep.Attempts = append(rep.Attempts, &pb.SequenceReport_Attempt{
		AttemptNumber:   int32(n),
		AttemptDeadline: dpb,
		ResponseTime:    rpb,
		AttemptDelay:    attDelay,
		Status:          st.Proto(),
	})

	return &empty.Empty{}, st.Err()
}

func (s *sequenceServerImpl) GetSequenceReport(ctx context.Context, in *pb.GetSequenceReportRequest) (*pb.SequenceReport, error) {
	name := in.GetName()
	if name == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"The field `name` is required.")
	}

	report, ok := s.reports.Load(name)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"The sequence report with %q does not exist.",
			name,
		)
	}

	return report.(*pb.SequenceReport), nil
}

func report(n string) string {
	return fmt.Sprintf("%s/sequenceReport", n)
}

func clone(s *pb.Sequence) *pb.Sequence {
	r := make([]*pb.Sequence_Response, len(s.GetResponses()))
	copy(r, s.GetResponses())

	return &pb.Sequence{
		Name:      s.GetName(),
		Responses: r,
	}
}

func (s *sequenceServerImpl) CreateStreamingSequence(ctx context.Context, in *pb.CreateStreamingSequenceRequest) (*pb.StreamingSequence, error) {
	seq := cloneStreamingSequence(in.GetStreamingsequence())

	// Assign Name.
	id := s.uid.Next()
	seq.Name = fmt.Sprintf("streamingsequences/%d", id)
	report := &pb.StreamingSequenceReport{
		Name: streamingReport(seq.GetName()),
	}

	s.streamingsequences.Store(seq.GetName(), seq)
	s.streamingreports.Store(report.GetName(), report)

	return seq, nil
}

func (s *sequenceServerImpl) AttemptStreamingSequence(in *pb.AttemptStreamingSequenceRequest, stream pb.SequenceService_AttemptStreamingSequenceServer) error {
	received := time.Now()
	name := in.GetName()
	if name == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `name` is required.")
	}

	// Retrieve Sequence and associated SequenceReport.
	i, ok := s.streamingsequences.Load(name)
	if !ok {
		return status.Errorf(
			codes.NotFound,
			"The Sequence with %q does not exist.",
			name,
		)
	}
	seq := i.(*pb.StreamingSequence)

	i, _ = s.streamingreports.Load(streamingReport(name))
	rep, _ := i.(*pb.StreamingSequenceReport)

	// Get the number of attempts, which coincides with this attempt's number.
	n := len(rep.Attempts)

	// Prepare the attempt response defined by the Sequence.
	st := status.New(codes.OK, "Successful attempt")
	sendStatusAtIndex := 0
	var delay time.Duration
	responses := seq.GetResponses()
	content := strings.Fields(seq.GetContent())
	if l := len(responses); l > 0 && n < l {
		resp := responses[n]
		delay = resp.GetDelay().AsDuration()
		st = status.FromProto(resp.GetStatus())
		sendStatusAtIndex = int(resp.SendStatusAtIndex) - len(strings.Fields(s.sentContent))

		if s.sentContent != "" {
			wordsWritten := strings.Fields(s.sentContent)
			content = content[len(wordsWritten):]
		}

	} else if n > l {
		st = status.New(codes.OutOfRange, "Attempt exceeded predefined responses")
	}

	for idx, word := range content {
		if idx >= sendStatusAtIndex {
			break
		}
		s.sentContent += word + " "
		err := stream.Send(&pb.AttemptStreamingSequenceResponse{Content: word})
		if err != nil {
			return err
		}
	}

	echoStreamingHeaders(stream)

	echoStreamingTrailers(stream)

	time.Sleep(delay)

	// Calculate the perceived delay since the last RPC attempt.
	attDelay := &duration.Duration{}
	if n > 0 {
		prev := rep.GetAttempts()[n-1]
		respTime := prev.GetResponseTime()
		d := received.Sub(respTime.AsTime())
		attDelay = ptypes.DurationProto(d)
	}

	// Clock the time that the server is sending the response
	responseTime := time.Now()
	rpb, err := ptypes.TimestampProto(responseTime)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	rep.Attempts = append(rep.Attempts, &pb.StreamingSequenceReport_Attempt{
		AttemptNumber: int32(n),
		ResponseTime:  rpb,
		AttemptDelay:  attDelay,
		Status:        st.Proto(),
		ContentSent:   s.sentContent,
	})

	if n+1 >= len(responses) {
		s.sentContent = ""
	}

	return st.Err()

}

func (s *sequenceServerImpl) GetStreamingSequenceReport(ctx context.Context, in *pb.GetStreamingSequenceReportRequest) (*pb.StreamingSequenceReport, error) {
	name := in.GetName()
	if name == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"The field `name` is required.")
	}

	report, ok := s.streamingreports.Load(name)
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"The sequence report with %q does not exist.",
			name,
		)
	}

	return report.(*pb.StreamingSequenceReport), nil
}

func cloneStreamingSequence(s *pb.StreamingSequence) *pb.StreamingSequence {
	r := make([]*pb.StreamingSequence_Response, len(s.GetResponses()))
	copy(r, s.GetResponses())

	return &pb.StreamingSequence{
		Name:      s.GetName(),
		Content:   s.GetContent(),
		Responses: r,
	}
}

func streamingReport(n string) string {
	return fmt.Sprintf("%s/streamingSequenceReport", n)
}
