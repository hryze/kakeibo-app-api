package infrastructure

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type UserRepository struct {
	*MySQLHandler
	*RedisHandler
}

func NewUserRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *UserRepository {
	userRepository := &UserRepository{mysqlHandler, redisHandler}
	return userRepository
}

func (u *UserRepository) FindID(signUpUser *model.SignUpUser) error {
	var dbID string
	query := "SELECT id FROM users WHERE id = ?"
	if err := u.MySQLHandler.conn.QueryRowx(query, signUpUser.ID).Scan(&dbID); err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) FindEmail(signUpUser *model.SignUpUser) error {
	var dbEmail string
	query := "SELECT email FROM users WHERE email = ?"
	if err := u.MySQLHandler.conn.QueryRowx(query, signUpUser.Email).Scan(&dbEmail); err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) CreateUser(signUpUser *model.SignUpUser) error {
	query := "INSERT INTO users(id, name, email, password) VALUES(?,?,?,?)"
	if _, err := u.MySQLHandler.conn.Exec(query, signUpUser.ID, signUpUser.Name, signUpUser.Email, signUpUser.Password); err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) FindUser(loginUser *model.LoginUser) (*model.LoginUser, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = ?"
	if err := u.MySQLHandler.conn.QueryRowx(query, loginUser.Email).StructScan(loginUser); err != nil {
		return nil, err
	}
	return loginUser, nil
}

func (u *UserRepository) SetSessionID(sessionID string, expiration int) error {
	conn := u.RedisHandler.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", sessionID, "ok", "EX", expiration)
	return err
}

func (u *UserRepository) DeleteSessionID(sessionID string) error {
	conn := u.RedisHandler.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", sessionID)
	return err
}
