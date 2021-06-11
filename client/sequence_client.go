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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	httptransport "google.golang.org/api/transport/http"
	locationpb "google.golang.org/genproto/googleapis/cloud/location"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

var newSequenceClientHook clientHook

// SequenceCallOptions contains the retry settings for each method of SequenceClient.
type SequenceCallOptions struct {
	CreateSequence     []gax.CallOption
	GetSequenceReport  []gax.CallOption
	AttemptSequence    []gax.CallOption
	ListLocations      []gax.CallOption
	GetLocation        []gax.CallOption
	SetIamPolicy       []gax.CallOption
	GetIamPolicy       []gax.CallOption
	TestIamPermissions []gax.CallOption
	ListOperations     []gax.CallOption
	GetOperation       []gax.CallOption
	DeleteOperation    []gax.CallOption
	CancelOperation    []gax.CallOption
}

func defaultSequenceGRPCClientOptions() []option.ClientOption {
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

func defaultSequenceCallOptions() *SequenceCallOptions {
	return &SequenceCallOptions{
		CreateSequence:    []gax.CallOption{},
		GetSequenceReport: []gax.CallOption{},
		AttemptSequence: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.Unavailable,
					codes.Unknown,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        3000 * time.Millisecond,
					Multiplier: 2.00,
				})
			}),
		},
		ListLocations:      []gax.CallOption{},
		GetLocation:        []gax.CallOption{},
		SetIamPolicy:       []gax.CallOption{},
		GetIamPolicy:       []gax.CallOption{},
		TestIamPermissions: []gax.CallOption{},
		ListOperations:     []gax.CallOption{},
		GetOperation:       []gax.CallOption{},
		DeleteOperation:    []gax.CallOption{},
		CancelOperation:    []gax.CallOption{},
	}
}

// internalSequenceClient is an interface that defines the methods availaible from Client Libraries Showcase API.
type internalSequenceClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	CreateSequence(context.Context, *genprotopb.CreateSequenceRequest, ...gax.CallOption) (*genprotopb.Sequence, error)
	GetSequenceReport(context.Context, *genprotopb.GetSequenceReportRequest, ...gax.CallOption) (*genprotopb.SequenceReport, error)
	AttemptSequence(context.Context, *genprotopb.AttemptSequenceRequest, ...gax.CallOption) error
	ListLocations(context.Context, *locationpb.ListLocationsRequest, ...gax.CallOption) *LocationIterator
	GetLocation(context.Context, *locationpb.GetLocationRequest, ...gax.CallOption) (*locationpb.Location, error)
	SetIamPolicy(context.Context, *iampb.SetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
	GetIamPolicy(context.Context, *iampb.GetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
	TestIamPermissions(context.Context, *iampb.TestIamPermissionsRequest, ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error)
	ListOperations(context.Context, *longrunningpb.ListOperationsRequest, ...gax.CallOption) *OperationIterator
	GetOperation(context.Context, *longrunningpb.GetOperationRequest, ...gax.CallOption) (*longrunningpb.Operation, error)
	DeleteOperation(context.Context, *longrunningpb.DeleteOperationRequest, ...gax.CallOption) error
	CancelOperation(context.Context, *longrunningpb.CancelOperationRequest, ...gax.CallOption) error
}

// SequenceClient is a client for interacting with Client Libraries Showcase API.
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type SequenceClient struct {
	// The internal transport-dependent client.
	internalClient internalSequenceClient

	// The call options for this service.
	CallOptions *SequenceCallOptions
}

// Wrapper methods routed to the internal client.

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *SequenceClient) Close() error {
	return c.internalClient.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *SequenceClient) setGoogleClientInfo(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *SequenceClient) Connection() *grpc.ClientConn {
	return c.internalClient.Connection()
}

// CreateSequence creates a sequence.
func (c *SequenceClient) CreateSequence(ctx context.Context, req *genprotopb.CreateSequenceRequest, opts ...gax.CallOption) (*genprotopb.Sequence, error) {
	return c.internalClient.CreateSequence(ctx, req, opts...)
}

