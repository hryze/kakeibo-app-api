package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"

	"github.com/hryze/kakeibo-app-api/user-rest-service/appcontext"
	"github.com/hryze/kakeibo-app-api/user-rest-service/interfaces/presenter"
	"github.com/hryze/kakeibo-app-api/user-rest-service/testutil"
	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase/output"
)

type mockGroupUsecase struct{}

func (u *mockGroupUsecase) FetchGroupList(in *input.AuthenticatedUser) (*output.GroupList, error) {
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

func (u *mockGroupUsecase) StoreGroup(authenticatedUser *input.AuthenticatedUser, group *input.Group) (*output.Group, error) {
	return &output.Group{
		GroupID:   1,
		GroupName: "group1",
	}, nil
}

func (u *mockGroupUsecase) UpdateGroupName(groupInput *input.Group) (*output.Group, error) {
	return &output.Group{
		GroupID:   1,
		GroupName: "group1",
	}, nil
}

func (u *mockGroupUsecase) StoreGroupUnapprovedUser(unapprovedUser *input.UnapprovedUser, group *input.Group) (*output.UnapprovedUser, error) {
	return &output.UnapprovedUser{
		GroupID:  1,
		UserID:   "userID1",
		UserName: "userName1",
	}, nil
}

func (u *mockGroupUsecase) DeleteGroupApprovedUser(authenticatedUser *input.AuthenticatedUser, group *input.Group) error {
	return nil
}

func (u *mockGroupUsecase) StoreGroupApprovedUser(authenticatedUser *input.AuthenticatedUser, group *input.Group) (*output.ApprovedUser, error) {
	return &output.ApprovedUser{
		GroupID:   1,
		UserID:    "userID1",
		UserName:  "userName1",
		ColorCode: "#FF0000",
	}, nil
}

func (u *mockGroupUsecase) DeleteGroupUnapprovedUser(authenticatedUser *input.AuthenticatedUser, group *input.Group) error {
	return nil
}

func (u *mockGroupUsecase) FetchApprovedUserIDList(group *input.Group) (*output.ApprovedUserIDList, error) {
	return &output.ApprovedUserIDList{
		"userID1",
		"userID2",
		"userID3",
	}, nil
}

func (u *mockGroupUsecase) VerifyGroupAffiliation(authenticatedUser *input.AuthenticatedUser, group *input.Group) error {
	return nil
}

func (u *mockGroupUsecase) VerifyGroupAffiliationForUsersList(approvedUsersList *input.ApprovedUsersList, group *input.Group) error {
	return nil
}

func Test_groupHandler_FetchGroupList(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodGet, "/groups", nil)
	w := httptest.NewRecorder()

	ctx := appcontext.SetUserID(r.Context(), "userID1")

	h.FetchGroupList(w, r.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &output.GroupList{}, &output.GroupList{})
}

func Test_groupHandler_StoreGroup(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	ctx := appcontext.SetUserID(r.Context(), "userID1")

	h.StoreGroup(w, r.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &output.Group{}, &output.Group{})
}

func Test_groupHandler_UpdateGroupName(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodPut, "/groups/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	ctx := appcontext.SetUserID(r.Context(), "userID1")

	h.UpdateGroupName(w, r.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &output.Group{}, &output.Group{})
}

func Test_groupHandler_StoreGroupUnapprovedUser(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodPost, "/groups/1/users", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	ctx := appcontext.SetUserID(r.Context(), "userID1")

	h.StoreGroupUnapprovedUser(w, r.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &output.UnapprovedUser{}, &output.UnapprovedUser{})
}

func Test_groupHandler_DeleteGroupApprovedUser(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodDelete, "/groups/1/users", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	ctx := appcontext.SetUserID(r.Context(), "userID1")

	h.DeleteGroupApprovedUser(w, r.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, presenter.NewSuccessString(""), presenter.NewSuccessString(""))
}

func Test_groupHandler_StoreGroupApprovedUser(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodPost, "/groups/1/users/approved", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	ctx := appcontext.SetUserID(r.Context(), "userID1")

	h.StoreGroupApprovedUser(w, r.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &output.ApprovedUser{}, &output.ApprovedUser{})
}

func Test_groupHandler_DeleteGroupUnapprovedUser(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodDelete, "/groups/1/users/unapproved", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	ctx := appcontext.SetUserID(r.Context(), "userID1")

	h.DeleteGroupUnapprovedUser(w, r.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, presenter.NewSuccessString(""), presenter.NewSuccessString(""))
}

func Test_groupHandler_FetchApprovedUserIDList(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodGet, "/groups/1/users", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	h.FetchApprovedUserIDList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &output.ApprovedUserIDList{}, &output.ApprovedUserIDList{})
}

func Test_groupHandler_VerifyGroupAffiliation(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodGet, "/groups/1/users/userID1/verify", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"user_id":  "userID1",
	})

	h.VerifyGroupAffiliation(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusOK, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_groupHandler_VerifyGroupAffiliationForUsersList(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodGet, "/groups/1/users/verify", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	h.VerifyGroupAffiliationForUsersList(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusOK, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
