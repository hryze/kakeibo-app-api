package config

import (
	"os"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLHandler struct {
	Conn *sqlx.DB
}

func NewMySQLHandler() (*MySQLHandler, error) {
	dsn := os.Getenv("MYSQL_DSN")

	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	mySQLHandler := new(MySQLHandler)
	mySQLHandler.Conn = conn

	return mySQLHandler, nil
}

type RedisHandler struct {
	Pool *redis.Pool
}

func NewRedisHandler() (*RedisHandler, error) {
	dsn := os.Getenv("REDIS_DSN")

	redisHandler := new(RedisHandler)
	redisHandler.Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", dsn)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}

	return redisHandler, nil
}