// GetSequenceReport retrieves a sequence.
func (c *SequenceClient) GetSequenceReport(ctx context.Context, req *genprotopb.GetSequenceReportRequest, opts ...gax.CallOption) (*genprotopb.SequenceReport, error) {
	return c.internalClient.GetSequenceReport(ctx, req, opts...)
}

// AttemptSequence attempts a sequence.
func (c *SequenceClient) AttemptSequence(ctx context.Context, req *genprotopb.AttemptSequenceRequest, opts ...gax.CallOption) error {
	return c.internalClient.AttemptSequence(ctx, req, opts...)
}

// ListLocations is a utility method from google.cloud.location.Locations.
func (c *SequenceClient) ListLocations(ctx context.Context, req *locationpb.ListLocationsRequest, opts ...gax.CallOption) *LocationIterator {
	return c.internalClient.ListLocations(ctx, req, opts...)
}

// GetLocation is a utility method from google.cloud.location.Locations.
func (c *SequenceClient) GetLocation(ctx context.Context, req *locationpb.GetLocationRequest, opts ...gax.CallOption) (*locationpb.Location, error) {
	return c.internalClient.GetLocation(ctx, req, opts...)
}

// SetIamPolicy is a utility method from google.iam.v1.IAMPolicy.
func (c *SequenceClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return c.internalClient.SetIamPolicy(ctx, req, opts...)
}

// GetIamPolicy is a utility method from google.iam.v1.IAMPolicy.
func (c *SequenceClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return c.internalClient.GetIamPolicy(ctx, req, opts...)
}

// TestIamPermissions is a utility method from google.iam.v1.IAMPolicy.
func (c *SequenceClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	return c.internalClient.TestIamPermissions(ctx, req, opts...)
}

// ListOperations is a utility method from google.longrunning.Operations.
func (c *SequenceClient) ListOperations(ctx context.Context, req *longrunningpb.ListOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	return c.internalClient.ListOperations(ctx, req, opts...)
}

// GetOperation is a utility method from google.longrunning.Operations.
func (c *SequenceClient) GetOperation(ctx context.Context, req *longrunningpb.GetOperationRequest, opts ...gax.CallOption) (*longrunningpb.Operation, error) {
	return c.internalClient.GetOperation(ctx, req, opts...)
}

// DeleteOperation is a utility method from google.longrunning.Operations.
func (c *SequenceClient) DeleteOperation(ctx context.Context, req *longrunningpb.DeleteOperationRequest, opts ...gax.CallOption) error {
	return c.internalClient.DeleteOperation(ctx, req, opts...)
}

// CancelOperation is a utility method from google.longrunning.Operations.
func (c *SequenceClient) CancelOperation(ctx context.Context, req *longrunningpb.CancelOperationRequest, opts ...gax.CallOption) error {
	return c.internalClient.CancelOperation(ctx, req, opts...)
}

