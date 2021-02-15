package userdomain

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/xerrors"
)

type Name string

const (
	minNameLength = 1
	maxNameLength = 50
)

var ErrInvalidUserName = xerrors.New("invalid user name")

func NewName(name string) (Name, error) {
	if n := utf8.RuneCountInString(name); n < minNameLength || n > maxNameLength {
		return "", xerrors.Errorf(
			"user name must be %d or more and %d or less: %s: %w",
			minNameLength, maxNameLength, name, ErrInvalidUserName,
		)
	}

	if strings.Contains(name, " ") ||
		strings.Contains(name, "　") {
		return "", xerrors.Errorf(
			"user name cannot contain spaces: %s: %w",
			name, ErrInvalidUserName,
		)
	}

	return Name(name), nil
}

func (n Name) Value() string {
	return string(n)
}
