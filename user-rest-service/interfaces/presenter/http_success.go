package presenter

import (
	"encoding/json"
	"net/http"
)

type successString struct {
	Message string `json:"message"`
}

func NewSuccessString(message string) *successString {
	return &successString{Message: message}
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
