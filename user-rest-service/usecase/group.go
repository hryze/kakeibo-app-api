package usecase

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/queryservice"
)

type GroupUsecase interface {
	FetchGroupList(in *input.AuthenticatedUser) (*output.GroupList, error)
}

type groupUsecase struct {
	groupQueryService queryservice.GroupQueryService
}

func NewGroupUsecase(groupQueryService queryservice.GroupQueryService) *groupUsecase {
	return &groupUsecase{
		groupQueryService: groupQueryService,
	}
}

func (u *groupUsecase) FetchGroupList(in *input.AuthenticatedUser) (*output.GroupList, error) {
	groupList, err := u.groupQueryService.FetchGroupList(in.UserID)
	if err != nil {
		return nil, err
	}

	return groupList, nil
}
