package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type MockHealthRepository struct{}

func (m MockHealthRepository) PingMySQL() error {
	return nil
}

func (m MockHealthRepository) PingRedis() error {
	return nil
}

func TestDBHandler_Readyz(t *testing.T) {
	h := DBHandler{
		HealthRepo: MockHealthRepository{},
	}

	r := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()

	h.Readyz(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusOK, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
