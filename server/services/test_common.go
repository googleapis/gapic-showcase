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

package services

import (
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

// Mock waiter type used in echo_service_test and operations_service_test to
// check that they defer to the waiter.
type mockWaiter struct {
	req *pb.WaitRequest
}

func (w *mockWaiter) Wait(req *pb.WaitRequest) *lropb.Operation {
	w.req = req
	return nil
}
