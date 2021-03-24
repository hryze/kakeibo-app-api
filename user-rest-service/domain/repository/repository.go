package repository

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type HealthRepository interface {
	PingMySQL() error
	PingRedis() error
}

type AuthRepository interface {
	GetUserID(sessionID string) (string, error)
}

type UserRepository interface {
	FindSignUpUserByUserID(userID string) (*model.SignUpUser, error)
	GetUser(userID string) (*model.LoginUser, error)
	AddSessionID(sessionID string, userID string, expiration int) error
	DeleteSessionID(sessionID string) error
}

type GroupRepository interface {
	GetGroup(groupID int) (*model.Group, error)
	PutGroup(group *model.Group, groupID int) error
	PostUnapprovedUser(unapprovedUser *model.UnapprovedUser, groupID int) (sql.Result, error)
	GetUnapprovedUser(groupUnapprovedUsersID int) (*model.UnapprovedUser, error)
	FindApprovedUser(groupID int, userID string) error
	FindUnapprovedUser(groupID int, userID string) error
	PostGroupApprovedUserAndDeleteGroupUnapprovedUser(groupID int, userID string, colorCode string) (sql.Result, error)
	GetApprovedUser(approvedUsersID int) (*model.ApprovedUser, error)
	DeleteGroupApprovedUser(groupID int, userID string) error
	DeleteGroupUnapprovedUser(groupID int, userID string) error
	FindApprovedUsersList(groupID int, groupUsersList []string) (model.GroupTasksUsersListReceiver, error)
	GetGroupUsersList(groupID int) ([]string, error)
}
