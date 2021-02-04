package errors

import (
	"encoding/json"
	"log"
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

type UserNotFoundError struct {
	Message string `json:"message"`
}

func (e *UserNotFoundError) Error() string {
	return e.Message
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
