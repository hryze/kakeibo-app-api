package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/hryze/kakeibo-app-api/todo-rest-service/config"
	"github.com/hryze/kakeibo-app-api/todo-rest-service/domain/model"
	"github.com/hryze/kakeibo-app-api/todo-rest-service/testutil"
)

type MockGroupTodoRepository struct{}

func (m MockGroupTodoRepository) GetDailyImplementationGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	return []model.GroupTodo{
		{ID: 3, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false, UserID: "userID1"},
	}, nil
}

func (m MockGroupTodoRepository) GetDailyDueGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	return []model.GroupTodo{
		{ID: 2, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true, UserID: "userID2"},
	}, nil
}

func (m MockGroupTodoRepository) GetMonthlyImplementationGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error) {
	return []model.GroupTodo{
		{ID: 1, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, TodoContent: "今月の予算を立てる", CompleteFlag: true, UserID: "userID1"},
		{ID: 2, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true, UserID: "userID2"},
		{ID: 3, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false, UserID: "userID1"},
	}, nil
}

func (m MockGroupTodoRepository) GetMonthlyDueGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error) {
	return []model.GroupTodo{
		{ID: 1, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, TodoContent: "今月の予算を立てる", CompleteFlag: true, UserID: "userID1"},
		{ID: 2, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true, UserID: "userID2"},
		{ID: 3, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false, UserID: "userID1"},
	}, nil
}

func (m MockGroupTodoRepository) GetExpiredGroupTodoList(dueDate time.Time, groupID int) (*model.ExpiredGroupTodoList, error) {
	return &model.ExpiredGroupTodoList{
		ExpiredGroupTodoList: []model.GroupTodo{
			{ID: 1, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, TodoContent: "今月の予算を立てる", CompleteFlag: false, UserID: "userID1"},
			{ID: 3, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false, UserID: "userID2"},
		},
	}, nil
}

func (m MockGroupTodoRepository) GetGroupTodo(groupTodoId int) (*model.GroupTodo, error) {
	return &model.GroupTodo{
		ID:                 1,
		PostedDate:         time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC),
		UpdatedDate:        time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC),
		ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)},
		DueDate:            model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)},
		TodoContent:        "今月の予算を立てる",
		CompleteFlag:       false,
		UserID:             "userID1",
	}, nil
}

func (m MockGroupTodoRepository) PostGroupTodo(groupTodo *model.GroupTodo, userID string, groupID int) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupTodoRepository) PutGroupTodo(groupTodo *model.GroupTodo, groupTodoID int) error {
	return nil
}

func (m MockGroupTodoRepository) DeleteGroupTodo(groupTodoID int) error {
	return nil
}

func (m MockGroupTodoRepository) SearchGroupTodoList(groupTodoSqlQuery string) ([]model.GroupTodo, error) {
	return []model.GroupTodo{
		{ID: 1, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, TodoContent: "今月の予算を立てる", CompleteFlag: true, UserID: "userID1"},
		{ID: 2, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), UpdatedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true, UserID: "userID2"},
	}, nil
}

func TestDBHandler_GetDailyGroupTodoList(t *testing.T) {
	h := DBHandler{
		AuthRepo:      MockAuthRepository{},
		GroupTodoRepo: MockGroupTodoRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/todo-list/2020-07-10", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"date":     "2020-07-10",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetDailyGroupTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTodoList{}, &model.GroupTodoList{})
}

func TestDBHandler_GetMonthlyGroupTodoList(t *testing.T) {
	h := DBHandler{
		AuthRepo:      MockAuthRepository{},
		GroupTodoRepo: MockGroupTodoRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/todo-list/2020-07", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "1",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetMonthlyGroupTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTodoList{}, &model.GroupTodoList{})
}

func TestDBHandler_GetExpiredGroupTodoList(t *testing.T) {
	h := DBHandler{
		AuthRepo:      MockAuthRepository{},
		GroupTodoRepo: MockGroupTodoRepository{},
		TimeManage:    MockTime{},
	}

	r := httptest.NewRequest("GET", "/groups/1/todo-list/expired", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetExpiredGroupTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ExpiredGroupTodoList{}, &model.ExpiredGroupTodoList{})
}

func TestDBHandler_PostGroupTodo(t *testing.T) {
	h := DBHandler{
		AuthRepo:      MockAuthRepository{},
		GroupTodoRepo: MockGroupTodoRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/todo-list", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupTodo(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupTodo{}, &model.GroupTodo{})
}

func TestDBHandler_PutGroupTodo(t *testing.T) {
	h := DBHandler{
		AuthRepo:      MockAuthRepository{},
		GroupTodoRepo: MockGroupTodoRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/todo-list/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "1",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroupTodo(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTodo{}, &model.GroupTodo{})
}

func TestDBHandler_DeleteGroupTodo(t *testing.T) {
	h := DBHandler{
		AuthRepo:      MockAuthRepository{},
		GroupTodoRepo: MockGroupTodoRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/1/todo-list/1", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "1",
	})

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteGroupTodo(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_SearchGroupTodoList(t *testing.T) {
	h := DBHandler{
		AuthRepo:      MockAuthRepository{},
		GroupTodoRepo: MockGroupTodoRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/todo-list/search", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	urlQuery := r.URL.Query()

	params := map[string]string{
		"date_type":     "due_date",
		"start_date":    "2020-07-05T00:00:00.0000",
		"end_date":      "2020-07-10T00:00:00.0000",
		"complete_flag": "true",
		"sort":          "due_date",
	}

	for k, v := range params {
		urlQuery.Add(k, v)
	}

	r.URL.RawQuery = urlQuery.Encode()

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.SearchGroupTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.SearchGroupTodoList{}, &model.SearchGroupTodoList{})
}
