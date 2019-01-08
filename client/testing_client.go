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

package client

import (
	"context"
	"time"

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// TestingCallOptions contains the retry settings for each method of TestingClient.
type TestingCallOptions struct {
	ReportSession []gax.CallOption
	DeleteTest []gax.CallOption
	RegisterTest []gax.CallOption
}

func defaultTestingClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("localhost:7469:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultTestingCallOptions() *TestingCallOptions {
	backoff := gax.Backoff{
		Initial: 100 * time.Millisecond,
		Max: time.Minute,
		Multiplier: 1.3,
	}

	nonidempotent := []gax.CallOption{
		gax.WithRetry(func() gax.Retryer {
			return gax.OnCodes([]codes.Code{
				codes.Unavailable,
			}, backoff)
		}),
	}

	return &TestingCallOptions{
		ReportSession: nonidempotent,
		DeleteTest: nonidempotent,
		RegisterTest: nonidempotent,
	}
}

// TestingClient is a client for interacting with  API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type TestingClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	testingClient genprotopb.TestingClient

	// The call options for this service.
	CallOptions *TestingCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewTestingClient creates a new testing client.
//
// A service to facilitate running discrete sets of tests
// against Showcase.
func NewTestingClient(ctx context.Context, opts ...option.ClientOption) (*TestingClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultTestingClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &TestingClient{
		conn:        conn,
		CallOptions: defaultTestingCallOptions(),

		testingClient: genprotopb.NewTestingClient(conn),
	}
	c.setGoogleClientInfo()

	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *TestingClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *TestingClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *TestingClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// ReportSession report on the status of a session.
// This generates a report detailing which tests have been completed,
// and an overall rollup.
func (c *TestingClient) ReportSession(ctx context.Context, req *genprotopb.ReportSessionRequest, opts ...gax.CallOption) (*genprotopb.ReportSessionResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ReportSession[0:len(c.CallOptions.ReportSession):len(c.CallOptions.ReportSession)], opts...)
	var resp *genprotopb.ReportSessionResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.testingClient.ReportSession(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteTest explicitly decline to implement a test.
//
// This removes the test from subsequent ListTests calls, and
// attempting to do the test will error.
//
// This method will error if attempting to delete a required test.
func (c *TestingClient) DeleteTest(ctx context.Context, req *genprotopb.DeleteTestRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteTest[0:len(c.CallOptions.DeleteTest):len(c.CallOptions.DeleteTest)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.testingClient.DeleteTest(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// RegisterTest register a response to a test.
//
// In cases where a test involves registering a final answer at the
// end of the test, this method provides the means to do so.
func (c *TestingClient) RegisterTest(ctx context.Context, req *genprotopb.RegisterTestRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.RegisterTest[0:len(c.CallOptions.RegisterTest):len(c.CallOptions.RegisterTest)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.testingClient.RegisterTest(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}
