package persistence

import (
	"database/sql"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/errors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/datasource"
)

type userRepository struct {
	*config.RedisHandler
	*config.MySQLHandler
}

func NewUserRepository(redisHandler *config.RedisHandler, mysqlHandler *config.MySQLHandler) *userRepository {
	return &userRepository{redisHandler, mysqlHandler}
}

func (r *userRepository) FindSignUpUserByUserID(userID string) (*userdomain.SignUpUser, error) {
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

	signUpUser := userdomain.NewSignUpUserFromDataSource(signUpUserDto.UserID, signUpUserDto.Name, signUpUserDto.Email, signUpUserDto.Password)

	return signUpUser, nil
}

func (r *userRepository) FindSignUpUserByEmail(email string) (*userdomain.SignUpUser, error) {
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

	signUpUser := userdomain.NewSignUpUserFromDataSource(signUpUserDto.UserID, signUpUserDto.Name, signUpUserDto.Email, signUpUserDto.Password)

	return signUpUser, nil
}

func (r *userRepository) CreateSignUpUser(signUpUser *userdomain.SignUpUser) error {
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

func (r *userRepository) DeleteSignUpUser(signUpUser *userdomain.SignUpUser) error {
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

	var loginUserDto datasource.LoginUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, email).StructScan(&loginUserDto); err != nil {
		return nil, err
	}

	loginUser := model.NewLoginUserFromDataSource(loginUserDto.UserID, loginUserDto.Name, loginUserDto.Email, loginUserDto.Password)

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
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) AddSessionID(sessionID string, userID string, expiration int) error {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", sessionID, userID, "EX", expiration)

	return err
}

func (r *userRepository) DeleteSessionID(sessionID string) error {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", sessionID)

	return err
}
