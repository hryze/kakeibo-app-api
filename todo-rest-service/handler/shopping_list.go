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

type CategoriesID struct {
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

func getShoppingItemCategoriesName(categoriesID CategoriesID) ([]byte, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/categories/name", accountHost)

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

func getShoppingItemCategoriesNameList(categoriesIdList []CategoriesID) ([]byte, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/categories/names", accountHost)

	requestBody, err := json.Marshal(&categoriesIdList)
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

func (h *DBHandler) GetShoppingDataByMonth(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingList, err := h.ShoppingListRepo.GetRegularShoppingList(userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(regularShoppingList.RegularShoppingList) != 0 {
		now := h.TimeManage.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if err = h.ShoppingListRepo.PutRegularShoppingList(regularShoppingList, userID, today); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		regularShoppingList, err = h.ShoppingListRepo.GetRegularShoppingList(userID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		categoriesIdList := make([]CategoriesID, len(regularShoppingList.RegularShoppingList))

		for i, regularShoppingItem := range regularShoppingList.RegularShoppingList {
			categoriesIdList[i] = CategoriesID{
				MediumCategoryID: regularShoppingItem.MediumCategoryID,
				CustomCategoryID: regularShoppingItem.CustomCategoryID,
			}
		}

		categoriesNameListBytes, err := getShoppingItemCategoriesNameList(categoriesIdList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		if err := json.Unmarshal(categoriesNameListBytes, &regularShoppingList.RegularShoppingList); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	firstDay, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	lastDay := time.Date(firstDay.Year(), firstDay.Month()+1, 1, 0, 0, 0, 0, firstDay.Location()).Add(-1 * time.Second)

	shoppingList, err := h.ShoppingListRepo.GetShoppingListByMonth(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	if len(shoppingList.ShoppingList) != 0 {
		categoriesIdList := make([]CategoriesID, len(shoppingList.ShoppingList))

		for i, shoppingItem := range shoppingList.ShoppingList {
			categoriesIdList[i] = CategoriesID{
				MediumCategoryID: shoppingItem.MediumCategoryID,
				CustomCategoryID: shoppingItem.CustomCategoryID,
			}
		}

		categoriesNameListBytes, err := getShoppingItemCategoriesNameList(categoriesIdList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		if err := json.Unmarshal(categoriesNameListBytes, &shoppingList.ShoppingList); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingData := model.ShoppingDataByMonth{
		RegularShoppingList: regularShoppingList,
		ShoppingList:        shoppingList,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetShoppingDataByCategories(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingList, err := h.ShoppingListRepo.GetRegularShoppingList(userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(regularShoppingList.RegularShoppingList) != 0 {
		now := h.TimeManage.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if err = h.ShoppingListRepo.PutRegularShoppingList(regularShoppingList, userID, today); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		regularShoppingList, err = h.ShoppingListRepo.GetRegularShoppingList(userID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		categoriesIdList := make([]CategoriesID, len(regularShoppingList.RegularShoppingList))

		for i, regularShoppingItem := range regularShoppingList.RegularShoppingList {
			categoriesIdList[i] = CategoriesID{
				MediumCategoryID: regularShoppingItem.MediumCategoryID,
				CustomCategoryID: regularShoppingItem.CustomCategoryID,
			}
		}

		categoriesNameListBytes, err := getShoppingItemCategoriesNameList(categoriesIdList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		if err := json.Unmarshal(categoriesNameListBytes, &regularShoppingList.RegularShoppingList); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	firstDay, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	lastDay := time.Date(firstDay.Year(), firstDay.Month()+1, 1, 0, 0, 0, 0, firstDay.Location()).Add(-1 * time.Second)

	shoppingList, err := h.ShoppingListRepo.GetShoppingListByCategories(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	if len(shoppingList.ShoppingList) != 0 {
		categoriesIdList := make([]CategoriesID, len(shoppingList.ShoppingList))

		for i, shoppingItem := range shoppingList.ShoppingList {
			categoriesIdList[i] = CategoriesID{
				MediumCategoryID: shoppingItem.MediumCategoryID,
				CustomCategoryID: shoppingItem.CustomCategoryID,
			}
		}

		categoriesNameListBytes, err := getShoppingItemCategoriesNameList(categoriesIdList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		if err := json.Unmarshal(categoriesNameListBytes, &shoppingList.ShoppingList); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	var shoppingListByCategories []model.ShoppingListByCategory

	for i, j := 0, 0; i < len(shoppingList.ShoppingList); i++ {
		firstIndexBigCategoryID := shoppingList.ShoppingList[j].BigCategoryID
		currentIndexBigCategoryID := shoppingList.ShoppingList[i].BigCategoryID

		if firstIndexBigCategoryID != currentIndexBigCategoryID {
			shoppingListByCategory := model.ShoppingListByCategory{
				BigCategoryName: shoppingList.ShoppingList[j].BigCategoryName,
				ShoppingList:    append(make([]model.ShoppingItem, 0, i-j), shoppingList.ShoppingList[j:i]...),
			}

			shoppingListByCategories = append(shoppingListByCategories, shoppingListByCategory)

			j = i
		} else if i == len(shoppingList.ShoppingList)-1 {
			shoppingListByCategory := model.ShoppingListByCategory{
				BigCategoryName: shoppingList.ShoppingList[j].BigCategoryName,
				ShoppingList:    append(make([]model.ShoppingItem, 0, i-j), shoppingList.ShoppingList[j:]...),
			}

			shoppingListByCategories = append(shoppingListByCategories, shoppingListByCategory)
		}
	}

	shoppingDataByCategories := model.ShoppingDataByCategories{
		RegularShoppingList:      regularShoppingList,
		ShoppingListByCategories: shoppingListByCategories,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingDataByCategories); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostRegularShoppingItem(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var regularShoppingItem model.RegularShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&regularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	regularShoppingItemResult, todayShoppingItemResult, laterThanTodayShoppingItemResult, err := h.ShoppingListRepo.PostRegularShoppingItem(&regularShoppingItem, userID, today)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingItemId, err := regularShoppingItemResult.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var todayShoppingItemID int64
	if todayShoppingItemResult != nil {
		todayShoppingItemID, err = todayShoppingItemResult.LastInsertId()
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	laterThanTodayShoppingItemID, err := laterThanTodayShoppingItemResult.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingItem, err = h.ShoppingListRepo.GetRegularShoppingItem(int(regularShoppingItemId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingList, err := h.ShoppingListRepo.GetShoppingListRelatedToRegularShoppingItem(int(todayShoppingItemID), int(laterThanTodayShoppingItemID))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesID := CategoriesID{
		MediumCategoryID: regularShoppingItem.MediumCategoryID,
		CustomCategoryID: regularShoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getShoppingItemCategoriesName(categoriesID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &regularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for i := 0; i < len(shoppingList.ShoppingList); i++ {
		if err := json.Unmarshal(categoriesNameBytes, &shoppingList.ShoppingList[i]); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingData := struct {
		RegularShoppingItem model.RegularShoppingItem `json:"regular_shopping_item"`
		model.ShoppingList
	}{
		RegularShoppingItem: regularShoppingItem,
		ShoppingList:        shoppingList,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&shoppingData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutRegularShoppingItem(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var regularShoppingItem model.RegularShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&regularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingItemID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"定期ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	todayShoppingItemResult, laterThanTodayShoppingItemResult, err := h.ShoppingListRepo.PutRegularShoppingItem(&regularShoppingItem, regularShoppingItemID, userID, today)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var todayShoppingItemID int64
	if todayShoppingItemResult != nil {
		todayShoppingItemID, err = todayShoppingItemResult.LastInsertId()
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	laterThanTodayShoppingItemID, err := laterThanTodayShoppingItemResult.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingItem, err = h.ShoppingListRepo.GetRegularShoppingItem(regularShoppingItemID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingList, err := h.ShoppingListRepo.GetShoppingListRelatedToRegularShoppingItem(int(todayShoppingItemID), int(laterThanTodayShoppingItemID))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesID := CategoriesID{
		MediumCategoryID: regularShoppingItem.MediumCategoryID,
		CustomCategoryID: regularShoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getShoppingItemCategoriesName(categoriesID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &regularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for i := 0; i < len(shoppingList.ShoppingList); i++ {
		if err := json.Unmarshal(categoriesNameBytes, &shoppingList.ShoppingList[i]); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingData := struct {
		RegularShoppingItem model.RegularShoppingItem `json:"regular_shopping_item"`
		model.ShoppingList
	}{
		RegularShoppingItem: regularShoppingItem,
		ShoppingList:        shoppingList,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteRegularShoppingItem(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingItemID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"定期ショッピングアイテムIDを正しく指定してください。"}))
		return
	}

	if err := h.ShoppingListRepo.DeleteRegularShoppingItem(regularShoppingItemID); err != nil {
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

	categoriesID := CategoriesID{
		MediumCategoryID: shoppingItem.MediumCategoryID,
		CustomCategoryID: shoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getShoppingItemCategoriesName(categoriesID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &shoppingItem); err != nil {
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

	categoriesID := CategoriesID{
		MediumCategoryID: shoppingItem.MediumCategoryID,
		CustomCategoryID: shoppingItem.CustomCategoryID,
	}

	categoriesNameBytes, err := getShoppingItemCategoriesName(categoriesID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := json.Unmarshal(categoriesNameBytes, &shoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

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
