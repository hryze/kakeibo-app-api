package infrastructure

import (
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/hryze/kakeibo-app-api/account-rest-service/config"
)

type MySQLHandler struct {
	conn *sqlx.DB
}

type RedisHandler struct {
	pool *redis.Pool
}

func NewMySQLHandler() (*MySQLHandler, error) {
	conn, err := sqlx.Open("mysql", config.Env.MySQL.Dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(config.Env.MySQL.MaxConn)
	conn.SetMaxIdleConns(config.Env.MySQL.MaxIdleConn)
	conn.SetConnMaxLifetime(config.Env.MySQL.MaxConnLifetime)

	return &MySQLHandler{
		conn: conn,
	}, nil
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
		pool: pool,
	}, nil
}
