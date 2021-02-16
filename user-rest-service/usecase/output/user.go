package output

type SignUpUser struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type LoginUser struct {
	UserID string     `json:"user_id"`
	Name   string     `json:"name"`
	Email  string     `json:"email"`
	Cookie CookieInfo `json:"-"`
}

type CookieInfo struct {
	SessionID string
}
