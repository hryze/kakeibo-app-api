package queryservice

import "github.com/hryze/kakeibo-app-api/user-rest-service/usecase/output"

type GroupQueryService interface {
	FetchGroupList(userID string) (*output.GroupList, error)
	FetchUnapprovedUser(groupID int, userID string) (*output.UnapprovedUser, error)
	FetchApprovedUser(groupID int, userID string) (*output.ApprovedUser, error)
}
