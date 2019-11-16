package service

import (
	"github.com/go-kit/kit/log"
	"github.com/mfslog/lab/go/kit/domain/dao"
	"github.com/mfslog/lab/go/kit/infrastructure/account"
)

type Service interface {
	Get(id int64)(*dao.Account,error)
}

type service struct{
	logger log.Logger
	repo  account.Repository
}

func NewService(logger log.Logger, repo account.Repository)Service{
	return  &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service)Get(id int64)(*dao.Account, error){
	return  s.repo.Get(id)
}
