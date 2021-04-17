package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockHealthUsecase struct{}

func (u *mockHealthUsecase) Readyz() error {
	return nil
}

func Test_healthHandler_Readyz(t *testing.T) {
	h := NewHealthHandler(&mockHealthUsecase{})

	r := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	h.Readyz(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusOK, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
