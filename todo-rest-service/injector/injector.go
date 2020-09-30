package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/handler"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/infrastructure"
)

func InjectMySQL(isLocal bool) *infrastructure.MySQLHandler {
	mySQLHandler, err := infrastructure.NewMySQLHandler(isLocal)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return mySQLHandler
}

func InjectRedis(isLocal bool) *infrastructure.RedisHandler {
	redisHandler, err := infrastructure.NewRedisHandler(isLocal)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return redisHandler
}

func InjectDBHandler(isLocal bool) *handler.DBHandler {
	return &handler.DBHandler{
		AuthRepo:       infrastructure.NewAuthRepository(InjectRedis(isLocal)),
		TodoRepo:       infrastructure.NewTodoRepository(InjectMySQL(isLocal)),
		GroupTodoRepo:  infrastructure.NewGroupTodoRepository(InjectMySQL(isLocal)),
		GroupTasksRepo: infrastructure.NewGroupTasksRepository(InjectMySQL(isLocal)),
		TimeManage:     handler.NewRealTime(),
	}
}
