package main

import (
	"context"
	pb "github.com/JerryZhou343/receivetime/genproto/github.com/JerryZhou343/lab/istio/receivetime"
	consulapi "github.com/armon/consul-api"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "google.golang.org/grpc/xds"
)

func init() {
	var err error
	err = os.MkdirAll("/data/log/receivetime", os.ModePerm)
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.OpenFile("/data/log/receivetime/receivetime.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	grpc_logrus.ReplaceGrpcLogger(log.NewEntry(log.StandardLogger()))
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
	log.Infof("have get request:%v", grpc_ctxtags.Extract(ctx).Values())
	//return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
	at := time.Now().Unix()
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

	return &pb.GetCurrentTimeReply{
		CurrentAt: at,
	}, nil
}



func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listener start...")
	s := grpc.NewServer()
	pb.RegisterTimeServerServer(s, &server{})
	var id string
	id , err = register()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			deregister(id)
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func register()(id string,err error){
	var client *consulapi.Client
	client,err = consulapi.NewClient(&consulapi.Config{
		Address:    "192.168.56.102:8500",
		WaitTime:   3,
	})
	if err != nil{
		log.Fatalf("new consul client failed. err:[%v]",err)
		return
	}
	id = uuid.NewV4().String()
	svcInfo := consulapi.AgentServiceRegistration{
		ID: id,
		Name:  "receivetime",
		Tags:  nil,
		Port:  50051,
		Check: nil,
	}

	err = client.Agent().ServiceRegister(&svcInfo)
	return
}

func deregister(id string)(err error){
	var client *consulapi.Client
	client,err = consulapi.NewClient(&consulapi.Config{
		Address:    "192.168.56.102:8500",
		WaitTime:   3,
	})
	if err != nil{
		log.Fatalf("new consul client failed. err:[%v]",err)
		return
	}

	err = client.Agent().ServiceDeregister(id)
	return
}
