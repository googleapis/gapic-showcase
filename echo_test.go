package feature_testing

import (
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"io"
	"strings"
	"testing"
)

func TestEcho_success(t *testing.T) {
	table := []string{"hello world", ""}

	server := EchoServer{}
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

	server := EchoServer{}
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
	pb.EchoService_ExpandServer
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
		&spb.Status{Code: int32(codes.OK)},
		&spb.Status{Code: int32(codes.InvalidArgument)},
		nil,
	}

	server := EchoServer{}
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

type mockCollectStream struct {
	reqs []*pb.EchoRequest
	exp  *string
	t    *testing.T
	pb.EchoService_CollectServer
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
		{[]string{"hola"}, strPtr("hola"), nil},
		{[]string{}, strPtr(""), nil},
	}

	server := &EchoServer{}
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

type mockChatStream struct {
	reqs []*pb.EchoRequest
	curr *pb.EchoRequest
	t    *testing.T
	pb.EchoService_ChatServer
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

func TestChatSuccess(t *testing.T) {
	tests := []struct {
		reqs []string
		err  *spb.Status
	}{
		{[]string{"Hello", "World"}, nil},
		{[]string{"Hello", "World"}, &spb.Status{Code: int32(codes.InvalidArgument)}},
		{[]string{}, &spb.Status{Code: int32(codes.InvalidArgument)}},
	}

	server := &EchoServer{}
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
