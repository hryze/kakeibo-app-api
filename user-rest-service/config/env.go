package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var Env ENV

func init() {
	env := os.Getenv("GO_ENV")

	if err := envconfig.Process(env, &Env); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type ENV struct {
	Server
	Cors
	Cookie
	MySQL
	Redis
	AccountApi
}

type Server struct {
	Port int `envconfig:"SERVER_PORT" required:"true"`
}

type Cors struct {
	AllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" required:"true"`
}

type Cookie struct {
	Domain string `envconfig:"COOKIE_DOMAIN" required:"true"`
}

type MySQL struct {
	Dsn             string        `envconfig:"MYSQL_DSN"               required:"true"`
	MaxConn         int           `envconfig:"MYSQL_MAX_CONN"          default:"25"`
	MaxIdleConn     int           `envconfig:"MYSQL_MAX_IDLE"          default:"25"`
	MaxConnLifetime time.Duration `envconfig:"MYSQL_MAX_CONN_LIFETIME" default:"300s"`
}

type Redis struct {
	Dsn             string        `envconfig:"REDIS_DSN"               required:"true"`
	MaxConn         int           `envconfig:"REDIS_MAX_CONN"          default:"25"`
	MaxIdleConn     int           `envconfig:"REDIS_MAX_IDLE"          default:"25"`
	MaxConnLifetime time.Duration `envconfig:"REDIS_MAX_CONN_LIFETIME" default:"300s"`
}

type AccountApi struct {
	Host string `envconfig:"ACCOUNT_HOST" required:"true"`
	Port int    `envconfig:"ACCOUNT_PORT" required:"true"`
}
