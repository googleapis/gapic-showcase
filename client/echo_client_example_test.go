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

// AUTO-GENERATED CODE. DO NOT EDIT.

package client_test

import (
	"context"
	"io"

	client "github.com/googleapis/gapic-showcase/client"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/api/iterator"
)

func ExampleNewEchoClient() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleEchoClient_Echo() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

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

func ExampleEchoClient_Wait() {
	ctx := context.Background()
	c, err := client.NewEchoClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

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
