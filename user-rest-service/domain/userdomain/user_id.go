package userdomain

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/xerrors"
)

type UserID string

const (
	minUserIDLength = 1
	maxUserIDLength = 10
)

var ErrInvalidUserID = xerrors.New("invalid user id")

func NewUserID(userID string) (UserID, error) {
	if n := utf8.RuneCountInString(userID); n < minUserIDLength || n > maxUserIDLength {
		return "", xerrors.Errorf(
			"user id must be %d or more and %d or less: %s: %w",
			minUserIDLength, maxUserIDLength, userID, ErrInvalidUserID,
		)
	}

	if strings.Contains(userID, " ") ||
		strings.Contains(userID, "　") {
		return "", xerrors.Errorf(
			"user id cannot contain spaces: %s: %w",
			userID, ErrInvalidUserID,
		)
	}

	return UserID(userID), nil
}

func (i UserID) Value() string {
	return string(i)
}

type UserIDList []UserID

func NewUserIDList(userIDList []string) (UserIDList, error) {
	userIDListVo := make(UserIDList, len(userIDList))
	for i, userID := range userIDList {
		userIDVo, err := NewUserID(userID)
		if err != nil {
			return nil, err
		}

		userIDListVo[i] = userIDVo
	}

	return userIDListVo, nil
}

func (il UserIDList) Value() []string {
	userIDList := make([]string, len(il))
	for i, userIDVo := range il {
		userIDList[i] = string(userIDVo)
	}

	return userIDList
}
