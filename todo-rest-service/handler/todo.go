package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
)

type NoContentMsg struct {
	Message string `json:"message"`
}

func (h *DBHandler) GetDailyTodoList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	date, err := time.Parse("2006-01-02", mux.Vars(r)["date"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	implementationTodoList, err := h.DBRepo.GetDailyImplementationTodoList(date, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dueTodoList, err := h.DBRepo.GetDailyDueTodoList(date, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(implementationTodoList) == 0 && len(dueTodoList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"今日実施予定todo、締切予定todoは登録されていません。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	todoList := model.NewTodoList(implementationTodoList, dueTodoList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&todoList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetMonthlyTodoList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	firstDay, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}
	lastDay := time.Date(firstDay.Year(), firstDay.Month()+1, 1, 0, 0, 0, 0, firstDay.Location()).Add(-1 * time.Second)

	implementationTodoList, err := h.DBRepo.GetMonthlyImplementationTodoList(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dueTodoList, err := h.DBRepo.GetMonthlyDueTodoList(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(implementationTodoList) == 0 && len(dueTodoList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"当月実施予定todoは登録されていません。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	todoList := model.NewTodoList(implementationTodoList, dueTodoList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&todoList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
