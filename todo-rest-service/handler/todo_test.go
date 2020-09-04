package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/testutil"

	"github.com/google/uuid"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/repository"
)

type MockTodoRepository struct {
	repository.TodoRepository
}

func (m MockTodoRepository) GetDailyImplementationTodoList(date time.Time, userID string) ([]model.Todo, error) {
	return []model.Todo{
		{ID: 3, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "電車定期券更新", CompleteFlag: true},
		{ID: 4, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false},
	}, nil
}

func (m MockTodoRepository) GetDailyDueTodoList(date time.Time, userID string) ([]model.Todo, error) {
	return []model.Todo{
		{ID: 2, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true},
		{ID: 3, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "電車定期券更新", CompleteFlag: true},
	}, nil
}

func (m MockTodoRepository) GetMonthlyImplementationTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error) {
	return []model.Todo{
		{ID: 1, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, TodoContent: "今月の予算を立てる", CompleteFlag: true},
		{ID: 2, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true},
		{ID: 3, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "電車定期券更新", CompleteFlag: true},
		{ID: 4, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false},
	}, nil
}

func (m MockTodoRepository) GetMonthlyDueTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error) {
	return []model.Todo{
		{ID: 1, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, TodoContent: "今月の予算を立てる", CompleteFlag: true},
		{ID: 2, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true},
		{ID: 3, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "電車定期券更新", CompleteFlag: true},
		{ID: 4, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false},
	}, nil
}

func (m MockTodoRepository) GetTodo(todoId int) (*model.Todo, error) {
	return &model.Todo{
		ID:                 1,
		PostedDate:         time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC),
		ImplementationDate: model.Date{Time: time.Date(2020, 7, 25, 0, 0, 0, 0, time.UTC)},
		DueDate:            model.Date{Time: time.Date(2020, 7, 30, 0, 0, 0, 0, time.UTC)},
		TodoContent:        "食器用洗剤2つ購入",
		CompleteFlag:       false,
	}, nil
}

func (m MockTodoRepository) PostTodo(todo *model.Todo, userID string) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockTodoRepository) PutTodo(todo *model.Todo, todoID int) error {
	return nil
}

func (m MockTodoRepository) DeleteTodo(todoID int) error {
	return nil
}

func (m MockTodoRepository) SearchTodoList(todoSqlQuery string) ([]model.Todo, error) {
	return []model.Todo{
		{ID: 1, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC)}, TodoContent: "今月の予算を立てる", CompleteFlag: true},
		{ID: 2, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true},
		{ID: 3, PostedDate: time.Date(2020, 9, 4, 17, 11, 0, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "電車定期券更新", CompleteFlag: true},
	}, nil
}

func TestDBHandler_GetDailyTodoList(t *testing.T) {
	h := DBHandler{
		AuthRepo: MockAuthRepository{},
		TodoRepo: MockTodoRepository{},
	}

	r := httptest.NewRequest("GET", "/todo-list/2020-07-10", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"date": "2020-07-10",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetDailyTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.TodoList{}, &model.TodoList{})
}

func TestDBHandler_GetMonthlyTodoList(t *testing.T) {
	h := DBHandler{
		AuthRepo: MockAuthRepository{},
		TodoRepo: MockTodoRepository{},
	}

	r := httptest.NewRequest("GET", "/todo-list/2020-07", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetMonthlyTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.TodoList{}, &model.TodoList{})
}

func TestDBHandler_PostTodo(t *testing.T) {
	h := DBHandler{
		AuthRepo: MockAuthRepository{},
		TodoRepo: MockTodoRepository{},
	}

	r := httptest.NewRequest("POST", "/todo-list", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostTodo(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.Todo{}, &model.Todo{})
}

func TestDBHandler_PutTodo(t *testing.T) {
	h := DBHandler{
		AuthRepo: MockAuthRepository{},
		TodoRepo: MockTodoRepository{},
	}

	r := httptest.NewRequest("PUT", "/todo-list/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutTodo(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.Todo{}, &model.Todo{})
}

func TestDBHandler_DeleteTodo(t *testing.T) {
	h := DBHandler{
		AuthRepo: MockAuthRepository{},
		TodoRepo: MockTodoRepository{},
	}

	r := httptest.NewRequest("DELETE", "/todo-list/1", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteTodo(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteTodoMsg{}, &DeleteTodoMsg{})
}

func TestDBHandler_SearchTodoList(t *testing.T) {
	h := DBHandler{
		AuthRepo: MockAuthRepository{},
		TodoRepo: MockTodoRepository{},
	}

	r := httptest.NewRequest("GET", "/todo-list/search", nil)
	w := httptest.NewRecorder()

	urlQuery := r.URL.Query()

	params := map[string]string{
		"date_type":     "implementation_date",
		"start_date":    "2020-07-05T00:00:00.0000",
		"end_date":      "2020-07-30T00:00:00.0000",
		"complete_flag": "true",
		"sort":          "due_date",
	}

	for k, v := range params {
		urlQuery.Add(k, v)
	}

	r.URL.RawQuery = urlQuery.Encode()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.SearchTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.SearchTodoList{}, &model.SearchTodoList{})
}
