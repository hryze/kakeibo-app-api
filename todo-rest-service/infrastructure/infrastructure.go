package infrastructure

type DBRepository struct {
	*AuthRepository
	*TodoRepository
}

func NewDBRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *DBRepository {
	return &DBRepository{
		&AuthRepository{redisHandler},
		&TodoRepository{mysqlHandler},
	}
}
