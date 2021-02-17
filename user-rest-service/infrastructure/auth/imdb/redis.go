package imdb

import (
	"github.com/garyburd/redigo/redis"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
)

type RedisHandler struct {
	Pool *redis.Pool
}

func NewRedisHandler() (*RedisHandler, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config.Env.Redis.Dsn)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
		MaxActive:       config.Env.Redis.MaxConn,
		MaxIdle:         config.Env.Redis.MaxIdleConn,
		MaxConnLifetime: config.Env.Redis.MaxConnLifetime,
		Wait:            true,
	}

	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return nil, err
	}

	return &RedisHandler{
		Pool: pool,
	}, nil
}
