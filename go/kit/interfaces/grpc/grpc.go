package grpc

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/mfslog/lab/go/kit/application"
	"github.com/mfslog/lab/go/kit/idl/account"
	"github.com/mfslog/lab/go/kit/idl/account/acctsvc"
)

type grpcServer struct {
	acc grpctransport.Handler
}

func NewGRPCServer(endpoint application.Set) acctsvc.AccountSvcServer {

	return &grpcServer{
		acc: grpctransport.NewServer(endpoint.AccApp, decodeGetAccountReq,encodeGetAccountRsp),
	}
}

func (g *grpcServer) GetAccount(ctx context.Context, req *account.GetAccountReq) (*account.GetAccountRsp, error) {
	_, rsp, err := g.acc.ServeGRPC(ctx, req)
	return rsp.(*account.GetAccountRsp), err
}

func decodeGetAccountReq(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeGetAccountRsp(_ context.Context, rsp interface{}) (interface{}, error) {
	return rsp, nil
}
