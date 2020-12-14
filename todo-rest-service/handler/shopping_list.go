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

type ShoppingItemCategoryID struct {
	MediumCategoryID model.NullInt64 `json:"medium_category_id"`
	CustomCategoryID model.NullInt64 `json:"custom_category_id"`
}

type RelatedTransaction struct {
	TransactionType  string           `json:"transaction_type"`
	TransactionDate  time.Time        `json:"transaction_date"`
	Shop             model.NullString `json:"shop"`
	Memo             string           `json:"memo"`
	Amount           int64            `json:"amount"`
	BigCategoryID    int              `json:"big_category_id"`
	MediumCategoryID model.NullInt64  `json:"medium_category_id"`
	CustomCategoryID model.NullInt64  `json:"custom_category_id"`
}

func getShoppingItemCategoryName(shoppingItem model.ShoppingItem) (model.ShoppingItem, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/categories/names", accountHost)

	shoppingItemCategoryID := ShoppingItemCategoryID{
		MediumCategoryID: shoppingItem.MediumCategoryID,
		CustomCategoryID: shoppingItem.CustomCategoryID,
	}

	requestBody, err := json.Marshal(&shoppingItemCategoryID)
	if err != nil {
		return shoppingItem, err
	}

	request, err := http.NewRequest(
		"GET",
		requestURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return shoppingItem, err
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
		return shoppingItem, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if err := json.NewDecoder(response.Body).Decode(&shoppingItem); err != nil {
		return shoppingItem, err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return shoppingItem, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return shoppingItem, nil
}

func postRelatedTransaction(shoppingItem model.ShoppingItem, cookie *http.Cookie) (model.ShoppingItem, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/transactions", accountHost)

	relatedTransaction := RelatedTransaction{
		TransactionType:  "expense",
		TransactionDate:  shoppingItem.ExpectedPurchaseDate.Time,
		Shop:             shoppingItem.Shop,
		Memo:             fmt.Sprintf("【買い物リスト】%s", shoppingItem.Purchase),
		Amount:           shoppingItem.Amount.Int64,
		BigCategoryID:    shoppingItem.BigCategoryID,
		MediumCategoryID: shoppingItem.MediumCategoryID,
		CustomCategoryID: shoppingItem.CustomCategoryID,
	}

	requestBody, err := json.Marshal(&relatedTransaction)
	if err != nil {
		return shoppingItem, err
	}

	request, err := http.NewRequest(
		"POST",
		requestURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return shoppingItem, err
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
		return shoppingItem, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if err := json.NewDecoder(response.Body).Decode(&shoppingItem.RelatedTransactionData); err != nil {
		return shoppingItem, err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return shoppingItem, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return shoppingItem, nil
}

func deleteRelatedTransaction(shoppingItem model.ShoppingItem, cookie *http.Cookie) (model.ShoppingItem, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/transactions/%d", accountHost, shoppingItem.RelatedTransactionData.ID.Int64)

	request, err := http.NewRequest(
		"DELETE",
		requestURL,
		nil,
	)
	if err != nil {
		return shoppingItem, err
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
		return shoppingItem, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode == http.StatusInternalServerError {
		return shoppingItem, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	shoppingItem.RelatedTransactionData = nil

	return shoppingItem, nil
}

func (h *DBHandler) PostShoppingItem(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var shoppingItem model.ShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&shoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	result, err := h.ShoppingListRepo.PostShoppingItem(&shoppingItem, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingItem, err = h.ShoppingListRepo.GetShoppingItem(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingItem, err = getShoppingItemCategoryName(shoppingItem)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&shoppingItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutShoppingItem(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var shoppingItem model.ShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&shoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingItem.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	if shoppingItem.CompleteFlag && shoppingItem.TransactionAutoAdd && shoppingItem.RelatedTransactionData == nil {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		shoppingItem, err = postRelatedTransaction(shoppingItem, cookie)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	} else if !shoppingItem.CompleteFlag && shoppingItem.RelatedTransactionData != nil {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		shoppingItem, err = deleteRelatedTransaction(shoppingItem, cookie)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	if _, err := h.ShoppingListRepo.PutShoppingItem(&shoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbShoppingItem, err := h.ShoppingListRepo.GetShoppingItem(shoppingItem.ID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingItem.PostedDate = dbShoppingItem.PostedDate
	shoppingItem.UpdatedDate = dbShoppingItem.UpdatedDate

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteShoppingItem(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingItemID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	if err := h.ShoppingListRepo.DeleteShoppingItem(shoppingItemID); err != nil {
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
