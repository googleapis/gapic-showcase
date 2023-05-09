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
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	durpb "github.com/golang/protobuf/ptypes/duration"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

func TestEcho_success(t *testing.T) {
	table := []string{"hello world", ""}

	server := NewEchoServer()
	for _, val := range table {
		in := &pb.EchoRequest{
			Response: &pb.EchoRequest_Content{Content: val},
			Severity: pb.Severity_CRITICAL,
		}
		mockStream := &mockUnaryStream{t: t}
		ctx := appendTestOutgoingMetadata(context.Background(), &mockSTS{t: t, stream: mockStream})
		out, err := server.Echo(ctx, in)
		if err != nil {
			t.Error(err)
		}
		if out.GetContent() != in.GetContent() {
			t.Errorf("Echo(%s) returned %s", in.GetContent(), out.GetContent())
		}
		if out.Severity != in.Severity {
			t.Errorf("Echo severity(%d) returned %d", in.Severity, out.Severity)
		}
		mockStream.verify(err != nil)
	}
	in := &pb.EchoRequest{
		Response: &pb.EchoRequest_Error{
			Error: &spb.Status{Code: int32(codes.OK)}}}

	mockStream := &mockUnaryStream{t: t}
	ctx := appendTestOutgoingMetadata(context.Background(), &mockSTS{t: t, stream: mockStream})
	_, err := server.Echo(ctx, in)
	if err != nil {
		t.Error(err)
	}
	mockStream.verify(err != nil)
}

func TestEcho_error(t *testing.T) {
	table := []codes.Code{codes.Canceled, codes.InvalidArgument}

	server := NewEchoServer()
	for _, val := range table {
		in := &pb.EchoRequest{
			Response: &pb.EchoRequest_Error{
				Error: &spb.Status{Code: int32(val)}}}
		out, err := server.Echo(context.Background(), in)
		if out != nil {
			t.Errorf("Echo called with code %d returned a non-nil response proto.", val)
		}
		if err == nil {
			t.Errorf("Echo called with code %d did not return an error.", val)
		}
		status, _ := status.FromError(err)
		if status.Code() != val {
			t.Errorf("Echo called with code %d returned an error with code %d", val, status.Code())
		}
	}
}

type mockSTS struct {
	stream grpc.ServerStream
	t      *testing.T
}

func (m *mockSTS) Method() string                  { return "" }
func (m *mockSTS) SetHeader(md metadata.MD) error  { m.stream.SetHeader(md); return nil }
func (m *mockSTS) SendHeader(md metadata.MD) error { return m.stream.SendHeader(md) }
func (m *mockSTS) SetTrailer(md metadata.MD) error { m.stream.SetTrailer(md); return nil }

type mockUnaryStream struct {
	head  []string
	trail []string
	t     *testing.T
	grpc.ServerStream
}

func (m *mockUnaryStream) Method() string                   { return "" }
func (m *mockUnaryStream) Send(resp *pb.EchoResponse) error { return nil }
func (m *mockUnaryStream) Context() context.Context         { return nil }
func (m *mockUnaryStream) SetTrailer(md metadata.MD) {
	m.trail = append(m.trail, md.Get("showcase-trailer")...)
}
func (m *mockUnaryStream) SetHeader(md metadata.MD) error {
	m.head = append(m.head, md.Get("x-goog-request-params")...)
	return nil
}
func (m *mockUnaryStream) verify(expectHeadersAndTrailers bool) {
	if expectHeadersAndTrailers && !reflect.DeepEqual([]string{"show", "case"}, m.trail) && !reflect.DeepEqual([]string{"showcaseHeader, anotherHeader"}, m.head) {
		m.t.Errorf("Unary stream did not get all expected headers and trailers.\nGot these headers: %+v\nGot these trailers: %+v", m.head, m.trail)
	}
}

type mockExpandStream struct {
	exp   []string
	head  []string
	trail []string
	t     *testing.T
	pb.Echo_ExpandServer
}

func (m *mockExpandStream) Send(resp *pb.EchoResponse) error {
	if resp.GetContent() != m.exp[0] {
		m.t.Errorf("Expand expected to send %s but sent %s", m.exp[0], resp.GetContent())
	}
	m.exp = m.exp[1:]
	return nil
}

