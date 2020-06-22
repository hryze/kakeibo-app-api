package model

type User struct {
	ID       string `json:"id"                 db:"id"       validate:"required,min=1,max=10,excludesall= "`
	Name     string `json:"name"               db:"name"     validate:"required,min=1,max=50,excludesall= "`
	Email    string `json:"email"              db:"email"    validate:"required,email,min=5,max=50,excludesall= "`
	Password string `json:"password,omitempty" db:"password" validate:"required,min=8,max=50,excludesall= "`
}
