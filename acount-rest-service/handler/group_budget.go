package handler

import (
	"encoding/json"
	"net/http"
)

type GroupID struct {
	GroupID int `json:"group_id"`
}

func (h *DBHandler) PostInitGroupStandardBudgets(w http.ResponseWriter, r *http.Request) {
	var groupID GroupID
	if err := json.NewDecoder(r.Body).Decode(&groupID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.DBRepo.PostInitGroupStandardBudgets(groupID.GroupID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