// sequenceGRPCClient is a client for interacting with Client Libraries Showcase API over gRPC transport.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type sequenceGRPCClient struct {
	// Connection pool of gRPC connections to the service.
	connPool gtransport.ConnPool

	// flag to opt out of default deadlines via GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE
	disableDeadlines bool

	// Points back to the CallOptions field of the containing SequenceClient
	CallOptions **SequenceCallOptions

	// The gRPC API client.
	sequenceClient genprotopb.SequenceServiceClient

	operationsClient longrunningpb.OperationsClient

	iamPolicyClient iampb.IAMPolicyClient

	locationsClient locationpb.LocationsClient

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewSequenceClient creates a new sequence service client based on gRPC.
// The returned client must be Closed when it is done being used to clean up its underlying connections.
func NewSequenceClient(ctx context.Context, opts ...option.ClientOption) (*SequenceClient, error) {
	clientOpts := defaultSequenceGRPCClientOptions()
	if newSequenceClientHook != nil {
		hookOpts, err := newSequenceClientHook(ctx, clientHookParams{})
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
	client := SequenceClient{CallOptions: defaultSequenceCallOptions()}

	c := &sequenceGRPCClient{
		connPool:         connPool,
		disableDeadlines: disableDeadlines,
		sequenceClient:   genprotopb.NewSequenceServiceClient(connPool),
		CallOptions:      &client.CallOptions,
		operationsClient: longrunningpb.NewOperationsClient(connPool),
		iamPolicyClient:  iampb.NewIAMPolicyClient(connPool),
		locationsClient:  locationpb.NewLocationsClient(connPool),
	}
	c.setGoogleClientInfo()

	client.internalClient = c

	return &client, nil
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *sequenceGRPCClient) Connection() *grpc.ClientConn {
	return c.connPool.Conn()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *sequenceGRPCClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *sequenceGRPCClient) Close() error {
	return c.connPool.Close()
}

// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type sequenceRESTClient struct {
	// The http endpoint to connect to.
	endpoint string

	// The http client.
	httpClient *http.Client

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewSequenceRESTClient creates a new sequence service rest client.
func NewSequenceRESTClient(ctx context.Context, opts ...option.ClientOption) (*SequenceClient, error) {
	clientOpts := defaultSequenceRESTClientOptions()
	httpClient, endpoint, err := httptransport.NewClient(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}

	c := &sequenceRESTClient{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
	c.setGoogleClientInfo()

	return &SequenceClient{internalClient: c, CallOptions: defaultSequenceCallOptions()}, nil
}

func defaultSequenceRESTClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("localhost:7469"),
		internaloption.WithDefaultMTLSEndpoint("localhost:7469"),
		internaloption.WithDefaultAudience("https://localhost/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
	}
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *sequenceRESTClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "rest", "UNKNOWN")
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *sequenceRESTClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *sequenceRESTClient) Connection() *grpc.ClientConn {
	return nil
}
func (c *sequenceGRPCClient) CreateSequence(ctx context.Context, req *genprotopb.CreateSequenceRequest, opts ...gax.CallOption) (*genprotopb.Sequence, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).CreateSequence[0:len((*c.CallOptions).CreateSequence):len((*c.CallOptions).CreateSequence)], opts...)
	var resp *genprotopb.Sequence
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.sequenceClient.CreateSequence(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *sequenceGRPCClient) GetSequenceReport(ctx context.Context, req *genprotopb.GetSequenceReportRequest, opts ...gax.CallOption) (*genprotopb.SequenceReport, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetSequenceReport[0:len((*c.CallOptions).GetSequenceReport):len((*c.CallOptions).GetSequenceReport)], opts...)
	var resp *genprotopb.SequenceReport
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.sequenceClient.GetSequenceReport(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *sequenceGRPCClient) AttemptSequence(ctx context.Context, req *genprotopb.AttemptSequenceRequest, opts ...gax.CallOption) error {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 10000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).AttemptSequence[0:len((*c.CallOptions).AttemptSequence):len((*c.CallOptions).AttemptSequence)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.sequenceClient.AttemptSequence(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

func (c *sequenceGRPCClient) ListLocations(ctx context.Context, req *locationpb.ListLocationsRequest, opts ...gax.CallOption) *LocationIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).ListLocations[0:len((*c.CallOptions).ListLocations):len((*c.CallOptions).ListLocations)], opts...)
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

func (c *sequenceGRPCClient) GetLocation(ctx context.Context, req *locationpb.GetLocationRequest, opts ...gax.CallOption) (*locationpb.Location, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).GetLocation[0:len((*c.CallOptions).GetLocation):len((*c.CallOptions).GetLocation)], opts...)
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

func (c *sequenceGRPCClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).SetIamPolicy[0:len((*c.CallOptions).SetIamPolicy):len((*c.CallOptions).SetIamPolicy)], opts...)
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

func (c *sequenceGRPCClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).GetIamPolicy[0:len((*c.CallOptions).GetIamPolicy):len((*c.CallOptions).GetIamPolicy)], opts...)
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

func (c *sequenceGRPCClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).TestIamPermissions[0:len((*c.CallOptions).TestIamPermissions):len((*c.CallOptions).TestIamPermissions)], opts...)
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

func (c *sequenceGRPCClient) ListOperations(ctx context.Context, req *longrunningpb.ListOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).ListOperations[0:len((*c.CallOptions).ListOperations):len((*c.CallOptions).ListOperations)], opts...)
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

