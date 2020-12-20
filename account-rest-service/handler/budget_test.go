package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/testutil"

	"github.com/google/go-cmp/cmp"
)

type MockBudgetsRepository struct{}

func (m MockBudgetsRepository) PostInitStandardBudgets(userID string) error {
	return nil
}

func (m MockBudgetsRepository) GetStandardBudgets(userID string) (*model.StandardBudgets, error) {
	return &model.StandardBudgets{
		StandardBudgets: []model.StandardBudgetByCategory{
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

func (m MockBudgetsRepository) PutStandardBudgets(standardBudgets *model.StandardBudgets, userID string) error {
	return nil
}

func (m MockBudgetsRepository) GetCustomBudgets(yearMonth time.Time, userID string) (*model.CustomBudgets, error) {
	return &model.CustomBudgets{
		CustomBudgets: []model.CustomBudgetByCategory{
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

func (m MockBudgetsRepository) PostCustomBudgets(customBudgets *model.CustomBudgets, yearMonth time.Time, userID string) error {
	return nil
}

func (m MockBudgetsRepository) PutCustomBudgets(customBudgets *model.CustomBudgets, yearMonth time.Time, userID string) error {
	return nil
}

func (m MockBudgetsRepository) DeleteCustomBudgets(yearMonth time.Time, userID string) error {
	return nil
}

func (m MockBudgetsRepository) GetMonthlyStandardBudget(userID string) (model.MonthlyBudget, error) {
	return model.MonthlyBudget{
		BudgetType:         "StandardBudget",
		MonthlyTotalBudget: 83600,
	}, nil
}

func (m MockBudgetsRepository) GetMonthlyCustomBudgets(year time.Time, userID string) ([]model.MonthlyBudget, error) {
	return []model.MonthlyBudget{
		{Month: model.Months{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)}, BudgetType: "CustomBudget", MonthlyTotalBudget: 88600},
		{Month: model.Months{Time: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC)}, BudgetType: "CustomBudget", MonthlyTotalBudget: 100000},
	}, nil
}

func TestDBHandler_PostInitStandardBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:    MockAuthRepository{},
		BudgetsRepo: MockBudgetsRepository{},
	}

	r := httptest.NewRequest("POST", "/standard-budgets", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	h.PostInitStandardBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusCreated, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func TestDBHandler_GetStandardBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		BudgetsRepo:      MockBudgetsRepository{},
		TransactionsRepo: MockTransactionsRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("GET", "/standard-budgets", nil)
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetStandardBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.StandardBudgets{}, &model.StandardBudgets{})
}

func TestDBHandler_PutStandardBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		BudgetsRepo:      MockBudgetsRepository{},
		TransactionsRepo: MockTransactionsRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("PUT", "/standard-budgets", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutStandardBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.StandardBudgets{}, &model.StandardBudgets{})
}

func TestDBHandler_GetCustomBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		BudgetsRepo:      MockBudgetsRepository{},
		TransactionsRepo: MockTransactionsRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("GET", "/custom-budgets/2020-07", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetCustomBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.CustomBudgets{}, &model.CustomBudgets{})
}

func TestDBHandler_PostCustomBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		BudgetsRepo:      MockBudgetsRepository{},
		TransactionsRepo: MockTransactionsRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("POST", "/custom-budgets/2020-07", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostCustomBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.CustomBudgets{}, &model.CustomBudgets{})
}

func TestDBHandler_PutCustomBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		BudgetsRepo:      MockBudgetsRepository{},
		TransactionsRepo: MockTransactionsRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("PUT", "/custom-budgets/2020-07", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutCustomBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.CustomBudgets{}, &model.CustomBudgets{})
}

func TestDBHandler_DeleteCustomBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:    MockAuthRepository{},
		BudgetsRepo: MockBudgetsRepository{},
	}

	r := httptest.NewRequest("DELETE", "/custom-budgets/2020-07", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteCustomBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_GetYearlyBudgets(t *testing.T) {
	h := DBHandler{
		AuthRepo:    MockAuthRepository{},
		BudgetsRepo: MockBudgetsRepository{},
	}

	r := httptest.NewRequest("GET", "/budgets/2020", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year": "2020",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetYearlyBudgets(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.YearlyBudget{}, &model.YearlyBudget{})
}
