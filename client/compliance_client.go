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

package client

import (
	"context"
	"fmt"
	"math"
	"net/url"

	"github.com/golang/protobuf/proto"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	locationpb "google.golang.org/genproto/googleapis/cloud/location"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var newComplianceGrpcClientHook clientHook

// ComplianceCallOptions contains the retry settings for each method of ComplianceClient.
type ComplianceCallOptions struct {
	RepeatDataBody                 []gax.CallOption
	RepeatDataBodyInfo             []gax.CallOption
	RepeatDataQuery                []gax.CallOption
	RepeatDataSimplePath           []gax.CallOption
	RepeatDataPathResource         []gax.CallOption
	RepeatDataPathTrailingResource []gax.CallOption
	ListLocations                  []gax.CallOption
	GetLocation                    []gax.CallOption
	SetIamPolicy                   []gax.CallOption
	GetIamPolicy                   []gax.CallOption
	TestIamPermissions             []gax.CallOption
	ListOperations                 []gax.CallOption
	GetOperation                   []gax.CallOption
	DeleteOperation                []gax.CallOption
	CancelOperation                []gax.CallOption
	WaitOperation                  []gax.CallOption
}

func defaultComplianceGrpcClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("localhost:7469"),
		internaloption.WithDefaultMTLSEndpoint("localhost:7469"),
		internaloption.WithDefaultAudience("https://localhost/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
		option.WithGRPCDialOption(grpc.WithDisableServiceConfig()),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultComplianceCallOptions() *ComplianceCallOptions {
	return &ComplianceCallOptions{
		RepeatDataBody:                 []gax.CallOption{},
		RepeatDataBodyInfo:             []gax.CallOption{},
		RepeatDataQuery:                []gax.CallOption{},
		RepeatDataSimplePath:           []gax.CallOption{},
		RepeatDataPathResource:         []gax.CallOption{},
		RepeatDataPathTrailingResource: []gax.CallOption{},
		ListLocations:                  []gax.CallOption{},
		GetLocation:                    []gax.CallOption{},
		SetIamPolicy:                   []gax.CallOption{},
		GetIamPolicy:                   []gax.CallOption{},
		TestIamPermissions:             []gax.CallOption{},
		ListOperations:                 []gax.CallOption{},
		GetOperation:                   []gax.CallOption{},
		DeleteOperation:                []gax.CallOption{},
		CancelOperation:                []gax.CallOption{},
		WaitOperation:                  []gax.CallOption{},
	}
}

// internalComplianceClient is an interface that defines the methods availaible from Client Libraries Showcase API.
type internalComplianceClient interface {
	RepeatDataBody(context.Context, *genprotopb.RepeatRequest, ...gax.CallOption) (*genprotopb.RepeatResponse, error)
	RepeatDataBodyInfo(context.Context, *genprotopb.RepeatRequest, ...gax.CallOption) (*genprotopb.RepeatResponse, error)
	RepeatDataQuery(context.Context, *genprotopb.RepeatRequest, ...gax.CallOption) (*genprotopb.RepeatResponse, error)
	RepeatDataSimplePath(context.Context, *genprotopb.RepeatRequest, ...gax.CallOption) (*genprotopb.RepeatResponse, error)
	RepeatDataPathResource(context.Context, *genprotopb.RepeatRequest, ...gax.CallOption) (*genprotopb.RepeatResponse, error)
	RepeatDataPathTrailingResource(context.Context, *genprotopb.RepeatRequest, ...gax.CallOption) (*genprotopb.RepeatResponse, error)
}

// ComplianceClient is a client for interacting with Client Libraries Showcase API.
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type ComplianceClient struct {
	// The internal transport-dependent client.
	internalClient internalComplianceClient

	// The call options for this service.
	CallOptions *ComplianceCallOptions
}

func (c *ComplianceClient) Close() error {
	return c.internalClient.Close()
}

func (c *ComplianceClient) RepeatDataBody(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	opts = append(c.CallOptions.RepeatDataBody[0:len(c.CallOptions.RepeatDataBody):len(c.CallOptions.RepeatDataBody)], opts...)
	return c.internalClient.RepeatDataBody(ctx, req, opts)
}
func (c *ComplianceClient) RepeatDataBodyInfo(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	opts = append(c.CallOptions.RepeatDataBodyInfo[0:len(c.CallOptions.RepeatDataBodyInfo):len(c.CallOptions.RepeatDataBodyInfo)], opts...)
	return c.internalClient.RepeatDataBodyInfo(ctx, req, opts)
}
func (c *ComplianceClient) RepeatDataQuery(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	opts = append(c.CallOptions.RepeatDataQuery[0:len(c.CallOptions.RepeatDataQuery):len(c.CallOptions.RepeatDataQuery)], opts...)
	return c.internalClient.RepeatDataQuery(ctx, req, opts)
}
func (c *ComplianceClient) RepeatDataSimplePath(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	opts = append(c.CallOptions.RepeatDataSimplePath[0:len(c.CallOptions.RepeatDataSimplePath):len(c.CallOptions.RepeatDataSimplePath)], opts...)
	return c.internalClient.RepeatDataSimplePath(ctx, req, opts)
}
func (c *ComplianceClient) RepeatDataPathResource(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	opts = append(c.CallOptions.RepeatDataPathResource[0:len(c.CallOptions.RepeatDataPathResource):len(c.CallOptions.RepeatDataPathResource)], opts...)
	return c.internalClient.RepeatDataPathResource(ctx, req, opts)
}
func (c *ComplianceClient) RepeatDataPathTrailingResource(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	opts = append(c.CallOptions.RepeatDataPathTrailingResource[0:len(c.CallOptions.RepeatDataPathTrailingResource):len(c.CallOptions.RepeatDataPathTrailingResource)], opts...)
	return c.internalClient.RepeatDataPathTrailingResource(ctx, req, opts)
}

// complianceGrpcClient is a client for interacting with Client Libraries Showcase API over gRPC transport.
// It satisfies the complianceinternalClient interface.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type complianceGrpcClient struct {
	// Connection pool of gRPC connections to the service.
	connPool gtransport.ConnPool

	// flag to opt out of default deadlines via GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE
	disableDeadlines bool

	// The gRPC API client.
	complianceClient genprotopb.ComplianceClient

	operationsClient longrunningpb.OperationsClient

	iamPolicyClient iampb.IAMPolicyClient

	locationsClient locationpb.LocationsClient

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewComplianceClient creates a new compliance client based on gRPC.
//
// This service is used to test that GAPICs can transcode proto3 requests to
// REST format correctly for various types of HTTP annotations.
func NewComplianceClient(ctx context.Context, opts ...option.ClientOption) (*ComplianceClient, error) {
	clientOpts := defaultComplianceGrpcClientOptions()
	if newComplianceGrpcClientHook != nil {
		hookOpts, err := newComplianceGrpcClientHook(ctx, clientHookParams{})
		if err != nil {
			return nil, err
		}
		clientOpts = append(clientOpts, hookOpts...)
	}

	disableDeadlines, err := checkDisableDeadlines()
	if err != nil {
		return nil, err
	}

	connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	c := &complianceGrpcClient{
		connPool:         connPool,
		disableDeadlines: disableDeadlines,
		complianceClient: genprotopb.NewComplianceClient(connPool),
	}
	c.setGoogleClientInfo()

	c.operationsClient = longrunningpb.NewOperationsClient(connPool)

	c.iamPolicyClient = iampb.NewIAMPolicyClient(connPool)

	c.locationsClient = locationpb.NewLocationsClient(connPool)

	return &ComplianceClient{internalComplianceClient: c, CallOptions: defaultComplianceCallOptions()}, nil
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *complianceGrpcClient) Connection() *grpc.ClientConn {
	return c.connPool.Conn()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *complianceGrpcClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *complianceGrpcClient) Close() error {
	return c.connPool.Close()
}

// RepeatDataBody this method echoes the ComplianceData request. This method exercises
// sending the entire request object in the REST body.
func (c *complianceGrpcClient) RepeatDataBody(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	var resp *genprotopb.RepeatResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.complianceClient.RepeatDataBody(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RepeatDataBodyInfo this method echoes the ComplianceData request. This method exercises
// sending the a message-type field in the REST body. Per AIP-127, only
// top-level, non-repeated fields can be sent this way.
func (c *complianceGrpcClient) RepeatDataBodyInfo(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	var resp *genprotopb.RepeatResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.complianceClient.RepeatDataBodyInfo(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RepeatDataQuery this method echoes the ComplianceData request. This method exercises
// sending all request fields as query parameters.
func (c *complianceGrpcClient) RepeatDataQuery(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	var resp *genprotopb.RepeatResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.complianceClient.RepeatDataQuery(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RepeatDataSimplePath this method echoes the ComplianceData request. This method exercises
// sending some parameters as “simple” path variables (i.e., of the form
// “/bar/{foo}” rather than “/{foo=bar/*}”), and the rest as query parameters.
func (c *complianceGrpcClient) RepeatDataSimplePath(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v&%s=%v&%s=%v&%s=%v", "info.f_string", url.QueryEscape(req.GetInfo().GetFString()), "info.f_int32", req.GetInfo().GetFInt32(), "info.f_double", url.QueryEscape(fmt.Sprintf("%g", req.GetInfo().GetFDouble())), "info.f_bool", req.GetInfo().GetFBool()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *genprotopb.RepeatResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.complianceClient.RepeatDataSimplePath(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RepeatDataPathResource same as RepeatDataSimplePath, but with a path resource.
func (c *complianceGrpcClient) RepeatDataPathResource(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v&%s=%v&%s=%v", "info.f_string", url.QueryEscape(req.GetInfo().GetFString()), "info.f_child.f_string", url.QueryEscape(req.GetInfo().GetFChild().GetFString()), "info.f_bool", req.GetInfo().GetFBool()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *genprotopb.RepeatResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.complianceClient.RepeatDataPathResource(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RepeatDataPathTrailingResource same as RepeatDataSimplePath, but with a trailing resource.
func (c *complianceGrpcClient) RepeatDataPathTrailingResource(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v&%s=%v", "info.f_string", url.QueryEscape(req.GetInfo().GetFString()), "info.f_child.f_string", url.QueryEscape(req.GetInfo().GetFChild().GetFString())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *genprotopb.RepeatResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.complianceClient.RepeatDataPathTrailingResource(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *complianceGrpcClient) ListLocations(ctx context.Context, req *locationpb.ListLocationsRequest, opts ...gax.CallOption) *LocationIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	it := &LocationIterator{}
	req = proto.Clone(req).(*locationpb.ListLocationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*locationpb.Location, string, error) {
		var resp *locationpb.ListLocationsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.locationsClient.ListLocations(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.GetLocations(), resp.GetNextPageToken(), nil
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
	it.pageInfo.MaxSize = int(req.GetPageSize())
	it.pageInfo.Token = req.GetPageToken()
	return it
}

func (c *complianceGrpcClient) GetLocation(ctx context.Context, req *locationpb.GetLocationRequest, opts ...gax.CallOption) (*locationpb.Location, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *locationpb.Location
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.locationsClient.GetLocation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *complianceGrpcClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "resource", url.QueryEscape(req.GetResource())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *iampb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.iamPolicyClient.SetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *complianceGrpcClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "resource", url.QueryEscape(req.GetResource())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *iampb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.iamPolicyClient.GetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *complianceGrpcClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "resource", url.QueryEscape(req.GetResource())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *iampb.TestIamPermissionsResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.iamPolicyClient.TestIamPermissions(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *complianceGrpcClient) ListOperations(ctx context.Context, req *longrunningpb.ListOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	it := &OperationIterator{}
	req = proto.Clone(req).(*longrunningpb.ListOperationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*longrunningpb.Operation, string, error) {
		var resp *longrunningpb.ListOperationsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.operationsClient.ListOperations(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.GetOperations(), resp.GetNextPageToken(), nil
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
	it.pageInfo.MaxSize = int(req.GetPageSize())
	it.pageInfo.Token = req.GetPageToken()
	return it
}

func (c *complianceGrpcClient) GetOperation(ctx context.Context, req *longrunningpb.GetOperationRequest, opts ...gax.CallOption) (*longrunningpb.Operation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.operationsClient.GetOperation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *complianceGrpcClient) DeleteOperation(ctx context.Context, req *longrunningpb.DeleteOperationRequest, opts ...gax.CallOption) error {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.operationsClient.DeleteOperation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

func (c *complianceGrpcClient) CancelOperation(ctx context.Context, req *longrunningpb.CancelOperationRequest, opts ...gax.CallOption) error {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.operationsClient.CancelOperation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

func (c *complianceGrpcClient) WaitOperation(ctx context.Context, req *longrunningpb.WaitOperationRequest, opts ...gax.CallOption) (*longrunningpb.Operation, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.operationsClient.WaitOperation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// LocationIterator manages a stream of *locationpb.Location.
type LocationIterator struct {
	items    []*locationpb.Location
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// Response is the raw response for the current page.
	// It must be cast to the RPC response type.
	// Calling Next() or InternalFetch() updates this value.
	Response interface{}

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*locationpb.Location, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *LocationIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *LocationIterator) Next() (*locationpb.Location, error) {
	var item *locationpb.Location
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *LocationIterator) bufLen() int {
	return len(it.items)
}

func (it *LocationIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

// OperationIterator manages a stream of *longrunningpb.Operation.
type OperationIterator struct {
	items    []*longrunningpb.Operation
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// Response is the raw response for the current page.
	// It must be cast to the RPC response type.
	// Calling Next() or InternalFetch() updates this value.
	Response interface{}

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*longrunningpb.Operation, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *OperationIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *OperationIterator) Next() (*longrunningpb.Operation, error) {
	var item *longrunningpb.Operation
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *OperationIterator) bufLen() int {
	return len(it.items)
}

func (it *OperationIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
