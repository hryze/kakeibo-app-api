package groupdomain

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"

type UnapprovedUser struct {
	groupID GroupID
	userID  userdomain.UserID
}

func NewUnapprovedUser(groupID GroupID, userID userdomain.UserID) *UnapprovedUser {
	return &UnapprovedUser{
		groupID: groupID,
		userID:  userID,
	}
}

func (u *UnapprovedUser) GroupID() GroupID {
	return u.groupID
}

func (u *UnapprovedUser) UserID() userdomain.UserID {
	return u.userID
}
