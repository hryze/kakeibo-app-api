package infrastructure

import (
	"github.com/garyburd/redigo/redis"
)

type AuthRepository struct {
	*RedisHandler
}

func (r *AuthRepository) GetUserID(sessionID string) (string, error) {
	conn := r.RedisHandler.pool.Get()
	defer conn.Close()
	userID, err := redis.String(conn.Do("GET", sessionID))
	if err != nil {
		return userID, err
	}
	return userID, nil
}
