// Copyright 2021 Google LLC
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

	"github.com/google/go-cmp/cmp"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/proto"
)

func TestComplianceRepeats(t *testing.T) {
	// Note that additional thorough test cases are exercised in
	// cmd/gapic-showcase/compliance_suite_test.go.
	server := NewComplianceServer()
	info := &pb.ComplianceData{
		FString:   "Terra Incognita",
		FInt32:    1,
		FSint32:   -2,
		FSfixed32: -300000000,
		FUint32:   5,
		FFixed32:  700000000,
		FInt64:    9,
		FSint64:   -1100000000,
		FSfixed64: -1300000000,
		FUint64:   1700000000000000000,
		FFixed64:  1900000000000000000,

		FDouble: 6.02e23,
		FFloat:  3.1415,
		FBool:   true,
		FBytes:  []byte("Lorem ipsum"),
	}
	request := &pb.RepeatRequest{Info: info}

	for idx, rpc := range [](func(ctx context.Context, in *pb.RepeatRequest) (*pb.RepeatResponse, error)){
		server.RepeatDataBody,
		server.RepeatDataQuery,
		server.RepeatDataSimplePath,
	} {
		response, err := rpc(context.Background(), request)
		if err != nil {
			t.Errorf("call %d: error: %s", idx, err)
		}
		if diff := cmp.Diff(response.GetInfo(), request.GetInfo(), cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("call %d: got=-, want=+:%s", idx, diff)
		}
	}
}
