package vo

import (
	"regexp"

	"golang.org/x/xerrors"
)

type Email string

const (
	maxEmailLength = 256
	emailFormat    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
)

var (
	emailRegex      = regexp.MustCompile(emailFormat)
	ErrInvalidEmail = xerrors.New("invalid email")
)

func NewEmail(email string) (Email, error) {
	if len(email) == 0 ||
		len(email) > maxEmailLength {
		return "", xerrors.Errorf(
			"email must be %d or less: %s: %w",
			maxEmailLength, email, ErrInvalidEmail,
		)
	}

	if ok := emailRegex.MatchString(email); !ok {
		return "", xerrors.Errorf(
			"invalid email format: %s: %w",
			email, ErrInvalidEmail,
		)
	}

	return Email(email), nil
}

func (e Email) Value() string {
	return string(e)
}
