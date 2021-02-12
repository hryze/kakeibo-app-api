package userdomain

import (
	"strings"
	"unicode/utf8"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
)

type Name string

const (
	minNameLength = 1
	maxNameLength = 50
)

func NewName(name string) (Name, error) {
	if utf8.RuneCountInString(name) < minNameLength ||
		utf8.RuneCountInString(name) > maxNameLength ||
		strings.Contains(name, " ") ||
		strings.Contains(name, "　") {
		return "", apierrors.ErrInvalidUserName
	}

	return Name(name), nil
}

func (n Name) Value() string {
	return string(n)
}
