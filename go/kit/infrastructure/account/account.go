package account

import "github.com/mfslog/lab/go/kit/domain/dao"

type Repository interface {
	Get(id int64) (*dao.Account, error)
}
