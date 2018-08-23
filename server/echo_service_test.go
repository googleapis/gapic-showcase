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

package server

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	durpb "github.com/golang/protobuf/ptypes/duration"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEcho_success(t *testing.T) {
	table := []string{"hello world", ""}

	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
	for _, val := range table {
		in := &pb.EchoRequest{Response: &pb.EchoRequest_Content{Content: val}}
		out, err := server.Echo(context.Background(), in)
		if err != nil {
			t.Error(err)
		}
		if out.GetContent() != in.GetContent() {
			t.Errorf("Echo(%s) returned %s", in.GetContent(), out.GetContent())
		}
	}
	in := &pb.EchoRequest{
		Response: &pb.EchoRequest_Error{
			Error: &spb.Status{Code: int32(codes.OK)}}}
	_, err := server.Echo(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
}

func TestEcho_error(t *testing.T) {
	table := []codes.Code{codes.Canceled, codes.InvalidArgument}

	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
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

type mockExpandStream struct {
	exp []string
	t   *testing.T
	pb.Echo_ExpandServer
}

func (m *mockExpandStream) Send(resp *pb.EchoResponse) error {
	if resp.GetContent() != m.exp[0] {
		m.t.Errorf("Expand expected to send %s but sent %s", m.exp[0], resp.GetContent())
	}
	m.exp = m.exp[1:]
	return nil
}

func (m *mockExpandStream) verify() {
	if len(m.exp) > 0 {
		m.t.Errorf("Exand did not stream all expected values. %d expected values remaining.", len(m.exp))
	}
}

func TestExpand(t *testing.T) {
	contentTable := []string{"Hello World", "Hola", ""}
	errTable := []*spb.Status{
		{Code: int32(codes.OK)},
		{Code: int32(codes.InvalidArgument)},
		nil,
	}

	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
	for _, c := range contentTable {
		for _, e := range errTable {
			stream := &mockExpandStream{exp: strings.Fields(c), t: t}
			err := server.Expand(&pb.ExpandRequest{Content: c, Error: e}, stream)
			status, _ := status.FromError(err)
			if int32(status.Code()) != e.GetCode() {
				t.Errorf("Expand expected stream to return status with code %d but code %d", status.Code(), e.GetCode())
			}
			stream.verify()
		}
	}
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
	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
	err := server.Expand(&pb.ExpandRequest{Content: "Hello World"}, stream)
	if e != err {
		t.Error("Expand expected to pass through stream errors.")
	}
}

type mockCollectStream struct {
	reqs []*pb.EchoRequest
	exp  *string
	t    *testing.T
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

	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
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
	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
	err := server.Collect(stream)
	if e != err {
		t.Error("Expand expected to pass through stream errors.")
	}
}

type mockChatStream struct {
	reqs []*pb.EchoRequest
	curr *pb.EchoRequest
	t    *testing.T
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

func TestChat(t *testing.T) {
	tests := []struct {
		reqs []string
		err  *spb.Status
	}{
		{[]string{"Hello", "World"}, nil},
		{[]string{"Hello", "World"}, &spb.Status{Code: int32(codes.InvalidArgument)}},
		{[]string{}, &spb.Status{Code: int32(codes.InvalidArgument)}},
	}

	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
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
	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
	err := server.Chat(stream)
	if e != err {
		t.Error("Expand expected to pass through stream errors.")
	}
}

func TestWaitSuccess(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		resp    string
	}{
		{1, int32(1000), "hello"},
		{10, int32(10), "world"},
	}
	for _, test := range tests {
		server := &echoServerImpl{sleepF: mockSleeper(test.seconds, test.nanos, t)}
		in := &pb.WaitRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.WaitRequest_Success{
				Success: &pb.WaitResponse{Content: test.resp},
			},
		}
		out, err := server.Wait(context.Background(), in)
		if err != nil {
			t.Error(err)
		}
		if out.GetContent() != test.resp {
			t.Errorf("Expected Wait test to return %s, but returned %s", out.GetContent(), test.resp)
		}
	}
}

func TestWaitError(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		code    codes.Code
	}{
		{0, int32(0), codes.InvalidArgument},
		{1000, int32(1000), codes.Unavailable},
	}

	for _, test := range tests {
		server := &echoServerImpl{sleepF: mockSleeper(test.seconds, test.nanos, t)}
		in := &pb.WaitRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.WaitRequest_Error{
				Error: status.New(test.code, "").Proto(),
			},
		}
		out, err := server.Wait(context.Background(), in)
		if out != nil {
			t.Errorf("Wait: Expected to error with code %d but returned success", test.code)
		}
		s, _ := status.FromError(err)
		if s.Code() != test.code {
			t.Errorf("Wait: Expected to error with code %d but errored with code %d", test.code, s.Code())
		}
	}
}

func mockSleeper(seconds int64, nanos int32, t *testing.T) func(d time.Duration) {
	return func(d time.Duration) {
		expected := time.Duration(seconds)*time.Second + time.Duration(nanos)*time.Nanosecond
		if d != expected {
			t.Errorf("Expected to sleep %d but was sleep was calledwith %d", expected, d)
		}
	}
}

func zeroNow() time.Time {
	return time.Unix(0, 0)
}

func TestPagination_invalidArgs(t *testing.T) {
	tests := []*pb.PaginationRequest{
		{PageSize: 0},
		{MaxResponse: -1},
		{PageToken: "-1"},
		{PageToken: "BOGUS"},
		{MaxResponse: 1, PageToken: "2"},
	}
	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
	for _, in := range tests {
		_, err := server.Pagination(context.Background(), in)
		s, _ := status.FromError(err)
		if s.Code() != codes.InvalidArgument {
			t.Errorf("Pagination with input '%s', expected to return code '%d' but returned code'%d",
				in.String(), codes.InvalidArgument, s.Code())
		}
	}
}

func TestPagination(t *testing.T) {
	tests := []struct {
		in  *pb.PaginationRequest
		out *pb.PaginationResponse
	}{
		{
			&pb.PaginationRequest{MaxResponse: 3},
			&pb.PaginationResponse{Responses: []int32{0, 1, 2}},
		},
		{
			&pb.PaginationRequest{PageSize: 3, MaxResponse: 2},
			&pb.PaginationResponse{Responses: []int32{0, 1}},
		},
		{
			&pb.PaginationRequest{PageSize: 3, MaxResponse: 10},
			&pb.PaginationResponse{Responses: []int32{0, 1, 2}, NextPageToken: "3"},
		},
		{
			&pb.PaginationRequest{PageSize: 3, MaxResponse: 10, PageToken: "3"},
			&pb.PaginationResponse{Responses: []int32{3, 4, 5}, NextPageToken: "6"},
		},
	}

	server := &echoServerImpl{
		sleepF:      nil,
		session: nil,
	}
	for _, test := range tests {
		out, err := server.Pagination(context.Background(), test.in)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(test.out.GetResponses(), out.GetResponses()) ||
			test.out.GetNextPageToken() != out.GetNextPageToken() {
			t.Errorf("Pagination with input '%s', expected '%s', but returned %s",
				test.in.String(), test.out.String(), out.String())
		}
	}
}
