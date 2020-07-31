package infrastructure

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type UserRepository struct {
	*MySQLHandler
	*RedisHandler
}

func (r *UserRepository) FindID(signUpUser *model.SignUpUser) error {
	var dbID string
	query := "SELECT user_id FROM users WHERE user_id = ?"
	if err := r.MySQLHandler.conn.QueryRowx(query, signUpUser.ID).Scan(&dbID); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindEmail(signUpUser *model.SignUpUser) error {
	var dbEmail string
	query := "SELECT email FROM users WHERE email = ?"
	if err := r.MySQLHandler.conn.QueryRowx(query, signUpUser.Email).Scan(&dbEmail); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateUser(signUpUser *model.SignUpUser) error {
	query := "INSERT INTO users(user_id, name, email, password) VALUES(?,?,?,?)"
	if _, err := r.MySQLHandler.conn.Exec(query, signUpUser.ID, signUpUser.Name, signUpUser.Email, signUpUser.Password); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(signUpUser *model.SignUpUser) error {
	query := `DELETE FROM users WHERE user_id = ?`
	_, err := r.MySQLHandler.conn.Exec(query, signUpUser.ID)
	return err
}

func (r *UserRepository) FindUser(loginUser *model.LoginUser) (*model.LoginUser, error) {
	query := "SELECT user_id, name, email, password FROM users WHERE email = ?"
	if err := r.MySQLHandler.conn.QueryRowx(query, loginUser.Email).StructScan(loginUser); err != nil {
		return nil, err
	}
	return loginUser, nil
}

func (r *UserRepository) SetSessionID(sessionID string, loginUserID string, expiration int) error {
	conn := r.RedisHandler.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", sessionID, loginUserID, "EX", expiration)
	return err
}

func (r *UserRepository) DeleteSessionID(sessionID string) error {
	conn := r.RedisHandler.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", sessionID)
	return err
}
