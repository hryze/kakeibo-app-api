package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/repository"

	"github.com/google/go-cmp/cmp"
)

type MockGroupBudgetsRepository struct {
	repository.GroupBudgetsRepository
}

func (m MockGroupBudgetsRepository) PostInitGroupStandardBudgets(groupID int) error {
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
