package model

type GroupList struct {
	GroupList []Group `json:"group_list"`
}

type Group struct {
	ID                       int      `json:"id"                                   db:"id"`
	GroupName                string   `json:"group_name"                           db:"group_name"`
	GroupUsersList           []string `json:"group_users_list,omitempty"`
	GroupUnapprovedUsersList []string `json:"group_unapproved_users_list,omitempty"`
}

type GroupUser struct {
	GroupID  int    `db:"group_id"`
	UserName string `db:"user_name"`
}

type GroupUnapprovedUser struct {
	GroupID  int    `db:"group_id"`
	UserName string `db:"user_name"`
}

func NewGroupList(groupList []Group) GroupList {
	return GroupList{
		GroupList: groupList,
	}
}
