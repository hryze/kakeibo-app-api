package input

type Group struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
}

type UnapprovedUser struct {
	UserID string `json:"user_id"`
}
