package feature_testing

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	featurepb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"github.com/grpc/grpc-go/status"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"time"
)

type operationInfo struct {
	name     string
	start    time.Time
	end      time.Time
	canceled bool
	resp     *featurepb.LongrunningTestResponse
	err      *statuspb.Status
}

type OperationStore struct {
	nowF  func() time.Time
	store map[string]*operationInfo
}

func (s *OperationStore) WithNowF(nowFunc func() time.Time) *OperationStore {
	return &OperationStore{
		nowF:  nowFunc,
		store: s.store,
	}
}

func (s *OperationStore) RegisterOp(op *featurepb.LongrunningTestRequest) (*lropb.Operation, error) {
	end, err := ptypes.Timestamp(op.CompletionTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Given operation completion time is invalid.")
	}
	now := s.nowF()
	name := fmt.Sprintf("lro-test-op-%d", now.UTC().Unix())
	s.store[name] = &operationInfo{
		name:     name,
		start:    now,
		end:      end,
		canceled: false,
		resp:     op.GetSuccess(),
		err:      op.GetError(),
	}
	return s.Get(name)
}

func (s *OperationStore) Get(name string) (*lropb.Operation, error) {
	op, ok := s.store[name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Operation '%s' not found.", name)
	}
	ret := &lropb.Operation{
		Name: op.name,
	}

	now := s.now()

	if op.canceled {
		ret.Result = &lropb.Operation_Error{
			Error: status.Newf(
				codes.Canceled,
				"Operation '%s' has been canceled.", name).Proto(),
		}
	} else if now.After(op.end) {
		if op.err != nil {
			ret.Result = &lropb.Operation_Error{Error: op.err}
		} else {
			resp, err := ptypes.MarshalAny(op.resp)
			if err != nil {
				return nil, err
			}
			ret.Result = &lropb.Operation_Response{Response: resp}
			ret.Done = true
			meta, err := ptypes.MarshalAny(&featurepb.LongrunningTestMetadata{TimeRemaining: ptypes.DurationProto(0)})
			if err != nil {
				return nil, err
			}
			ret.Metadata = meta
		}
	} else {
		meta, err := ptypes.MarshalAny(
			&featurepb.LongrunningTestMetadata{
				TimeRemaining: ptypes.DurationProto(now.Sub(op.end))})
		if err != nil {
			return nil, err
		}
		ret.Metadata = meta
	}
	return ret, nil
}

func (s *OperationStore) Cancel(name string) error {
	op, ok := s.store[name]
	if !ok {
		return status.Errorf(codes.NotFound, "Operation '%s' not found.", name)
	}
	op.canceled = true
	s.store[name] = op
	return nil
}

func (s *OperationStore) now() time.Time {
	if s.nowF != nil {
		return s.nowF()
	}
	return time.Now()
}
