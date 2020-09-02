package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/testutil"

	"github.com/google/go-cmp/cmp"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/repository"
)

type MockBudgetsRepository struct {
	repository.BudgetsRepository
}

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
		AuthRepo:    MockAuthRepository{},
		BudgetsRepo: MockBudgetsRepository{},
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
