package handler

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"

	"github.com/hryze/kakeibo-app-api/todo-rest-service/config"
	"github.com/hryze/kakeibo-app-api/todo-rest-service/domain/model"
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

type GroupRegularShoppingItemValidationErrorMsg struct {
	Message []string `json:"message"`
}

func (e *GroupRegularShoppingItemValidationErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}

type GroupShoppingItemValidationErrorMsg struct {
	Message []string `json:"message"`
}

func (e *GroupShoppingItemValidationErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}

func validateGroupRegularShoppingItem(groupRegularShoppingItem model.GroupRegularShoppingItem) error {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(groupRegularShoppingItemValidateValuer, model.Date{}, model.NullString{}, model.NullInt{}, model.NullInt64{})
	if err := validate.RegisterValidation("blank", blankValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("date_range", dateRangeValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("either_id", eitherCategoryIDValidation); err != nil {
		return err
	}

	err := validate.Struct(groupRegularShoppingItem)
	if err == nil {
		return nil
	}

	var groupRegularShoppingItemValidationErrorMsg GroupRegularShoppingItemValidationErrorMsg
	for _, err := range err.(validator.ValidationErrors) {
		var errorMessage string

		fieldName := err.Field()
		switch fieldName {
		case "ExpectedPurchaseDate":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "購入予定日が選択されていません。"
			case "date_range":
				errorMessage = "購入予定日は今日以降の日付を選択してください。"
			}
		case "CycleType":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "購入周期タイプが選択されていません。"
			case "oneof":
				errorMessage = "購入周期タイプを正しく選択してください。"
			}
		case "Cycle":
			errorMessage = "購入周期は1以上の正の整数を入力してください。"
		case "Purchase":
			tagName := err.Tag()
			switch tagName {
			case "max":
				errorMessage = "定期購入品は50文字以内で入力してください。"
			case "blank":
				errorMessage = "定期購入品の文字列先頭か末尾に空白がないか確認してください。"
			}
		case "Shop":
			tagName := err.Tag()
			switch tagName {
			case "max":
				errorMessage = "店名は20文字以内で入力してください。"
			case "blank":
				errorMessage = "店名の文字列先頭か末尾に空白がないか確認してください。"
			}
		case "Amount":
			errorMessage = "金額は1以上の正の整数を入力してください。"
		case "BigCategoryID":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "大カテゴリーが選択されていません。"
			case "min", "max":
				errorMessage = "大カテゴリーを正しく選択してください。"
			case "either_id":
				errorMessage = "中カテゴリーを正しく選択してください。"
			}
		case "MediumCategoryID":
			errorMessage = "中カテゴリーを正しく選択してください。"
		case "CustomCategoryID":
			errorMessage = "中カテゴリーを正しく選択してください。"
		case "PaymentUserID":
			errorMessage = "支払ユーザーを正しく選択してください。"
		}

		groupRegularShoppingItemValidationErrorMsg.Message = append(groupRegularShoppingItemValidationErrorMsg.Message, errorMessage)
	}

	return &groupRegularShoppingItemValidationErrorMsg
}

func groupRegularShoppingItemValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}

	return nil
}

func validateGroupShoppingItem(groupShoppingItem model.GroupShoppingItem) error {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(groupShoppingItemValidateValuer, model.Date{}, model.NullString{}, model.NullInt64{})
	if err := validate.RegisterValidation("blank", blankValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("date_range", dateRangeValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("either_id", eitherCategoryIDValidation); err != nil {
		return err
	}

	err := validate.Struct(groupShoppingItem)
	if err == nil {
		return nil
	}

	var groupShoppingItemValidationErrorMsg GroupShoppingItemValidationErrorMsg
	for _, err := range err.(validator.ValidationErrors) {
		var errorMessage string

		fieldName := err.Field()
		switch fieldName {
		case "ExpectedPurchaseDate":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "購入予定日が選択されていません。"
			case "date_range":
				errorMessage = "購入予定日は今日以降の日付を選択してください。"
			}
		case "Purchase":
			tagName := err.Tag()
			switch tagName {
			case "max":
				errorMessage = "購入品は50文字以内で入力してください。"
			case "blank":
				errorMessage = "購入品の文字列先頭か末尾に空白がないか確認してください。"
			}
		case "Shop":
			tagName := err.Tag()
			switch tagName {
			case "max":
				errorMessage = "店名は20文字以内で入力してください。"
			case "blank":
				errorMessage = "店名の文字列先頭か末尾に空白がないか確認してください。"
			}
		case "Amount":
			errorMessage = "金額は1以上の正の整数を入力してください。"
		case "BigCategoryID":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "大カテゴリーが選択されていません。"
			case "min", "max":
				errorMessage = "大カテゴリーを正しく選択してください。"
			case "either_id":
				errorMessage = "中カテゴリーを正しく選択してください。"
			}
		case "MediumCategoryID":
			errorMessage = "中カテゴリーを正しく選択してください。"
		case "CustomCategoryID":
			errorMessage = "中カテゴリーを正しく選択してください。"
		case "PaymentUserID":
			errorMessage = "支払ユーザーを正しく選択してください。"
		}

		groupShoppingItemValidationErrorMsg.Message = append(groupShoppingItemValidationErrorMsg.Message, errorMessage)
	}

	return &groupShoppingItemValidationErrorMsg
}

func groupShoppingItemValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}

	return nil
}

func getGroupShoppingItemCategoriesName(categoriesID CategoriesID, groupID int) ([]byte, error) {
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/categories/name",
		config.Env.AccountApi.Host, config.Env.AccountApi.Port, groupID,
	)

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

func getGroupShoppingItemCategoriesNameList(categoriesIdList []CategoriesID, groupID int) ([]byte, error) {
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/categories/names",
		config.Env.AccountApi.Host, config.Env.AccountApi.Port, groupID,
	)

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

func getGroupShoppingItemRelatedTransactionDataList(transactionIdList []int64, groupID int) ([]*model.GroupTransactionData, error) {
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/transactions/related-shopping-list",
		config.Env.AccountApi.Host, config.Env.AccountApi.Port, groupID,
	)

	requestBody, err := json.Marshal(&transactionIdList)
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

	var groupShoppingItemRelatedTransactionDataList []*model.GroupTransactionData
	if err := json.NewDecoder(response.Body).Decode(&groupShoppingItemRelatedTransactionDataList); err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return nil, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return groupShoppingItemRelatedTransactionDataList, nil
}

func postGroupRelatedTransaction(groupShoppingItem model.GroupShoppingItem, groupID int, cookie *http.Cookie) (model.GroupShoppingItem, error) {
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/transactions",
		config.Env.AccountApi.Host, config.Env.AccountApi.Port, groupID,
	)

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
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/transactions/%d",
		config.Env.AccountApi.Host, config.Env.AccountApi.Port, groupID, groupShoppingItem.RelatedTransactionData.ID.Int64,
	)

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

func generateGroupRegularShoppingList(groupRegularShoppingList model.GroupRegularShoppingList, groupID int) (model.GroupRegularShoppingList, error) {
	categoriesIdList := make([]CategoriesID, len(groupRegularShoppingList.GroupRegularShoppingList))

	for i, groupRegularShoppingItem := range groupRegularShoppingList.GroupRegularShoppingList {
		categoriesIdList[i] = CategoriesID{
			MediumCategoryID: groupRegularShoppingItem.MediumCategoryID,
			CustomCategoryID: groupRegularShoppingItem.CustomCategoryID,
		}
	}

	categoriesNameListBytes, err := getGroupShoppingItemCategoriesNameList(categoriesIdList, groupID)
	if err != nil {
		return groupRegularShoppingList, err
	}

	if err := json.Unmarshal(categoriesNameListBytes, &groupRegularShoppingList.GroupRegularShoppingList); err != nil {
		return groupRegularShoppingList, err
	}

	return groupRegularShoppingList, nil
}

func generateGroupShoppingList(groupShoppingList model.GroupShoppingList, groupID int) (model.GroupShoppingList, error) {
	categoriesIdList := make([]CategoriesID, len(groupShoppingList.GroupShoppingList))
	var transactionIdList []int64

	for i, groupShoppingItem := range groupShoppingList.GroupShoppingList {
		categoriesIdList[i] = CategoriesID{
			MediumCategoryID: groupShoppingItem.MediumCategoryID,
			CustomCategoryID: groupShoppingItem.CustomCategoryID,
		}

		if groupShoppingItem.RelatedTransactionData != nil {
			transactionIdList = append(transactionIdList, groupShoppingItem.RelatedTransactionData.ID.Int64)
		}
	}

	categoriesNameListBytes, err := getGroupShoppingItemCategoriesNameList(categoriesIdList, groupID)
	if err != nil {
		return groupShoppingList, err
	}

	if err := json.Unmarshal(categoriesNameListBytes, &groupShoppingList.GroupShoppingList); err != nil {
		return groupShoppingList, err
	}

	if len(transactionIdList) != 0 {
		groupShoppingItemRelatedTransactionDataList, err := getGroupShoppingItemRelatedTransactionDataList(transactionIdList, groupID)
		if err != nil {
			return groupShoppingList, err
		}

		for _, groupShoppingItemRelatedTransactionData := range groupShoppingItemRelatedTransactionDataList {
			for i, groupShoppingItem := range groupShoppingList.GroupShoppingList {
				if groupShoppingItem.RelatedTransactionData != nil && groupShoppingItemRelatedTransactionData.ID.Int64 == groupShoppingItem.RelatedTransactionData.ID.Int64 {
					groupShoppingList.GroupShoppingList[i].RelatedTransactionData = groupShoppingItemRelatedTransactionData
				}
			}
		}
	}

	return groupShoppingList, nil
}

