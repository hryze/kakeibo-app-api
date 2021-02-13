package vo

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

type Password string

const (
	minPasswordLength  = 8
	maxPasswordLength  = 50
	hashPasswordLength = 60
)

var ErrInvalidPassword = xerrors.New("invalid password")

func NewPassword(password string) (Password, error) {
	if l := len(password); l < minPasswordLength || l > maxPasswordLength {
		return "", xerrors.Errorf(
			"password must be %d or more and %d or less: %s: %w",
			minPasswordLength, maxPasswordLength, password, ErrInvalidPassword,
		)
	}

	if strings.Contains(password, " ") ||
		strings.Contains(password, "　") {
		return "", xerrors.Errorf(
			"password cannot contain spaces: %s: %w",
			password, ErrInvalidPassword,
		)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", xerrors.Errorf("can't generate hash password: %s", password)
	}

	return Password(string(hashPassword)), nil
}

func NewHashPassword(hashPassword string) (Password, error) {
	if len(hashPassword) != hashPasswordLength {
		return "", xerrors.Errorf(
			"hash password must be %d characters: %s: %w",
			hashPasswordLength, hashPassword, ErrInvalidPassword,
		)
	}

	return Password(string(hashPassword)), nil
}

func (p Password) Value() string {
	return string(p)
}
