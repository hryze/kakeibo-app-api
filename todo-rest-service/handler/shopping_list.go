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
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"

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

type RegularShoppingItemValidationErrorMsg struct {
	Message []string `json:"message"`
}

func (e *RegularShoppingItemValidationErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}

type ShoppingItemValidationErrorMsg struct {
	Message []string `json:"message"`
}

func (e *ShoppingItemValidationErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}

func validateRegularShoppingItem(regularShoppingItem model.RegularShoppingItem) error {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(regularShoppingItemValidateValuer, model.Date{}, model.NullString{}, model.NullInt{}, model.NullInt64{})
	if err := validate.RegisterValidation("blank", blankValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("date_range", dateRangeValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("either_id", eitherCategoryIDValidation); err != nil {
		return err
	}

	err := validate.Struct(regularShoppingItem)
	if err == nil {
		return nil
	}

	var regularShoppingItemValidationErrorMsg RegularShoppingItemValidationErrorMsg
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
		}

		regularShoppingItemValidationErrorMsg.Message = append(regularShoppingItemValidationErrorMsg.Message, errorMessage)
	}

	return &regularShoppingItemValidationErrorMsg
}

func regularShoppingItemValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}

	return nil
}

func validateShoppingItem(shoppingItem model.ShoppingItem) error {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(shoppingItemValidateValuer, model.Date{}, model.NullString{}, model.NullInt64{})
	if err := validate.RegisterValidation("blank", blankValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("date_range", dateRangeValidation); err != nil {
		return err
	}

	if err := validate.RegisterValidation("either_id", eitherCategoryIDValidation); err != nil {
		return err
	}

	err := validate.Struct(shoppingItem)
	if err == nil {
		return nil
	}

	var shoppingItemValidationErrorMsg ShoppingItemValidationErrorMsg
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
		}

		shoppingItemValidationErrorMsg.Message = append(shoppingItemValidationErrorMsg.Message, errorMessage)
	}

	return &shoppingItemValidationErrorMsg
}

func shoppingItemValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}

	return nil
}

func dateRangeValidation(fl validator.FieldLevel) bool {
	flDate, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	var today time.Time

	switch item := fl.Parent().Interface().(type) {
	case model.RegularShoppingItem:
		today = item.Today
	case model.ShoppingItem:
		today = item.Today
	case model.GroupRegularShoppingItem:
		today = item.Today
	case model.GroupShoppingItem:
		today = item.Today
	default:
		return false
	}

	return !today.After(flDate)
}

