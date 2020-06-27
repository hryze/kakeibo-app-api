package infrastructure

import (
	"database/sql"
	"errors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(sqlHandler SQLHandler) *UserRepository {
	userRepository := UserRepository{sqlHandler}
	return &userRepository
}

func (u *UserRepository) FindID(signUpUser *model.SignUpUser) (string, error) {
	var dbID string
	query := "SELECT id FROM users WHERE id = ?"
	if err := u.SQLHandler.DB.QueryRowx(query, signUpUser.ID).Scan(&dbID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbID, nil
		} else if err != nil {
			return dbID, err
		}
	}
	return dbID, nil
}

func (u *UserRepository) CreateUser(signUpUser *model.SignUpUser) error {
	query := "INSERT INTO users(id, name, email, password) VALUES(?,?,?,?)"
	if _, err := u.SQLHandler.DB.Exec(query, signUpUser.ID, signUpUser.Name, signUpUser.Email, signUpUser.Password); err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) FindUser(loginUser *model.LoginUser) (*model.LoginUser, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = ?"
	if err := u.SQLHandler.DB.QueryRowx(query, loginUser.Email).StructScan(loginUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
	}
	return loginUser, nil
}
