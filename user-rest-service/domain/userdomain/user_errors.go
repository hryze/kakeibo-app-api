package userdomain

import (
	"encoding/json"
	"log"
)

type ValidationError struct {
	UserID   string `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (e *ValidationError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}

func (e *ValidationError) IsEmpty() bool {
	if e.UserID != "" ||
		e.Name != "" ||
		e.Email != "" ||
		e.Password != "" {
		return false
	}

	return true
}

type ConflictError struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

func (e *ConflictError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}
