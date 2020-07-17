package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"
)

func (h *DBHandler) PostTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var transactionReceiver model.TransactionReceiver
	if err := json.NewDecoder(r.Body).Decode(&transactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	result, err := h.DBRepo.PostTransaction(&transactionReceiver, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	lastInsertId, err := result.LastInsertId()

	var transactionSender model.TransactionSender
	dbTransactionSender, err := h.DBRepo.GetTransaction(&transactionSender, int(lastInsertId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dbTransactionSender); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
