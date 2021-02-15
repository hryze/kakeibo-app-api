package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/db"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/interfaces/handler"
)

func InjectMySQL() *db.MySQLHandler {
	mySQLHandler, err := db.NewMySQLHandler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return mySQLHandler
}

func InjectRedis() *db.RedisHandler {
	redisHandler, err := db.NewRedisHandler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return redisHandler
}

func InjectDBHandler() *handler.DBHandler {
	return &handler.DBHandler{
		HealthRepo: infrastructure.NewHealthRepository(InjectRedis(), InjectMySQL()),
		AuthRepo:   infrastructure.NewAuthRepository(InjectRedis()),
		UserRepo:   infrastructure.NewUserRepository(InjectRedis(), InjectMySQL()),
		GroupRepo:  infrastructure.NewGroupRepository(InjectMySQL()),
	}
}
