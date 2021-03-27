package groupdomain

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"

type Repository interface {
	StoreGroupAndApprovedUser(group *Group, userID userdomain.UserID) (*Group, error)
	DeleteGroupAndApprovedUser(group *Group) error
}
