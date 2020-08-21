package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/handler"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/infrastructure"
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

func InjectTransactionsRepository() *infrastructure.TransactionsRepository {
	mySQLHandler := InjectMySQL()
	return infrastructure.NewTransactionsRepository(mySQLHandler)
}

func InjectCategoriesRepository() *infrastructure.CategoriesRepository {
	mySQLHandler := InjectMySQL()
	return infrastructure.NewCategoriesRepository(mySQLHandler)
}

func InjectBudgetsRepository() *infrastructure.BudgetsRepository {
	mySQLHandler := InjectMySQL()
	return infrastructure.NewBudgetsRepository(mySQLHandler)
}

func InjectGroupTransactionsRepository() *infrastructure.GroupTransactionsRepository {
	mySQLHandler := InjectMySQL()
	return infrastructure.NewGroupTransactionsRepository(mySQLHandler)
}

func InjectGroupCategoriesRepository() *infrastructure.GroupCategoriesRepository {
	mySQLHandler := InjectMySQL()
	return infrastructure.NewGroupCategoriesRepository(mySQLHandler)
}

func InjectGroupBudgetsRepository() *infrastructure.GroupBudgetsRepository {
	mySQLHandler := InjectMySQL()
	return infrastructure.NewGroupBudgetsRepository(mySQLHandler)
}

func InjectDBHandler() *handler.DBHandler {
	return &handler.DBHandler{
		AuthRepo:              InjectAuthRepository(),
		TransactionsRepo:      InjectTransactionsRepository(),
		CategoriesRepo:        InjectCategoriesRepository(),
		BudgetsRepo:           InjectBudgetsRepository(),
		GroupTransactionsRepo: InjectGroupTransactionsRepository(),
		GroupCategoriesRepo:   InjectGroupCategoriesRepository(),
		GroupBudgetsRepo:      InjectGroupBudgetsRepository(),
	}
}
