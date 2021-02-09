package vo

import (
	"regexp"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
)

type Email string

const (
	maxEmailLength = 256
	emailFormat    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
)

var emailRegex = regexp.MustCompile(emailFormat)

func NewEmail(email string) (Email, error) {
	if len(email) == 0 ||
		len(email) > maxEmailLength ||
		!emailRegex.MatchString(email) {
		return "", apierrors.ErrInvalidEmail
	}

	return Email(email), nil
}

func (e Email) Value() string {
	return string(e)
}
