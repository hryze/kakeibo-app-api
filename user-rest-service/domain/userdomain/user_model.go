package userdomain

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/vo"

type SignUpUser struct {
	userID   UserID
	name     Name
	email    vo.Email
	password vo.Password
}

func NewSignUpUser(userID UserID, name Name, email vo.Email, password vo.Password) *SignUpUser {
	return &SignUpUser{
		userID:   userID,
		name:     name,
		email:    email,
		password: password,
	}
}

func NewSignUpUserFromDataSource(userID UserID, name Name, email vo.Email) *SignUpUser {
	return &SignUpUser{
		userID: userID,
		name:   name,
		email:  email,
	}
}

func (u *SignUpUser) UserID() UserID {
	return u.userID
}

func (u *SignUpUser) Name() Name {
	return u.name
}

func (u *SignUpUser) Email() vo.Email {
	return u.email
}

func (u *SignUpUser) Password() vo.Password {
	return u.password
}
