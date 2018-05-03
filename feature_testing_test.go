package feature_testing

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	durpb "github.com/golang/protobuf/ptypes/duration"
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"reflect"
	"testing"
	"time"
)

func TestTimeoutTestSuccess(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		resp    string
	}{
		{1, int32(1000), "hello"},
		{10, int32(10), "world"},
	}
	for _, test := range tests {
		server := NewFeatureTestingServer(nil).WithSleepFunc(mockSleeper(test.seconds, test.nanos, t))
		in := &pb.TimeoutTestRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.TimeoutTestRequest_Success{
				Success: &pb.TimeoutTestResponse{Content: test.resp},
			},
		}
		out, err := server.TimeoutTest(context.Background(), in)
		if err != nil {
			t.Error(err)
		}
		if out.GetContent() != test.resp {
			t.Errorf("Expected Timeout test to return %s, but returned %s", out.GetContent(), test.resp)
		}
	}
}

func TestTimeoutTestError(t *testing.T) {
	tests := []struct {
		seconds int64
		nanos   int32
		code    codes.Code
	}{
		{0, int32(0), codes.InvalidArgument},
		{1000, int32(1000), codes.Unavailable},
	}

	for _, test := range tests {
		server := NewFeatureTestingServer(nil).WithSleepFunc(mockSleeper(test.seconds, test.nanos, t))
		in := &pb.TimeoutTestRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.TimeoutTestRequest_Error{
				Error: status.New(test.code, "").Proto(),
			},
		}
		out, err := server.TimeoutTest(context.Background(), in)
		if out != nil {
			t.Errorf("TimeoutTest: Expected to error with code %d but returned success", test.code)
		}
		s, _ := status.FromError(err)
		if s.Code() != test.code {
			t.Errorf("TimeoutTest: Expected to error with code %d but errored with code %d", test.code, s.Code())
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

func TestSetupRetryTest(t *testing.T) {
	tests := []struct {
		in  []codes.Code
		out codes.Code
	}{
		{[]codes.Code{codes.OK}, codes.OK},
		{nil, codes.InvalidArgument},
		{[]codes.Code{}, codes.InvalidArgument},
	}

	for _, test := range tests {
		server := (&FeatureTestingServer{}).WithNowFunc(zeroNow)
		var resps []*spb.Status
		if test.in != nil {
			resps = []*spb.Status{}
			for _, code := range test.in {
				resps = append(resps, &spb.Status{Code: int32(code)})
			}
		}
		in := &pb.SetupRetryTestRequest{Responses: resps}
		out, err := server.SetupRetryTest(context.Background(), in)
		if out != nil && out.GetId() != "retry-test-0" {
			t.Errorf("Expected SetupRetryTest to return the ID 'retry-test-0', but returned '%s'", out.GetId())
		}
		s, _ := status.FromError(err)
		if s.Code() != test.out {
			t.Errorf("Expected SetupRetryTest to return with code '%d', but returned code '%d'", test.out, s.Code())
		}
	}
}

func TestRetryTest(t *testing.T) {
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
			[]codes.Code{codes.Unavailable, codes.OK},
			[]codes.Code{codes.Unavailable, codes.OK},
		},
	}

	for _, test := range tests {
		server := (&FeatureTestingServer{}).WithNowFunc(zeroNow)
		resps := []*spb.Status{}
		for _, code := range test.in {
			resps = append(resps, &spb.Status{Code: int32(code)})
		}
		in := &pb.SetupRetryTestRequest{Responses: resps}
		out, err := server.SetupRetryTest(context.Background(), in)
		if err != nil {
			t.Errorf("SetupRetryTest failed to setup.")
		}
		for _, expected := range test.out {
			_, err := server.RetryTest(context.Background(), out)
			s, _ := status.FromError(err)
			if expected != s.Code() {
				t.Errorf("RetryTest expected to return code '%d', but returned '%d'", expected, s.Code())
			}
		}
	}
}

