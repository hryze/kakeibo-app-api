package groupdomain

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/xerrors"
)

type GroupName string

const (
	minNameLength = 1
	maxNameLength = 20
)

var (
	ErrCharacterCountGroupName = xerrors.New("invalid character count for group name")
	ErrPrefixSpaceGroupName    = xerrors.New("invalid prefix space in group name")
	ErrSuffixSpaceGroupName    = xerrors.New("invalid suffix space in group name")
)

func NewGroupName(groupName string) (GroupName, error) {
	if n := utf8.RuneCountInString(groupName); n < minNameLength || n > maxNameLength {
		return "", xerrors.Errorf(
			"group name must be %d or more and %d or less: %s: %w",
			minNameLength, maxNameLength, groupName, ErrCharacterCountGroupName,
		)
	}

	if strings.HasPrefix(groupName, " ") ||
		strings.HasPrefix(groupName, "　") {
		return "", xerrors.Errorf(
			"group name prefix cannot contain spaces: %s: %w",
			groupName, ErrPrefixSpaceGroupName,
		)
	}

	if strings.HasSuffix(groupName, " ") ||
		strings.HasSuffix(groupName, "　") {
		return "", xerrors.Errorf(
			"group name suffix cannot contain spaces: %s: %w",
			groupName, ErrSuffixSpaceGroupName,
		)
	}

	return GroupName(groupName), nil
}

func (n GroupName) Value() string {
	return string(n)
}
