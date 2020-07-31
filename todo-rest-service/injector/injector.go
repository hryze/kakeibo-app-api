package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/handler"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/infrastructure"
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

func InjectDBRepository() *infrastructure.DBRepository {
	mySQLHandler := InjectMySQL()
	redisHandler := InjectRedis()
	return infrastructure.NewDBRepository(mySQLHandler, redisHandler)
}

func InjectDBHandler() *handler.DBHandler {
	DBRepository := InjectDBRepository()
	return handler.NewDBHandler(DBRepository)
}
