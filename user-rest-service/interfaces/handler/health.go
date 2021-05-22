package handler

import (
	"net/http"

	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase"
)

type healthHandler struct {
	healthUsecase usecase.HealthUsecase
}

func NewHealthHandler(healthUsecase usecase.HealthUsecase) *healthHandler {
	return &healthHandler{
		healthUsecase: healthUsecase,
	}
}

func (h *healthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if err := h.healthUsecase.Readyz(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
