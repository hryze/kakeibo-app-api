package output

import "time"

type SignUpUser struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type LoginUser struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	SessionID string    `json:"-"`
	Expires   time.Time `json:"-"`
}
