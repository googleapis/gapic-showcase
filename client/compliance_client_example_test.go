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

func ExampleNewComplianceClient() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	// TODO: Use client.
	_ = c
}

func ExampleComplianceClient_RepeatDataBody() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_RepeatDataBodyInfo() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_RepeatDataQuery() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_RepeatDataSimplePath() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_RepeatDataPathResource() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_RepeatDataPathTrailingResource() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_RepeatDataBodyPut() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataBodyPut(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceClient_RepeatDataBodyPatch() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &genprotopb.RepeatRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RepeatDataBodyPatch(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleComplianceClient_ListLocations() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_GetLocation() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_SetIamPolicy() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_GetIamPolicy() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_TestIamPermissions() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_ListOperations() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_GetOperation() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

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

func ExampleComplianceClient_DeleteOperation() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &longrunningpb.DeleteOperationRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteOperation(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleComplianceClient_CancelOperation() {
	ctx := context.Background()
	c, err := client.NewComplianceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &longrunningpb.CancelOperationRequest{
		// TODO: Fill request struct fields.
	}
	err = c.CancelOperation(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}