func (m *mockExpandStream) Context() context.Context {
	return appendTestOutgoingMetadata(context.Background(), &mockSTS{stream: m, t: m.t})
}

func (m *mockExpandStream) SetTrailer(md metadata.MD) {
	m.trail = append(m.trail, md.Get("showcase-trailer")...)
}

func (m *mockExpandStream) SetHeader(md metadata.MD) error {
	m.head = append(m.head, md.Get("x-goog-request-params")...)
	return nil
}

func (m *mockExpandStream) verify(expectHeadersAndTrailers bool) {
	if len(m.exp) > 0 {
		m.t.Errorf("Expand did not stream all expected values. %d expected values remaining.", len(m.exp))
	}
	if expectHeadersAndTrailers && !reflect.DeepEqual([]string{"show", "case"}, m.trail) && !reflect.DeepEqual([]string{"showcaseHeader", "anotherHeader"}, m.head) {
		m.t.Errorf("Expand did not get all expected headers and trailers.\nGot these headers: %+v\nGot these trailers: %+v", m.head, m.trail)
	}
}

func TestExpand(t *testing.T) {
	contentTable := []string{"Hello World", "Hola", ""}
	errTable := []*spb.Status{
		{Code: int32(codes.OK)},
		{Code: int32(codes.InvalidArgument)},
		nil,
	}

	server := NewEchoServer()
	for _, c := range contentTable {
		for _, e := range errTable {
			stream := &mockExpandStream{exp: strings.Fields(c), t: t}
			err := server.Expand(&pb.ExpandRequest{Content: c, Error: e}, stream)
			status, _ := status.FromError(err)
			if int32(status.Code()) != e.GetCode() {
				t.Errorf("Expand expected stream to return status with code %d but code %d", status.Code(), e.GetCode())
			}
			stream.verify(e == nil)
		}
	}
}

func TestExpandWithWaitTime(t *testing.T) {
	server := NewEchoServer()
	//This stream should take at least 300ms to complete because there are 7 messages, and we wait 50ms between sending each message.
	content := "This stream should take 300ms to complete"
	stream := &mockExpandStream{exp: strings.Fields(content), t: t}
	streamWaitTime := durationpb.New(time.Duration(50) * time.Millisecond)
	start := time.Now()

	err := server.Expand(&pb.ExpandRequest{Content: content, StreamWaitTime: streamWaitTime}, stream)

	actualTimeSpent := int(time.Since(start).Milliseconds())
	expectedTimeSent := 300
	if actualTimeSpent < expectedTimeSent {
		t.Error("Expand stream should take at least " + strconv.Itoa(expectedTimeSent) + "ms to complete, but it only took " + strconv.Itoa(actualTimeSpent) + "ms")
	}
	stream.verify(err == nil)
}

type errorExpandStream struct {
	err error
	pb.Echo_ExpandServer
}

func (s *errorExpandStream) Send(resp *pb.EchoResponse) error {
	return s.err
}

func TestExpand_streamErr(t *testing.T) {
	e := errors.New("Test Error")
	stream := &errorExpandStream{err: e}
	server := NewEchoServer()
	err := server.Expand(&pb.ExpandRequest{Content: "Hello World"}, stream)
	if e != err {
		t.Error("Expand expected to pass through stream errors.")
	}
}

type mockCollectStream struct {
	reqs  []*pb.EchoRequest
	head  []string
	trail []string
	exp   *string
	t     *testing.T
	pb.Echo_CollectServer
}

func (m *mockCollectStream) SendAndClose(r *pb.EchoResponse) error {
	if m.exp == nil {
		m.t.Errorf("Collect Stream SendAndClose called unexpectedly")
	}
	if r.GetContent() != *m.exp {
		m.t.Errorf("Collect expected to return '%s', but returned '%s'", *m.exp, r.GetContent())
	}
	return nil
}

func (m *mockCollectStream) Recv() (*pb.EchoRequest, error) {
	if len(m.reqs) > 0 {
		ret := m.reqs[0]
		m.reqs = m.reqs[1:]
		return ret, nil
	}
	return nil, io.EOF
}

func (m *mockCollectStream) Context() context.Context {
	return appendTestOutgoingMetadata(context.Background(), &mockSTS{stream: m, t: m.t})
}

