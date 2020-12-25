package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type GroupRelatedTransaction struct {
	TransactionType  string           `json:"transaction_type"`
	TransactionDate  time.Time        `json:"transaction_date"`
	Shop             model.NullString `json:"shop"`
	Memo             string           `json:"memo"`
	Amount           int64            `json:"amount"`
	PaymentUserID    string           `json:"payment_user_id"`
	BigCategoryID    int              `json:"big_category_id"`
	MediumCategoryID model.NullInt64  `json:"medium_category_id"`
	CustomCategoryID model.NullInt64  `json:"custom_category_id"`
}

func getGroupShoppingItemCategoriesName(categoriesID CategoriesID, groupID int) ([]byte, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/groups/%d/categories/name", accountHost, groupID)

	requestBody, err := json.Marshal(&categoriesID)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		"GET",
		requestURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 60 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode == http.StatusInternalServerError {
		return nil, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	categoriesNameBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return categoriesNameBytes, nil
}

func postGroupRelatedTransaction(groupShoppingItem model.GroupShoppingItem, groupID int, cookie *http.Cookie) (model.GroupShoppingItem, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/groups/%d/transactions", accountHost, groupID)

	groupRelatedTransaction := GroupRelatedTransaction{
		TransactionType:  "expense",
		TransactionDate:  groupShoppingItem.ExpectedPurchaseDate.Time,
		Shop:             groupShoppingItem.Shop,
		Memo:             fmt.Sprintf("【買い物リスト】%s", groupShoppingItem.Purchase),
		Amount:           groupShoppingItem.Amount.Int64,
		PaymentUserID:    groupShoppingItem.PaymentUserID.String,
		BigCategoryID:    groupShoppingItem.BigCategoryID,
		MediumCategoryID: groupShoppingItem.MediumCategoryID,
		CustomCategoryID: groupShoppingItem.CustomCategoryID,
	}

	requestBody, err := json.Marshal(&groupRelatedTransaction)
	if err != nil {
		return groupShoppingItem, err
	}

	request, err := http.NewRequest(
		"POST",
		requestURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return groupShoppingItem, err
	}

	request.AddCookie(cookie)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 60 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return groupShoppingItem, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if err := json.NewDecoder(response.Body).Decode(&groupShoppingItem.RelatedTransactionData); err != nil {
		return groupShoppingItem, err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return groupShoppingItem, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return groupShoppingItem, nil
}

func deleteGroupRelatedTransaction(groupShoppingItem model.GroupShoppingItem, groupID int, cookie *http.Cookie) (model.GroupShoppingItem, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/groups/%d/transactions/%d", accountHost, groupID, groupShoppingItem.RelatedTransactionData.ID.Int64)

	request, err := http.NewRequest(
		"DELETE",
		requestURL,
		nil,
	)
	if err != nil {
		return groupShoppingItem, err
	}

	request.AddCookie(cookie)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 60 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return groupShoppingItem, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode == http.StatusInternalServerError {
		return groupShoppingItem, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	groupShoppingItem.RelatedTransactionData = nil

	return groupShoppingItem, nil
}

func (h *DBHandler) PostGroupRegularShoppingItem(w http.ResponseWriter, r *http.Request) {
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

	var groupRegularShoppingItem model.GroupRegularShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&groupRegularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	groupRegularShoppingItemResult, todayGroupShoppingItemResult, laterThanTodayGroupShoppingItemResult, err := h.GroupShoppingListRepo.PostGroupRegularShoppingItem(&groupRegularShoppingItem, groupID, today)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupRegularShoppingItemID, err := groupRegularShoppingItemResult.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var todayGroupShoppingItemID int64
	if todayGroupShoppingItemResult != nil {
		todayGroupShoppingItemID, err = todayGroupShoppingItemResult.LastInsertId()
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	laterThanTodayGroupShoppingItemID, err := laterThanTodayGroupShoppingItemResult.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupRegularShoppingItem, err = h.GroupShoppingListRepo.GetGroupRegularShoppingItem(int(groupRegularShoppingItemID))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupShoppingList, err := h.GroupShoppingListRepo.GetGroupShoppingListRelatedToGroupRegularShoppingItem(int(todayGroupShoppingItemID), int(laterThanTodayGroupShoppingItemID))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesID := CategoriesID{
		MediumCategoryID: groupRegularShoppingItem.MediumCategoryID,
		CustomCategoryID: groupRegularShoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getGroupShoppingItemCategoriesName(categoriesID, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &groupRegularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for i := 0; i < len(groupShoppingList.GroupShoppingList); i++ {
		if err := json.Unmarshal(categoriesNameBytes, &groupShoppingList.GroupShoppingList[i]); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingData := struct {
		GroupRegularShoppingItem model.GroupRegularShoppingItem `json:"regular_shopping_item"`
		model.GroupShoppingList
	}{
		GroupRegularShoppingItem: groupRegularShoppingItem,
		GroupShoppingList:        groupShoppingList,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&shoppingData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutGroupRegularShoppingItem(w http.ResponseWriter, r *http.Request) {
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

	var groupRegularShoppingItem model.GroupRegularShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&groupRegularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupRegularShoppingItemID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"定期ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	todayGroupShoppingItemResult, laterThanTodayGroupShoppingItemResult, err := h.GroupShoppingListRepo.PutGroupRegularShoppingItem(&groupRegularShoppingItem, groupRegularShoppingItemID, groupID, today)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var todayGroupShoppingItemID int64
	if todayGroupShoppingItemResult != nil {
		todayGroupShoppingItemID, err = todayGroupShoppingItemResult.LastInsertId()
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	laterThanTodayGroupShoppingItemID, err := laterThanTodayGroupShoppingItemResult.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupRegularShoppingItem, err = h.GroupShoppingListRepo.GetGroupRegularShoppingItem(groupRegularShoppingItemID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupShoppingList, err := h.GroupShoppingListRepo.GetGroupShoppingListRelatedToGroupRegularShoppingItem(int(todayGroupShoppingItemID), int(laterThanTodayGroupShoppingItemID))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesID := CategoriesID{
		MediumCategoryID: groupRegularShoppingItem.MediumCategoryID,
		CustomCategoryID: groupRegularShoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getGroupShoppingItemCategoriesName(categoriesID, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &groupRegularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for i := 0; i < len(groupShoppingList.GroupShoppingList); i++ {
		if err := json.Unmarshal(categoriesNameBytes, &groupShoppingList.GroupShoppingList[i]); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingData := struct {
		GroupRegularShoppingItem model.GroupRegularShoppingItem `json:"regular_shopping_item"`
		model.GroupShoppingList
	}{
		GroupRegularShoppingItem: groupRegularShoppingItem,
		GroupShoppingList:        groupShoppingList,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupRegularShoppingItem(w http.ResponseWriter, r *http.Request) {
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

	groupRegularShoppingItemID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"定期ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	if err := h.GroupShoppingListRepo.DeleteGroupRegularShoppingItem(groupRegularShoppingItemID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"定期ショッピングアイテムを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroupShoppingItem(w http.ResponseWriter, r *http.Request) {
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

	var groupShoppingItem model.GroupShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&groupShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	result, err := h.GroupShoppingListRepo.PostGroupShoppingItem(&groupShoppingItem, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupShoppingItem, err = h.GroupShoppingListRepo.GetGroupShoppingItem(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesID := CategoriesID{
		MediumCategoryID: groupShoppingItem.MediumCategoryID,
		CustomCategoryID: groupShoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getGroupShoppingItemCategoriesName(categoriesID, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &groupShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&groupShoppingItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutGroupShoppingItem(w http.ResponseWriter, r *http.Request) {
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

	var groupShoppingItem model.GroupShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&groupShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupShoppingItem.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	if groupShoppingItem.CompleteFlag && groupShoppingItem.TransactionAutoAdd && groupShoppingItem.RelatedTransactionData == nil {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupShoppingItem, err = postGroupRelatedTransaction(groupShoppingItem, groupID, cookie)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	} else if !groupShoppingItem.CompleteFlag && groupShoppingItem.RelatedTransactionData != nil {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupShoppingItem, err = deleteGroupRelatedTransaction(groupShoppingItem, groupID, cookie)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	if _, err := h.GroupShoppingListRepo.PutGroupShoppingItem(&groupShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbShoppingItem, err := h.GroupShoppingListRepo.GetGroupShoppingItem(groupShoppingItem.ID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupShoppingItem.PostedDate = dbShoppingItem.PostedDate
	groupShoppingItem.UpdatedDate = dbShoppingItem.UpdatedDate

	categoriesID := CategoriesID{
		MediumCategoryID: groupShoppingItem.MediumCategoryID,
		CustomCategoryID: groupShoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getGroupShoppingItemCategoriesName(categoriesID, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &groupShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupShoppingItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupShoppingItem(w http.ResponseWriter, r *http.Request) {
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

	groupShoppingItemID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	if err := h.GroupShoppingListRepo.DeleteGroupShoppingItem(groupShoppingItemID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"ショッピングアイテムを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
