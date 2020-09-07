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
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/testutil"
)

var count int

type MockGroupTasksRepository struct{}

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

func (m MockGroupTasksRepository) GetGroupTasksUser(groupTasksUser model.GroupTasksUser, groupID int) (*model.GroupTasksUser, error) {
	if count == 0 {
		count++
		return nil, sql.ErrNoRows
	}

	count = 0

	return &model.GroupTasksUser{
		ID:        4,
		UserID:    "userID4",
		GroupID:   1,
		TasksList: make([]model.GroupTask, 0),
	}, nil
}

func (m MockGroupTasksRepository) PostGroupTasksUser(groupTasksUser model.GroupTasksUser, groupID int) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupTasksRepository) GetGroupTasksList(groupID int) ([]model.GroupTask, error) {
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
			ID:               3,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(2020, 9, 2, 0, 0, 0, 0, time.UTC), Valid: true}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "every", Valid: true}},
			Cycle:            model.NullInt{Int: 7, Valid: true},
			TaskName:         "トイレ掃除",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 4, Valid: true},
		},
		{
			ID:               4,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), Valid: false}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Cycle:            model.NullInt{Int: 0, Valid: false},
			TaskName:         "台所掃除",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 0, Valid: false},
		},
		{
			ID:               5,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(2020, 8, 31, 0, 0, 0, 0, time.UTC), Valid: true}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "consecutive", Valid: true}},
			Cycle:            model.NullInt{Int: 7, Valid: true},
			TaskName:         "風呂掃除",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 3, Valid: true},
		},
	}, nil
}

func (m MockGroupTasksRepository) GetGroupTask(groupTasksID int) (*model.GroupTask, error) {
	if groupTasksID == 1 {
		return &model.GroupTask{
			ID:               1,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), Valid: false}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Cycle:            model.NullInt{Int: 0, Valid: false},
			TaskName:         "料理",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 0, Valid: false},
		}, nil
	}

	if count == 0 {
		count++
		return &model.GroupTask{
			ID:               2,
			BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), Valid: false}},
			CycleType:        model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Cycle:            model.NullInt{Int: 0, Valid: false},
			TaskName:         "洗濯",
			GroupID:          1,
			GroupTasksUserID: model.NullInt{Int: 0, Valid: false},
		}, nil
	}

	count = 0

	return &model.GroupTask{
		ID:               2,
		BaseDate:         model.NullTime{NullTime: sql.NullTime{Time: time.Date(2020, 9, 3, 0, 0, 0, 0, time.UTC), Valid: true}},
		CycleType:        model.NullString{NullString: sql.NullString{String: "every", Valid: true}},
		Cycle:            model.NullInt{Int: 3, Valid: true},
		TaskName:         "洗濯",
		GroupID:          1,
		GroupTasksUserID: model.NullInt{Int: 1, Valid: true},
	}, nil
}

func (m MockGroupTasksRepository) PostGroupTask(groupTask model.GroupTask, groupID int) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupTasksRepository) PutGroupTask(groupTask *model.GroupTask, groupTasksID int) error {
	return nil
}

func (m MockGroupTasksRepository) DeleteGroupTask(groupTasksID int) error {
	return nil
}

type MockTime struct{}

func (m MockTime) Now() time.Time {
	return time.Date(2020, 9, 6, 0, 0, 0, 0, time.UTC)
}

func TestDBHandler_GetGroupTasksListForEachUser(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:       MockAuthRepository{},
		GroupTasksRepo: MockGroupTasksRepository{},
		TimeManage:     MockTime{},
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

func TestDBHandler_PostGroupTasksUser(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:       MockAuthRepository{},
		GroupTasksRepo: MockGroupTasksRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/tasks/users", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupTasksUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupTasksUser{}, &model.GroupTasksUser{})
}

func TestDBHandler_GetGroupTasksList(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:       MockAuthRepository{},
		GroupTasksRepo: MockGroupTasksRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/tasks", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetGroupTasksList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTasksList{}, &model.GroupTasksList{})
}

func TestDBHandler_PostGroupTask(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:       MockAuthRepository{},
		GroupTasksRepo: MockGroupTasksRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/tasks", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupTask(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupTask{}, &model.GroupTask{})
}

func TestDBHandler_PutGroupTask(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:       MockAuthRepository{},
		GroupTasksRepo: MockGroupTasksRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/tasks/2", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "2",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroupTask(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTask{}, &model.GroupTask{})
}

func TestDBHandler_DeleteGroupTask(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:       MockAuthRepository{},
		GroupTasksRepo: MockGroupTasksRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/1/tasks/2", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteGroupTask(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteGroupTaskMsg{}, &DeleteGroupTaskMsg{})
}
