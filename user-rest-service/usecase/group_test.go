package usecase

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/groupdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
)

type mockGroupRepository struct{}

func (r *mockGroupRepository) StoreGroupAndApprovedUser(group *groupdomain.Group, userID userdomain.UserID) (*groupdomain.Group, error) {
	groupID, _ := groupdomain.NewGroupID(1)
	group = groupdomain.NewGroup(groupID, group.GroupName())

	return group, nil
}

func (r *mockGroupRepository) DeleteGroupAndApprovedUser(group *groupdomain.Group) error {
	return nil
}

func (r *mockGroupRepository) UpdateGroupName(group *groupdomain.Group) error {
	return nil
}

func (r *mockGroupRepository) FindGroupByID(groupID *groupdomain.GroupID) (*groupdomain.Group, error) {
	groupName, _ := groupdomain.NewGroupName("group1")
	group := groupdomain.NewGroup(*groupID, groupName)

	return group, nil
}

type mockGroupQueryService struct{}

func (u *mockGroupQueryService) FetchGroupList(userID string) (*output.GroupList, error) {
	return &output.GroupList{
		ApprovedGroupList: []output.ApprovedGroup{
			{
				GroupID:   1,
				GroupName: "group1",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 1, UserID: "userID1", UserName: "userName1", ColorCode: "#FF0000"},
					{GroupID: 1, UserID: "userID2", UserName: "userName2", ColorCode: "#00FFFF"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 1, UserID: "userID3", UserName: "userName3"},
				},
			},
			{
				GroupID:   2,
				GroupName: "group2",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 2, UserID: "userID1", UserName: "userName1", ColorCode: "#FF0000"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 2, UserID: "userID3", UserName: "userName3"},
					{GroupID: 2, UserID: "userID4", UserName: "userName4"},
				},
			},
			{
				GroupID:   3,
				GroupName: "group3",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 3, UserID: "userID1", UserName: "userName1", ColorCode: "#FF0000"},
					{GroupID: 3, UserID: "userID2", UserName: "userName2", ColorCode: "#00FFFF"},
				},
				UnapprovedUsersList: make([]output.UnapprovedUser, 0),
			},
		},
		UnapprovedGroupList: []output.UnapprovedGroup{
			{
				GroupID:   4,
				GroupName: "group4",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 4, UserID: "userID2", UserName: "userName2", ColorCode: "#FF0000"},
					{GroupID: 4, UserID: "userID4", UserName: "userName4", ColorCode: "#00FFFF"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 4, UserID: "userID1", UserName: "userName1"},
					{GroupID: 4, UserID: "userID3", UserName: "userName3"},
				},
			},
			{
				GroupID:   5,
				GroupName: "group5",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 5, UserID: "userID4", UserName: "userName4", ColorCode: "#FF0000"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 5, UserID: "userID1", UserName: "userName1"},
				},
			},
		},
	}, nil
}

func Test_groupUsecase_FetchGroupList(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{})

	in := input.AuthenticatedUser{
		UserID: "userID1",
	}

	gotOut, err := u.FetchGroupList(&in)
	if err != nil {
		t.Errorf("unexpected error by groupUsecase.FetchGroupList '%#v'", err)
	}

	wantOut := &output.GroupList{
		ApprovedGroupList: []output.ApprovedGroup{
			{
				GroupID:   1,
				GroupName: "group1",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 1, UserID: "userID1", UserName: "userName1", ColorCode: "#FF0000"},
					{GroupID: 1, UserID: "userID2", UserName: "userName2", ColorCode: "#00FFFF"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 1, UserID: "userID3", UserName: "userName3"},
				},
			},
			{
				GroupID:   2,
				GroupName: "group2",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 2, UserID: "userID1", UserName: "userName1", ColorCode: "#FF0000"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 2, UserID: "userID3", UserName: "userName3"},
					{GroupID: 2, UserID: "userID4", UserName: "userName4"},
				},
			},
			{
				GroupID:   3,
				GroupName: "group3",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 3, UserID: "userID1", UserName: "userName1", ColorCode: "#FF0000"},
					{GroupID: 3, UserID: "userID2", UserName: "userName2", ColorCode: "#00FFFF"},
				},
				UnapprovedUsersList: make([]output.UnapprovedUser, 0),
			},
		},
		UnapprovedGroupList: []output.UnapprovedGroup{
			{
				GroupID:   4,
				GroupName: "group4",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 4, UserID: "userID2", UserName: "userName2", ColorCode: "#FF0000"},
					{GroupID: 4, UserID: "userID4", UserName: "userName4", ColorCode: "#00FFFF"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 4, UserID: "userID1", UserName: "userName1"},
					{GroupID: 4, UserID: "userID3", UserName: "userName3"},
				},
			},
			{
				GroupID:   5,
				GroupName: "group5",
				ApprovedUsersList: []output.ApprovedUser{
					{GroupID: 5, UserID: "userID4", UserName: "userName4", ColorCode: "#FF0000"},
				},
				UnapprovedUsersList: []output.UnapprovedUser{
					{GroupID: 5, UserID: "userID1", UserName: "userName1"},
				},
			},
		},
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_groupUsecase_StoreGroup(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{})

	authenticatedUser := input.AuthenticatedUser{
		UserID: "userID1",
	}

	groupInput := input.Group{
		GroupName: "group1",
	}

	gotOut, err := u.StoreGroup(&authenticatedUser, &groupInput)
	if err != nil {
		t.Errorf("unexpected error by groupUsecase.StoreGroup '%#v'", err)
	}

	wantOut := &output.Group{
		GroupID:   1,
		GroupName: "group1",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_groupUsecase_UpdateGroupName(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{})

	groupInput := input.Group{
		GroupID:   1,
		GroupName: "group2",
	}

	gotOut, err := u.UpdateGroupName(&groupInput)
	if err != nil {
		t.Errorf("unexpected error by groupUsecase.UpdateGroupName '%#v'", err)
	}

	wantOut := &output.Group{
		GroupID:   1,
		GroupName: "group2",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
