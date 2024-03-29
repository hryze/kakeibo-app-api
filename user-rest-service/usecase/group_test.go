package usecase

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hryze/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/hryze/kakeibo-app-api/user-rest-service/domain/groupdomain"
	"github.com/hryze/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase/output"
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

func (r *mockGroupRepository) StoreUnapprovedUser(unapprovedUser *groupdomain.UnapprovedUser) error {
	return nil
}

func (r *mockGroupRepository) DeleteApprovedUser(approvedUser *groupdomain.ApprovedUser) error {
	return nil
}

func (r *mockGroupRepository) StoreApprovedUser(approvedUser *groupdomain.ApprovedUser) error {
	return nil
}

func (r *mockGroupRepository) DeleteUnapprovedUser(unapprovedUser *groupdomain.UnapprovedUser) error {
	return nil
}

func (r *mockGroupRepository) FindGroupByID(groupID *groupdomain.GroupID) (*groupdomain.Group, error) {
	groupName, _ := groupdomain.NewGroupName("group1")
	group := groupdomain.NewGroup(*groupID, groupName)

	return group, nil
}

func (r *mockGroupRepository) FindApprovedUser(groupID groupdomain.GroupID, userID userdomain.UserID) (*groupdomain.ApprovedUser, error) {
	if groupID.Value() == 1 {
		return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
	}

	colorCode, _ := groupdomain.NewColorCode("#FF0000")
	approvedUser := groupdomain.NewApprovedUser(groupID, userID, colorCode)

	return approvedUser, nil
}

func (r *mockGroupRepository) FindUnapprovedUser(groupID groupdomain.GroupID, userID userdomain.UserID) (*groupdomain.UnapprovedUser, error) {
	if groupID.Value() == 1 {
		return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
	}

	unapprovedUser := groupdomain.NewUnapprovedUser(groupID, userID)

	return unapprovedUser, nil
}

func (r *mockGroupRepository) FindApprovedUsersList(groupID groupdomain.GroupID, userIDList userdomain.UserIDList) ([]groupdomain.ApprovedUser, error) {
	approvedUserList := make([]groupdomain.ApprovedUser, len(userIDList))
	for i, userID := range userIDList {
		if i == 0 {
			colorCode, _ := groupdomain.NewColorCode("#FF0000")
			approvedUserList[i] = *groupdomain.NewApprovedUser(groupID, userID, colorCode)
			continue
		}

		colorCode, _ := groupdomain.NewColorCodeToUser(userIDList[:i])
		approvedUserList[i] = *groupdomain.NewApprovedUser(groupID, userID, colorCode)
	}

	return approvedUserList, nil
}

func (r *mockGroupRepository) FetchApprovedUserIDList(groupID groupdomain.GroupID) ([]userdomain.UserID, error) {
	return []userdomain.UserID{
		"userID1",
		"userID2",
		"userID3",
	}, nil
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

func (u *mockGroupQueryService) FetchUnapprovedUser(groupID int, userID string) (*output.UnapprovedUser, error) {
	return &output.UnapprovedUser{
		GroupID:  1,
		UserID:   "userID1",
		UserName: "userName1",
	}, nil
}

func (u *mockGroupQueryService) FetchApprovedUser(groupID int, userID string) (*output.ApprovedUser, error) {
	return &output.ApprovedUser{
		GroupID:   2,
		UserID:    "userID4",
		UserName:  "userName4",
		ColorCode: "#8000FF",
	}, nil
}

func Test_groupUsecase_FetchGroupList(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

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
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

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
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	groupInput := input.Group{
		GroupID:   1,
		GroupName: "group1",
	}

	gotOut, err := u.UpdateGroupName(&groupInput)
	if err != nil {
		t.Errorf("unexpected error by groupUsecase.UpdateGroupName '%#v'", err)
	}

	wantOut := &output.Group{
		GroupID:   1,
		GroupName: "group1",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_groupUsecase_StoreGroupUnapprovedUser(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	unapprovedUser := input.UnapprovedUser{
		UserID: "userID1",
	}

	groupInput := input.Group{
		GroupID: 1,
	}

	gotOut, err := u.StoreGroupUnapprovedUser(&unapprovedUser, &groupInput)
	if err != nil {
		t.Errorf("unexpected error by groupUsecase.StoreGroupUnapprovedUser '%#v'", err)
	}

	wantOut := &output.UnapprovedUser{
		GroupID:  1,
		UserID:   "userID1",
		UserName: "userName1",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_groupUsecase_DeleteGroupApprovedUser(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	authenticatedUser := input.AuthenticatedUser{
		UserID: "userID1",
	}

	groupInput := input.Group{
		GroupID: 2,
	}

	if err := u.DeleteGroupApprovedUser(&authenticatedUser, &groupInput); err != nil {
		t.Errorf("unexpected error by groupUsecase.DeleteGroupApprovedUser '%#v'", err)
	}
}

func Test_groupUsecase_StoreGroupApprovedUser(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	authenticatedUser := input.AuthenticatedUser{
		UserID: "userID1",
	}

	groupInput := input.Group{
		GroupID: 2,
	}

	gotOut, err := u.StoreGroupApprovedUser(&authenticatedUser, &groupInput)
	if err != nil {
		t.Errorf("unexpected error by groupUsecase.StoreGroupApprovedUser '%#v'", err)
	}

	wantOut := &output.ApprovedUser{
		GroupID:   2,
		UserID:    "userID4",
		UserName:  "userName4",
		ColorCode: "#8000FF",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_groupUsecase_DeleteGroupUnapprovedUser(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	authenticatedUser := input.AuthenticatedUser{
		UserID: "userID1",
	}

	groupInput := input.Group{
		GroupID: 2,
	}

	if err := u.DeleteGroupUnapprovedUser(&authenticatedUser, &groupInput); err != nil {
		t.Errorf("unexpected error by groupUsecase.DeleteGroupUnapprovedUser '%#v'", err)
	}
}

func Test_groupUsecase_FetchApprovedUserIDList(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	groupInput := input.Group{
		GroupID: 2,
	}

	gotOut, err := u.FetchApprovedUserIDList(&groupInput)
	if err != nil {
		t.Errorf("unexpected error by groupUsecase.FetchApprovedUserIDList '%#v'", err)
	}

	wantOut := &output.ApprovedUserIDList{
		"userID1",
		"userID2",
		"userID3",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_groupUsecase_VerifyGroupAffiliation(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	authenticatedUser := input.AuthenticatedUser{
		UserID: "userID1",
	}

	groupInput := input.Group{
		GroupID: 2,
	}

	if err := u.VerifyGroupAffiliation(&authenticatedUser, &groupInput); err != nil {
		t.Errorf("unexpected error by groupUsecase.VerifyGroupAffiliation '%#v'", err)
	}
}

func Test_groupUsecase_VerifyGroupAffiliationForUsersList(t *testing.T) {
	u := NewGroupUsecase(&mockGroupRepository{}, &mockGroupQueryService{}, &mockAccountApi{}, &mockUserRepository{})

	authenticatedUser := input.ApprovedUsersList{
		UserIDList: []string{
			"userID1",
			"userID2",
			"userID3",
		},
	}

	groupInput := input.Group{
		GroupID: 1,
	}

	if err := u.VerifyGroupAffiliationForUsersList(&authenticatedUser, &groupInput); err != nil {
		t.Errorf("unexpected error by groupUsecase.VerifyGroupAffiliationForUsersList '%#v'", err)
	}
}
