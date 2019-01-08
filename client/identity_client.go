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
	"math"
	"time"

	"github.com/golang/protobuf/proto"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// IdentityCallOptions contains the retry settings for each method of IdentityClient.
type IdentityCallOptions struct {
	CreateUser []gax.CallOption
	GetUser []gax.CallOption
	UpdateUser []gax.CallOption
	DeleteUser []gax.CallOption
	ListUsers []gax.CallOption
}

func defaultIdentityClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("localhost:7469:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultIdentityCallOptions() *IdentityCallOptions {
	backoff := gax.Backoff{
		Initial: 100 * time.Millisecond,
		Max: time.Minute,
		Multiplier: 1.3,
	}

	idempotent := []gax.CallOption{
		gax.WithRetry(func() gax.Retryer {
			return gax.OnCodes([]codes.Code{
				codes.Aborted,
				codes.Internal,
				codes.Unavailable,
				codes.Unknown,
			}, backoff)
		}),
	}

	nonidempotent := []gax.CallOption{
		gax.WithRetry(func() gax.Retryer {
			return gax.OnCodes([]codes.Code{
				codes.Unavailable,
			}, backoff)
		}),
	}

	return &IdentityCallOptions{
		GetUser: idempotent,
		ListUsers: idempotent,
		CreateUser: nonidempotent,
		UpdateUser: nonidempotent,
		DeleteUser: nonidempotent,
	}
}

// IdentityClient is a client for interacting with  API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type IdentityClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	identityClient genprotopb.IdentityClient

	// The call options for this service.
	CallOptions *IdentityCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewIdentityClient creates a new identity client.
//
// A simple identity service.
func NewIdentityClient(ctx context.Context, opts ...option.ClientOption) (*IdentityClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultIdentityClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &IdentityClient{
		conn:        conn,
		CallOptions: defaultIdentityCallOptions(),

		identityClient: genprotopb.NewIdentityClient(conn),
	}
	c.setGoogleClientInfo()

	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *IdentityClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *IdentityClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *IdentityClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// CreateUser creates a user.
func (c *IdentityClient) CreateUser(ctx context.Context, req *genprotopb.CreateUserRequest, opts ...gax.CallOption) (*genprotopb.User, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateUser[0:len(c.CallOptions.CreateUser):len(c.CallOptions.CreateUser)], opts...)
	var resp *genprotopb.User
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.identityClient.CreateUser(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUser retrieves the User with the given uri.
func (c *IdentityClient) GetUser(ctx context.Context, req *genprotopb.GetUserRequest, opts ...gax.CallOption) (*genprotopb.User, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetUser[0:len(c.CallOptions.GetUser):len(c.CallOptions.GetUser)], opts...)
	var resp *genprotopb.User
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.identityClient.GetUser(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateUser updates a user.
func (c *IdentityClient) UpdateUser(ctx context.Context, req *genprotopb.UpdateUserRequest, opts ...gax.CallOption) (*genprotopb.User, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateUser[0:len(c.CallOptions.UpdateUser):len(c.CallOptions.UpdateUser)], opts...)
	var resp *genprotopb.User
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.identityClient.UpdateUser(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteUser deletes a user, their profile, and all of their authored messages.
func (c *IdentityClient) DeleteUser(ctx context.Context, req *genprotopb.DeleteUserRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteUser[0:len(c.CallOptions.DeleteUser):len(c.CallOptions.DeleteUser)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.identityClient.DeleteUser(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// ListUsers lists all users.
func (c *IdentityClient) ListUsers(ctx context.Context, req *genprotopb.ListUsersRequest, opts ...gax.CallOption) *UserIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListUsers[0:len(c.CallOptions.ListUsers):len(c.CallOptions.ListUsers)], opts...)
	it := &UserIterator{}
	req = proto.Clone(req).(*genprotopb.ListUsersRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*genprotopb.User, string, error) {
		var resp *genprotopb.ListUsersResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.identityClient.ListUsers(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Users, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.PageSize)
	return it
}

// UserIterator manages a stream of *genprotopb.User.
type UserIterator struct {
	items    []*genprotopb.User
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*genprotopb.User, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *UserIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *UserIterator) Next() (*genprotopb.User, error) {
	var item *genprotopb.User
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *UserIterator) bufLen() int {
	return len(it.items)
}

func (it *UserIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
