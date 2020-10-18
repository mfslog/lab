package main

import (
	"context"
	"fmt"
	pb "github.com/JerryZhou343/echo/genproto/github.com/JerryZhou343/lab/istio/echo"
	"github.com/JerryZhou343/echo/genproto/github.com/JerryZhou343/lab/istio/receivetime"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"os"
	//_ "google.golang.org/grpc/xds"
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
	receiveTarget ="eds_experimental:///receivetime-server"+port
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
	meta,ok := metadata.FromIncomingContext(ctx)
	if ok{
		traceID := meta.Get("x-b3-traceid")
		requestID := meta.Get("x-request-id")
		spanID := meta.Get("x-b3-spanid")
		parentspanid := meta.Get("x-b3-parentspanid")
		sampled := meta.Get("x-b3-sampled")
		flags := meta.Get("x-b3-flags")
		spanCtx := meta.Get("x-ot-span-context")
		log.Infof("traceID:%v, requestID:%v, spanID:%v,parentSpanID:%v, sampled:%v, flags:%v,spanCtx:%v",
			traceID,requestID,spanID, parentspanid,sampled,flags,spanCtx)
	}else{
		log.Infof("not ok")
	}

	req := receivetime.GetCurrentTimeRequest{}
	//cc,err = grpc.Dial(receiveTarget,grpc.WithInsecure())

		serviceConfigStr := fmt.Sprintf(`
{
  "loadBalancingConfig":[
    {"eds_experimental":{ "EDSServiceName": "%s" }}
  ]
}`, "receivetime-server"+port)

		cc, err = grpc.Dial("anything",
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(serviceConfigStr))

	if err != nil{
		log.Errorf("Failed to dial receivetime server. scheme:[%s] err:[%v]",receiveTarget,err)
		return nil,err
	}
	client  = receivetime.NewTimeServerClient(cc)
	rsp ,err = client.GetCurrentTime(ctx,&req)
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
	s := grpc.NewServer( grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_logrus.StreamServerInterceptor(log.NewEntry(log.StandardLogger())),
		grpc_recovery.StreamServerInterceptor(),
	)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			grpc_recovery.UnaryServerInterceptor(),
		)),
		)
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


