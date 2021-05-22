package usecase

import (
	"github.com/hryze/kakeibo-app-api/user-rest-service/domain/healthdomain"
	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase/sessionstore"
)

type HealthUsecase interface {
	Readyz() error
}

type healthUsecase struct {
	healthRepository healthdomain.Repository
	sessionStore     sessionstore.SessionStore
}

func NewHealthUsecase(healthRepository healthdomain.Repository, sessionStore sessionstore.SessionStore) *healthUsecase {
	return &healthUsecase{
		healthRepository: healthRepository,
		sessionStore:     sessionStore,
	}
}

func (u *healthUsecase) Readyz() error {
	if err := u.healthRepository.PingDataStore(); err != nil {
		return err
	}

	if err := u.sessionStore.PingSessionStore(); err != nil {
		return err
	}

	return nil
}
