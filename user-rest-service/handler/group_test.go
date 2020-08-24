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
		{1, "group1", make([]model.ApprovedUser, 0), make([]model.UnapprovedUser, 0)},
		{2, "group2", make([]model.ApprovedUser, 0), make([]model.UnapprovedUser, 0)},
		{3, "group3", make([]model.ApprovedUser, 0), make([]model.UnapprovedUser, 0)},
	}, nil
}

func (t MockGroupRepository) GetUnApprovedGroupList(userID string) ([]model.UnapprovedGroup, error) {
	return []model.UnapprovedGroup{
		{4, "group4", make([]model.ApprovedUser, 0), make([]model.UnapprovedUser, 0)},
		{5, "group5", make([]model.ApprovedUser, 0), make([]model.UnapprovedUser, 0)},
	}, nil
}

func (t MockGroupRepository) GetApprovedUsersList(approvedGroupIDList []interface{}) ([]model.ApprovedUser, error) {
	return []model.ApprovedUser{
		{1, "userID1", "userName1"},
		{1, "userID2", "userName2"},
		{2, "userID1", "userName1"},
		{3, "userID1", "userName1"},
		{3, "userID2", "userName2"},
		{4, "userID2", "userName2"},
		{4, "userID4", "userName4"},
		{5, "userID4", "userName4"},
	}, nil
}

func (t MockGroupRepository) GetUnapprovedUsersList(unapprovedGroupIDList []interface{}) ([]model.UnapprovedUser, error) {
	return []model.UnapprovedUser{
		{1, "userID3", "userName3"},
		{2, "userID3", "userName3"},
		{2, "userID4", "userName4"},
		{4, "userID1", "userName1"},
		{4, "userID3", "userName3"},
		{5, "userID1", "userName1"},
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