func generateGroupShoppingListByCategories(groupShoppingList model.GroupShoppingList) []model.GroupShoppingListByCategory {
	groupShoppingListByCategories := make([]model.GroupShoppingListByCategory, 0)

	for i, j := 0, 0; i < len(groupShoppingList.GroupShoppingList); i++ {
		firstIndexBigCategoryID := groupShoppingList.GroupShoppingList[j].BigCategoryID
		currentIndexBigCategoryID := groupShoppingList.GroupShoppingList[i].BigCategoryID

		if firstIndexBigCategoryID != currentIndexBigCategoryID {
			groupShoppingListByCategory := model.GroupShoppingListByCategory{
				BigCategoryName:   groupShoppingList.GroupShoppingList[j].BigCategoryName,
				GroupShoppingList: append(make([]model.GroupShoppingItem, 0, i-j), groupShoppingList.GroupShoppingList[j:i]...),
			}

			groupShoppingListByCategories = append(groupShoppingListByCategories, groupShoppingListByCategory)

			j = i
		}

		if i == len(groupShoppingList.GroupShoppingList)-1 {
			if firstIndexBigCategoryID == currentIndexBigCategoryID {
				groupShoppingListByCategory := model.GroupShoppingListByCategory{
					BigCategoryName:   groupShoppingList.GroupShoppingList[j].BigCategoryName,
					GroupShoppingList: append(make([]model.GroupShoppingItem, 0, i-j), groupShoppingList.GroupShoppingList[j:]...),
				}

				groupShoppingListByCategories = append(groupShoppingListByCategories, groupShoppingListByCategory)
			} else if i == j {
				groupShoppingListByCategory := model.GroupShoppingListByCategory{
					BigCategoryName:   groupShoppingList.GroupShoppingList[i].BigCategoryName,
					GroupShoppingList: []model.GroupShoppingItem{groupShoppingList.GroupShoppingList[i]},
				}

				groupShoppingListByCategories = append(groupShoppingListByCategories, groupShoppingListByCategory)
			}
		}
	}

	return groupShoppingListByCategories
}

