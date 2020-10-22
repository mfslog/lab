package main

import (
	"context"
	"fmt"
	pb "github.com/JerryZhou343/echo/genproto/github.com/JerryZhou343/lab/istio/echo"
	"github.com/JerryZhou343/echo/genproto/github.com/JerryZhou343/lab/istio/receivetime"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	consulapi "github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	_ "google.golang.org/grpc/xds"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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
	grpc_logrus.ReplaceGrpcLogger(log.NewEntry(log.StandardLogger()))
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
}`, "outbound_.50051_._.receivetime.default.svc.cluster.local")

		cc, err = grpc.Dial("anything",
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(serviceConfigStr))

	if err != nil{
		log.Errorf("Failed to dial receivetime server. config:[%s] err:[%v]",serviceConfigStr,err)
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
	var id string
	id , err = register()
	log.Infof("register id:[%s]",id)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		log.Infof("get a signal %s", s.String())
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
		Kind:              "",
		ID:                id,
		Name:              "echo",
		Tags:              nil,
		Port:              50051,
		Address:           InternalIP(),
		TaggedAddresses:   nil,
		EnableTagOverride: false,
		Meta:              nil,
		Weights:           nil,
		Check:             nil,
		Checks:            nil,
		Proxy:             nil,
		Connect:           nil,
		Namespace:         "",
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

func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if !isUp(inter.Flags) {
			continue
		}
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}


// isUp Interface is up
func isUp(v net.Flags) bool {
	return v&net.FlagUp == net.FlagUp
}
