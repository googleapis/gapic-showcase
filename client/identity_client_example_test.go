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

	client "github.com/googleapis/gapic-showcase/client"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/api/iterator"
)

func ExampleNewIdentityClient() {
	ctx := context.Background()
	c, err := client.NewIdentityClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleIdentityClient_CreateUser() {
	ctx := context.Background()
	c, err := client.NewIdentityClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.CreateUserRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateUser(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleIdentityClient_GetUser() {
	ctx := context.Background()
	c, err := client.NewIdentityClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.GetUserRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetUser(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleIdentityClient_UpdateUser() {
	ctx := context.Background()
	c, err := client.NewIdentityClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.UpdateUserRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateUser(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleIdentityClient_DeleteUser() {
	ctx := context.Background()
	c, err := client.NewIdentityClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.DeleteUserRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteUser(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleIdentityClient_ListUsers() {
	ctx := context.Background()
	c, err := client.NewIdentityClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.ListUsersRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListUsers(ctx, req)
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
