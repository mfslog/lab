package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	pb "github.com/JerryZhou343/receivetime/genproto/github.com/JerryZhou343/lab/istio/receivetime"
	"google.golang.org/grpc"
	"time"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedTimeServerServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.GetCurrentTimeRequest) (*pb.GetCurrentTimeReply, error) {
	logrus.Infof("have get request")
	//return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
	at := time.Now().Unix()
	return &pb.GetCurrentTimeReply{
		CurrentAt:            at,
	},nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}