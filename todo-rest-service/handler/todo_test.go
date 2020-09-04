package handler

import (
	"net/http"
	"net/http/httptest"
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
