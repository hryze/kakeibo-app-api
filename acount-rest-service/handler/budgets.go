package handler

import (
	"encoding/json"
	"net/http"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"
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

func (h *DBHandler) GetStandardBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	standardBudgets, err := h.DBRepo.GetStandardBudgets(userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&standardBudgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutStandardBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var standardBudgets model.StandardBudgets
	if err := json.NewDecoder(r.Body).Decode(&standardBudgets); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.DBRepo.PutStandardBudgets(&standardBudgets, userID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&standardBudgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
