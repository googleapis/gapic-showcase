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

// Code generated by protoc-gen-go_gapic. DO NOT EDIT.

package client_test

import (
	"context"

	client "github.com/googleapis/gapic-showcase/client"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/api/iterator"
	locationpb "google.golang.org/genproto/googleapis/cloud/location"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
)

func ExampleNewComplianceGrpcClient() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleComplianceGrpcClient_RepeatDataBody() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataBody(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_RepeatDataBodyInfo() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataBodyInfo(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_RepeatDataQuery() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataQuery(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_RepeatDataSimplePath() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataSimplePath(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_RepeatDataPathResource() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataPathResource(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_RepeatDataPathTrailingResource() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataPathTrailingResource(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_ListLocations() {
	// import locationpb "google.golang.org/genproto/googleapis/cloud/location"
	// import "google.golang.org/api/iterator"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &locationpb.ListLocationsRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListLocations(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		// TODO: Use resp.
		_ = resp
	}
}

func ExampleComplianceGrpcClient_GetLocation() {
	// import locationpb "google.golang.org/genproto/googleapis/cloud/location"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &locationpb.GetLocationRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetLocation(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_SetIamPolicy() {
	// import iampb "google.golang.org/genproto/googleapis/iam/v1"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &iampb.SetIamPolicyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.SetIamPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_GetIamPolicy() {
	// import iampb "google.golang.org/genproto/googleapis/iam/v1"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &iampb.GetIamPolicyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetIamPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_TestIamPermissions() {
	// import iampb "google.golang.org/genproto/googleapis/iam/v1"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &iampb.TestIamPermissionsRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.TestIamPermissions(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_ListOperations() {
	// import longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	// import "google.golang.org/api/iterator"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &longrunningpb.ListOperationsRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListOperations(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		// TODO: Use resp.
		_ = resp
	}
}

func ExampleComplianceGrpcClient_GetOperation() {
	// import longrunningpb "google.golang.org/genproto/googleapis/longrunning"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &longrunningpb.GetOperationRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetOperation(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceGrpcClient_DeleteOperation() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &longrunningpb.DeleteOperationRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteOperation(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleComplianceGrpcClient_CancelOperation() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &longrunningpb.CancelOperationRequest{
		// TODO: Fill request struct fields.
	}
	err = c.CancelOperation(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleComplianceGrpcClient_WaitOperation() {
	// import longrunningpb "google.golang.org/genproto/googleapis/longrunning"

	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &longrunningpb.WaitOperationRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.WaitOperation(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
