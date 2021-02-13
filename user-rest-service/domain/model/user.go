package model

import "encoding/json"

type SignUpUser struct {
	ID       string `json:"id"                 db:"user_id"  validate:"required,min=1,max=10,excludesall= "`
	Name     string `json:"name"               db:"name"     validate:"required,min=1,max=50,excludesall= "`
	Email    string `json:"email"              db:"email"    validate:"required,email,min=5,max=50,excludesall= "`
	Password string `json:"password,omitempty" db:"password" validate:"required,min=8,max=50,excludesall= "`
}

type LoginUser struct {
	ID       string `json:"id"                 db:"user_id"`
	Name     string `json:"name"               db:"name"`
	Email    string `json:"email"              db:"email"    validate:"required,email,min=5,max=50,excludesall= "`
	Password string `json:"password,omitempty" db:"password" validate:"required,min=8,max=50,excludesall= "`
}

func (u SignUpUser) ShowUser() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return string(b), err
	}

	return string(b), nil
}

func (u LoginUser) ShowUser() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return string(b), err
	}

	return string(b), nil
}
