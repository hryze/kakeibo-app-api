package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/datasource"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/errors"
)

type userRepository struct {
	*config.RedisHandler
	*config.MySQLHandler
}

func NewUserRepository(redisHandler *config.RedisHandler, mysqlHandler *config.MySQLHandler) *userRepository {
	return &userRepository{redisHandler, mysqlHandler}
}

func (r *userRepository) FindSignUpUserByUserID(userID string) (*model.SignUpUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email,
            password
        FROM
            users
        WHERE
            user_id = ?`

	var signUpUserDto datasource.SignUpUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, userID).StructScan(&signUpUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	signUpUser := model.NewSignUpUserFromDataSource(signUpUserDto.UserID, signUpUserDto.Name, signUpUserDto.Email, signUpUserDto.Password)

	return signUpUser, nil
}

func (r *userRepository) FindSignUpUserByEmail(email string) (*model.SignUpUser, error) {
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

	var signUpUserDto datasource.SignUpUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, email).StructScan(&signUpUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	signUpUser := model.NewSignUpUserFromDataSource(signUpUserDto.UserID, signUpUserDto.Name, signUpUserDto.Email, signUpUserDto.Password)

	return signUpUser, nil
}

func (r *userRepository) CreateSignUpUser(signUpUser *model.SignUpUser) error {
	query := `
        INSERT INTO users
            (user_id, name, email, password)
        VALUES
            (?,?,?,?)`

	if _, err := r.MySQLHandler.Conn.Exec(query, signUpUser.UserID(), signUpUser.Name(), signUpUser.Email(), signUpUser.Password()); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteSignUpUser(signUpUser *model.SignUpUser) error {
	query := `
        DELETE
        FROM
            users
        WHERE
            user_id = ?`

	_, err := r.MySQLHandler.Conn.Exec(query, signUpUser.UserID())

	return err
}

func (r *userRepository) FindLoginUserByEmail(email string) (*model.LoginUser, error) {
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

	var loginUser *model.LoginUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, email).StructScan(loginUser); err != nil {
		return nil, err
	}

	return loginUser, nil
}

func (r *userRepository) GetUser(userID string) (*model.LoginUser, error) {
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
	if err := r.MySQLHandler.Conn.QueryRowx(query, userID).StructScan(&user); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) AddSessionID(sessionID string, loginUserID string, expiration int) error {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", sessionID, loginUserID, "EX", expiration)

	return err
}

func (r *userRepository) DeleteSessionID(sessionID string) error {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", sessionID)

	return err
}
