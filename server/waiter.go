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
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	lropb "google.golang.org/genproto/googleapis/longrunning"
)

var waiterSingleton Waiter = &waiterImpl{
	nowF: time.Now,
}

// GetWaiterInstance returns the waiter singleton.
func GetWaiterInstance() Waiter {
	return waiterSingleton
}

// Waiter handles the echo.Wait method for both the LRO service and the echo service.
type Waiter interface {
	Wait(req *pb.WaitRequest) *lropb.Operation
}

type waiterImpl struct {
	nowF func() time.Time
}

func (w *waiterImpl) Wait(req *pb.WaitRequest) *lropb.Operation {
	endTime := time.Unix(0, 0).UTC()
	if ttl := req.GetTtl(); ttl != nil {
		duration, _ := ptypes.Duration(ttl)
		endTime = w.nowF().Add(duration)
	}
	if end := req.GetEndTime(); end != nil {
		endTime, _ = ptypes.Timestamp(end)
	}
	endTimeProto, _ := ptypes.TimestampProto(endTime)
	req.End = &pb.WaitRequest_EndTime{
		EndTime: endTimeProto,
	}

	done := w.nowF().After(endTime)
	reqBytes, _ := proto.Marshal(req)
	name := fmt.Sprintf(
		"operations/google.showcase.v1alpha3.Echo/Wait/%s",
		base64.StdEncoding.EncodeToString(reqBytes))
	answer := &lropb.Operation{
		Name: name,
		Done: done,
	}

	if done && (req.GetError() != nil) {
		answer.Result = &lropb.Operation_Error{Error: req.GetError()}
	}

	if done && (req.GetSuccess() != nil) {
		resp, _ := ptypes.MarshalAny(req.GetSuccess())
		answer.Result = &lropb.Operation_Response{Response: resp}
	}

	if !done {
		meta, _ := ptypes.MarshalAny(&pb.WaitMetadata{EndTime: endTimeProto})
		answer.Metadata = meta
	}

	return answer
}
