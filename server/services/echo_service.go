// Copyright 2018 Google LLC
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

package services

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	errdetails "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// NewEchoServer returns a new EchoServer for the Showcase API.
func NewEchoServer() pb.EchoServer {
	return &echoServerImpl{waiter: server.GetWaiterInstance()}
}

type echoServerImpl struct {
	waiter server.Waiter
}

func (s *echoServerImpl) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	err := status.ErrorProto(in.GetError())
	if err != nil {
		return nil, err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	requestHeaders := make(map[string]*pb.EchoResponse_RepeatedValues)
	if ok {
		for k, v := range md {
			requestHeaders[k] = &pb.EchoResponse_RepeatedValues{HeaderValues: v}
		}
	}
	echoHeaders(ctx)
	echoTrailers(ctx)
	return &pb.EchoResponse{Content: in.GetContent(), Severity: in.GetSeverity(), RequestId: in.GetRequestId(), OtherRequestId: in.GetOtherRequestId(), HttpRequestHeaderValue: requestHeaders}, nil
}

func (s *echoServerImpl) EchoErrorDetails(ctx context.Context, in *pb.EchoErrorDetailsRequest) (*pb.EchoErrorDetailsResponse, error) {
	var singleDetailError *pb.EchoErrorDetailsResponse_SingleDetail
	singleDetailText := in.GetSingleDetailText()
	if len(singleDetailText) > 0 {
		singleErrorInfo := &errdetails.ErrorInfo{Reason: singleDetailText}
		singleMarshalledError, err := anypb.New(singleErrorInfo)
		if err != nil {
			return nil, fmt.Errorf("failure with single error detail in EchoErrorDetails: %w", err)
		}
		singleDetailError = &pb.EchoErrorDetailsResponse_SingleDetail{
			Error: &pb.ErrorWithSingleDetail{Details: singleMarshalledError},
		}
	}

	var multipleDetailsError *pb.EchoErrorDetailsResponse_MultipleDetails
	multipleDetailText := in.GetMultiDetailText()
	if len(multipleDetailText) > 0 {
		details := []*anypb.Any{}
		for idx, text := range multipleDetailText {
			errorInfo := &errdetails.ErrorInfo{
				Reason: text,
			}
			marshalledError, err := anypb.New(errorInfo)
			if err != nil {
				return nil, fmt.Errorf("failure in EchoErrorDetails[%d]: %w", idx, err)
			}

			details = append(details, marshalledError)
		}

		multipleDetailsError = &pb.EchoErrorDetailsResponse_MultipleDetails{
			Error: &pb.ErrorWithMultipleDetails{Details: details},
		}
	}

	echoHeaders(ctx)
	echoTrailers(ctx)
	response := &pb.EchoErrorDetailsResponse{
		SingleDetail:    singleDetailError,
		MultipleDetails: multipleDetailsError,
	}
	return response, nil
}

func (s *echoServerImpl) Expand(in *pb.ExpandRequest, stream pb.Echo_ExpandServer) error {
	for _, word := range strings.Fields(in.GetContent()) {
		err := stream.Send(&pb.EchoResponse{Content: word})
		if err != nil {
			return err
		}
		time.Sleep(in.GetStreamWaitTime().AsDuration())
	}
	echoStreamingHeaders(stream)
	if in.GetError() != nil {
		return status.ErrorProto(in.GetError())
	}
	echoStreamingTrailers(stream)
	return nil
}

func (s *echoServerImpl) Collect(stream pb.Echo_CollectServer) error {
	var resp []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			echoStreamingHeaders(stream)
			echoStreamingTrailers(stream)
			return stream.SendAndClose(&pb.EchoResponse{Content: strings.Join(resp, " ")})
		}
		if err != nil {
			return err
		}
		s := status.ErrorProto(req.GetError())
		if s != nil {
			return s
		}
		if req.GetContent() != "" {
			resp = append(resp, req.GetContent())
		}
	}
}

func (s *echoServerImpl) Chat(stream pb.Echo_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Echo headers and trailers when the stream ends
			echoStreamingHeaders(stream)
			echoStreamingTrailers(stream)
			return nil
		}
		if err != nil {
			return err
		}

		s := status.ErrorProto(req.GetError())
		if s != nil {
			return s
		}
		stream.Send(&pb.EchoResponse{Content: req.GetContent()})
	}
}

func (s *echoServerImpl) PagedExpandLegacy(ctx context.Context, in *pb.PagedExpandLegacyRequest) (*pb.PagedExpandResponse, error) {
	req := &pb.PagedExpandRequest{
		Content:   in.Content,
		PageSize:  in.MaxResults,
		PageToken: in.PageToken,
	}
	return s.PagedExpand(ctx, req)
}

