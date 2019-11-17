package application

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/mfslog/lab/go/kit/domain/dao"
	"github.com/mfslog/lab/go/kit/domain/service"
	"github.com/mfslog/lab/go/kit/idl/account"
)

type Set struct {
	AccApp endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger) Set {
	var accApp endpoint.Endpoint
	{
		accApp = MakeAccApp(svc)
	}
	return Set{
		AccApp: accApp,
	}
}

func MakeAccApp(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(account.GetAccountReq)
		do, err := s.Get(req.ID)
		return DO2DTO(do), err
	}
}

func DO2DTO(src *dao.Account) *account.Account {
	if src == nil {
		return nil
	}

	return &account.Account{
		Id:       src.ID,
		Name:     src.Name,
		Email:    src.Email,
		Password: src.Password,
	}
}
