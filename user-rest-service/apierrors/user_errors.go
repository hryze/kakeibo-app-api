package apierrors

import (
	"encoding/json"
	"log"

	"golang.org/x/xerrors"
)

var (
	ErrUserNotFound    = xerrors.New("user not found")
	ErrInvalidUserID   = xerrors.New("invalid user id")
	ErrInvalidUserName = xerrors.New("invalid user name")
	ErrInvalidEmail    = xerrors.New("invalid email")
	ErrInvalidPassword = xerrors.New("invalid password")
)

type UserValidationError struct {
	UserID   string `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (e *UserValidationError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}

type UserConflictError struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

func (e *UserConflictError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}
