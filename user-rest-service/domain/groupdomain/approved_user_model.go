package groupdomain

import "github.com/hryze/kakeibo-app-api/user-rest-service/domain/userdomain"

type ApprovedUser struct {
	groupID   GroupID
	userID    userdomain.UserID
	colorCode ColorCode
}

func NewApprovedUser(groupID GroupID, userID userdomain.UserID, colorCode ColorCode) *ApprovedUser {
	return &ApprovedUser{
		groupID:   groupID,
		userID:    userID,
		colorCode: colorCode,
	}
}

func (u *ApprovedUser) GroupID() GroupID {
	return u.groupID
}

func (u *ApprovedUser) UserID() userdomain.UserID {
	return u.userID
}

func (u *ApprovedUser) ColorCode() ColorCode {
	return u.colorCode
}
