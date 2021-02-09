package db

import (
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisHandler struct {
	Pool *redis.Pool
}

func NewRedisHandler() (*RedisHandler, error) {
	dsn := os.Getenv("REDIS_DSN")

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", dsn)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
		MaxActive:       25,
		MaxIdle:         25,
		MaxConnLifetime: 300 * time.Second,
		Wait:            true,
	}

	return &RedisHandler{
		Pool: pool,
	}, nil
}
