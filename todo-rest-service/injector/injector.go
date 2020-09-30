package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/handler"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/infrastructure"
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
		AuthRepo:       infrastructure.NewAuthRepository(InjectRedis(env)),
		TodoRepo:       infrastructure.NewTodoRepository(InjectMySQL(env)),
		GroupTodoRepo:  infrastructure.NewGroupTodoRepository(InjectMySQL(env)),
		GroupTasksRepo: infrastructure.NewGroupTasksRepository(InjectMySQL(env)),
		TimeManage:     handler.NewRealTime(),
	}
}
