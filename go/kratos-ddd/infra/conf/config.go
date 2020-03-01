package conf

import (
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
)

// Config .
type Config struct {
	Redis  *redis.Config
	MySQL  *sql.Config
}

func NewConf()(conf *Config,err error){
	conf = &Config{
		Redis:  &redis.Config{},
		MySQL:  &sql.Config{},
	}
	err = paladin.Get("db.toml").UnmarshalTOML(conf.MySQL)
	return
}
