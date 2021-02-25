package userdomain

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/vo"

type LoginUser struct {
	userID   UserID
	name     Name
	email    vo.Email
	password vo.Password
}

func NewLoginUser(email vo.Email, password vo.Password) *LoginUser {
	return &LoginUser{
		email:    email,
		password: password,
	}
}

func NewLoginUserWithHashPassword(userID UserID, name Name, email vo.Email, hashPassword vo.Password) *LoginUser {
	return &LoginUser{
		userID:   userID,
		name:     name,
		email:    email,
		password: hashPassword,
	}
}

func NewLoginUserWithNoPassword(userID UserID, name Name, email vo.Email) *LoginUser {
	return &LoginUser{
		userID: userID,
		name:   name,
		email:  email,
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
