package persistence

import (
	"context"
	"time"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/rdb"
)

type healthRepository struct {
	*rdb.MySQLHandler
}

func NewHealthRepository(mysqlHandler *rdb.MySQLHandler) *healthRepository {
	return &healthRepository{mysqlHandler}
}

func (r *healthRepository) PingDataStore() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := r.MySQLHandler.Conn.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
