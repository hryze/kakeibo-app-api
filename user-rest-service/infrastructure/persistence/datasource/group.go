package datasource

type Group struct {
	GroupID   int    `db:"id"`
	GroupName string `db:"group_name"`
}
