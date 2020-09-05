package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/repository"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/testutil"
)

type MockGroupTodoRepository struct {
	repository.GroupTodoRepository
}

func (m MockGroupTodoRepository) GetDailyImplementationGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	return []model.GroupTodo{
		{ID: 3, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)}, TodoContent: "醤油購入", CompleteFlag: false, UserID: "userID1"},
	}, nil
}

func (m MockGroupTodoRepository) GetDailyDueGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	return []model.GroupTodo{
		{ID: 2, PostedDate: time.Date(2020, 9, 5, 1, 29, 8, 0, time.UTC), ImplementationDate: model.Date{Time: time.Date(2020, 7, 9, 0, 0, 0, 0, time.UTC)}, DueDate: model.Date{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)}, TodoContent: "コストコ鶏肉セール 5パック購入", CompleteFlag: true, UserID: "userID2"},
	}, nil
}

func TestDBHandler_GetDailyGroupTodoList(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

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
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetDailyGroupTodoList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTodoList{}, &model.GroupTodoList{})
}
