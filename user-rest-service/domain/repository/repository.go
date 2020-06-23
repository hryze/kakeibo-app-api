package repository

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type UserRepository interface {
	FindID(user *model.User) (bool, error)
	CreateUser(user *model.User) error
}