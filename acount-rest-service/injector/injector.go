package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/handler"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/infrastructure"
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
		AuthRepo:              infrastructure.NewAuthRepository(InjectRedis(isLocal)),
		TransactionsRepo:      infrastructure.NewTransactionsRepository(InjectMySQL(isLocal)),
		CategoriesRepo:        infrastructure.NewCategoriesRepository(InjectMySQL(isLocal)),
		BudgetsRepo:           infrastructure.NewBudgetsRepository(InjectMySQL(isLocal)),
		GroupTransactionsRepo: infrastructure.NewGroupTransactionsRepository(InjectMySQL(isLocal)),
		GroupCategoriesRepo:   infrastructure.NewGroupCategoriesRepository(InjectMySQL(isLocal)),
		GroupBudgetsRepo:      infrastructure.NewGroupBudgetsRepository(InjectMySQL(isLocal)),
	}
}
