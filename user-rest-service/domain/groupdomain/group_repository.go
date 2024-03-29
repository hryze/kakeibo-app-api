package groupdomain

import "github.com/hryze/kakeibo-app-api/user-rest-service/domain/userdomain"

type Repository interface {
	StoreGroupAndApprovedUser(group *Group, userID userdomain.UserID) (*Group, error)
	DeleteGroupAndApprovedUser(group *Group) error
	UpdateGroupName(group *Group) error
	StoreUnapprovedUser(unapprovedUser *UnapprovedUser) error
	DeleteApprovedUser(approvedUser *ApprovedUser) error
	StoreApprovedUser(approvedUser *ApprovedUser) error
	DeleteUnapprovedUser(unapprovedUser *UnapprovedUser) error
	FindGroupByID(groupID *GroupID) (*Group, error)
	FindApprovedUser(groupID GroupID, userID userdomain.UserID) (*ApprovedUser, error)
	FindUnapprovedUser(groupID GroupID, userID userdomain.UserID) (*UnapprovedUser, error)
	FindApprovedUsersList(groupID GroupID, userIDList userdomain.UserIDList) ([]ApprovedUser, error)
	FetchApprovedUserIDList(groupID GroupID) ([]userdomain.UserID, error)
}
