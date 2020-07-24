package infrastructure

type DBRepository struct {
	*AuthRepository
	*CategoriesRepository
	*TransactionsRepository
	*BudgetsRepository
}

func NewDBRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *DBRepository {
	return &DBRepository{
		&AuthRepository{redisHandler},
		&CategoriesRepository{mysqlHandler},
		&TransactionsRepository{mysqlHandler},
		&BudgetsRepository{mysqlHandler},
	}
}
