package infrastructure

import (
	"context"
	"time"
)

type HealthRepository struct {
	*RedisHandler
	*MySQLHandler
}

func NewHealthRepository(redisHandler *RedisHandler, mysqlHandler *MySQLHandler) *HealthRepository {
	return &HealthRepository{redisHandler, mysqlHandler}
}

func (r *HealthRepository) PingMySQL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := r.MySQLHandler.conn.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (r *HealthRepository) PingRedis() error {
	conn := r.RedisHandler.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return err
	}

	return nil
}
