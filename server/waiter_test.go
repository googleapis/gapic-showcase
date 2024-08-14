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
	"encoding/base64"
	"strings"
	"testing"
	"time"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/googleapis/rpc/status"
)

func TestGetWaiterInstance(t *testing.T) {
	waiter := GetWaiterInstance()
	if waiter != waiterSingleton {
		t.Error("GetWaiterInstance: Expected to get waiter singleton.")
	}
}

func TestWait_pending(t *testing.T) {
	now := time.Unix(1, 0)
	endTime := time.Unix(2, 0)
	ttl := endTime.Sub(now)
	nowF := func() time.Time { return time.Unix(1, 0) }
	endTimeProto := timestampProto(endTime)

	tests := []*pb.WaitRequest{
		&pb.WaitRequest{
			End: &pb.WaitRequest_EndTime{
				EndTime: endTimeProto,
			},
		},
		&pb.WaitRequest{
			End: &pb.WaitRequest_Ttl{
				Ttl: ptypes.DurationProto(ttl),
			},
		},
	}

	for _, req := range tests {
		waiter := &waiterImpl{nowF: nowF}
		op := waiter.Wait(req)

		if op.Done {
			t.Errorf("Wait() for %q expectee done=false got done=true", req)
		}

		checkName(t, req, op)

		if op.Metadata == nil {
			t.Errorf("Wait() for %q expected metadata, got nil", req)
		}

		meta := &pb.WaitMetadata{}
		ptypes.UnmarshalAny(op.Metadata, meta)
		if !proto.Equal(endTimeProto, meta.EndTime) {
			t.Errorf(
				"Wait for %q expected metadata with Endtime=%q, got %q",
				req,
				endTimeProto,
				meta.EndTime)
		}
	}
}

func TestWait_success(t *testing.T) {
	nowF := func() time.Time { return time.Unix(3, 0) }
	endTime := timestampProto(time.Unix(2, 0))
	success := &pb.WaitResponse{Content: "Hello World!"}
	req := &pb.WaitRequest{
		End: &pb.WaitRequest_EndTime{
			EndTime: endTime,
		},
		Response: &pb.WaitRequest_Success{Success: success},
	}

	waiter := &waiterImpl{nowF: nowF}
	op := waiter.Wait(req)

	checkName(t, req, op)

	if !op.Done {
		t.Errorf("Wait() for %q expected done=true got done=false", req)
	}
	if op.Metadata != nil {
		t.Errorf("Wait() for %q expected nil metadata, got %q", req, op.Metadata)
	}
	if op.GetError() != nil {
		t.Errorf("Wait() expected op.Error=nil, got %q", op.GetError())
	}
	if op.GetResponse() == nil {
		t.Error("Wait() expected op.Response!=nil")
	}
	resp := &pb.WaitResponse{}
	ptypes.UnmarshalAny(op.GetResponse(), resp)
	if !proto.Equal(resp, success) {
		t.Errorf("Wait() expected op.GetResponse()=%q, got %q", success, resp)
	}
}

func TestWait_error(t *testing.T) {
	nowF := func() time.Time { return time.Unix(3, 0) }
	endTime := timestampProto(time.Unix(2, 0))
	expErr := &status.Status{Code: int32(1), Message: "Error!"}
	req := &pb.WaitRequest{
		End: &pb.WaitRequest_EndTime{
			EndTime: endTime,
		},
		Response: &pb.WaitRequest_Error{Error: expErr},
	}

	waiter := &waiterImpl{nowF: nowF}
	op := waiter.Wait(req)

	checkName(t, req, op)

	if !op.Done {
		t.Errorf("Wait() for %q expected done=true got done=false", req)
	}
	if op.Metadata != nil {
		t.Errorf("Wait() for %q expected nil metadata, got %q", req, op.Metadata)
	}
	if op.GetResponse() != nil {
		t.Errorf("Wait() expected op.Response=nil, got %q", op.GetResponse())
	}
	if !proto.Equal(expErr, op.GetError()) {
		t.Errorf("Wait() expected op.Error=%q, got %q", expErr, op.GetError())
	}
}

func timestampProto(t time.Time) *timestamp.Timestamp {
	ts, _ := ptypes.TimestampProto(t)
	return ts
}

func checkName(t *testing.T, req *pb.WaitRequest, op *lropb.Operation) {
	if !strings.HasPrefix(op.Name, "operations/google.showcase.v1beta1.Echo/Wait/") {
		t.Errorf(
			"Wait() expected op.Name prefex 'operations/google.showcase.v1beta1.Echo/Wait/', got: %s'",
			op.Name)
	}
	nameProto := &pb.WaitRequest{}
	encodedBytes := strings.TrimPrefix(
		op.Name,
		"operations/google.showcase.v1beta1.Echo/Wait/")
	bytes, _ := base64.StdEncoding.DecodeString(encodedBytes)
	proto.Unmarshal(bytes, nameProto)
	if !proto.Equal(nameProto, req) {
		t.Errorf(
			"Wait() for %q expected unmarshalled name=%q, got name=%q",
			req,
			req,
			nameProto)
	}
}
