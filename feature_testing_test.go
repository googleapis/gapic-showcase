package feature_testing

import (
	durpb "github.com/golang/protobuf/ptypes/duration"
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"testing"
	"time"
)

func mockSleeper(seconds int64, nanos int32, t *testing.T) func(d time.Duration) {
	return func(d time.Duration) {
		expected := time.Duration(seconds)*time.Second + time.Duration(nanos)*time.Nanosecond
		if d != expected {
			t.Errorf("Expected to sleep %d but was sleep was calledwith %d", expected, d)
		}
	}
}

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
