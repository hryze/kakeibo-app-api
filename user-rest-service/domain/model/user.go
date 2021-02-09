package model

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
)

const (
	minUserIDLength   = 1
	maxUserIDLength   = 10
	minNameLength     = 1
	maxNameLength     = 50
	minEmailLength    = 5
	maxEmailLength    = 50
	minPasswordLength = 8
	maxPasswordLength = 50
	emailRegexString  = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

var emailRegex = regexp.MustCompile(emailRegexString)

type SignUpUser struct {
	userID   string
	name     string
	email    string
	password string
}

func (u *SignUpUser) UserID() string {
	return u.userID
}

func (u *SignUpUser) Email() string {
	return u.email
}

func (u *SignUpUser) Name() string {
	return u.name
}

func (u *SignUpUser) Password() string {
	return u.password
}

func (u *SignUpUser) SetPassword(hashStr string) {
	u.password = hashStr
}

func NewSignUpUser(userID, name, email, password string) (*SignUpUser, error) {
	var userValidationError apierrors.UserValidationError

	if utf8.RuneCountInString(userID) < minUserIDLength ||
		utf8.RuneCountInString(userID) > maxUserIDLength ||
		strings.Contains(userID, " ") ||
		strings.Contains(userID, "　") {
		userValidationError.UserID = "ユーザーIDを正しく入力してください"
	}

	if utf8.RuneCountInString(name) < minNameLength ||
		utf8.RuneCountInString(name) > maxNameLength ||
		strings.Contains(name, " ") ||
		strings.Contains(name, "　") {
		userValidationError.Name = "名前を正しく入力してください"
	}

	if len(email) < minEmailLength ||
		len(email) > maxEmailLength ||
		strings.Contains(email, " ") ||
		strings.Contains(email, "　") ||
		!emailRegex.MatchString(email) {
		userValidationError.Email = "メールアドレスを正しく入力してください"
	}

	if len(password) < minPasswordLength ||
		len(password) > maxPasswordLength ||
		strings.Contains(password, " ") ||
		strings.Contains(password, "　") {
		userValidationError.Password = "パスワードを正しく入力してください"
	}

	if userValidationError.UserID != "" ||
		userValidationError.Name != "" ||
		userValidationError.Email != "" ||
		userValidationError.Password != "" {
		return nil, &userValidationError
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
		userID:   userID,
		name:     name,
		email:    email,
		password: password,
	}
}

type LoginUser struct {
	userID   string
	name     string
	email    string
	password string
}

func (u *LoginUser) UserID() string {
	return u.userID
}

func (u *LoginUser) Email() string {
	return u.email
}

func (u *LoginUser) Name() string {
	return u.name
}

func (u *LoginUser) Password() string {
	return u.password
}

func (u *LoginUser) SetPassword(hashStr string) {
	u.password = hashStr
}

func NewLoginUser(email, password string) (*LoginUser, error) {
	var userValidationError apierrors.UserValidationError

	if len(email) < minEmailLength ||
		len(email) > maxEmailLength ||
		strings.Contains(email, " ") ||
		strings.Contains(email, "　") ||
		!emailRegex.MatchString(email) {
		userValidationError.Email = "メールアドレスを正しく入力してください"
	}

	if len(password) < minPasswordLength ||
		len(password) > maxPasswordLength ||
		strings.Contains(password, " ") ||
		strings.Contains(password, "　") {
		userValidationError.Password = "パスワードを正しく入力してください"
	}

	if userValidationError.Email != "" ||
		userValidationError.Password != "" {
		return nil, &userValidationError
	}

	return &LoginUser{
		email:    email,
		password: password,
	}, nil
}

func NewLoginUserFromDataSource(userID, name, email, password string) *LoginUser {
	return &LoginUser{
		userID:   userID,
		name:     name,
		email:    email,
		password: password,
	}
}
