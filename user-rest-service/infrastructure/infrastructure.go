package infrastructure

import (
	"database/sql"
	"errors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
)

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(sqlHandler SQLHandler) repository.UserRepository {
	userRepository := UserRepository{sqlHandler}
	return &userRepository
}

func (u *UserRepository) FindID(user *model.User) (bool, error) {
	var dbID string
	if err := u.SQLHandler.DB.QueryRowx("SELECT id FROM users WHERE id = ?", user.ID).Scan(&dbID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		} else if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (u *UserRepository) CreateUser(user *model.User) error {
	if _, err := u.SQLHandler.DB.Exec("INSERT INTO users(id, name, email, password) VALUES(?,?,?,?)", user.ID, user.Name, user.Email, user.Password); err != nil {
		return err
	}
	return nil
}
