package main

import (
	"context"
	"fmt"
	pb "github.com/JerryZhou343/echo/genproto/github.com/JerryZhou343/lab/istio/echo"
	"github.com/JerryZhou343/echo/genproto/github.com/JerryZhou343/lab/istio/receivetime"
	"google.golang.org/grpc"
	"net"
	"os"
	log "github.com/sirupsen/logrus"
)
func init() {
	var err error
	err = os.MkdirAll("/data/log/echo",os.ModePerm)
	if err != nil{
		panic(err)
	}
	var f *os.File
	f, err = os.OpenFile("/data/log/echo/echo.log",os.O_CREATE|os.O_RDWR|os.O_TRUNC,0666)
	if err != nil{
		panic(err)
	}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}


const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}


const (
	receiveTarget ="xds:///receivetime-svc"
)


// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Debugf("Received: %v", in.GetName())
	var (
		cc *grpc.ClientConn
		err error
		client receivetime.TimeServerClient
		rsp *receivetime.GetCurrentTimeReply
	)
	req := receivetime.GetCurrentTimeRequest{}
	cc,err = grpc.Dial(receiveTarget,grpc.WithInsecure())
	if err != nil{
		log.Errorf("Failed to dial receivetime server. scheme:[%s] err:[%v]",receiveTarget,err)
		return nil,err
	}
	client  = receivetime.NewTimeServerClient(cc)
	rsp ,err = client.GetCurrentTime(context.Background(),&req)
	if err != nil{
		log.Errorf("Failed to get current time. err:[%v]",err)
		return  nil,err
	}
	return &pb.HelloReply{Message: "Hello " + in.GetName() + fmt.Sprintf(" current time:%d",rsp.CurrentAt)}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Info("listener start....")
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


