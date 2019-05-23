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

// Code generated by protoc-gen-go_gapic. DO NOT EDIT.

package client

import (
	"context"
	"fmt"
	"math"
	"time"

	"cloud.google.com/go/longrunning"
	lroauto "cloud.google.com/go/longrunning/autogen"
	"github.com/golang/protobuf/proto"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// MessagingCallOptions contains the retry settings for each method of MessagingClient.
type MessagingCallOptions struct {
	CreateRoom   []gax.CallOption
	GetRoom      []gax.CallOption
	UpdateRoom   []gax.CallOption
	DeleteRoom   []gax.CallOption
	ListRooms    []gax.CallOption
	CreateBlurb  []gax.CallOption
	GetBlurb     []gax.CallOption
	UpdateBlurb  []gax.CallOption
	DeleteBlurb  []gax.CallOption
	ListBlurbs   []gax.CallOption
	SearchBlurbs []gax.CallOption
	StreamBlurbs []gax.CallOption
	SendBlurbs   []gax.CallOption
	Connect      []gax.CallOption
}

func defaultMessagingClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("localhost:7469"),
		option.WithGRPCDialOption(grpc.WithDisableServiceConfig()),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultMessagingCallOptions() *MessagingCallOptions {
	backoff := gax.Backoff{
		Initial:    100 * time.Millisecond,
		Max:        time.Minute,
		Multiplier: 1.3,
	}

	idempotent := []gax.CallOption{
		gax.WithRetry(func() gax.Retryer {
			return gax.OnCodes([]codes.Code{
				codes.Aborted,
				codes.Unavailable,
				codes.Unknown,
			}, backoff)
		}),
	}

	return &MessagingCallOptions{
		GetRoom:    idempotent,
		ListRooms:  idempotent,
		GetBlurb:   idempotent,
		ListBlurbs: idempotent,
	}
}

// MessagingClient is a client for interacting with .
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type MessagingClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	messagingClient genprotopb.MessagingClient

	// LROClient is used internally to handle longrunning operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient *lroauto.OperationsClient

	// The call options for this service.
	CallOptions *MessagingCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewMessagingClient creates a new messaging client.
//
// A simple messaging service that implements chat rooms and profile posts.
//
// This messaging service showcases the features that API clients
// generated by gapic-generators implement.
func NewMessagingClient(ctx context.Context, opts ...option.ClientOption) (*MessagingClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultMessagingClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &MessagingClient{
		conn:        conn,
		CallOptions: defaultMessagingCallOptions(),

		messagingClient: genprotopb.NewMessagingClient(conn),
	}
	c.setGoogleClientInfo()

	c.LROClient, err = lroauto.NewOperationsClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		// This error "should not happen", since we are just reusing old connection
		// and never actually need to dial.
		// If this does happen, we could leak conn. However, we cannot close conn:
		// If the user invoked the function with option.WithGRPCConn,
		// we would close a connection that's still in use.
		// TODO(pongad): investigate error conditions.
		return nil, err
	}
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *MessagingClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *MessagingClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *MessagingClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// CreateRoom creates a room.
func (c *MessagingClient) CreateRoom(ctx context.Context, req *genprotopb.CreateRoomRequest, opts ...gax.CallOption) (*genprotopb.Room, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateRoom[0:len(c.CallOptions.CreateRoom):len(c.CallOptions.CreateRoom)], opts...)
	var resp *genprotopb.Room
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.CreateRoom(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetRoom retrieves the Room with the given resource name.
func (c *MessagingClient) GetRoom(ctx context.Context, req *genprotopb.GetRoomRequest, opts ...gax.CallOption) (*genprotopb.Room, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("name=%v", req.GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetRoom[0:len(c.CallOptions.GetRoom):len(c.CallOptions.GetRoom)], opts...)
	var resp *genprotopb.Room
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.GetRoom(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateRoom updates a room.
func (c *MessagingClient) UpdateRoom(ctx context.Context, req *genprotopb.UpdateRoomRequest, opts ...gax.CallOption) (*genprotopb.Room, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateRoom[0:len(c.CallOptions.UpdateRoom):len(c.CallOptions.UpdateRoom)], opts...)
	var resp *genprotopb.Room
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.UpdateRoom(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteRoom deletes a room and all of its blurbs.
func (c *MessagingClient) DeleteRoom(ctx context.Context, req *genprotopb.DeleteRoomRequest, opts ...gax.CallOption) error {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("name=%v", req.GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.DeleteRoom[0:len(c.CallOptions.DeleteRoom):len(c.CallOptions.DeleteRoom)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.messagingClient.DeleteRoom(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// ListRooms lists all chat rooms.
func (c *MessagingClient) ListRooms(ctx context.Context, req *genprotopb.ListRoomsRequest, opts ...gax.CallOption) *RoomIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListRooms[0:len(c.CallOptions.ListRooms):len(c.CallOptions.ListRooms)], opts...)
	it := &RoomIterator{}
	req = proto.Clone(req).(*genprotopb.ListRoomsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*genprotopb.Room, string, error) {
		var resp *genprotopb.ListRoomsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.messagingClient.ListRooms(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.Rooms, resp.NextPageToken, nil
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
	it.pageInfo.Token = req.PageToken
	return it
}

// CreateBlurb creates a blurb. If the parent is a room, the blurb is understood to be a
// message in that room. If the parent is a profile, the blurb is understood
// to be a post on the profile.
func (c *MessagingClient) CreateBlurb(ctx context.Context, req *genprotopb.CreateBlurbRequest, opts ...gax.CallOption) (*genprotopb.Blurb, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("parent=%v", req.GetParent()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.CreateBlurb[0:len(c.CallOptions.CreateBlurb):len(c.CallOptions.CreateBlurb)], opts...)
	var resp *genprotopb.Blurb
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.CreateBlurb(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBlurb retrieves the Blurb with the given resource name.
func (c *MessagingClient) GetBlurb(ctx context.Context, req *genprotopb.GetBlurbRequest, opts ...gax.CallOption) (*genprotopb.Blurb, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("name=%v", req.GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetBlurb[0:len(c.CallOptions.GetBlurb):len(c.CallOptions.GetBlurb)], opts...)
	var resp *genprotopb.Blurb
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.GetBlurb(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateBlurb updates a blurb.
func (c *MessagingClient) UpdateBlurb(ctx context.Context, req *genprotopb.UpdateBlurbRequest, opts ...gax.CallOption) (*genprotopb.Blurb, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateBlurb[0:len(c.CallOptions.UpdateBlurb):len(c.CallOptions.UpdateBlurb)], opts...)
	var resp *genprotopb.Blurb
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.UpdateBlurb(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteBlurb deletes a blurb.
func (c *MessagingClient) DeleteBlurb(ctx context.Context, req *genprotopb.DeleteBlurbRequest, opts ...gax.CallOption) error {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("name=%v", req.GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.DeleteBlurb[0:len(c.CallOptions.DeleteBlurb):len(c.CallOptions.DeleteBlurb)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.messagingClient.DeleteBlurb(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// ListBlurbs lists blurbs for a specific chat room or user profile depending on the
// parent resource name.
func (c *MessagingClient) ListBlurbs(ctx context.Context, req *genprotopb.ListBlurbsRequest, opts ...gax.CallOption) *BlurbIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("parent=%v", req.GetParent()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ListBlurbs[0:len(c.CallOptions.ListBlurbs):len(c.CallOptions.ListBlurbs)], opts...)
	it := &BlurbIterator{}
	req = proto.Clone(req).(*genprotopb.ListBlurbsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*genprotopb.Blurb, string, error) {
		var resp *genprotopb.ListBlurbsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.messagingClient.ListBlurbs(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.Blurbs, resp.NextPageToken, nil
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
	it.pageInfo.Token = req.PageToken
	return it
}

// SearchBlurbs this method searches through all blurbs across all rooms and profiles
// for blurbs containing to words found in the query. Only posts that
// contain an exact match of a queried word will be returned.
func (c *MessagingClient) SearchBlurbs(ctx context.Context, req *genprotopb.SearchBlurbsRequest, opts ...gax.CallOption) (*SearchBlurbsOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("parent=%v", req.GetParent()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.SearchBlurbs[0:len(c.CallOptions.SearchBlurbs):len(c.CallOptions.SearchBlurbs)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.SearchBlurbs(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &SearchBlurbsOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// StreamBlurbs this returns a stream that emits the blurbs that are created for a
// particular chat room or user profile.
func (c *MessagingClient) StreamBlurbs(ctx context.Context, req *genprotopb.StreamBlurbsRequest, opts ...gax.CallOption) (genprotopb.Messaging_StreamBlurbsClient, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("name=%v", req.GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.StreamBlurbs[0:len(c.CallOptions.StreamBlurbs):len(c.CallOptions.StreamBlurbs)], opts...)
	var resp genprotopb.Messaging_StreamBlurbsClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.StreamBlurbs(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SendBlurbs this is a stream to create multiple blurbs. If an invalid blurb is
// requested to be created, the stream will close with an error.
func (c *MessagingClient) SendBlurbs(ctx context.Context, opts ...gax.CallOption) (genprotopb.Messaging_SendBlurbsClient, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.SendBlurbs[0:len(c.CallOptions.SendBlurbs):len(c.CallOptions.SendBlurbs)], opts...)
	var resp genprotopb.Messaging_SendBlurbsClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.SendBlurbs(ctx, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Connect this method starts a bidirectional stream that receives all blurbs that
// are being created after the stream has started and sends requests to create
// blurbs. If an invalid blurb is requested to be created, the stream will
// close with an error.
func (c *MessagingClient) Connect(ctx context.Context, opts ...gax.CallOption) (genprotopb.Messaging_ConnectClient, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.Connect[0:len(c.CallOptions.Connect):len(c.CallOptions.Connect)], opts...)
	var resp genprotopb.Messaging_ConnectClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.messagingClient.Connect(ctx, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SearchBlurbsOperation manages a long-running operation from SearchBlurbs.
type SearchBlurbsOperation struct {
	lro *longrunning.Operation
}

// SearchBlurbsOperation returns a new SearchBlurbsOperation from a given name.
// The name must be that of a previously created SearchBlurbsOperation, possibly from a different process.
func (c *MessagingClient) SearchBlurbsOperation(name string) *SearchBlurbsOperation {
	return &SearchBlurbsOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *SearchBlurbsOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*genprotopb.SearchBlurbsResponse, error) {
	var resp genprotopb.SearchBlurbsResponse
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *SearchBlurbsOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*genprotopb.SearchBlurbsResponse, error) {
	var resp genprotopb.SearchBlurbsResponse
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *SearchBlurbsOperation) Metadata() (*genprotopb.SearchBlurbsMetadata, error) {
	var meta genprotopb.SearchBlurbsMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *SearchBlurbsOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *SearchBlurbsOperation) Name() string {
	return op.lro.Name()
}

// BlurbIterator manages a stream of *genprotopb.Blurb.
type BlurbIterator struct {
	items    []*genprotopb.Blurb
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// Response is the raw response for the current page.
	// Calling Next() or InternalFetch() updates this value.
	Response *genprotopb.ListBlurbsResponse

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*genprotopb.Blurb, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *BlurbIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *BlurbIterator) Next() (*genprotopb.Blurb, error) {
	var item *genprotopb.Blurb
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *BlurbIterator) bufLen() int {
	return len(it.items)
}

func (it *BlurbIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

// RoomIterator manages a stream of *genprotopb.Room.
type RoomIterator struct {
	items    []*genprotopb.Room
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// Response is the raw response for the current page.
	// Calling Next() or InternalFetch() updates this value.
	Response *genprotopb.ListRoomsResponse

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*genprotopb.Room, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *RoomIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *RoomIterator) Next() (*genprotopb.Room, error) {
	var item *genprotopb.Room
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *RoomIterator) bufLen() int {
	return len(it.items)
}

func (it *RoomIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
