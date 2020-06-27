package repository

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type UserRepository interface {
	FindID(user *model.SignUpUser) (string, error)
	CreateUser(user *model.SignUpUser) error
	FindUser(user *model.LoginUser) (*model.LoginUser, error)
}
