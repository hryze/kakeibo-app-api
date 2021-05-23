package injector

import (
	"fmt"
	"os"

	"github.com/hryze/kakeibo-app-api/account-rest-service/handler"
	"github.com/hryze/kakeibo-app-api/account-rest-service/infrastructure"
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
		HealthRepo:            infrastructure.NewHealthRepository(InjectRedis(), InjectMySQL()),
		AuthRepo:              infrastructure.NewAuthRepository(InjectRedis()),
		TransactionsRepo:      infrastructure.NewTransactionsRepository(InjectMySQL()),
		CategoriesRepo:        infrastructure.NewCategoriesRepository(InjectMySQL()),
		BudgetsRepo:           infrastructure.NewBudgetsRepository(InjectMySQL()),
		GroupTransactionsRepo: infrastructure.NewGroupTransactionsRepository(InjectMySQL()),
		GroupCategoriesRepo:   infrastructure.NewGroupCategoriesRepository(InjectMySQL()),
		GroupBudgetsRepo:      infrastructure.NewGroupBudgetsRepository(InjectMySQL()),
		TimeManage:            handler.NewRealTime(),
	}
}
