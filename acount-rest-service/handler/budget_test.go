package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
