package userdomain

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/vo"

type Repository interface {
	FindSignUpUserByUserID(userID UserID) (*SignUpUser, error)
	FindSignUpUserByEmail(email vo.Email) (*SignUpUser, error)
	CreateSignUpUser(user *SignUpUser) error
	DeleteSignUpUser(signUpUser *SignUpUser) error
	FindLoginUserByEmail(email vo.Email) (*LoginUser, error)
	GetUser(userID string) (*LoginUser, error)
	AddSessionID(sessionID string, userID UserID, expiration int) error
	DeleteSessionID(sessionID string) error
}
