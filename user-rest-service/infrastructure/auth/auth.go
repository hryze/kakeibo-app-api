package auth

import (
	"github.com/garyburd/redigo/redis"
	"golang.org/x/xerrors"

	"github.com/hryze/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/hryze/kakeibo-app-api/user-rest-service/config"
	"github.com/hryze/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/hryze/kakeibo-app-api/user-rest-service/infrastructure/auth/imdb"
	"github.com/hryze/kakeibo-app-api/user-rest-service/interfaces/presenter"
)

type sessionStore struct {
	*imdb.RedisHandler
}

func NewSessionStore(redisHandler *imdb.RedisHandler) *sessionStore {
	return &sessionStore{redisHandler}
}

func (s *sessionStore) PingSessionStore() error {
	conn := s.RedisHandler.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return err
	}

	return nil
}

func (s *sessionStore) StoreUserBySessionID(sessionID string, userID userdomain.UserID) error {
	conn := s.RedisHandler.Pool.Get()
	defer conn.Close()

	expirationS := int(config.Env.Cookie.Expiration.Seconds())

	if _, err := conn.Do("SET", sessionID, userID.Value(), "EX", expirationS); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}

func (s *sessionStore) DeleteUserBySessionID(sessionID string) error {
	conn := s.RedisHandler.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", sessionID); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}

func (s *sessionStore) FetchUserByUserID(sessionID string) (userdomain.UserID, error) {
	conn := s.RedisHandler.Pool.Get()
	defer conn.Close()

	userID, err := redis.String(conn.Do("GET", sessionID))
	if err != nil {
		if xerrors.Is(err, redis.ErrNil) {
			return "", apierrors.NewAuthenticationError(apierrors.NewErrorString("このページを表示するにはログインが必要です"))
		}

		return "", apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	var userValidationError presenter.UserValidationError

	userIDVo, err := userdomain.NewUserID(userID)
	if err != nil {
		userValidationError.UserID = "ユーザーIDが正しくありません"
	}

	if !userValidationError.IsEmpty() {
		return "", apierrors.NewBadRequestError(&userValidationError)
	}

	return userIDVo, nil
}
