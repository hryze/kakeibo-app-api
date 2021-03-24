package groupdomain

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"

type Repository interface {
	StoreGroupAndApprovedUser(groupName GroupName, userID userdomain.UserID) (*Group, error)
	DeleteGroupAndApprovedUser(groupID GroupID) error
}
