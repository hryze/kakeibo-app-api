package db

import (
	"os"
	"time"

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

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(300 * time.Second)

	return &MySQLHandler{
		Conn: conn,
	}, nil
}
