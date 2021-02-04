package infrastructure

import (
	"context"
	"time"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
)

type HealthRepository struct {
	*config.RedisHandler
	*config.MySQLHandler
}

func NewHealthRepository(redisHandler *config.RedisHandler, mysqlHandler *config.MySQLHandler) *HealthRepository {
	return &HealthRepository{redisHandler, mysqlHandler}
}

func (r *HealthRepository) PingMySQL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := r.MySQLHandler.Conn.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (r *HealthRepository) PingRedis() error {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return err
	}

	return nil
}