func (m *mockCollectStream) SetHeader(md metadata.MD) error {
	m.head = append(m.head, md.Get("x-goog-request-params")...)
	return nil
}

func (m *mockCollectStream) SetTrailer(md metadata.MD) {
	m.trail = append(m.trail, md.Get("showcase-trailer")...)
}

func (m *mockCollectStream) verify(expectHeadersAndTrailers bool) {
	if expectHeadersAndTrailers && !reflect.DeepEqual([]string{"show", "case"}, m.trail) && !reflect.DeepEqual([]string{"showcaseHeader", "anotherHeader"}, m.head) {
		m.t.Errorf("Collect did not get all expected trailers.\nGot these headers: %+v\nGot these trailers: %+v", m.head, m.trail)
	}
}

func TestCollect(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	tests := []struct {
		reqs []string
		exp  *string
		err  *spb.Status
	}{
		{[]string{"Hello", "", "World"}, strPtr("Hello World"), nil},
		{[]string{"Hello", "World"}, strPtr("Hello World"), &spb.Status{Code: int32(codes.OK)}},
		{[]string{"Hello", "World"}, nil, &spb.Status{Code: int32(codes.InvalidArgument)}},
		{[]string{}, nil, &spb.Status{Code: int32(codes.InvalidArgument)}},
		{[]string{}, strPtr(""), nil},
	}

	server := NewEchoServer()
	for _, test := range tests {
		reqs := []*pb.EchoRequest{}
		for _, req := range test.reqs {
			reqs = append(reqs, &pb.EchoRequest{Response: &pb.EchoRequest_Content{Content: req}})
		}
		if test.err != nil {
			reqs = append(reqs, &pb.EchoRequest{Response: &pb.EchoRequest_Error{Error: test.err}})
		}
		mockStream := &mockCollectStream{reqs: reqs, exp: test.exp, t: t}
		err := server.Collect(mockStream)
		expCode := status.FromProto(test.err).Code()
		s, _ := status.FromError(err)
		if expCode != s.Code() {
			t.Errorf("Collect expected to return with code %d, but returned %d", expCode, s.Code())
		}
		mockStream.verify(test.err == nil)
	}
}

type errorCollectStream struct {
	err error
	pb.Echo_CollectServer
}

func (s *errorCollectStream) Recv() (*pb.EchoRequest, error) {
	return nil, s.err
}

func TestCollect_streamErr(t *testing.T) {
	e := errors.New("Test Error")
	stream := &errorCollectStream{err: e}
	server := NewEchoServer()
	err := server.Collect(stream)
	if e != err {
		t.Error("Collect expected to pass through stream errors.")
	}
}

type mockChatStream struct {
	reqs  []*pb.EchoRequest
	head  []string
	trail []string
	curr  *pb.EchoRequest
	t     *testing.T
	pb.Echo_ChatServer
}

func (m *mockChatStream) Recv() (*pb.EchoRequest, error) {
	if len(m.reqs) > 0 {
		m.curr = m.reqs[0]
		m.reqs = m.reqs[1:]
		return m.curr, nil
	}
	return nil, io.EOF
}

func (m *mockChatStream) Send(r *pb.EchoResponse) error {
	if m.curr == nil {
		m.t.Errorf("Chat unexpectedly tried to send content.")
	}
	if r.GetContent() != m.curr.GetContent() {
		m.t.Errorf("Chat expected to send content %s, but sent %s", m.curr.GetContent(), r.GetContent())
		m.curr = nil
	}
	return nil
}

func (m *mockChatStream) Context() context.Context {
	return appendTestOutgoingMetadata(context.Background(), &mockSTS{stream: m, t: m.t})
}

func (m *mockChatStream) SetHeader(md metadata.MD) error {
	m.head = append(m.head, md.Get("x-goog-request-params")...)
	return nil
}

func (m *mockChatStream) SetTrailer(md metadata.MD) {
	m.trail = append(m.trail, md.Get("showcase-trailer")...)
}

func (m *mockChatStream) verify(expectHeadersAndTrailers bool) {
	if expectHeadersAndTrailers && !reflect.DeepEqual([]string{"show", "case"}, m.trail) && !reflect.DeepEqual([]string{"showcaseHeader", "anotherHeader"}, m.head) {
		m.t.Errorf("Chat did not get all expected trailers.\nGot these headers: %+v\nGot these trailers: %+v", m.head, m.trail)
	}
}

