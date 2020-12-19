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

func InjectDBHandler() *handler.DBHandler {
	return &handler.DBHandler{
		HealthRepo:       infrastructure.NewHealthRepository(InjectRedis(), InjectMySQL()),
		AuthRepo:         infrastructure.NewAuthRepository(InjectRedis()),
		TodoRepo:         infrastructure.NewTodoRepository(InjectMySQL()),
		ShoppingListRepo: infrastructure.NewShoppingListRepository(InjectMySQL()),
		GroupTodoRepo:    infrastructure.NewGroupTodoRepository(InjectMySQL()),
		GroupTasksRepo:   infrastructure.NewGroupTasksRepository(InjectMySQL()),
		TimeManage:       handler.NewRealTime(),
	}
}
