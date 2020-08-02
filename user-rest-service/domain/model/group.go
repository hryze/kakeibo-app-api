package model

type GroupList struct {
	GroupList []Group `json:"group_list"`
}

type Group struct {
	GroupID                  int                   `json:"group_id"                              db:"id"`
	GroupName                string                `json:"group_name"                            db:"group_name"`
	GroupUsersList           []GroupUser           `json:"group_users_list,omitempty"`
	GroupUnapprovedUsersList []GroupUnapprovedUser `json:"group_unapproved_users_list,omitempty"`
}

type GroupUser struct {
	GroupID  int    `json:"group_id"  db:"group_id"`
	UserID   string `json:"user_id"   db:"user_id"`
	UserName string `json:"user_name" db:"user_name"`
}

type GroupUnapprovedUser struct {
	GroupID  int    `json:"group_id"  db:"group_id"`
	UserID   string `json:"user_id"   db:"user_id"`
	UserName string `json:"user_name" db:"user_name"`
}

func NewGroupList(groupList []Group) GroupList {
	return GroupList{
		GroupList: groupList,
	}
}