func (s *echoServerImpl) PagedExpand(ctx context.Context, in *pb.PagedExpandRequest) (*pb.PagedExpandResponse, error) {
	words := strings.Fields(in.GetContent())

	start, end, nextToken, err := processPageTokens(len(words), in.GetPageSize(), in.GetPageToken())
	if err != nil {
		return nil, err
	}

	responses := []*pb.EchoResponse{}
	for _, word := range words[start:end] {
		responses = append(responses, &pb.EchoResponse{Content: word})
	}

	echoHeaders(ctx)
	echoTrailers(ctx)
	return &pb.PagedExpandResponse{
		Responses:     responses,
		NextPageToken: nextToken,
	}, nil
}

func (s *echoServerImpl) PagedExpandLegacyMapped(ctx context.Context, in *pb.PagedExpandRequest) (*pb.PagedExpandLegacyMappedResponse, error) {
	words := strings.Fields(in.GetContent())
	start, end, nextToken, err := processPageTokens(len(words), in.GetPageSize(), in.GetPageToken())
	if err != nil {
		return nil, err
	}

	// Construct a map with the following properties:
	//
	// 1. The map has a one-rune string key corresponding to the first rune of EVERY word in words.
	// 2. The value corresponding to a given rune key is a list of only those words between
	// `start` and `end` whose first rune is that key.
	// 3. Consequently, initial runes that only appear outside the [start,end) range will have
	// empty list entries, even if they are non-empty in subsequent pages.
	alphabetized := make(map[string]*pb.PagedExpandResponseList, 255) //assume most input is ASCII
	for idx, word := range words {
		initialRune, _ := utf8.DecodeRuneInString(word)
		key := string(initialRune) // enforces #1
		prev, ok := alphabetized[key]
		if !ok {
			prev = &pb.PagedExpandResponseList{} // enforces #3
			alphabetized[key] = prev
		}
		if int32(idx) >= start && int32(idx) < end { // enforces #2
			prev.Words = append(prev.Words, word)
		}
	}

	echoHeaders(ctx)
	echoTrailers(ctx)
	return &pb.PagedExpandLegacyMappedResponse{
		Alphabetized:  alphabetized,
		NextPageToken: nextToken,
	}, nil
}

func processPageTokens(numElements int, pageSize int32, pageToken string) (start, end int32, nextToken string, err error) {
	if pageSize < 0 {
		return 0, 0, "", status.Error(codes.InvalidArgument, "the page size provided must not be negative.")
	}

	if pageToken != "" {
		token, err := strconv.Atoi(pageToken)
		token32 := int32(token)
		if err != nil || token32 < 0 || token32 >= int32(numElements) {
			return 0, 0, "", status.Errorf(
				codes.InvalidArgument,
				"invalid page token: %s. Token must be within the range [0, %d)",
				pageToken,
				numElements)
		}
		start = token32
	}

	if pageSize == 0 {
		pageSize = int32(numElements)
	}
	end = min(start+pageSize, int32(numElements))

	if end < int32(numElements) {
		nextToken = strconv.Itoa(int(end))
	}

	return start, end, nextToken, nil
}

func min(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func (s *echoServerImpl) Wait(ctx context.Context, in *pb.WaitRequest) (*lropb.Operation, error) {
	echoHeaders(ctx)
	echoTrailers(ctx)
	return s.waiter.Wait(in), nil
}

func (s *echoServerImpl) Block(ctx context.Context, in *pb.BlockRequest) (*pb.BlockResponse, error) {
	d, _ := ptypes.Duration(in.GetResponseDelay())
	time.Sleep(d)
	if in.GetError() != nil {
		return nil, status.ErrorProto(in.GetError())
	}
	echoHeaders(ctx)
	echoTrailers(ctx)
	return in.GetSuccess(), nil
}

// echo any provided headers in the metadata
func echoHeaders(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	values := md.Get("x-goog-request-params")
	for _, value := range values {
		header := metadata.Pairs("x-goog-request-params", value)
		grpc.SetHeader(ctx, header)
	}
}

func echoStreamingHeaders(stream grpc.ServerStream) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return
	}
	values := md.Get("x-goog-request-params")
	for _, value := range values {
		header := metadata.Pairs("x-goog-request-params", value)
		if stream.SetHeader(header) != nil {
			return
		}
	}
}

// echo any provided trailing metadata
func echoTrailers(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	values := md.Get("showcase-trailer")
	for _, value := range values {
		trailer := metadata.Pairs("showcase-trailer", value)
		grpc.SetTrailer(ctx, trailer)
	}
}

func echoStreamingTrailers(stream grpc.ServerStream) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return
	}

	values := md.Get("showcase-trailer")
	for _, value := range values {
		trailer := metadata.Pairs("showcase-trailer", value)
		stream.SetTrailer(trailer)
	}
}
