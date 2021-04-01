package datasource

type Group struct {
	GroupID   int    `db:"id"`
	GroupName string `db:"group_name"`
}

type ApprovedUser struct {
	GroupID   int    `db:"group_id"`
	UserID    string `db:"user_id"`
	ColorCode string `db:"color_code"`
}

type UnapprovedUser struct {
	GroupID int    `db:"group_id"`
	UserID  string `db:"user_id"`
}
