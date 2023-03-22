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
	"strings"
	"testing"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockStreamSequence struct {
	exp   []string
	head  []string
	trail []string
	t     *testing.T
	pb.SequenceService_AttemptStreamingSequenceServer
}

func TestAttemptStreamingSequence(t *testing.T) {
	s := NewSequenceServer()
	stream := &mockStreamSequence{exp: strings.Fields("10"), t: t}
	attemptRequest := &pb.AttemptStreamingSequenceRequest{Name: "sequences/0"}
	err := s.AttemptStreamingSequence(attemptRequest,stream)
	if c := status.Code(err); c != codes.NotFound {
		t.Errorf("%s: expected error to be %s but was %s", t.Name(), codes.NotFound, c)
	}
}
