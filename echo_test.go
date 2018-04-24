package feature_testing

import (
	pb "github.com/googleapis/feature-testing-server/genproto/google/example/feature_testing/v1"
	"golang.org/x/net/context"
	"testing"
)

func TestEchoServerEcho(t *testing.T) {
	server := EchoServer{}
	table := []string{"hello", "world"}

	for _, val := range table {
		in := &pb.EchoMessage{Content: val}
		out, err := server.Echo(context.TODO(), in)
		if err != nil {
			t.Error(err)
		}
		if out.GetContent() != in.GetContent() {
			t.Errorf("Echo(%s) returned %s", in.GetContent(), out.GetContent())
		}
	}
}
