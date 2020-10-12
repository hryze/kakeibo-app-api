package infrastructure

import (
	"os"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLHandler struct {
	conn *sqlx.DB
}

type RedisHandler struct {
	pool *redis.Pool
}

func NewMySQLHandler() (*MySQLHandler, error) {
	dsn := os.Getenv("MYSQL_DSN")

	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	mySQLHandler := new(MySQLHandler)
	mySQLHandler.conn = conn

	return mySQLHandler, nil
}

func NewRedisHandler() (*RedisHandler, error) {
	dsn := os.Getenv("REDIS_DSN")
	//password := os.Getenv("REDIS_AUTH")

	redisHandler := new(RedisHandler)
	redisHandler.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", dsn)
			if err != nil {
				return nil, err
			}
			//if _, err := conn.Do("AUTH", password); err != nil {
			//	return nil, err
			//}
			return conn, nil
		},
	}

	return redisHandler, nil
}