func (h *DBHandler) GetDailyGroupShoppingDataByDay(w http.ResponseWriter, r *http.Request) {
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

	groupRegularShoppingList, err := h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupRegularShoppingList.GroupRegularShoppingList) != 0 {
		now := h.TimeManage.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if err = h.GroupShoppingListRepo.PutGroupRegularShoppingList(groupRegularShoppingList, groupID, today); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = generateGroupRegularShoppingList(groupRegularShoppingList, groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	date, err := time.Parse("2006-01-02", mux.Vars(r)["date"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"日付を正しく指定してください。"}))
		return
	}

	groupShoppingList, err := h.GroupShoppingListRepo.GetDailyGroupShoppingListByDay(date, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupShoppingList.GroupShoppingList) != 0 {
		groupShoppingList, err = generateGroupShoppingList(groupShoppingList, groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	groupShoppingData := model.GroupShoppingDataByDay{
		GroupRegularShoppingList: groupRegularShoppingList,
		GroupShoppingList:        groupShoppingList,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupShoppingData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetDailyGroupShoppingDataByCategory(w http.ResponseWriter, r *http.Request) {
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

	groupRegularShoppingList, err := h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupRegularShoppingList.GroupRegularShoppingList) != 0 {
		now := h.TimeManage.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if err = h.GroupShoppingListRepo.PutGroupRegularShoppingList(groupRegularShoppingList, groupID, today); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = generateGroupRegularShoppingList(groupRegularShoppingList, groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	date, err := time.Parse("2006-01-02", mux.Vars(r)["date"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"日付を正しく指定してください。"}))
		return
	}

	groupShoppingList, err := h.GroupShoppingListRepo.GetDailyGroupShoppingListByCategory(date, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupShoppingList.GroupShoppingList) != 0 {
		groupShoppingList, err = generateGroupShoppingList(groupShoppingList, groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	groupShoppingListByCategories := generateGroupShoppingListByCategories(groupShoppingList)

	shoppingDataByCategories := model.GroupShoppingDataByCategory{
		GroupRegularShoppingList:    groupRegularShoppingList,
		GroupShoppingListByCategory: groupShoppingListByCategories,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingDataByCategories); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetMonthlyGroupShoppingDataByDay(w http.ResponseWriter, r *http.Request) {
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

	groupRegularShoppingList, err := h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupRegularShoppingList.GroupRegularShoppingList) != 0 {
		now := h.TimeManage.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if err = h.GroupShoppingListRepo.PutGroupRegularShoppingList(groupRegularShoppingList, groupID, today); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = generateGroupRegularShoppingList(groupRegularShoppingList, groupID)
		if err != nil {
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

	groupShoppingList, err := h.GroupShoppingListRepo.GetMonthlyGroupShoppingListByDay(firstDay, lastDay, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupShoppingList.GroupShoppingList) != 0 {
		groupShoppingList, err = generateGroupShoppingList(groupShoppingList, groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	groupShoppingData := model.GroupShoppingDataByDay{
		GroupRegularShoppingList: groupRegularShoppingList,
		GroupShoppingList:        groupShoppingList,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupShoppingData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetMonthlyGroupShoppingDataByCategory(w http.ResponseWriter, r *http.Request) {
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

	groupRegularShoppingList, err := h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupRegularShoppingList.GroupRegularShoppingList) != 0 {
		now := h.TimeManage.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		if err = h.GroupShoppingListRepo.PutGroupRegularShoppingList(groupRegularShoppingList, groupID, today); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = h.GroupShoppingListRepo.GetGroupRegularShoppingList(groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		groupRegularShoppingList, err = generateGroupRegularShoppingList(groupRegularShoppingList, groupID)
		if err != nil {
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

	groupShoppingList, err := h.GroupShoppingListRepo.GetMonthlyGroupShoppingListByCategory(firstDay, lastDay, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupShoppingList.GroupShoppingList) != 0 {
		groupShoppingList, err = generateGroupShoppingList(groupShoppingList, groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	groupShoppingListByCategories := generateGroupShoppingListByCategories(groupShoppingList)

	shoppingDataByCategories := model.GroupShoppingDataByCategory{
		GroupRegularShoppingList:    groupRegularShoppingList,
		GroupShoppingListByCategory: groupShoppingListByCategories,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingDataByCategories); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetExpiredGroupShoppingList(w http.ResponseWriter, r *http.Request) {
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

	now := h.TimeManage.Now()
	dueDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	expiredGroupShoppingList, err := h.GroupShoppingListRepo.GetExpiredGroupShoppingList(dueDate, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(expiredGroupShoppingList.ExpiredGroupShoppingList) != 0 {
		categoriesIdList := make([]CategoriesID, len(expiredGroupShoppingList.ExpiredGroupShoppingList))

		for i, groupShoppingItem := range expiredGroupShoppingList.ExpiredGroupShoppingList {
			categoriesIdList[i] = CategoriesID{
				MediumCategoryID: groupShoppingItem.MediumCategoryID,
				CustomCategoryID: groupShoppingItem.CustomCategoryID,
			}
		}

		categoriesNameListBytes, err := getGroupShoppingItemCategoriesNameList(categoriesIdList, groupID)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		if err := json.Unmarshal(categoriesNameListBytes, &expiredGroupShoppingList.ExpiredGroupShoppingList); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&expiredGroupShoppingList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	groupRegularShoppingItem.Today = today

	if err := validateGroupRegularShoppingItem(groupRegularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

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

	groupShoppingList, err := h.GroupShoppingListRepo.GetGroupShoppingListRelatedToPostedGroupRegularShoppingItem(int(todayGroupShoppingItemID), int(laterThanTodayGroupShoppingItemID))
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

	if err := validateGroupRegularShoppingItem(groupRegularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if err := h.GroupShoppingListRepo.PutGroupRegularShoppingItem(&groupRegularShoppingItem, groupRegularShoppingItemID, groupID, today); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupRegularShoppingItem, err = h.GroupShoppingListRepo.GetGroupRegularShoppingItem(groupRegularShoppingItemID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupShoppingList, err := h.GroupShoppingListRepo.GetGroupShoppingListRelatedToUpdatedGroupRegularShoppingItem(groupRegularShoppingItemID)
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

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	groupShoppingItem.Today = today

	if err := validateGroupShoppingItem(groupShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
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

	if err := validateGroupShoppingItem(groupShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if groupShoppingItem.CompleteFlag && groupShoppingItem.TransactionAutoAdd && groupShoppingItem.RelatedTransactionData == nil {
		cookie, err := r.Cookie(config.Env.Cookie.Name)
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
		cookie, err := r.Cookie(config.Env.Cookie.Name)
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

func (h *DBHandler) PutGroupShoppingListCustomCategoryIdToMediumCategoryId(w http.ResponseWriter, r *http.Request) {
	categoriesID := struct {
		MediumCategoryID int `json:"medium_category_id"`
		CustomCategoryID int `json:"custom_category_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&categoriesID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.GroupShoppingListRepo.PutGroupShoppingListCustomCategoryIdToMediumCategoryId(categoriesID.MediumCategoryID, categoriesID.CustomCategoryID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
}
