package handler

import "net/http"

func (h *DBHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if err := h.HealthRepo.PingMySQL(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := h.HealthRepo.PingRedis(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}
