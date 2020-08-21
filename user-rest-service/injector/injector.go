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

func InjectAuthRepository() *infrastructure.AuthRepository {
	redisHandler := InjectRedis()
	return infrastructure.NewAuthRepository(redisHandler)
}

func InjectUserRepository() *infrastructure.UserRepository {
	mySQLHandler := InjectMySQL()
	redisHandler := InjectRedis()
	return infrastructure.NewUserRepository(mySQLHandler, redisHandler)
}

func InjectGroupRepository() *infrastructure.GroupRepository {
	mySQLHandler := InjectMySQL()
	return infrastructure.NewGroupRepository(mySQLHandler)
}

func InjectDBHandler() *handler.DBHandler {
	authRepo := InjectAuthRepository()
	userRepo := InjectUserRepository()
	groupRepo := InjectGroupRepository()
	return handler.NewDBHandler(authRepo, userRepo, groupRepo)
}
