package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Env ENV

func init() {
	isLocal := flag.Bool("local", false, "Please specify -local flag")
	flag.Parse()

	if *isLocal {
		if err := godotenv.Load("../../development.env"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if err := envconfig.Process("", &Env); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type ENV struct {
	Server     server
	Cors       cors
	Cookie     cookie
	MySQL      mysql
	Redis      redis
	AccountApi accountApi
}

type server struct {
	Port int `envconfig:"SERVER_PORT" required:"true"`
}

type cors struct {
	AllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" required:"true"`
}

type cookie struct {
	Domain string `envconfig:"COOKIE_DOMAIN" required:"true"`
}

type mysql struct {
	Dsn             string        `envconfig:"MYSQL_DSN"               required:"true"`
	MaxConn         int           `envconfig:"MYSQL_MAX_CONN"          default:"25"`
	MaxIdleConn     int           `envconfig:"MYSQL_MAX_IDLE"          default:"25"`
	MaxConnLifetime time.Duration `envconfig:"MYSQL_MAX_CONN_LIFETIME" default:"300s"`
}

type redis struct {
	Dsn             string        `envconfig:"REDIS_DSN"               required:"true"`
	MaxConn         int           `envconfig:"REDIS_MAX_CONN"          default:"25"`
	MaxIdleConn     int           `envconfig:"REDIS_MAX_IDLE"          default:"25"`
	MaxConnLifetime time.Duration `envconfig:"REDIS_MAX_CONN_LIFETIME" default:"300s"`
}

type accountApi struct {
	Host string `envconfig:"ACCOUNT_HOST" required:"true"`
	Port int    `envconfig:"ACCOUNT_PORT" required:"true"`
}
