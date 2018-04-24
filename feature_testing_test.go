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
		server := NewFeatureTestingServer(nil).WithSleepFunc(func(d time.Duration) {
			expected := time.Duration(test.seconds)*time.Second + time.Duration(test.nanos)*time.Nanosecond
			if d != expected {
				t.Errorf("Expected to sleep %d but was sleep was calledwith %d", expected, d)
			}
		})
		in := &pb.TimeoutTestRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.TimeoutTestRequest_Success{
				Success: &pb.TimeoutTestResponse{Content: test.resp},
			},
		}
		out, err := server.TimeoutTest(context.TODO(), in)
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
		server := NewFeatureTestingServer(nil).WithSleepFunc(func(d time.Duration) {
			expected := time.Duration(test.seconds)*time.Second + time.Duration(test.nanos)*time.Nanosecond
			if d != expected {
				t.Errorf("Expected to sleep %d but was sleep was called with %d", expected, d)
			}
		})
		in := &pb.TimeoutTestRequest{
			ResponseDelay: &durpb.Duration{
				Seconds: test.seconds,
				Nanos:   test.nanos,
			},
			Response: &pb.TimeoutTestRequest_Error{
				Error: status.New(test.code, "").Proto(),
			},
		}
		out, err := server.TimeoutTest(context.TODO(), in)
		if out != nil {
			t.Errorf("TimeoutTest: Expected to error with code %d but returned success", test.code)
		}
		s, _ := status.FromError(err)
		if s.Code() != test.code {
			t.Errorf("TimeoutTest: Expected to error with code %d but errored with code %d", test.code, s.Code())
		}
	}
}
