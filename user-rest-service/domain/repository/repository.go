package repository

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type DBRepository interface {
	UserRepository
}

type UserRepository interface {
	FindID(signUpUser *model.SignUpUser) error
	FindEmail(signUpUser *model.SignUpUser) error
	CreateUser(user *model.SignUpUser) error
	DeleteUser(signUpUser *model.SignUpUser) error
	FindUser(user *model.LoginUser) (*model.LoginUser, error)
	SetSessionID(sessionID string, loginUserID string, expiration int) error
	DeleteSessionID(sessionID string) error
}
