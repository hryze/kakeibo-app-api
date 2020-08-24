package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/testutil"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
)

type MockGroupRepository struct{ repository.GroupRepository }

func (t MockGroupRepository) GetApprovedGroupList(userID string) ([]model.ApprovedGroup, error) {
	return []model.ApprovedGroup{
		{GroupID: 1, GroupName: "group1", ApprovedUsersList: make([]model.ApprovedUser, 0), UnapprovedUsersList: make([]model.UnapprovedUser, 0)},
		{GroupID: 2, GroupName: "group2", ApprovedUsersList: make([]model.ApprovedUser, 0), UnapprovedUsersList: make([]model.UnapprovedUser, 0)},
		{GroupID: 3, GroupName: "group3", ApprovedUsersList: make([]model.ApprovedUser, 0), UnapprovedUsersList: make([]model.UnapprovedUser, 0)},
	}, nil
}

func (t MockGroupRepository) GetUnApprovedGroupList(userID string) ([]model.UnapprovedGroup, error) {
	return []model.UnapprovedGroup{
		{GroupID: 4, GroupName: "group4", ApprovedUsersList: make([]model.ApprovedUser, 0), UnapprovedUsersList: make([]model.UnapprovedUser, 0)},
		{GroupID: 5, GroupName: "group5", ApprovedUsersList: make([]model.ApprovedUser, 0), UnapprovedUsersList: make([]model.UnapprovedUser, 0)},
	}, nil
}

func (t MockGroupRepository) GetApprovedUsersList(approvedGroupIDList []interface{}) ([]model.ApprovedUser, error) {
	return []model.ApprovedUser{
		{GroupID: 1, UserID: "userID1", UserName: "userName1"},
		{GroupID: 1, UserID: "userID2", UserName: "userName2"},
		{GroupID: 2, UserID: "userID1", UserName: "userName1"},
		{GroupID: 3, UserID: "userID1", UserName: "userName1"},
		{GroupID: 3, UserID: "userID2", UserName: "userName2"},
		{GroupID: 4, UserID: "userID2", UserName: "userName2"},
		{GroupID: 4, UserID: "userID4", UserName: "userName4"},
		{GroupID: 5, UserID: "userID4", UserName: "userName4"},
	}, nil
}

func (t MockGroupRepository) GetUnapprovedUsersList(unapprovedGroupIDList []interface{}) ([]model.UnapprovedUser, error) {
	return []model.UnapprovedUser{
		{GroupID: 1, UserID: "userID3", UserName: "userName3"},
		{GroupID: 2, UserID: "userID3", UserName: "userName3"},
		{GroupID: 2, UserID: "userID4", UserName: "userName4"},
		{GroupID: 4, UserID: "userID1", UserName: "userName1"},
		{GroupID: 4, UserID: "userID3", UserName: "userName3"},
		{GroupID: 5, UserID: "userID1", UserName: "userName1"},
	}, nil
}

func TestDBHandler_GetGroupList(t *testing.T) {
	h := DBHandler{
		AuthRepo:  MockAuthRepository{},
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("GET", "/groups", nil)
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetGroupList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, "./testdata/group/get_group_list/response.json.golden")
}
