package server

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	durpb "github.com/golang/protobuf/ptypes/duration"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/grpc/grpc-go/status"

	"golang.org/x/net/context"

	lropb "google.golang.org/genproto/googleapis/longrunning"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

func TestEcho_success(t *testing.T) {
	table := []string{"hello world", ""}

	server := NewShowcaseServer(nil)
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

	server := NewShowcaseServer(nil)
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
	pb.Showcase_ExpandServer
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

	server := NewShowcaseServer(nil)
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
	pb.Showcase_ExpandServer
}

func (s *errorExpandStream) Send(resp *pb.EchoResponse) error {
	return s.err
}

func TestExpand_streamErr(t *testing.T) {
	e := errors.New("Test Error")
	stream := &errorExpandStream{err: e}
	server := NewShowcaseServer(nil)
	err := server.Expand(&pb.ExpandRequest{Content: "Hello World"}, stream)
	if e != err {
		t.Error("Expand expected to pass through stream errors.")
	}
}

type mockCollectStream struct {
	reqs []*pb.EchoRequest
	exp  *string
	t    *testing.T
	pb.Showcase_CollectServer
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

	server := NewShowcaseServer(nil)
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
	pb.Showcase_CollectServer
}

func (s *errorCollectStream) Recv() (*pb.EchoRequest, error) {
	return nil, s.err
}

func TestCollect_streamErr(t *testing.T) {
	e := errors.New("Test Error")
	stream := &errorCollectStream{err: e}
	server := NewShowcaseServer(nil)
	err := server.Collect(stream)
	if e != err {
		t.Error("Expand expected to pass through stream errors.")
	}
}

type mockChatStream struct {
	reqs []*pb.EchoRequest
	curr *pb.EchoRequest
	t    *testing.T
	pb.Showcase_ChatServer
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

	server := NewShowcaseServer(nil)
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
	pb.Showcase_ChatServer
}

func (s *errorChatStream) Recv() (*pb.EchoRequest, error) {
	return nil, s.err
}

func TestChat_streamErr(t *testing.T) {
	e := errors.New("Test Error")
	stream := &errorChatStream{err: e}
	server := NewShowcaseServer(nil)
	err := server.Chat(stream)
	if e != err {
		t.Error("Expand expected to pass through stream errors.")
	}
}

func TestTimeoutSuccess(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		resp    string
	}{
		{1, int32(1000), "hello"},
		{10, int32(10), "world"},
	}
	for _, test := range tests {
		server := &showcaseServerImpl{sleepF: mockSleeper(test.seconds, test.nanos, t)}
		in := &pb.TimeoutRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.TimeoutRequest_Success{
				Success: &pb.TimeoutResponse{Content: test.resp},
			},
		}
		out, err := server.Timeout(context.Background(), in)
		if err != nil {
			t.Error(err)
		}
		if out.GetContent() != test.resp {
			t.Errorf("Expected Timeout test to return %s, but returned %s", out.GetContent(), test.resp)
		}
	}
}

