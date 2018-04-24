package feature_testing

import (
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"golang.org/x/net/context"
	"io"
	"strings"
)

type EchoServer struct{}

func (s *EchoServer) Echo(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	return in, nil
}

func (s *EchoServer) Expand(in *pb.EchoMessage, stream pb.EchoService_ExpandServer) error {
	for _, word := range strings.Fields(in.Content) {
		err := stream.Send(&pb.EchoMessage{Content: word})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *EchoServer) Collect(stream pb.EchoService_CollectServer) error {
	var resp []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&pb.EchoMessage{Content: strings.Join(resp, " ")})
		}
		if err != nil {
			return err
		}
		resp = append(resp, req.Content)
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
		stream.Send(req)
	}
}
