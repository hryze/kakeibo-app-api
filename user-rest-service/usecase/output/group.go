package output

type GroupList struct {
	ApprovedGroupList   []ApprovedGroup   `json:"approved_group_list"`
	UnapprovedGroupList []UnapprovedGroup `json:"unapproved_group_list"`
}

type ApprovedGroup struct {
	GroupID             int              `json:"group_id"               db:"group_id"`
	GroupName           string           `json:"group_name"             db:"group_name"`
	ApprovedUsersList   []ApprovedUser   `json:"approved_users_list"`
	UnapprovedUsersList []UnapprovedUser `json:"unapproved_users_list"`
}

type UnapprovedGroup struct {
	GroupID             int              `json:"group_id"               db:"group_id"`
	GroupName           string           `json:"group_name"             db:"group_name"`
	ApprovedUsersList   []ApprovedUser   `json:"approved_users_list"`
	UnapprovedUsersList []UnapprovedUser `json:"unapproved_users_list"`
}

type ApprovedUser struct {
	GroupID   int    `json:"group_id"   db:"group_id"`
	UserID    string `json:"user_id"    db:"user_id"`
	UserName  string `json:"user_name"  db:"user_name"`
	ColorCode string `json:"color_code" db:"color_code"`
}

type UnapprovedUser struct {
	GroupID  int    `json:"group_id"  db:"group_id"`
	UserID   string `json:"user_id"   db:"user_id"`
	UserName string `json:"user_name" db:"user_name"`
}
