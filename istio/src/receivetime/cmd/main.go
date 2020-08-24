package main

import (
	"context"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"net"
	pb "github.com/JerryZhou343/receivetime/genproto/github.com/JerryZhou343/lab/istio/receivetime"
	"google.golang.org/grpc"
	"os"
	"time"
)


func init() {
	var err error
	err = os.MkdirAll("/data/log/receivetime",os.ModePerm)
	if err != nil{
		panic(err)
	}
	var f *os.File
	f, err = os.OpenFile("/data/log/receivetime/receivetime.log",os.O_CREATE|os.O_RDWR|os.O_TRUNC,0666)
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