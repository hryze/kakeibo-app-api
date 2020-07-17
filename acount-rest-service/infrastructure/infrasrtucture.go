package infrastructure

type DBRepository struct {
	*AuthRepository
	*CategoriesRepository
}

func NewDBRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *DBRepository {
	DBRepository := &DBRepository{
		&AuthRepository{redisHandler},
		&CategoriesRepository{mysqlHandler},
	}
	return DBRepository
}
