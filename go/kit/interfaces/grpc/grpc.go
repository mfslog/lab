package grpc

import (
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct{
	acc grpctransport.Handler
}

func NewGRPCServer(endpoint endpoint.Endpoint)