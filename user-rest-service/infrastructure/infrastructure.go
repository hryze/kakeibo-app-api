package infrastructure

type DBRepository struct {
	*AuthRepository
	*UserRepository
	*GroupRepository
}

func NewDBRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *DBRepository {
	return &DBRepository{
		&AuthRepository{redisHandler},
		&UserRepository{mysqlHandler, redisHandler},
		&GroupRepository{mysqlHandler},
	}
}