func TestPaginationTest_invalidArgs(t *testing.T) {
	tests := []*pb.PaginationTestRequest{
		&pb.PaginationTestRequest{PageSize: 0},
		&pb.PaginationTestRequest{PageSizeOverride: -1},
		&pb.PaginationTestRequest{MaxResponse: -1},
		&pb.PaginationTestRequest{PageToken: "-1"},
		&pb.PaginationTestRequest{PageToken: "BOGUS"},
		&pb.PaginationTestRequest{MaxResponse: 1, PageToken: "2"},
	}
	server := &FeatureTestingServer{}
	for _, in := range tests {
		_, err := server.PaginationTest(context.Background(), in)
		s, _ := status.FromError(err)
		if s.Code() != codes.InvalidArgument {
			t.Errorf("PaginationTest with input '%s', expected to return code '%d' but returned code'%d",
				in.String(), codes.InvalidArgument, s.Code())
		}
	}
}

func TestPaginationTest(t *testing.T) {
	tests := []struct {
		in  *pb.PaginationTestRequest
		out *pb.PaginationTestResponse
	}{
		{
			&pb.PaginationTestRequest{MaxResponse: 3},
			&pb.PaginationTestResponse{Responses: []int32{0, 1, 2}},
		},
		{
			&pb.PaginationTestRequest{PageSize: 3, MaxResponse: 2},
			&pb.PaginationTestResponse{Responses: []int32{0, 1}},
		},
		{
			&pb.PaginationTestRequest{PageSize: 3, MaxResponse: 10},
			&pb.PaginationTestResponse{Responses: []int32{0, 1, 2}, NextPageToken: "3"},
		},
		{
			&pb.PaginationTestRequest{PageSize: 3, MaxResponse: 10, PageToken: "3"},
			&pb.PaginationTestResponse{Responses: []int32{3, 4, 5}, NextPageToken: "6"},
		},
		{
			&pb.PaginationTestRequest{PageSize: 3, MaxResponse: 10, PageSizeOverride: 5},
			&pb.PaginationTestResponse{Responses: []int32{0, 1, 2, 3, 4}, NextPageToken: "5"},
		},
	}

	server := &FeatureTestingServer{}
	for _, test := range tests {
		out, err := server.PaginationTest(context.Background(), test.in)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(test.out.GetResponses(), out.GetResponses()) ||
			test.out.GetNextPageToken() != out.GetNextPageToken() {
			t.Errorf("PaginationTest with input '%s', expected '%s', but returned %s",
				test.in.String(), test.out.String(), out.String())
		}
	}
}

type mockOperationStore struct {
	t   *testing.T
	req *pb.LongrunningTestRequest
	OperationStore
}

func (m *mockOperationStore) RegisterOp(in *pb.LongrunningTestRequest) (*lropb.Operation, error) {
	m.req = in
	return &lropb.Operation{}, nil
}

func TestLongrunningTest(t *testing.T) {
	mockStore := &mockOperationStore{t: t}
	server := NewFeatureTestingServer(mockStore)
	stamp, _ := ptypes.TimestampProto(time.Unix(0, 0))
	in := &pb.LongrunningTestRequest{
		CompletionTime: stamp,
		Response: &pb.LongrunningTestRequest_Success{
			Success: &pb.LongrunningTestResponse{Content: "Hello World"},
		},
	}
	_, err := server.LongrunningTest(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
	if mockStore.req != in {
		t.Error("LongrunningTest expected to register the an operation with the store.")
	}
	if !proto.Equal(mockStore.req, in) {
		t.Error("LongrunningTest unexpectedly altered the input registering with the operation store.")
	}

}

func TestParameterFlatteningTest(t *testing.T) {
	in := &pb.ParameterFlatteningTestMessage{
		Content:         "hello world",
		RepeatedContent: []string{"hello", "world"},
		Nested:          &pb.ParameterFlatteningTestMessage{Content: "hola"},
	}
	server := &FeatureTestingServer{}
	out, err := server.ParameterFlatteningTest(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
	if in != out {
		t.Errorf("ParameterFlatteningTest expected to pass back the input.")
	}
	if !proto.Equal(in, out) {
		t.Errorf("ParameterFlatteningTest unexpectedly altered the input proto.")
	}
}

func TestResourceNameTest(t *testing.T) {
	in := &pb.ResourceNameTestMessage{
		SingleTemplate:    "/hello/world",
		MultipleTemplates: "/hola/world",
	}
	server := &FeatureTestingServer{}
	out, err := server.ResourceNameTest(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
	if in != out {
		t.Errorf("ResourceNameTest expected to pass back the input.")
	}
	if !proto.Equal(in, out) {
		t.Errorf("ResourceNameTest unexpectedly altered the input proto.")
	}
}
