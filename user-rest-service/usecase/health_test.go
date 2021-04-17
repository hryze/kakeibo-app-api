package usecase

import (
	"testing"
)

type mockHealthRepository struct{}

func (r *mockHealthRepository) PingDataStore() error {
	return nil
}

func Test_healthUsecase_Readyz(t *testing.T) {
	u := NewHealthUsecase(&mockHealthRepository{}, &mockSessionStore{})

	if err := u.Readyz(); err != nil {
		t.Errorf("unexpected error by healthUsecase.Readyz '%#v'", err)
	}
}
