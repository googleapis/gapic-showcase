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
)

func ExampleNewSequenceClient() {
	ctx := context.Background()
	c, err := client.NewSequenceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleSequenceClient_CreateSequence() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewSequenceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.CreateSequenceRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateSequence(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleSequenceClient_GetSequenceReport() {
	// import genprotopb "github.com/googleapis/gapic-showcase/server/genproto"

	ctx := context.Background()
	c, err := client.NewSequenceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.GetSequenceReportRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetSequenceReport(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleSequenceClient_AttemptSequence() {
	ctx := context.Background()
	c, err := client.NewSequenceClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &genprotopb.AttemptSequenceRequest{
		// TODO: Fill request struct fields.
	}
	err = c.AttemptSequence(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}
