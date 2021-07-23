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
	"io"

	client "github.com/googleapis/gapic-showcase/client"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/api/iterator"
	locationpb "google.golang.org/genproto/googleapis/cloud/location"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
)

func ExampleNewEchoClient() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	// TODO: Use client.
	_ = c
}

func ExampleEchoClient_Echo() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &genprotopb.EchoRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.Echo(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleEchoClient_Chat() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()
	stream, err := c.Chat(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	go func() {
		reqs := []*genprotopb.EchoRequest{
			// TODO: Create requests.
		}
		for _, req := range reqs {
			if err := stream.Send(req); err != nil {
				// TODO: Handle error.
			}
		}
		stream.CloseSend()
	}()
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// TODO: handle error.
		}
		// TODO: Use resp.
		_ = resp
	}
}

func ExampleEchoClient_PagedExpand() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &genprotopb.PagedExpandRequest{
		// TODO: Fill request struct fields.
	}
	it := c.PagedExpand(ctx, req)
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

func ExampleEchoClient_PagedExpandLegacy() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &genprotopb.PagedExpandLegacyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.PagedExpandLegacy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleEchoClient_Wait() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &genprotopb.WaitRequest{
		// TODO: Fill request struct fields.
	}
	op, err := c.Wait(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleEchoClient_Block() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &genprotopb.BlockRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.Block(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleEchoClient_ListLocations() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_GetLocation() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_SetIamPolicy() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_GetIamPolicy() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_TestIamPermissions() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_ListOperations() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_GetOperation() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_DeleteOperation() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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

func ExampleEchoClient_CancelOperation() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
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
