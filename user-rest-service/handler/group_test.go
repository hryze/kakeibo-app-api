package handler

import (
	"database/sql"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/testutil"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
)

type MockGroupRepository struct {
	repository.GroupRepository
}

type MockSqlResult struct {
	sql.Result
}

func (r MockSqlResult) LastInsertId() (int64, error) {
	return 1, nil
}

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

func (t MockGroupRepository) PostGroupAndApprovedUser(group *model.Group, userID string) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (t MockGroupRepository) DeleteGroupAndApprovedUser(groupID int, userID string) error {
	return nil
}

func (t MockGroupRepository) GetGroup(groupID int) (*model.Group, error) {
	return &model.Group{
		GroupID:   1,
		GroupName: "group1",
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

func TestDBHandler_PostGroup(t *testing.T) {
	postInitGroupStandardBudgetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: postInitGroupStandardBudgetsHandler},
	}
	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:  MockAuthRepository{},
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("POST", "/groups", strings.NewReader(testutil.GetJsonFromTestData(t, "./testdata/group/post_group/request.json")))
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroup(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, "./testdata/group/post_group/response.json.golden")
}
