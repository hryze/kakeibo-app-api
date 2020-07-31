package infrastructure

type DBRepository struct {
	*UserRepository
}

func NewDBRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *DBRepository {
	return &DBRepository{
		&UserRepository{mysqlHandler, redisHandler},
	}
}