func TestChat(t *testing.T) {
	tests := []struct {
		reqs []string
		err  *spb.Status
	}{
		{[]string{"Hello", "World"}, nil},
		{[]string{"Hello", "World"}, &spb.Status{Code: int32(codes.InvalidArgument)}},
		{[]string{}, &spb.Status{Code: int32(codes.InvalidArgument)}},
	}

	server := NewEchoServer()
	for _, test := range tests {
		reqs := []*pb.EchoRequest{}
		for _, req := range test.reqs {
			reqs = append(reqs, &pb.EchoRequest{Response: &pb.EchoRequest_Content{Content: req}})
		}
		if test.err != nil {
			reqs = append(reqs, &pb.EchoRequest{Response: &pb.EchoRequest_Error{Error: test.err}})
		}

		mockStream := &mockChatStream{reqs: reqs, t: t}
		err := server.Chat(mockStream)
		expCode := status.FromProto(test.err).Code()
		s, _ := status.FromError(err)
		if expCode != s.Code() {
			t.Errorf("Chat expected to return status with code %d, but returned %d", expCode, s.Code())
		}
		mockStream.verify(test.err == nil)
	}
}

type errorChatStream struct {
	err error
	pb.Echo_ChatServer
}

func (s *errorChatStream) Recv() (*pb.EchoRequest, error) {
	return nil, s.err
}

func TestChat_streamErr(t *testing.T) {
	e := errors.New("Test Error")
	stream := &errorChatStream{err: e}
	server := NewEchoServer()
	err := server.Chat(stream)
	if e != err {
		t.Error("Chat expected to pass through stream errors.")
	}
}

func TestPagedExpand_invalidArgs(t *testing.T) {
	tests := []*pb.PagedExpandRequest{
		{PageSize: -1},
		{PageToken: "-1"},
		{PageToken: "BOGUS"},
		{Content: "one", PageToken: "1"},
		{Content: "one", PageToken: "2"},
	}
	server := NewEchoServer()
	for _, in := range tests {
		_, err := server.PagedExpand(context.Background(), in)
		s, _ := status.FromError(err)
		if s.Code() != codes.InvalidArgument {
			t.Errorf("PagedExpand() expected error code: %d, got error code %d",
				codes.InvalidArgument, s.Code())
		}
	}
}

func TestPagedExpand(t *testing.T) {
	tests := []struct {
		in  *pb.PagedExpandRequest
		out *pb.PagedExpandResponse
	}{
		{
			&pb.PagedExpandRequest{Content: "Hello world!"},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "Hello"},
					{Content: "world!"},
				},
			},
		},
		{
			&pb.PagedExpandRequest{PageSize: 3, Content: "Hello world!"},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "Hello"},
					{Content: "world!"},
				},
			},
		},
		{
			&pb.PagedExpandRequest{
				PageSize: 3,
				Content:  "The rain in Spain falls mainly on the plain!",
			},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "The"},
					{Content: "rain"},
					{Content: "in"},
				},
				NextPageToken: "3",
			},
		},
		{
			&pb.PagedExpandRequest{
				PageSize:  3,
				PageToken: "3",
				Content:   "The rain in Spain falls mainly on the plain!",
			},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "Spain"},
					{Content: "falls"},
					{Content: "mainly"},
				},
				NextPageToken: "6",
			},
		},
	}

	server := NewEchoServer()
	for _, test := range tests {
		mockStream := &mockUnaryStream{t: t}
		ctx := appendTestOutgoingMetadata(context.Background(), &mockSTS{t: t, stream: mockStream})
		out, err := server.PagedExpand(ctx, test.in)
		if err != nil {
			t.Error(err)
		}
		if !proto.Equal(test.out, out) {
			t.Errorf("PagedExpand with input %q, expected: %q, got: %q",
				test.in.String(), test.out.String(), out.String())
		}
		mockStream.verify(err == nil)
	}
}

// NOTE: The TestPagedExpandLegacy*() tests mirror the TestPagedExpand*() tests.

