package repository

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type UserRepository interface {
	FindID(signUpUser *model.SignUpUser) error
	CreateUser(user *model.SignUpUser) error
	FindUser(user *model.LoginUser) (*model.LoginUser, error)
	SetSessionID(sessionID string, expiration int) error
}
