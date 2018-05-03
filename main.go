package feature_testing

import (
	"log"
	"net"

	featurepb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	featurepb.RegisterEchoServiceServer(s, &EchoServer{})
	opStore := &OperationStoreImpl{}
	featurepb.RegisterFeatureTestingServiceServer(s, NewFeatureTestingServer(opStore))
	lropb.RegisterOperationsServer(s, NewOperationsServer(opStore))
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
