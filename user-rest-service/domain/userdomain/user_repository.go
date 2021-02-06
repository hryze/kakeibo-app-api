package userdomain

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"

type Repository interface {
	FindSignUpUserByUserID(userID string) (*SignUpUser, error)
	FindSignUpUserByEmail(email string) (*SignUpUser, error)
	CreateSignUpUser(user *SignUpUser) error
	DeleteSignUpUser(signUpUser *SignUpUser) error
	FindLoginUserByEmail(email string) (*model.LoginUser, error)
	GetUser(userID string) (*model.LoginUser, error)
	AddSessionID(sessionID string, userID string, expiration int) error
	DeleteSessionID(sessionID string) error
}