func eitherCategoryIDValidation(fl validator.FieldLevel) bool {
	var mediumCategoryIDValid bool
	var customCategoryIDValid bool

	switch item := fl.Parent().Interface().(type) {
	case model.RegularShoppingItem:
		mediumCategoryIDValid = item.MediumCategoryID.Valid
		customCategoryIDValid = item.CustomCategoryID.Valid
	case model.ShoppingItem:
		mediumCategoryIDValid = item.MediumCategoryID.Valid
		customCategoryIDValid = item.CustomCategoryID.Valid
	case model.GroupRegularShoppingItem:
		mediumCategoryIDValid = item.MediumCategoryID.Valid
		customCategoryIDValid = item.CustomCategoryID.Valid
	case model.GroupShoppingItem:
		mediumCategoryIDValid = item.MediumCategoryID.Valid
		customCategoryIDValid = item.CustomCategoryID.Valid
	default:
		return false
	}

	if mediumCategoryIDValid && customCategoryIDValid {
		return false
	}

	if mediumCategoryIDValid || customCategoryIDValid {
		return true
	}

	return false
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

func getShoppingItemRelatedTransactionDataList(transactionIdList []int64) ([]*model.TransactionData, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/transactions/related-shopping-list", accountHost)

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

	var shoppingItemRelatedTransactionDataList []*model.TransactionData
	if err := json.NewDecoder(response.Body).Decode(&shoppingItemRelatedTransactionDataList); err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return nil, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return shoppingItemRelatedTransactionDataList, nil
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

func generateRegularShoppingList(regularShoppingList model.RegularShoppingList) (model.RegularShoppingList, error) {
	categoriesIdList := make([]CategoriesID, len(regularShoppingList.RegularShoppingList))

	for i, regularShoppingItem := range regularShoppingList.RegularShoppingList {
		categoriesIdList[i] = CategoriesID{
			MediumCategoryID: regularShoppingItem.MediumCategoryID,
			CustomCategoryID: regularShoppingItem.CustomCategoryID,
		}
	}

	categoriesNameListBytes, err := getShoppingItemCategoriesNameList(categoriesIdList)
	if err != nil {
		return regularShoppingList, err
	}

	if err := json.Unmarshal(categoriesNameListBytes, &regularShoppingList.RegularShoppingList); err != nil {
		return regularShoppingList, err
	}

	return regularShoppingList, nil
}

func generateShoppingList(shoppingList model.ShoppingList) (model.ShoppingList, error) {
	categoriesIdList := make([]CategoriesID, len(shoppingList.ShoppingList))
	var transactionIdList []int64

	for i, shoppingItem := range shoppingList.ShoppingList {
		categoriesIdList[i] = CategoriesID{
			MediumCategoryID: shoppingItem.MediumCategoryID,
			CustomCategoryID: shoppingItem.CustomCategoryID,
		}

		if shoppingItem.RelatedTransactionData != nil {
			transactionIdList = append(transactionIdList, shoppingItem.RelatedTransactionData.ID.Int64)
		}
	}

	categoriesNameListBytes, err := getShoppingItemCategoriesNameList(categoriesIdList)
	if err != nil {
		return shoppingList, err
	}

	if err := json.Unmarshal(categoriesNameListBytes, &shoppingList.ShoppingList); err != nil {
		return shoppingList, err
	}

	if len(transactionIdList) != 0 {
		shoppingItemRelatedTransactionDataList, err := getShoppingItemRelatedTransactionDataList(transactionIdList)
		if err != nil {
			return shoppingList, err
		}

		for _, shoppingItemRelatedTransactionData := range shoppingItemRelatedTransactionDataList {
			for i, shoppingItem := range shoppingList.ShoppingList {
				if shoppingItem.RelatedTransactionData != nil && shoppingItemRelatedTransactionData.ID.Int64 == shoppingItem.RelatedTransactionData.ID.Int64 {
					shoppingList.ShoppingList[i].RelatedTransactionData = shoppingItemRelatedTransactionData
				}
			}
		}
	}

	return shoppingList, nil
}

func generateShoppingListByCategories(shoppingList model.ShoppingList) []model.ShoppingListByCategory {
	shoppingListByCategories := make([]model.ShoppingListByCategory, 0)

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
		}

		if i == len(shoppingList.ShoppingList)-1 {
			if firstIndexBigCategoryID == currentIndexBigCategoryID {
				shoppingListByCategory := model.ShoppingListByCategory{
					BigCategoryName: shoppingList.ShoppingList[j].BigCategoryName,
					ShoppingList:    append(make([]model.ShoppingItem, 0, i-j), shoppingList.ShoppingList[j:]...),
				}

				shoppingListByCategories = append(shoppingListByCategories, shoppingListByCategory)
			} else if i == j {
				shoppingListByCategory := model.ShoppingListByCategory{
					BigCategoryName: shoppingList.ShoppingList[i].BigCategoryName,
					ShoppingList:    []model.ShoppingItem{shoppingList.ShoppingList[i]},
				}

				shoppingListByCategories = append(shoppingListByCategories, shoppingListByCategory)
			}
		}
	}

	return shoppingListByCategories
}

