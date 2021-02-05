package datasource

type SignUpUser struct {
	UserID   string `db:"user_id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type LoginUser struct {
	UserID   string `db:"user_id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}
