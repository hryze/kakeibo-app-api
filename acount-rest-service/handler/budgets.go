package handler

import (
	"encoding/json"
	"net/http"
)

type UserID struct {
	UserID string `json:"user_id"`
}

func (h *DBHandler) PostInitStandardBudgets(w http.ResponseWriter, r *http.Request) {
	var userID UserID
	if err := json.NewDecoder(r.Body).Decode(&userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.DBRepo.PostInitStandardBudgets(userID.UserID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
