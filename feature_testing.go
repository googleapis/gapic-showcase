package feature_testing

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"strconv"
	"time"
)

type FeatureTestingServer struct {
	retryStore     map[string][]*statuspb.Status
	operationStore *OperationStore
	nowF           func() time.Time
	sleepF         func(time.Duration)
}

func NewFeatureTestingServer(opStore *OperationStore) *FeatureTestingServer {
	return &FeatureTestingServer{operationStore: opStore}
}

func (s *FeatureTestingServer) WithNowFunc(nowFunc func() time.Time) *FeatureTestingServer {
	return &FeatureTestingServer{
		retryStore:     s.retryStore,
		operationStore: s.operationStore,
		nowF:           nowFunc,
		sleepF:         s.sleepF,
	}
}

func (s *FeatureTestingServer) WithSleepFunc(sleepFunc func(time.Duration)) *FeatureTestingServer {
	return &FeatureTestingServer{
		retryStore:     s.retryStore,
		operationStore: s.operationStore,
		nowF:           s.nowF,
		sleepF:         sleepFunc,
	}
}

func (s *FeatureTestingServer) TimeoutTest(ctx context.Context, in *pb.TimeoutTestRequest) (*pb.TimeoutTestResponse, error) {
	d, _ := ptypes.Duration(in.GetResponseDelay())
	s.sleep(d)
	if in.GetError() != nil {
		return nil, status.ErrorProto(in.GetError())
	}
	return in.GetSuccess(), nil
}

func (s *FeatureTestingServer) SetupRetryTest(ctx context.Context, in *pb.SetupRetryTestRequest) (*pb.RetryTestId, error) {
	if in.GetResponses() == nil {
		return nil, status.Error(codes.InvalidArgument, "A list of responses must be specified.")
	}
	id := fmt.Sprintf("retry-test-%d", s.now().UTC().Unix())
	s.retryStore[id] = in.GetResponses()
	return &pb.RetryTestId{Id: id}, nil
}

func (s *FeatureTestingServer) RetryTest(ctx context.Context, in *pb.RetryTestId) (*empty.Empty, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "An Id must be specified.")
	}
	resps, ok := s.retryStore[in.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "RetryTest with Id: %s was not found.", in.GetId())
	}
	resp, resps := resps[0], resps[1:]
	if len(resps) == 0 {
		delete(s.retryStore, in.GetId())
	} else {
		s.retryStore[in.GetId()] = resps
	}

	if status.FromProto(resp).Code() == codes.OK {
		return &empty.Empty{}, nil
	}
	return nil, status.ErrorProto(resp)
}

func (s *FeatureTestingServer) LongrunningTest(ctx context.Context, in *pb.LongrunningTestRequest) (*lropb.Operation, error) {
	return s.operationStore.RegisterOp(in)
}

func (s *FeatureTestingServer) PaginationTest(ctx context.Context, in *pb.PaginationTestRequest) (*pb.PaginationTestResponse, error) {
	if in.GetPageSize() < 0 || in.GetPageSizeOverride() < 0 {
		return nil, status.Error(codes.InvalidArgument, "The page size provided must not be negative.")
	}

	if in.GetMaxResponse() < 0 {
		return nil, status.Error(codes.InvalidArgument, "The maximum response provided must not be negative.")
	}

	start := int32(0)
	if in.GetPageToken() != "" {
		token, err := strconv.Atoi(in.GetPageToken())
		token32 := int32(token)
		if err != nil || token32 < 0 || token32 > in.GetMaxResponse() {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid page token: %s. Token must be within the range [0, request.MaxResponse]", in.GetPageToken())
		}
		start = token32
	}

	actualSize := in.GetPageSize()
	if in.GetPageSizeOverride() > 0 {
		actualSize = in.GetPageSizeOverride()
	}

	end := start + actualSize
	if actualSize == 0 {
		end = in.GetMaxResponse()
	}
	if end > in.GetMaxResponse() {
		end = in.GetMaxResponse()
	}

	nextToken := ""
	if end < in.GetMaxResponse() {
		nextToken = strconv.Itoa(int(end))
	}

	page := []int32{}
	for i := start; i <= end; i++ {
		page = append(page, i)
	}

	return &pb.PaginationTestResponse{
		Responses:     page,
		NextPageToken: nextToken,
	}, nil
}

func (s *FeatureTestingServer) ParameterFlatteningTest(ctx context.Context, in *pb.ParameterFlatteningTestMessage) (*pb.ParameterFlatteningTestMessage, error) {
	return in, nil
}

func (s *FeatureTestingServer) ResourceNameTest(ctx context.Context, in *pb.ResourceNameTestMessage) (*pb.ResourceNameTestMessage, error) {
	return in, nil
}

func (s *FeatureTestingServer) sleep(d time.Duration) {
	if s.sleepF != nil {
		s.sleepF(d)
	} else {
		time.Sleep(d)
	}
}

func (s FeatureTestingServer) now() time.Time {
	if s.nowF != nil {
		return s.nowF()
	}
	return time.Now()
}