func TestTimeoutError(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		code    codes.Code
	}{
		{0, int32(0), codes.InvalidArgument},
		{1000, int32(1000), codes.Unavailable},
	}

	for _, test := range tests {
		server := &showcaseServerImpl{sleepF: mockSleeper(test.seconds, test.nanos, t)}
		in := &pb.TimeoutRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.TimeoutRequest_Error{
				Error: status.New(test.code, "").Proto(),
			},
		}
		out, err := server.Timeout(context.Background(), in)
		if out != nil {
			t.Errorf("Timeout: Expected to error with code %d but returned success", test.code)
		}
		s, _ := status.FromError(err)
		if s.Code() != test.code {
			t.Errorf("Timeout: Expected to error with code %d but errored with code %d", test.code, s.Code())
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

func TestSetupRetry(t *testing.T) {
	tests := []struct {
		in  []codes.Code
		out codes.Code
	}{
		{[]codes.Code{codes.OK}, codes.OK},
		{nil, codes.InvalidArgument},
		{[]codes.Code{}, codes.InvalidArgument},
	}

	for _, test := range tests {
		server := (&showcaseServerImpl{nowF: zeroNow})
		var resps []*spb.Status
		if test.in != nil {
			resps = []*spb.Status{}
			for _, code := range test.in {
				resps = append(resps, &spb.Status{Code: int32(code)})
			}
		}
		in := &pb.SetupRetryRequest{Responses: resps}
		out, err := server.SetupRetry(context.Background(), in)
		if out != nil && out.GetId() != "retry-test-0" {
			t.Errorf("Expected SetupRetry to return the ID 'retry-test-0', but returned '%s'", out.GetId())
		}
		s, _ := status.FromError(err)
		if s.Code() != test.out {
			t.Errorf("Expected SetupRetry to return with code '%d', but returned code '%d'", test.out, s.Code())
		}
	}
}

func TestRetry(t *testing.T) {
	tests := []struct {
		in  []codes.Code
		out []codes.Code
	}{
		{
			[]codes.Code{codes.OK},
			[]codes.Code{codes.OK},
		},
		{
			[]codes.Code{codes.OK, codes.Unavailable},
			[]codes.Code{codes.OK, codes.NotFound},
		},
		{
			[]codes.Code{codes.Unavailable, codes.Unavailable},
			[]codes.Code{codes.Unavailable, codes.Unavailable},
		},
		{
			[]codes.Code{codes.Unavailable, codes.OK},
			[]codes.Code{codes.Unavailable, codes.OK},
		},
	}

	for _, test := range tests {
		server := &showcaseServerImpl{nowF: zeroNow}
		resps := []*spb.Status{}
		for _, code := range test.in {
			resps = append(resps, &spb.Status{Code: int32(code)})
		}
		in := &pb.SetupRetryRequest{Responses: resps}
		out, err := server.SetupRetry(context.Background(), in)
		if err != nil {
			t.Errorf("SetupRetry failed to setup.")
		}
		for _, expected := range test.out {
			_, err := server.Retry(context.Background(), out)
			s, _ := status.FromError(err)
			if expected != s.Code() {
				t.Errorf("Retry expected to return code '%d', but returned '%d'", expected, s.Code())
			}
		}
	}
}

func TestRetry_noId(t *testing.T) {
	server := NewShowcaseServer(nil)
	_, err := server.Retry(context.Background(), &pb.RetryId{})
	s, _ := status.FromError(err)
	if codes.InvalidArgument != s.Code() {
		t.Error("Retry expected to return InvalidArgument if no id was specified.")
	}
}

func TestPagination_invalidArgs(t *testing.T) {
	tests := []*pb.PaginationRequest{
		{PageSize: 0},
		{PageSizeOverride: -1},
		{MaxResponse: -1},
		{PageToken: "-1"},
		{PageToken: "BOGUS"},
		{MaxResponse: 1, PageToken: "2"},
	}
	server := NewShowcaseServer(nil)
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
		{
			&pb.PaginationRequest{PageSize: 3, MaxResponse: 10, PageSizeOverride: 5},
			&pb.PaginationResponse{Responses: []int32{0, 1, 2, 3, 4}, NextPageToken: "5"},
		},
	}

	server := NewShowcaseServer(nil)
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

type mockOperationStore struct {
	t   *testing.T
	req *pb.LongrunningRequest
	OperationStore
}

func (m *mockOperationStore) RegisterOp(in *pb.LongrunningRequest) (*lropb.Operation, error) {
	m.req = in
	return &lropb.Operation{}, nil
}

func TestLongrunning(t *testing.T) {
	mockStore := &mockOperationStore{t: t}
	server := NewShowcaseServer(mockStore)
	stamp, _ := ptypes.TimestampProto(time.Unix(0, 0))
	in := &pb.LongrunningRequest{
		CompletionTime: stamp,
		Response: &pb.LongrunningRequest_Success{
			Success: &pb.LongrunningResponse{Content: "Hello World"},
		},
	}
	_, err := server.Longrunning(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
	if mockStore.req != in {
		t.Error("Longrunning expected to register the an operation with the store.")
	}
	if !proto.Equal(mockStore.req, in) {
		t.Error("Longrunning unexpectedly altered the input registering with the operation store.")
	}

}

func TestParameterFlattening(t *testing.T) {
	in := &pb.ParameterFlatteningMessage{
		Content:         "hello world",
		RepeatedContent: []string{"hello", "world"},
		Nested:          &pb.ParameterFlatteningMessage{Content: "hola"},
	}
	server := NewShowcaseServer(nil)
	out, err := server.ParameterFlattening(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
	if in != out {
		t.Errorf("ParameterFlattening expected to pass back the input.")
	}
	if !proto.Equal(in, out) {
		t.Errorf("ParameterFlattening unexpectedly altered the input proto.")
	}
}

func TestResourceName(t *testing.T) {
	in := &pb.ResourceNameMessage{
		SingleTemplate:    "/hello/world",
		MultipleTemplates: "/hola/world",
	}
	server := NewShowcaseServer(nil)
	out, err := server.ResourceName(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
	if in != out {
		t.Errorf("ResourceName expected to pass back the input.")
	}
	if !proto.Equal(in, out) {
		t.Errorf("ResourceName unexpectedly altered the input proto.")
	}
}
