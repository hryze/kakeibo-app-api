package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure"
)

func InjectMySQL() *infrastructure.MySQLHandler {
	mySQLHandler, err := infrastructure.NewMySQLHandler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return mySQLHandler
}

func InjectRedis() *infrastructure.RedisHandler {
	redisHandler, err := infrastructure.NewRedisHandler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return redisHandler
}

func InjectDBHandler() *handler.DBHandler {
	return &handler.DBHandler{
		AuthRepo:  infrastructure.NewAuthRepository(InjectRedis()),
		UserRepo:  infrastructure.NewUserRepository(InjectRedis(), InjectMySQL()),
		GroupRepo: infrastructure.NewGroupRepository(InjectMySQL()),
	}
}
