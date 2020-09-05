package handler

import (
	"database/sql"
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

type MockGroupTasksRepository struct {
	repository.GroupTasksRepository
}

func (m MockGroupTasksRepository) GetGroupTasksUsersList(groupID int) ([]model.GroupTasksUser, error) {
	return []model.GroupTasksUser{
		{ID: 1, UserID: "userID1", GroupID: 1, TasksList: make([]model.GroupTask, 0)},
		{ID: 2, UserID: "userID2", GroupID: 1, TasksList: make([]model.GroupTask, 0)},
		{ID: 3, UserID: "userID3", GroupID: 1, TasksList: make([]model.GroupTask, 0)},
	}, nil
}

func (m MockGroupTasksRepository) GetGroupTasksListAssignedToUser(groupID int) ([]model.GroupTask, error) {
	return []model.GroupTask{
		{
			ID:               1,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(2020, 9, 5, 0, 0, 0, 0, time.UTC), Valid: true}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "every", Valid: true}},
			Cycle:            model.NullInt{Int: 1, Valid: true},
			TaskName:         "料理",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 2, Valid: true},
		},
		{
			ID:               2,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(2020, 9, 3, 0, 0, 0, 0, time.UTC), Valid: true}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "every", Valid: true}},
			Cycle:            model.NullInt{Int: 3, Valid: true},
			TaskName:         "洗濯",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 1, Valid: true},
		},
		{
			ID:               5,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(2020, 8, 31, 0, 0, 0, 0, time.UTC), Valid: true}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "consecutive", Valid: true}},
			Cycle:            model.NullInt{Int: 7, Valid: true},
			TaskName:         "風呂掃除",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 2, Valid: true},
		},
	}, nil
}

func TestDBHandler_GetGroupTasksListForEachUser(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:       MockAuthRepository{},
		GroupTasksRepo: MockGroupTasksRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/tasks/users", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetGroupTasksListForEachUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTasksListForEachUser{}, &model.GroupTasksListForEachUser{})
}
