package userdomain

import (
	"strings"
	"unicode/utf8"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/vo"
)

type SignUpUser struct {
	userID   UserID
	name     string
	email    vo.Email
	password vo.Password
}

const (
	minNameLength = 1
	maxNameLength = 50
)

func NewSignUpUser(userID UserID, name string, email vo.Email, password vo.Password) (*SignUpUser, error) {
	if utf8.RuneCountInString(name) < minNameLength ||
		utf8.RuneCountInString(name) > maxNameLength ||
		strings.Contains(name, " ") ||
		strings.Contains(name, "　") {
		return nil, apierrors.ErrInvalidUserName
	}

	return &SignUpUser{
		userID:   userID,
		name:     name,
		email:    email,
		password: password,
	}, nil
}

func NewSignUpUserFromDataSource(userID, name, email, password string) *SignUpUser {
	return &SignUpUser{
		userID:   UserID(userID),
		name:     name,
		email:    vo.Email(email),
		password: vo.Password(password),
	}
}

func (u *SignUpUser) UserID() UserID {
	return u.userID
}

func (u *SignUpUser) Name() string {
	return u.name
}

func (u *SignUpUser) Email() vo.Email {
	return u.email
}

func (u *SignUpUser) Password() vo.Password {
	return u.password
}
