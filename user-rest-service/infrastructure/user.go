package infrastructure

import (
	"fmt"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type UserRepository struct {
	*RedisHandler
	*MySQLHandler
}

func NewUserRepository(redisHandler *RedisHandler, mysqlHandler *MySQLHandler) *UserRepository {
	return &UserRepository{redisHandler, mysqlHandler}
}

func (r *UserRepository) FindUserID(userID string) error {
	query := `
        SELECT
            user_id
        FROM
            users
        WHERE 
            user_id = ?`

	var dbUserID string
	if err := r.MySQLHandler.conn.QueryRowx(query, userID).Scan(&dbUserID); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindEmail(email string) error {
	query := `
        SELECT
            email
        FROM
            users
        WHERE
            email = ?`

	var dbEmail string
	if err := r.MySQLHandler.conn.QueryRowx(query, email).Scan(&dbEmail); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CreateUser(signUpUser *model.SignUpUser) error {
	query := `
        INSERT INTO users
            (user_id, name, email, password)
        VALUES
            (?,?,?,?)`

	if _, err := r.MySQLHandler.conn.Exec(query, signUpUser.ID, signUpUser.Name, signUpUser.Email, signUpUser.Password); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(signUpUser *model.SignUpUser) error {
	query := `
        DELETE
        FROM
            users
        WHERE
            user_id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, signUpUser.ID)

	return err
}

func (r *UserRepository) FindUser(loginUser *model.LoginUser) (*model.LoginUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email,
            password
        FROM 
            users
        WHERE
            email = ?`

	if err := r.MySQLHandler.conn.QueryRowx(query, loginUser.Email).StructScan(loginUser); err != nil {
		return nil, err
	}

	return loginUser, nil
}

func (r *UserRepository) GetUser(userID string) (*model.LoginUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email
        FROM 
            users
        WHERE
            user_id = ?`

	var user model.LoginUser
	if err := r.MySQLHandler.conn.QueryRowx(query, userID).StructScan(&user); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
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
