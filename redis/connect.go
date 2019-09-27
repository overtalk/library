package redis

import (
	"errors"
	"github.com/go-redis/redis"
)

type Config struct {
	Addr     []string `json:"redis_addrs" env:"REDIS_ADDRS" envDefault:"127.0.0.1:6379"`
	Pwd      string   `json:"redis_pwd" env:"REDIS_PWD"`
	PoolSize int      `json:"redis_pool_size" env:"REDIS_POOL_SIZE" envDefault:"1000"`
	DB       int      `json:"redis_db" env:"REDIS_DB"` // 单机模式下选择使用哪个DB，集群模式下无效
}

func (c Config) Connect() (redis.Cmdable, error) {
	addrNum := len(c.Addr)
	if addrNum == 0 {
		return nil, errors.New("redis addr is absent")
	}

	if addrNum > 1 {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    c.Addr,
			Password: c.Pwd,
			PoolSize: c.PoolSize,
		}), nil
	}
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr[0],
		Password: c.Pwd,
		PoolSize: c.PoolSize,
		DB:       c.DB,
	}), nil
}