func (h *DBHandler) GetDailyShoppingDataByDay(w http.ResponseWriter, r *http.Request) {
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

		regularShoppingList, err = generateRegularShoppingList(regularShoppingList)
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

	shoppingList, err := h.ShoppingListRepo.GetDailyShoppingListByDay(date, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(shoppingList.ShoppingList) != 0 {
		shoppingList, err = generateShoppingList(shoppingList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingData := model.ShoppingDataByDay{
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

func (h *DBHandler) GetDailyShoppingDataByCategory(w http.ResponseWriter, r *http.Request) {
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

		regularShoppingList, err = generateRegularShoppingList(regularShoppingList)
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

	shoppingList, err := h.ShoppingListRepo.GetDailyShoppingListByCategory(date, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(shoppingList.ShoppingList) != 0 {
		shoppingList, err = generateShoppingList(shoppingList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingListByCategories := generateShoppingListByCategories(shoppingList)

	shoppingDataByCategories := model.ShoppingDataByCategory{
		RegularShoppingList:    regularShoppingList,
		ShoppingListByCategory: shoppingListByCategories,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingDataByCategories); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetMonthlyShoppingDataByDay(w http.ResponseWriter, r *http.Request) {
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

		regularShoppingList, err = generateRegularShoppingList(regularShoppingList)
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

	shoppingList, err := h.ShoppingListRepo.GetMonthlyShoppingListByDay(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(shoppingList.ShoppingList) != 0 {
		shoppingList, err = generateShoppingList(shoppingList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingData := model.ShoppingDataByDay{
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

func (h *DBHandler) GetMonthlyShoppingDataByCategory(w http.ResponseWriter, r *http.Request) {
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

		regularShoppingList, err = generateRegularShoppingList(regularShoppingList)
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

	shoppingList, err := h.ShoppingListRepo.GetMonthlyShoppingListByCategory(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(shoppingList.ShoppingList) != 0 {
		shoppingList, err = generateShoppingList(shoppingList)
		if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	shoppingListByCategories := generateShoppingListByCategories(shoppingList)

	shoppingDataByCategories := model.ShoppingDataByCategory{
		RegularShoppingList:    regularShoppingList,
		ShoppingListByCategory: shoppingListByCategories,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&shoppingDataByCategories); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetExpiredShoppingList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	now := h.TimeManage.Now()
	dueDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	expiredShoppingList, err := h.ShoppingListRepo.GetExpiredShoppingList(dueDate, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(expiredShoppingList.ExpiredShoppingList) != 0 {
		categoriesIdList := make([]CategoriesID, len(expiredShoppingList.ExpiredShoppingList))

		for i, shoppingItem := range expiredShoppingList.ExpiredShoppingList {
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

		if err := json.Unmarshal(categoriesNameListBytes, &expiredShoppingList.ExpiredShoppingList); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&expiredShoppingList); err != nil {
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

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	shoppingItem.Today = today

	if err := validateShoppingItem(shoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
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

	if err := validateShoppingItem(shoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
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

	regularShoppingItem.Today = today

	if err := validateRegularShoppingItem(regularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

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

	shoppingList, err := h.ShoppingListRepo.GetShoppingListRelatedToPostedRegularShoppingItem(int(todayShoppingItemID), int(laterThanTodayShoppingItemID))
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

	if err := validateRegularShoppingItem(regularShoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if err := h.ShoppingListRepo.PutRegularShoppingItem(&regularShoppingItem, regularShoppingItemID, userID, today); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	regularShoppingItem, err = h.ShoppingListRepo.GetRegularShoppingItem(regularShoppingItemID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingList, err := h.ShoppingListRepo.GetShoppingListRelatedToUpdatedRegularShoppingItem(regularShoppingItemID)
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

func (h *DBHandler) PutShoppingListCustomCategoryIdToMediumCategoryId(w http.ResponseWriter, r *http.Request) {
	categoriesID := struct {
		MediumCategoryID int `json:"medium_category_id"`
		CustomCategoryID int `json:"custom_category_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&categoriesID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.ShoppingListRepo.PutShoppingListCustomCategoryIdToMediumCategoryId(categoriesID.MediumCategoryID, categoriesID.CustomCategoryID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
}