func (c *sequenceGRPCClient) GetOperation(ctx context.Context, req *longrunningpb.GetOperationRequest, opts ...gax.CallOption) (*longrunningpb.Operation, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).GetOperation[0:len((*c.CallOptions).GetOperation):len((*c.CallOptions).GetOperation)], opts...)
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

func (c *sequenceGRPCClient) DeleteOperation(ctx context.Context, req *longrunningpb.DeleteOperationRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).DeleteOperation[0:len((*c.CallOptions).DeleteOperation):len((*c.CallOptions).DeleteOperation)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.operationsClient.DeleteOperation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

func (c *sequenceGRPCClient) CancelOperation(ctx context.Context, req *longrunningpb.CancelOperationRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).CancelOperation[0:len((*c.CallOptions).CancelOperation):len((*c.CallOptions).CancelOperation)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.operationsClient.CancelOperation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// CreateSequence creates a sequence.
func (c *sequenceRESTClient) CreateSequence(ctx context.Context, req *genprotopb.CreateSequenceRequest, opts ...gax.CallOption) (*genprotopb.Sequence, error) {
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}

	body := req.GetSequence()
	queryParamStrs := []string{}
	if body.GetName() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("name=%v", body.GetName()))
	}
	if body.GetResponses() != nil {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("responses=%v", body.GetResponses()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s/v1beta1/sequences",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode >= http.StatusOK {
		return nil, fmt.Errorf(httpRsp.Status)
	}

	buf, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}
	rsp := &genprotopb.Sequence{}

	return rsp, unm.Unmarshal(buf, rsp)
}

