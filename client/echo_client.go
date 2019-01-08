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
	"google.golang.org/grpc/metadata"
)

// EchoCallOptions contains the retry settings for each method of EchoClient.
type EchoCallOptions struct {
	Echo        []gax.CallOption
	Expand      []gax.CallOption
	Collect     []gax.CallOption
	Chat        []gax.CallOption
	PagedExpand []gax.CallOption
	Wait        []gax.CallOption
}

func defaultEchoClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("localhost:7469:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultEchoCallOptions() *EchoCallOptions {
	return &EchoCallOptions{}
}

// EchoClient is a client for interacting with  API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type EchoClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	echoClient genprotopb.EchoClient

	// LROClient is used internally to handle longrunning operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient *lroauto.OperationsClient

	// The call options for this service.
	CallOptions *EchoCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewEchoClient creates a new echo client.
//
// This service is used showcase the four main types of rpcs - unary, server
// side streaming, client side streaming, and bidirectional streaming. This
// service also exposes methods that explicitly implement server delay, and
// paginated calls.
func NewEchoClient(ctx context.Context, opts ...option.ClientOption) (*EchoClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultEchoClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &EchoClient{
		conn:        conn,
		CallOptions: defaultEchoCallOptions(),

		echoClient: genprotopb.NewEchoClient(conn),
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
func (c *EchoClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *EchoClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *EchoClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Echo this method simply echos the request. This method is showcases unary rpcs.
func (c *EchoClient) Echo(ctx context.Context, req *genprotopb.EchoRequest, opts ...gax.CallOption) (*genprotopb.EchoResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.Echo[0:len(c.CallOptions.Echo):len(c.CallOptions.Echo)], opts...)
	var resp *genprotopb.EchoResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.echoClient.Echo(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Expand this method split the given content into words and will pass each word back
// through the stream. This method showcases server-side streaming rpcs.
func (c *EchoClient) Expand(ctx context.Context, req *genprotopb.ExpandRequest, opts ...gax.CallOption) (genprotopb.Echo_ExpandClient, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.Expand[0:len(c.CallOptions.Expand):len(c.CallOptions.Expand)], opts...)
	var resp genprotopb.Echo_ExpandClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.echoClient.Expand(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Collect this method will collect the words given to it. When the stream is closed
// by the client, this method will return the a concatenation of the strings
// passed to it. This method showcases client-side streaming rpcs.
func (c *EchoClient) Collect(ctx context.Context, opts ...gax.CallOption) (genprotopb.Echo_CollectClient, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.Collect[0:len(c.CallOptions.Collect):len(c.CallOptions.Collect)], opts...)
	var resp genprotopb.Echo_CollectClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.echoClient.Collect(ctx, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Chat this method, upon receiving a request on the stream, the same content will
// be passed  back on the stream. This method showcases bidirectional
// streaming rpcs.
func (c *EchoClient) Chat(ctx context.Context, opts ...gax.CallOption) (genprotopb.Echo_ChatClient, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.Chat[0:len(c.CallOptions.Chat):len(c.CallOptions.Chat)], opts...)
	var resp genprotopb.Echo_ChatClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.echoClient.Chat(ctx, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PagedExpand this is similar to the Expand method but instead of returning a stream of
// expanded words, this method returns a paged list of expanded words.
func (c *EchoClient) PagedExpand(ctx context.Context, req *genprotopb.PagedExpandRequest, opts ...gax.CallOption) *EchoResponseIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.PagedExpand[0:len(c.CallOptions.PagedExpand):len(c.CallOptions.PagedExpand)], opts...)
	it := &EchoResponseIterator{}
	req = proto.Clone(req).(*genprotopb.PagedExpandRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*genprotopb.EchoResponse, string, error) {
		var resp *genprotopb.PagedExpandResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.echoClient.PagedExpand(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Responses, resp.NextPageToken, nil
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

// Wait this method will wait the requested amount of and then return.
// This method showcases how a client handles a request timing out.
func (c *EchoClient) Wait(ctx context.Context, req *genprotopb.WaitRequest, opts ...gax.CallOption) (*WaitOperation, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.Wait[0:len(c.CallOptions.Wait):len(c.CallOptions.Wait)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.echoClient.Wait(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &WaitOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// WaitOperation manages a long-running operation from Wait.
type WaitOperation struct {
	lro *longrunning.Operation
}

// WaitOperation returns a new WaitOperation from a given name.
// The name must be that of a previously created WaitOperation, possibly from a different process.
func (c *EchoClient) WaitOperation(name string) *WaitOperation {
	return &WaitOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *WaitOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*genprotopb.WaitResponse, error) {
	var resp genprotopb.WaitResponse
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
func (op *WaitOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*genprotopb.WaitResponse, error) {
	var resp genprotopb.WaitResponse
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
func (op *WaitOperation) Metadata() (*genprotopb.WaitMetadata, error) {
	var meta genprotopb.WaitMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *WaitOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *WaitOperation) Name() string {
	return op.lro.Name()
}

// EchoResponseIterator manages a stream of *genprotopb.EchoResponse.
type EchoResponseIterator struct {
	items    []*genprotopb.EchoResponse
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*genprotopb.EchoResponse, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *EchoResponseIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *EchoResponseIterator) Next() (*genprotopb.EchoResponse, error) {
	var item *genprotopb.EchoResponse
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *EchoResponseIterator) bufLen() int {
	return len(it.items)
}

func (it *EchoResponseIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
