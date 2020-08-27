package main

import (
	"context"
	envoy_tracer "github.com/JerryZhou343/receivetime/envoy-tracer"
	pb "github.com/JerryZhou343/receivetime/genproto/github.com/JerryZhou343/lab/istio/receivetime"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
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
func (s *server) GetCurrentTime(ctx context.Context, in *pb.GetCurrentTimeRequest) (*pb.GetCurrentTimeReply, error) {
	log.Infof("have get request:%v",grpc_ctxtags.Extract(ctx).Values())
	//return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
	at := time.Now().Unix()
	return &pb.GetCurrentTimeReply{
		CurrentAt:            at,
	},nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listener start...")
	s := grpc.NewServer( grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(envoy_tracer.EnvoyTracer{})),
		//grpc_prometheus.StreamServerInterceptor,
		//grpc_zap.StreamServerInterceptor(zapLogger),
		grpc_logrus.StreamServerInterceptor(log.NewEntry(log.StandardLogger())),
		//grpc_auth.StreamServerInterceptor(myAuthFunction),
		grpc_recovery.StreamServerInterceptor(),
	)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(envoy_tracer.EnvoyTracer{})),
			//grpc_prometheus.UnaryServerInterceptor,
			//grpc_zap.UnaryServerInterceptor(zapLogger),
			grpc_logrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			//grpc_auth.UnaryServerInterceptor(myAuthFunction),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	pb.RegisterTimeServerServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}