// GetSequenceReport retrieves a sequence.
func (c *sequenceRESTClient) GetSequenceReport(ctx context.Context, req *genprotopb.GetSequenceReportRequest, opts ...gax.CallOption) (*genprotopb.SequenceReport, error) {
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}

	queryParamStrs := []string{}
	if req.GetName() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("name=%v", req.GetName()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s/v1beta1/{name=sequences/*/sequenceReport}",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode >= http.StatusOK {
		return nil, fmt.Errorf(httpRsp.Status)
	}

	buf, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}
	rsp := &genprotopb.SequenceReport{}

	return rsp, unm.Unmarshal(buf, rsp)
}

// AttemptSequence attempts a sequence.
func (c *sequenceRESTClient) AttemptSequence(ctx context.Context, req *genprotopb.AttemptSequenceRequest, opts ...gax.CallOption) error {
	// The default (false) for the other options are fine.
	// Field names should be lowerCamel, not snake.
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true, UseProtoNames: false}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return err
	}

	queryParamStrs := []string{}
	if req.GetName() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("name=%v", req.GetName()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s/v1beta1/{name=sequences/*}",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonReq))
	if err != nil {
		return err
	}

	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf(httpRsp.Status)
	}

	return nil
}

// ListLocations is a utility method from google.cloud.location.Locations.
func (c *sequenceRESTClient) ListLocations(ctx context.Context, req *locationpb.ListLocationsRequest, opts ...gax.CallOption) *LocationIterator {
	it := &LocationIterator{}
	req = proto.Clone(req).(*locationpb.ListLocationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*locationpb.Location, string, error) {
		return nil, "", fmt.Errorf("ListLocations not yet supported for REST clients")
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

// GetLocation is a utility method from google.cloud.location.Locations.
func (c *sequenceRESTClient) GetLocation(ctx context.Context, req *locationpb.GetLocationRequest, opts ...gax.CallOption) (*locationpb.Location, error) {
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}

	queryParamStrs := []string{}
	if req.GetName() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("name=%v", req.GetName()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "", url, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode >= http.StatusOK {
		return nil, fmt.Errorf(httpRsp.Status)
	}

	buf, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}
	rsp := &locationpb.Location{}

	return rsp, unm.Unmarshal(buf, rsp)
}

// SetIamPolicy is a utility method from google.iam.v1.IAMPolicy.
func (c *sequenceRESTClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}

	queryParamStrs := []string{}
	if req.GetPolicy() != nil {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("policy=%v", req.GetPolicy()))
	}
	if req.GetResource() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("resource=%v", req.GetResource()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "", url, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode >= http.StatusOK {
		return nil, fmt.Errorf(httpRsp.Status)
	}

	buf, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}
	rsp := &iampb.Policy{}

	return rsp, unm.Unmarshal(buf, rsp)
}

// GetIamPolicy is a utility method from google.iam.v1.IAMPolicy.
func (c *sequenceRESTClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}

	queryParamStrs := []string{}
	if req.GetOptions() != nil {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("options=%v", req.GetOptions()))
	}
	if req.GetResource() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("resource=%v", req.GetResource()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "", url, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode >= http.StatusOK {
		return nil, fmt.Errorf(httpRsp.Status)
	}

	buf, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}
	rsp := &iampb.Policy{}

	return rsp, unm.Unmarshal(buf, rsp)
}

// TestIamPermissions is a utility method from google.iam.v1.IAMPolicy.
func (c *sequenceRESTClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}

	queryParamStrs := []string{}
	if req.GetPermissions() != nil {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("permissions=%v", req.GetPermissions()))
	}
	if req.GetResource() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("resource=%v", req.GetResource()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "", url, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode >= http.StatusOK {
		return nil, fmt.Errorf(httpRsp.Status)
	}

	buf, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}
	rsp := &iampb.TestIamPermissionsResponse{}

	return rsp, unm.Unmarshal(buf, rsp)
}

// ListOperations is a utility method from google.longrunning.Operations.
func (c *sequenceRESTClient) ListOperations(ctx context.Context, req *longrunningpb.ListOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	it := &OperationIterator{}
	req = proto.Clone(req).(*longrunningpb.ListOperationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*longrunningpb.Operation, string, error) {
		return nil, "", fmt.Errorf("ListOperations not yet supported for REST clients")
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

// GetOperation is a utility method from google.longrunning.Operations.
func (c *sequenceRESTClient) GetOperation(ctx context.Context, req *longrunningpb.GetOperationRequest, opts ...gax.CallOption) (*longrunningpb.Operation, error) {
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}

	queryParamStrs := []string{}
	if req.GetName() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("name=%v", req.GetName()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "", url, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode >= http.StatusOK {
		return nil, fmt.Errorf(httpRsp.Status)
	}

	buf, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}
	rsp := &longrunningpb.Operation{}

	return rsp, unm.Unmarshal(buf, rsp)
}

// DeleteOperation is a utility method from google.longrunning.Operations.
func (c *sequenceRESTClient) DeleteOperation(ctx context.Context, req *longrunningpb.DeleteOperationRequest, opts ...gax.CallOption) error {
	// The default (false) for the other options are fine.
	// Field names should be lowerCamel, not snake.
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true, UseProtoNames: false}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return err
	}

	queryParamStrs := []string{}
	if req.GetName() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("name=%v", req.GetName()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "", url, bytes.NewReader(jsonReq))
	if err != nil {
		return err
	}

	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf(httpRsp.Status)
	}

	return nil
}

// CancelOperation is a utility method from google.longrunning.Operations.
func (c *sequenceRESTClient) CancelOperation(ctx context.Context, req *longrunningpb.CancelOperationRequest, opts ...gax.CallOption) error {
	// The default (false) for the other options are fine.
	// Field names should be lowerCamel, not snake.
	m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true, UseProtoNames: false}
	jsonReq, err := m.Marshal(req)
	if err != nil {
		return err
	}

	queryParamStrs := []string{}
	if req.GetName() != "" {
		queryParamStrs = append(queryParamStrs, fmt.Sprintf("name=%v", req.GetName()))
	}
	query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")

	url := fmt.Sprintf("%s",
		c.endpoint,
	)

	if query != "" {
		url += "?" + query
	}

	httpReq, err := http.NewRequestWithContext(ctx, "", url, bytes.NewReader(jsonReq))
	if err != nil {
		return err
	}

	// Set the headers
	for k, v := range c.xGoogMetadata {
		httpReq.Header[k] = v
	}
	httpReq.Header["Content-Type"] = []string{"application/json"}

	httpRsp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf(httpRsp.Status)
	}

	return nil
}
