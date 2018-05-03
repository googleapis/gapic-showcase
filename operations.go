package feature_testing

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
)

type OperationsServer struct {
	store OperationStore
}

func NewOperationsServer(opStore OperationStore) *OperationsServer {
	return &OperationsServer{store: opStore}
}

func (s *OperationsServer) GetOperation(ctx context.Context, in *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	return s.store.Get(in.GetName())
}

func (s *OperationsServer) CancelOperation(ctx context.Context, in *longrunning.CancelOperationRequest) (*empty.Empty, error) {
	err := s.store.Cancel(in.GetName())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OperationsServer) ListOperations(ctx context.Context, in *longrunning.ListOperationsRequest) (*longrunning.ListOperationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.ListOperations is unimplemented.")
}

func (s *OperationsServer) DeleteOperation(ctx context.Context, in *longrunning.DeleteOperationRequest) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "google.longrunning.DeleteOperation is unimplemented.")
}
