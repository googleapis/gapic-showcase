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
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (s *sequenceServerImpl) CreateStreamingSequence(ctx context.Context, in *pb.CreateSequenceRequest) (*pb.Sequence, error) {
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

func (s *sequenceServerImpl) AttemptStreamingSequence(in *pb.AttemptStreamingSequenceRequest, stream pb.SequenceService_AttemptStreamingSequenceServer) error {
	received := time.Now()
	name := in.GetName()
	if name == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `name` is required.")
	}

	// Retrieve Sequence and associated SequenceReport.
	i, ok := s.sequences.Load(name)
	if !ok {
		return status.Errorf(
			codes.NotFound,
			"The Sequence with %q does not exist.",
			name,
		)
	}
	seq := i.(*pb.Sequence)

	i, _ = s.reports.Load(report(name))
	rep, _ := i.(*pb.SequenceReport)

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

	r, _ := strconv.Atoi(in.GetContent())

	for end := time.Now().Add(delay*time.Nanosecond); ; {
		if time.Now().After(end) {
			break
		}

		for number := 0; number < r; number++ {
			err := stream.Send(&pb.AttemptStreamingSequenceResponse{Content: strconv.Itoa(number)})
			time.Sleep(delay)
			if err != nil {
				return err
			}
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

	rep.Attempts = append(rep.Attempts, &pb.SequenceReport_Attempt{
		AttemptNumber: int32(n),
		ResponseTime:  rpb,
		AttemptDelay:  attDelay,
		Status:        st.Proto(),
	})

	return st.Err()
}

func (s *sequenceServerImpl) GetStreamingSequenceReport(ctx context.Context, in *pb.GetSequenceReportRequest) (*pb.SequenceReport, error) {
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

