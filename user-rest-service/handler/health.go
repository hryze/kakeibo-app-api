package handler

import (
	"log"
	"net/http"
)

func (h *DBHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if err := h.HealthRepo.PingMySQL(); err != nil {
		log.Printf("mysql ping error: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := h.HealthRepo.PingRedis(); err != nil {
		log.Printf("redis ping error: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}
