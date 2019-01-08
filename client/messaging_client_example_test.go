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

func ExampleNewMessagingClient() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleMessagingClient_CreateRoom() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.CreateRoomRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateRoom(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleMessagingClient_GetRoom() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.GetRoomRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetRoom(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleMessagingClient_UpdateRoom() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.UpdateRoomRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateRoom(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleMessagingClient_DeleteRoom() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.DeleteRoomRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteRoom(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleMessagingClient_ListRooms() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.ListRoomsRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListRooms(ctx, req)
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

func ExampleMessagingClient_CreateBlurb() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.CreateBlurbRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateBlurb(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleMessagingClient_GetBlurb() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.GetBlurbRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetBlurb(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleMessagingClient_UpdateBlurb() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.UpdateBlurbRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateBlurb(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleMessagingClient_DeleteBlurb() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.DeleteBlurbRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteBlurb(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleMessagingClient_ListBlurbs() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.ListBlurbsRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListBlurbs(ctx, req)
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

func ExampleMessagingClient_SearchBlurbs() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.SearchBlurbsRequest{
		// TODO: Fill request struct fields.
	}
	op, err := c.SearchBlurbs(ctx, req)
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

func ExampleMessagingClient_Connect() {
	ctx := context.Background()
	c, err := client.NewMessagingClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	stream, err := c.Connect(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	go func() {
		reqs := []*genprotopb.ConnectRequest{
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
