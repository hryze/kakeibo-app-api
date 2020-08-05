package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
)

func (h *DBHandler) GetGroupCategoriesList(w http.ResponseWriter, r *http.Request) {
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

	if err := getVerifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	groupBigCategoriesList, err := h.DBRepo.GetGroupBigCategoriesList()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	groupMediumCategoriesList, err := h.DBRepo.GetGroupMediumCategoriesList()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	groupCustomCategoriesList, err := h.DBRepo.GetGroupCustomCategoriesList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	for i, groupBigCategory := range groupBigCategoriesList {
		for _, groupCustomCategory := range groupCustomCategoriesList {
			if groupBigCategory.TransactionType == "income" && groupBigCategory.ID == groupCustomCategory.BigCategoryID {
				groupBigCategoriesList[i].IncomeAssociatedCategoriesList = append(groupBigCategoriesList[i].IncomeAssociatedCategoriesList, groupCustomCategory)
			} else if groupBigCategory.TransactionType == "expense" && groupBigCategory.ID == groupCustomCategory.BigCategoryID {
				groupBigCategoriesList[i].ExpenseAssociatedCategoriesList = append(groupBigCategoriesList[i].ExpenseAssociatedCategoriesList, groupCustomCategory)
			}
		}
	}
	for i, groupBigCategory := range groupBigCategoriesList {
		for _, groupMediumCategory := range groupMediumCategoriesList {
			if groupBigCategory.TransactionType == "income" && groupBigCategory.ID == groupMediumCategory.BigCategoryID {
				groupBigCategoriesList[i].IncomeAssociatedCategoriesList = append(groupBigCategoriesList[i].IncomeAssociatedCategoriesList, groupMediumCategory)
			} else if groupBigCategory.TransactionType == "expense" && groupBigCategory.ID == groupMediumCategory.BigCategoryID {
				groupBigCategoriesList[i].ExpenseAssociatedCategoriesList = append(groupBigCategoriesList[i].ExpenseAssociatedCategoriesList, groupMediumCategory)
			}
		}
	}
	var groupCategoriesList model.GroupCategoriesList
	for _, groupBigCategory := range groupBigCategoriesList {
		if groupBigCategory.TransactionType == "income" {
			groupCategoriesList.GroupIncomeBigCategoriesList = append(groupCategoriesList.GroupIncomeBigCategoriesList, model.NewGroupIncomeBigCategory(&groupBigCategory))
		} else if groupBigCategory.TransactionType == "expense" {
			groupCategoriesList.GroupExpenseBigCategoriesList = append(groupCategoriesList.GroupExpenseBigCategoriesList, model.NewGroupExpenseBigCategory(&groupBigCategory))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupCategoriesList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
