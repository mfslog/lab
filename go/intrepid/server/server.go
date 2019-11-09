package server

import "github.com/mfslog/lab/go/intrepid/transport"

type Server interface {
	Init()
	Start()
	Stop()
	Handle()
}



type rpcServer struct{
	opts Options
	transport transport.Transport
	dispatch *dispatch
}

func (r *rpcServer)Init(){}

func (r *rpcServer)Start(){}

func (r *rpcServer)Stop(){}

func (r *rpcServer)Handle(){}