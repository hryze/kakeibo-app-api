package model

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
	GroupID  int    `json:"group_id"  db:"group_id"`
	UserID   string `json:"user_id"   db:"user_id"`
	UserName string `json:"user_name" db:"user_name"`
}

type UnapprovedUser struct {
	GroupID  int    `json:"group_id"  db:"group_id"`
	UserID   string `json:"user_id"   db:"user_id"`
	UserName string `json:"user_name" db:"user_name"`
}

type Group struct {
	GroupID   int    `json:"group_id"   db:"id"`
	GroupName string `json:"group_name" db:"group_name"`
}

func NewGroupList(approvedGroupList []ApprovedGroup, unapprovedGroupList []UnapprovedGroup) GroupList {
	return GroupList{
		ApprovedGroupList:   approvedGroupList,
		UnapprovedGroupList: unapprovedGroupList,
	}
}
