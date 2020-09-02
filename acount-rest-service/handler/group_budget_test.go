package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/repository"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/testutil"

	"github.com/google/go-cmp/cmp"
)

type MockGroupBudgetsRepository struct {
	repository.GroupBudgetsRepository
}

func (m MockGroupBudgetsRepository) PostInitGroupStandardBudgets(groupID int) error {
	return nil
}

func (m MockGroupBudgetsRepository) GetGroupStandardBudgets(groupID int) (*model.GroupStandardBudgets, error) {
	return &model.GroupStandardBudgets{
		GroupStandardBudgets: []model.GroupStandardBudgetByCategory{
			{BigCategoryID: 2, BigCategoryName: "食費", Budget: 25000},
			{BigCategoryID: 3, BigCategoryName: "日用品", Budget: 5000},
			{BigCategoryID: 4, BigCategoryName: "趣味・娯楽", Budget: 4500},
			{BigCategoryID: 5, BigCategoryName: "交際費", Budget: 1000},
			{BigCategoryID: 6, BigCategoryName: "交通費", Budget: 1000},
			{BigCategoryID: 7, BigCategoryName: "衣服・美容", Budget: 0},
			{BigCategoryID: 8, BigCategoryName: "健康・医療", Budget: 4900},
			{BigCategoryID: 9, BigCategoryName: "通信費", Budget: 4400},
			{BigCategoryID: 10, BigCategoryName: "教養・教育", Budget: 10000},
			{BigCategoryID: 11, BigCategoryName: "住宅", Budget: 15000},
			{BigCategoryID: 12, BigCategoryName: "水道・光熱費", Budget: 3000},
			{BigCategoryID: 13, BigCategoryName: "自動車", Budget: 0},
			{BigCategoryID: 14, BigCategoryName: "保険", Budget: 9800},
			{BigCategoryID: 15, BigCategoryName: "税金・社会保険", Budget: 0},
			{BigCategoryID: 16, BigCategoryName: "現金・カード", Budget: 0},
			{BigCategoryID: 17, BigCategoryName: "その他", Budget: 0},
		},
	}, nil
}

func (m MockGroupBudgetsRepository) PutGroupStandardBudgets(groupStandardBudgets *model.GroupStandardBudgets, groupID int) error {
	return nil
}

func (m MockGroupBudgetsRepository) GetGroupCustomBudgets(yearMonth time.Time, groupID int) (*model.GroupCustomBudgets, error) {
	return &model.GroupCustomBudgets{
		GroupCustomBudgets: []model.GroupCustomBudgetByCategory{
			{BigCategoryID: 2, BigCategoryName: "食費", Budget: 30000},
			{BigCategoryID: 3, BigCategoryName: "日用品", Budget: 5000},
			{BigCategoryID: 4, BigCategoryName: "趣味・娯楽", Budget: 4500},
			{BigCategoryID: 5, BigCategoryName: "交際費", Budget: 1000},
			{BigCategoryID: 6, BigCategoryName: "交通費", Budget: 1000},
			{BigCategoryID: 7, BigCategoryName: "衣服・美容", Budget: 0},
			{BigCategoryID: 8, BigCategoryName: "健康・医療", Budget: 4900},
			{BigCategoryID: 9, BigCategoryName: "通信費", Budget: 4400},
			{BigCategoryID: 10, BigCategoryName: "教養・教育", Budget: 10000},
			{BigCategoryID: 11, BigCategoryName: "住宅", Budget: 15000},
			{BigCategoryID: 12, BigCategoryName: "水道・光熱費", Budget: 3000},
			{BigCategoryID: 13, BigCategoryName: "自動車", Budget: 0},
			{BigCategoryID: 14, BigCategoryName: "保険", Budget: 9800},
			{BigCategoryID: 15, BigCategoryName: "税金・社会保険", Budget: 0},
			{BigCategoryID: 16, BigCategoryName: "現金・カード", Budget: 0},
			{BigCategoryID: 17, BigCategoryName: "その他", Budget: 0},
		},
	}, nil
}

func (m MockGroupBudgetsRepository) PostGroupCustomBudgets(groupCustomBudgets *model.GroupCustomBudgets, yearMonth time.Time, groupID int) error {
	return nil
}

func (m MockGroupBudgetsRepository) PutGroupCustomBudgets(groupCustomBudgets *model.GroupCustomBudgets, yearMonth time.Time, groupID int) error {
	return nil
}

func TestDBHandler_PostInitGroupStandardBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		GroupBudgetsRepo: MockGroupBudgetsRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/standard-budgets", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	h.PostInitGroupStandardBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusCreated, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func TestDBHandler_GetGroupStandardBudgets(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		GroupBudgetsRepo: MockGroupBudgetsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/standard-budgets", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetGroupStandardBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupStandardBudgets{}, &model.GroupStandardBudgets{})
}

func TestDBHandler_PutGroupStandardBudgets(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		GroupBudgetsRepo: MockGroupBudgetsRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/standard-budgets", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroupStandardBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupStandardBudgets{}, &model.GroupStandardBudgets{})
}

func TestDBHandler_GetGroupCustomBudgets(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		GroupBudgetsRepo: MockGroupBudgetsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/custom-budgets/2020-07", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "1",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetGroupCustomBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupCustomBudgets{}, &model.GroupCustomBudgets{})
}

func TestDBHandler_PostGroupCustomBudgets(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		GroupBudgetsRepo: MockGroupBudgetsRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/custom-budgets/2020-07", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "1",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupCustomBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupCustomBudgets{}, &model.GroupCustomBudgets{})
}

func TestDBHandler_PutGroupCustomBudgets(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		GroupBudgetsRepo: MockGroupBudgetsRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/custom-budgets/2020-07", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "1",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroupCustomBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupCustomBudgets{}, &model.GroupCustomBudgets{})
}
