package vo

import (
	"strings"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/errors"

	"golang.org/x/crypto/bcrypt"
)

type Password string

const (
	minPasswordLength = 8
	maxPasswordLength = 50
)

func NewPassword(password string) (Password, error) {
	if len(password) < minPasswordLength ||
		len(password) > maxPasswordLength ||
		strings.Contains(password, " ") ||
		strings.Contains(password, "　") {
		return "", errors.ErrInvalidPassword
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return Password(string(hashPassword)), nil
}

func (p Password) Value() string {
	return string(p)
}
