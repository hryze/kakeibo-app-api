package infrastructure

import (
	"github.com/garyburd/redigo/redis"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/db"
)

type AuthRepository struct {
	*db.RedisHandler
}

func NewAuthRepository(redisHandler *db.RedisHandler) *AuthRepository {
	return &AuthRepository{redisHandler}
}

func (r *AuthRepository) GetUserID(sessionID string) (string, error) {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	userID, err := redis.String(conn.Do("GET", sessionID))
	if err != nil {
		return userID, err
	}

	return userID, nil
}
