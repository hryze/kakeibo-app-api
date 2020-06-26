package infrastructure

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type SQLHandler struct {
	DB *sqlx.DB
}

func NewSQLHandler() (*SQLHandler, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, err
	}
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	SQLHandler := new(SQLHandler)
	SQLHandler.DB = db

	return SQLHandler, nil
}
