package model

type Group struct {
	ID        int    `json:"id"         db:"id"`
	GroupName string `json:"group_name" db:"group_name"`
}
