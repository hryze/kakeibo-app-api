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

type LoginUser struct {
	userID   UserID
	name     Name
	email    vo.Email
	password vo.Password
}

func NewLoginUser(email vo.Email, password vo.Password) (*LoginUser, error) {
	return &LoginUser{
		email:    email,
		password: password,
	}, nil
}

func NewLoginUserFromDataSource(userID UserID, name Name, email vo.Email, hashPassword vo.Password) *LoginUser {
	return &LoginUser{
		userID:   userID,
		name:     name,
		email:    email,
		password: hashPassword,
	}
}

func (u *LoginUser) UserID() UserID {
	return u.userID
}

func (u *LoginUser) Name() Name {
	return u.name
}

func (u *LoginUser) Email() vo.Email {
	return u.email
}

func (u *LoginUser) Password() vo.Password {
	return u.password
}
