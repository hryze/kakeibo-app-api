package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/testutil"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
)

type MockGroupRepository struct{}

type MockUserRepositoryForGroup struct {
	repository.UserRepository
}

type MockSqlResult struct {
	sql.Result
}

func (t MockUserRepositoryForGroup) FindSignUpUserByUserID(userID string) (*model.SignUpUser, error) {
	signUpUser := &model.SignUpUser{
		ID:       "testID",
		Name:     "testName",
		Email:    "test@icloud.com",
		Password: "$2a$10$teJL.9I0QfBESpaBIwlbl.VkivuHEOKhy674CW6J.4k3AnfEpcYLy",
	}

	return signUpUser, nil
}

func (r MockSqlResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (t MockGroupRepository) GetGroup(groupID int) (*model.Group, error) {
	return &model.Group{
		GroupID:   1,
		GroupName: "group1",
	}, nil
}

func (t MockGroupRepository) PutGroup(group *model.Group, groupID int) error {
	return nil
}

func (t MockGroupRepository) PostUnapprovedUser(unapprovedUser *model.UnapprovedUser, groupID int) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (t MockGroupRepository) GetUnapprovedUser(groupUnapprovedUsersID int) (*model.UnapprovedUser, error) {
	return &model.UnapprovedUser{
		GroupID:  1,
		UserID:   "userID2",
		UserName: "userName2",
	}, nil
}

func (t MockGroupRepository) FindApprovedUser(groupID int, userID string) error {
	if groupID == 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (t MockGroupRepository) FindUnapprovedUser(groupID int, userID string) error {
	if groupID == 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (t MockGroupRepository) PostGroupApprovedUserAndDeleteGroupUnapprovedUser(groupID int, userID string, colorCode string) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (t MockGroupRepository) GetApprovedUser(approvedUsersID int) (*model.ApprovedUser, error) {
	return &model.ApprovedUser{
		GroupID:   2,
		UserID:    "userID1",
		UserName:  "userName1",
		ColorCode: "#FF0000",
	}, nil
}

func (t MockGroupRepository) DeleteGroupApprovedUser(groupID int, userID string) error {
	return nil
}

func (t MockGroupRepository) DeleteGroupUnapprovedUser(groupID int, userID string) error {
	return nil
}

func (t MockGroupRepository) FindApprovedUsersList(groupID int, groupUsersList []string) (model.GroupTasksUsersListReceiver, error) {
	return model.GroupTasksUsersListReceiver{
		GroupUsersList: []string{
			"userID4",
			"userID5",
			"userID6",
		},
	}, nil
}

func (t MockGroupRepository) GetGroupUsersList(groupID int) ([]string, error) {
	return []string{"userID1", "userID4", "userID5", "userID3", "userID2"}, nil
}

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

func Test_groupHandler_FetchGroupList(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodGet, "/groups", nil)
	w := httptest.NewRecorder()

	context.Set(r, config.Env.RequestCtx.UserID, "userID1")

	h.FetchGroupList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &output.GroupList{}, &output.GroupList{})
}

func Test_groupHandler_StoreGroup(t *testing.T) {
	h := NewGroupHandler(&mockGroupUsecase{})

	r := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	context.Set(r, config.Env.RequestCtx.UserID, "userID1")

	h.StoreGroup(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &output.Group{}, &output.Group{})
}

func TestDBHandler_PutGroup(t *testing.T) {
	h := DBHandler{
		AuthRepo:  MockAuthRepository{},
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroup(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.Group{}, &model.Group{})
}

func TestDBHandler_PostGroupUnapprovedUser(t *testing.T) {
	h := DBHandler{
		AuthRepo:  MockAuthRepository{},
		UserRepo:  MockUserRepositoryForGroup{},
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/users", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupUnapprovedUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.UnapprovedUser{}, &model.UnapprovedUser{})
}

func TestDBHandler_DeleteGroupApprovedUser(t *testing.T) {
	h := DBHandler{
		AuthRepo:  MockAuthRepository{},
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/2/users", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "2",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteGroupApprovedUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_PostGroupApprovedUser(t *testing.T) {
	h := DBHandler{
		AuthRepo:  MockAuthRepository{},
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/2/users/approved", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "2",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupApprovedUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.ApprovedUser{}, &model.ApprovedUser{})
}

func TestDBHandler_DeleteGroupUnapprovedUser(t *testing.T) {
	h := DBHandler{
		AuthRepo:  MockAuthRepository{},
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/2/users/unapproved", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "2",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteGroupUnapprovedUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_VerifyGroupAffiliation(t *testing.T) {
	h := DBHandler{
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/2/users/userID1/verify", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "2",
		"user_id":  "userID1",
	})

	h.VerifyGroupAffiliation(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusOK, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func TestDBHandler_VerifyGroupAffiliationOfUsersList(t *testing.T) {
	h := DBHandler{
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/2/users/verify", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "2",
	})

	h.VerifyGroupAffiliationOfUsersList(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusOK, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func TestDBHandler_GetGroupUserIDList(t *testing.T) {
	h := DBHandler{
		GroupRepo: MockGroupRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/2/users", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "2",
	})

	h.GetGroupUserIDList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &[]string{}, &[]string{})
}
