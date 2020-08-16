package model

type GroupTasksUser struct {
	ID      int    `json:"id"       db:"id"`
	UserID  string `json:"user_id"  db:"user_id"`
	GroupID int    `json:"group_id" db:"group_id"`
}
