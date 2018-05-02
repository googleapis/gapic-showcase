package feature_testing

import (
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"github.com/grpc/grpc-go/status"
	"golang.org/x/net/context"
	"io"
	"strings"
)

type EchoServer struct{}

func (s *EchoServer) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	err := status.ErrorProto(in.GetError())
	if err != nil {
		return nil, err
	}
	return &pb.EchoResponse{Content: in.GetContent()}, nil
}

func (s *EchoServer) Expand(in *pb.ExpandRequest, stream pb.EchoService_ExpandServer) error {
	for _, word := range strings.Fields(in.GetContent()) {
		err := stream.Send(&pb.EchoResponse{Content: word})
		if err != nil {
			return err
		}
	}
	if in.GetError() != nil {
		return status.ErrorProto(in.GetError())
	}
	return nil
}

func (s *EchoServer) Collect(stream pb.EchoService_CollectServer) error {
	var resp []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.EchoResponse{Content: strings.Join(resp, " ")})
		}
		if err != nil {
			return err
		}
		s := status.ErrorProto(req.GetError())
		if s != nil {
			return s
		}
		if req.GetContent() != "" {
			resp = append(resp, req.GetContent())
		}
	}
}

func (s *EchoServer) Chat(stream pb.EchoService_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s := status.ErrorProto(req.GetError())
		if s != nil {
			return s
		}
		stream.Send(&pb.EchoResponse{Content: req.GetContent()})
	}
}
