package infrastructure

type DBRepository struct {
	*AuthRepository
	*TodoRepository
	*GroupTodoRepository
	*GroupTasksRepository
}

func NewDBRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *DBRepository {
	return &DBRepository{
		&AuthRepository{redisHandler},
		&TodoRepository{mysqlHandler},
		&GroupTodoRepository{mysqlHandler},
		&GroupTasksRepository{mysqlHandler},
	}
}