func TestPagedExpandLegacy_invalidArgs(t *testing.T) {
	tests := []*pb.PagedExpandLegacyRequest{
		{MaxResults: -1},
		{PageToken: "-1"},
		{PageToken: "BOGUS"},
		{Content: "one", PageToken: "1"},
		{Content: "one", PageToken: "2"},
	}
	server := NewEchoServer()
	for _, in := range tests {
		_, err := server.PagedExpandLegacy(context.Background(), in)
		s, _ := status.FromError(err)
		if s.Code() != codes.InvalidArgument {
			t.Errorf("PagedExpandLegacy() expected error code: %d, got error code %d",
				codes.InvalidArgument, s.Code())
		}
	}
}

func TestPagedExpandLegacy(t *testing.T) {
	tests := []struct {
		in  *pb.PagedExpandLegacyRequest
		out *pb.PagedExpandResponse
	}{
		{
			&pb.PagedExpandLegacyRequest{Content: "Hello world!"},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "Hello"},
					{Content: "world!"},
				},
			},
		},
		{
			&pb.PagedExpandLegacyRequest{MaxResults: 3, Content: "Hello world!"},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "Hello"},
					{Content: "world!"},
				},
			},
		},
		{
			&pb.PagedExpandLegacyRequest{
				MaxResults: 3,
				Content:    "The rain in Spain falls mainly on the plain!",
			},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "The"},
					{Content: "rain"},
					{Content: "in"},
				},
				NextPageToken: "3",
			},
		},
		{
			&pb.PagedExpandLegacyRequest{
				MaxResults: 3,
				PageToken:  "3",
				Content:    "The rain in Spain falls mainly on the plain!",
			},
			&pb.PagedExpandResponse{
				Responses: []*pb.EchoResponse{
					{Content: "Spain"},
					{Content: "falls"},
					{Content: "mainly"},
				},
				NextPageToken: "6",
			},
		},
	}

	server := NewEchoServer()
	for _, test := range tests {
		mockStream := &mockUnaryStream{t: t}
		ctx := appendTestOutgoingMetadata(context.Background(), &mockSTS{t: t, stream: mockStream})
		out, err := server.PagedExpandLegacy(ctx, test.in)
		if err != nil {
			t.Error(err)
		}
		if !proto.Equal(test.out, out) {
			t.Errorf("PagedExpandLegacy with input %q, expected: %q, got: %q",
				test.in.String(), test.out.String(), out.String())
		}
		mockStream.verify(err == nil)
	}
}

func TestPagedExpandLegacyMapped_invalidArgs(t *testing.T) {
	tests := []*pb.PagedExpandRequest{
		{PageSize: -1},
		{PageToken: "-1"},
		{PageToken: "BOGUS"},
		{Content: "one", PageToken: "1"},
		{Content: "one", PageToken: "2"},
	}
	server := NewEchoServer()
	for _, in := range tests {
		_, err := server.PagedExpandLegacyMapped(context.Background(), in)
		s, _ := status.FromError(err)
		if s.Code() != codes.InvalidArgument {
			t.Errorf("PagedExpandLegacyMapped() expected error code: %d, got error code %d",
				codes.InvalidArgument, s.Code())
		}
	}
}

