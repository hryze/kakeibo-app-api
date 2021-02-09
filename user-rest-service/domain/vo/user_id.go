package vo

import (
	"strings"
	"unicode/utf8"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
)

type UserID string

const (
	minUserIDLength = 1
	maxUserIDLength = 10
)

func NewUserID(userID string) (UserID, error) {
	if utf8.RuneCountInString(userID) < minUserIDLength ||
		utf8.RuneCountInString(userID) > maxUserIDLength ||
		strings.Contains(userID, " ") ||
		strings.Contains(userID, "　") {
		return "", apierrors.ErrInvalidUserID
	}

	return UserID(userID), nil
}

func (i UserID) Value() string {
	return string(i)
}
