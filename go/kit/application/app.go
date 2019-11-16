package application

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/mfslog/lab/go/kit/domain/service"
)

type Set struct{
	AccApp endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger)Set{
	var accApp endpoint.Endpoint
	{

	}
	return
}

type AccGetRequest struct{
	id int64
}




func MakeAccApp(s service.Service)endpoint.Endpoint{
	return func(ctx context.Context, request interface{})(response interface{}, err error){
		req := request.(AccGetRequest)
		do, err := s.Get(req.id)
		return
	}
}