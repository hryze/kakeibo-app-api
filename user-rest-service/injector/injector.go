package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure"
)

func InjectMySQL(env string) *infrastructure.MySQLHandler {
	mySQLHandler, err := infrastructure.NewMySQLHandler(env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return mySQLHandler
}

func InjectRedis(env string) *infrastructure.RedisHandler {
	redisHandler, err := infrastructure.NewRedisHandler(env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return redisHandler
}

func InjectDBHandler(env string) *handler.DBHandler {
	return &handler.DBHandler{
		AuthRepo:  infrastructure.NewAuthRepository(InjectRedis(env)),
		UserRepo:  infrastructure.NewUserRepository(InjectRedis(env), InjectMySQL(env)),
		GroupRepo: infrastructure.NewGroupRepository(InjectMySQL(env)),
	}
}
