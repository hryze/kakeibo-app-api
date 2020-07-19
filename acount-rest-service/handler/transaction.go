package handler

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"

	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"
)

type DeleteTransactionMsg struct {
	Message string `json:"message"`
}

type TransactionValidationErrorMsg struct {
	Message []string `json:"message"`
}

func (e *TransactionValidationErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func validateTransaction(transactionReceiver *model.TransactionReceiver) error {
	var transactionValidationErrorMsg TransactionValidationErrorMsg

	validate := validator.New()
	validate.RegisterCustomTypeFunc(ValidateValuer, model.Date{}, model.NullString{}, model.NullInt64{})
	validate.RegisterValidation("blank", blankValidation)
	validate.RegisterValidation("date", dateValidation)
	validate.RegisterValidation("either_id", eitherIDValidation)
	err := validate.Struct(transactionReceiver)
	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		var errorMessage string

		fieldName := err.Field()
		switch fieldName {
		case "TransactionType":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "取引タイプが選択されていません。"
			case "oneof":
				errorMessage = "取引タイプを正しく選択してください。"
			}
		case "TransactionDate":
			errorMessage = "日付を正しく選択してください。"
		case "Shop":
			tagName := err.Tag()
			switch tagName {
			case "max":
				errorMessage = "店名は20文字以内で入力してください。"
			case "blank":
				errorMessage = "店名の文字列先頭か末尾に空白がないか確認してください。"
			}
		case "Memo":
			tagName := err.Tag()
			switch tagName {
			case "max":
				errorMessage = "メモは50文字以内で入力してください"
			case "blank":
				errorMessage = "メモの文字列先頭か末尾に空白がないか確認してください。"
			}
		case "Amount":
			errorMessage = "金額が入力されていません。"
		case "BigCategoryID":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "カテゴリーが選択されていません。"
			case "min", "max":
				errorMessage = "カテゴリーを正しく選択してください。"
			case "either_id":
				errorMessage = "中カテゴリーを正しく選択してください。"
			}
		case "MediumCategoryID":
			errorMessage = "中カテゴリーを正しく選択してください。"
		case "CustomCategoryID":
			errorMessage = "中カテゴリーを正しく選択してください。"
		}
		transactionValidationErrorMsg.Message = append(transactionValidationErrorMsg.Message, errorMessage)
	}

	return &transactionValidationErrorMsg
}

func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}

func blankValidation(fl validator.FieldLevel) bool {
	text := fl.Field().String()

	if strings.HasPrefix(text, " ") || strings.HasPrefix(text, "　") || strings.HasSuffix(text, " ") || strings.HasSuffix(text, "　") {
		return false
	}

	return true
}

func dateValidation(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	stringDate := date.String()
	trimDate := strings.Trim(string(stringDate), "\"")[:10]

	minDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

	dateTime, err := time.Parse("2006-01-02", trimDate)
	if err != nil {
		return false
	}
	if dateTime.Before(minDate) || dateTime.After(maxDate) {
		return false
	}

	return true
}

func eitherIDValidation(fl validator.FieldLevel) bool {
	transactionReceiver, ok := fl.Parent().Interface().(*model.TransactionReceiver)
	if !ok {
		return false
	}

	if transactionReceiver.MediumCategoryID.Valid && transactionReceiver.CustomCategoryID.Valid {
		return false
	}

	if transactionReceiver.CustomCategoryID.Valid {
		return true
	}

	if transactionReceiver.MediumCategoryID.Valid {
		return true
	}

	return false
}

func (h *DBHandler) GetMonthlyTransactionsList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbTransactionsList, err := h.DBRepo.GetMonthlyTransactionsList(userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	transactionsList := model.NewTransactionsList(dbTransactionsList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&transactionsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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

	if err := validateTransaction(&transactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	result, err := h.DBRepo.PostTransaction(&transactionReceiver, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

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

func (h *DBHandler) PutTransaction(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
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

	transactionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.DBRepo.PutTransaction(&transactionReceiver, transactionID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var transactionSender model.TransactionSender
	dbTransactionSender, err := h.DBRepo.GetTransaction(&transactionSender, transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbTransactionSender); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	transactionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.DBRepo.DeleteTransaction(transactionID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteTransactionMsg{"トランザクションを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
