package main

import (
	"context"
	"github.com/micro/go-micro"
	"github.com/mfslog/lab/go/micro/server/hello"
)

func main(){
	service := micro.NewService(micro.Name("go.micro.srv.hello"))

	service.Init()

	hello.RegisterHelloHandler(service.Server(), )
}


type Handler struct{}


func (h *Handler)SayHello(ctx context.Context, req *hello.SayHelloReq,rsp *hello.SayHelloRsp)error{

}