func TestPagedExpandLegacyMapped(t *testing.T) {
	text := "It was the best of times, it was the worst of times"
	tests := []struct {
		in  *pb.PagedExpandRequest
		out *pb.PagedExpandLegacyMappedResponse
	}{
		{
			&pb.PagedExpandRequest{Content: text},
			&pb.PagedExpandLegacyMappedResponse{
				Alphabetized: map[string]*pb.PagedExpandResponseList{
					"b": {Words: []string{"best"}},
					"I": {Words: []string{"It"}},
					"i": {Words: []string{"it"}},
					"o": {Words: []string{"of", "of"}},
					"t": {Words: []string{"the", "times,", "the", "times"}},
					"w": {Words: []string{"was", "was", "worst"}},
				},
			},
		},
		{
			&pb.PagedExpandRequest{PageSize: 1, Content: text},
			&pb.PagedExpandLegacyMappedResponse{
				Alphabetized: map[string]*pb.PagedExpandResponseList{
					"I": {Words: []string{"It"}},
					"b": {Words: []string{}},
					"i": {Words: []string{}},
					"o": {Words: []string{}},
					"t": {Words: []string{}},
					"w": {Words: []string{}},
				},
				NextPageToken: "1",
			},
		},
		{
			&pb.PagedExpandRequest{PageSize: 4, PageToken: "2", Content: text},
			&pb.PagedExpandLegacyMappedResponse{
				Alphabetized: map[string]*pb.PagedExpandResponseList{
					"b": {Words: []string{"best"}},
					"I": {Words: []string{}},
					"i": {Words: []string{}},
					"o": {Words: []string{"of"}},
					"t": {Words: []string{"the", "times,"}},
					"w": {Words: []string{}},
				},
				NextPageToken: "6",
			},
		},
		{
			&pb.PagedExpandRequest{PageSize: 4, PageToken: "8", Content: text},
			&pb.PagedExpandLegacyMappedResponse{
				Alphabetized: map[string]*pb.PagedExpandResponseList{
					"b": {Words: []string{}},
					"I": {Words: []string{}},
					"i": {Words: []string{}},
					"o": {Words: []string{"of"}},
					"t": {Words: []string{"the", "times"}},
					"w": {Words: []string{"worst"}},
				},
			},
		},
	}

	server := NewEchoServer()
	for _, test := range tests {
		mockStream := &mockUnaryStream{t: t}
		ctx := appendTestOutgoingMetadata(context.Background(), &mockSTS{t: t, stream: mockStream})
		out, err := server.PagedExpandLegacyMapped(ctx, test.in)
		if err != nil {
			t.Error(err)
		}
		if !proto.Equal(test.out, out) {
			t.Errorf("PagedExpandLegacyMapped with input %q:\n  expected: %#v\n       got: %#v\n",
				test.in.String(), test.out.String(), out.String())
		}
		mockStream.verify(err == nil)
	}
}

func TestWait(t *testing.T) {
	endTime, _ := ptypes.TimestampProto(time.Now())
	req := &pb.WaitRequest{End: &pb.WaitRequest_EndTime{EndTime: endTime}}
	waiter := &mockWaiter{}
	server := &echoServerImpl{waiter: waiter}
	mockStream := &mockUnaryStream{t: t}
	ctx := appendTestOutgoingMetadata(context.Background(), &mockSTS{t: t, stream: mockStream})
	server.Wait(ctx, req)
	if !proto.Equal(waiter.req, req) {
		t.Error("Expected echo.Wait to defer to waiter.")
	}
	mockStream.verify(true)
}

func TestBlockSuccess(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		resp    string
	}{
		{1, int32(1000), "hello"},
		{5, int32(10), "world"},
	}
	for _, test := range tests {
		waiter := &mockWaiter{}
		server := &echoServerImpl{waiter: waiter}
		in := &pb.BlockRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.BlockRequest_Success{
				Success: &pb.BlockResponse{Content: test.resp},
			},
		}
		mockStream := &mockUnaryStream{t: t}
		ctx := appendTestOutgoingMetadata(context.Background(), &mockSTS{t: t, stream: mockStream})
		out, err := server.Block(ctx, in)
		if err != nil {
			t.Error(err)
		}
		if out.GetContent() != test.resp {
			t.Errorf("Expected Wait test to return %s, but returned %s", out.GetContent(), test.resp)
		}
		mockStream.verify(err == nil)
	}
}

func TestBlockError(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		code    codes.Code
	}{
		{0, int32(0), codes.InvalidArgument},
		{2, int32(1000), codes.Unavailable},
	}

	for _, test := range tests {
		waiter := &mockWaiter{}
		server := &echoServerImpl{waiter: waiter}
		in := &pb.BlockRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.BlockRequest_Error{
				Error: status.New(test.code, "").Proto(),
			},
		}
		out, err := server.Block(context.Background(), in)
		if out != nil {
			t.Errorf("Block: Expected to error with code %d but returned success", test.code)
		}
		s, _ := status.FromError(err)
		if s.Code() != test.code {
			t.Errorf("Block: Expected to error with code %d but errored with code %d", test.code, s.Code())
		}
	}
}

func appendTestOutgoingMetadata(ctx context.Context, stream grpc.ServerTransportStream) context.Context {
	ctx = grpc.NewContextWithServerTransportStream(ctx, stream)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("showcase-trailer", "show", "showcase-trailer", "case", "trailer", "trail", "x-goog-request-params", "showcaseHeader", "x-goog-request-params", "anotherHeader", "header", "head"))
	return ctx
}
