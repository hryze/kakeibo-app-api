package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

func (h *DBHandler) PostInitGroupStandardBudgets(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.GroupBudgetsRepo.PostInitGroupStandardBudgets(groupID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *DBHandler) GetGroupStandardBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	groupStandardBudgets, err := h.GroupBudgetsRepo.GetGroupStandardBudgets(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	firstDayOfLastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	lastDayOfLastMonth := firstDayOfLastMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	groupTransactionTotalAmountByBigCategoryList, err := h.GroupTransactionsRepo.GetMonthlyGroupTransactionTotalAmountByBigCategory(groupID, firstDayOfLastMonth, lastDayOfLastMonth)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for _, groupTransactionTotalAmountByBigCategory := range groupTransactionTotalAmountByBigCategoryList {
		for i, groupStandardBudgetByCategory := range groupStandardBudgets.GroupStandardBudgets {
			if groupTransactionTotalAmountByBigCategory.BigCategoryID == groupStandardBudgetByCategory.BigCategoryID {
				groupStandardBudgets.GroupStandardBudgets[i].LastMonthExpenses = groupTransactionTotalAmountByBigCategory.TotalAmount

				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupStandardBudgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutGroupStandardBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	var groupStandardBudgets model.GroupStandardBudgets
	if err := json.NewDecoder(r.Body).Decode(&groupStandardBudgets); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateBudgets(groupStandardBudgets); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.GroupBudgetsRepo.PutGroupStandardBudgets(&groupStandardBudgets, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupStandardBudgets, err := h.GroupBudgetsRepo.GetGroupStandardBudgets(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	firstDayOfLastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	lastDayOfLastMonth := firstDayOfLastMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	groupTransactionTotalAmountByBigCategoryList, err := h.GroupTransactionsRepo.GetMonthlyGroupTransactionTotalAmountByBigCategory(groupID, firstDayOfLastMonth, lastDayOfLastMonth)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for _, groupTransactionTotalAmountByBigCategory := range groupTransactionTotalAmountByBigCategoryList {
		for i, groupStandardBudgetByCategory := range dbGroupStandardBudgets.GroupStandardBudgets {
			if groupTransactionTotalAmountByBigCategory.BigCategoryID == groupStandardBudgetByCategory.BigCategoryID {
				dbGroupStandardBudgets.GroupStandardBudgets[i].LastMonthExpenses = groupTransactionTotalAmountByBigCategory.TotalAmount

				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&dbGroupStandardBudgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetGroupCustomBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	yearMonth, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	dbGroupCustomBudgets, err := h.GroupBudgetsRepo.GetGroupCustomBudgets(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	firstDayOfLastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	lastDayOfLastMonth := firstDayOfLastMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	groupTransactionTotalAmountByBigCategoryList, err := h.GroupTransactionsRepo.GetMonthlyGroupTransactionTotalAmountByBigCategory(groupID, firstDayOfLastMonth, lastDayOfLastMonth)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for _, groupTransactionTotalAmountByBigCategory := range groupTransactionTotalAmountByBigCategoryList {
		for i, groupStandardBudgetByCategory := range dbGroupCustomBudgets.GroupCustomBudgets {
			if groupTransactionTotalAmountByBigCategory.BigCategoryID == groupStandardBudgetByCategory.BigCategoryID {
				dbGroupCustomBudgets.GroupCustomBudgets[i].LastMonthExpenses = groupTransactionTotalAmountByBigCategory.TotalAmount

				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&dbGroupCustomBudgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroupCustomBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	yearMonth, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	var groupCustomBudgets model.GroupCustomBudgets
	if err := json.NewDecoder(r.Body).Decode(&groupCustomBudgets); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateBudgets(groupCustomBudgets); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.GroupBudgetsRepo.PostGroupCustomBudgets(&groupCustomBudgets, yearMonth, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupCustomBudgets, err := h.GroupBudgetsRepo.GetGroupCustomBudgets(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	firstDayOfLastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	lastDayOfLastMonth := firstDayOfLastMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	groupTransactionTotalAmountByBigCategoryList, err := h.GroupTransactionsRepo.GetMonthlyGroupTransactionTotalAmountByBigCategory(groupID, firstDayOfLastMonth, lastDayOfLastMonth)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for _, groupTransactionTotalAmountByBigCategory := range groupTransactionTotalAmountByBigCategoryList {
		for i, groupStandardBudgetByCategory := range dbGroupCustomBudgets.GroupCustomBudgets {
			if groupTransactionTotalAmountByBigCategory.BigCategoryID == groupStandardBudgetByCategory.BigCategoryID {
				dbGroupCustomBudgets.GroupCustomBudgets[i].LastMonthExpenses = groupTransactionTotalAmountByBigCategory.TotalAmount

				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&dbGroupCustomBudgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutGroupCustomBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	yearMonth, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	var groupCustomBudgets model.GroupCustomBudgets
	if err := json.NewDecoder(r.Body).Decode(&groupCustomBudgets); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateBudgets(groupCustomBudgets); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.GroupBudgetsRepo.PutGroupCustomBudgets(&groupCustomBudgets, yearMonth, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupCustomBudgets, err := h.GroupBudgetsRepo.GetGroupCustomBudgets(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	firstDayOfLastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	lastDayOfLastMonth := firstDayOfLastMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	groupTransactionTotalAmountByBigCategoryList, err := h.GroupTransactionsRepo.GetMonthlyGroupTransactionTotalAmountByBigCategory(groupID, firstDayOfLastMonth, lastDayOfLastMonth)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for _, groupTransactionTotalAmountByBigCategory := range groupTransactionTotalAmountByBigCategoryList {
		for i, groupStandardBudgetByCategory := range dbGroupCustomBudgets.GroupCustomBudgets {
			if groupTransactionTotalAmountByBigCategory.BigCategoryID == groupStandardBudgetByCategory.BigCategoryID {
				dbGroupCustomBudgets.GroupCustomBudgets[i].LastMonthExpenses = groupTransactionTotalAmountByBigCategory.TotalAmount

				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&dbGroupCustomBudgets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupCustomBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	yearMonth, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	if err := h.GroupBudgetsRepo.DeleteGroupCustomBudgets(yearMonth, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"カスタム予算を削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetYearlyGroupBudgets(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	year, err := time.Parse("2006", mux.Vars(r)["year"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年を正しく指定してください。"}))
		return
	}

	monthlyGroupStandardBudget, err := h.GroupBudgetsRepo.GetMonthlyGroupStandardBudget(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	monthlyGroupCustomBudgets, err := h.GroupBudgetsRepo.GetMonthlyGroupCustomBudgets(year, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	yearlyBudget := model.NewYearlyGroupBudget(year)

	for i, j := 0, 0; i < len(yearlyBudget.GroupMonthlyBudgets); i++ {
		if j < len(monthlyGroupCustomBudgets) && time.Month(i)+1 == monthlyGroupCustomBudgets[j].Month.Month() {
			yearlyBudget.YearlyTotalBudget += monthlyGroupCustomBudgets[j].MonthlyTotalBudget
			yearlyBudget.GroupMonthlyBudgets[i] = monthlyGroupCustomBudgets[j]

			j++
			continue
		}

		monthlyGroupStandardBudget.Month.Time = year.AddDate(0, i, 0)
		yearlyBudget.YearlyTotalBudget += monthlyGroupStandardBudget.MonthlyTotalBudget
		yearlyBudget.GroupMonthlyBudgets[i] = monthlyGroupStandardBudget
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(yearlyBudget); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
