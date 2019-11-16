package mgo

import (
	"github.com/mfslog/lab/go/kit/domain/dao"
	"github.com/mfslog/lab/go/kit/infrastructure/account"
	"gopkg.in/mgo.v2"
)

type repository struct{
	collection string
	ses *mgo.Session
}

func NewRepository(ses *mgo.Session) account.Repository {
	return &repository{
		collection: "account",
		ses:        ses,
	}
}

func (r *repository)Get(id int64)(*dao.Account,error){

	return nil,nil
}
