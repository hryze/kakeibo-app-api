package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/repository"
)

type DBHandler struct {
	DBRepo repository.DBRepository
}

type DeleteCustomCategoryMsg struct {
	Message string `json:"message"`
}

type HTTPError struct {
	Status       int   `json:"status"`
	ErrorMessage error `json:"error"`
}

type AuthenticationErrorMsg struct {
	Message string `json:"message"`
}

type ValidationErrorMsg struct {
	Message string `json:"message"`
}
type ConflictErrorMsg struct {
	Message string `json:"message"`
}

type InternalServerErrorMsg struct {
	Message string `json:"message"`
}

func NewDBHandler(DBRepo repository.DBRepository) *DBHandler {
	DBHandler := DBHandler{DBRepo: DBRepo}
	return &DBHandler
}

func NewHTTPError(status int, err error) error {
	switch status {
	case http.StatusBadRequest:
		return &HTTPError{
			Status:       status,
			ErrorMessage: err.(*ValidationErrorMsg),
		}
	case http.StatusConflict:
		return &HTTPError{
			Status:       status,
			ErrorMessage: &ConflictErrorMsg{"中カテゴリーの登録に失敗しました。 同じカテゴリー名が既に存在していないか確認してください。"},
		}
	case http.StatusUnauthorized:
		return &HTTPError{
			Status:       status,
			ErrorMessage: &AuthenticationErrorMsg{"このページを表示するにはログインが必要です"},
		}
	default:
		return &HTTPError{
			Status:       status,
			ErrorMessage: &InternalServerErrorMsg{"500 Internal Server Error"},
		}
	}
}

func (e *HTTPError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func (e *ValidationErrorMsg) Error() string {
	return e.Message
}

func (e *ConflictErrorMsg) Error() string {
	return e.Message
}

func (e *AuthenticationErrorMsg) Error() string {
	return e.Message
}

func (e *InternalServerErrorMsg) Error() string {
	return e.Message
}

func validateCustomCategory(customCategory *model.CustomCategory) error {
	if strings.HasPrefix(customCategory.Name, " ") || strings.HasPrefix(customCategory.Name, "　") {
		return &ValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列先頭に空白がないか確認してください。"}
	}
	if strings.HasSuffix(customCategory.Name, " ") || strings.HasSuffix(customCategory.Name, "　") {
		return &ValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列末尾に空白がないか確認してください。"}
	}
	if utf8.RuneCountInString(customCategory.Name) > 9 {
		return &ValidationErrorMsg{"カテゴリー名は9文字以下で入力してください。"}
	}
	return nil
}

func checkForUniqueCustomCategory(h *DBHandler, customCategory *model.CustomCategory, userID string) error {
	if err := h.DBRepo.FindCustomCategory(customCategory, userID); err != nil {
		return err
	}
	return nil
}

func responseByJSON(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		httpError, ok := err.(*HTTPError)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(httpError.Status)
		if err := json.NewEncoder(w).Encode(httpError); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func verifySessionID(h *DBHandler, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	sessionID := cookie.Value
	userID, err := h.DBRepo.GetUserID(sessionID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (h *DBHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	bigCategoriesList, err := h.DBRepo.GetBigCategoriesList()
	if err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	mediumCategoriesList, err := h.DBRepo.GetMediumCategoriesList()
	if err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategoriesList, err := h.DBRepo.GetCustomCategoriesList(userID)
	if err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	for i, bigCategory := range bigCategoriesList {
		for _, customCategory := range customCategoriesList {
			if bigCategory.ID == customCategory.BigCategoryID {
				bigCategoriesList[i].AssociatedCategoriesList = append(bigCategoriesList[i].AssociatedCategoriesList, customCategory)
			}
		}
	}
	for i, bigCategory := range bigCategoriesList {
		for _, mediumCategory := range mediumCategoriesList {
			if bigCategory.ID == mediumCategory.BigCategoryID {
				bigCategoriesList[i].AssociatedCategoriesList = append(bigCategoriesList[i].AssociatedCategoriesList, mediumCategory)
			}
		}
	}
	CategoriesList := model.NewCategoriesList(bigCategoriesList)
	responseByJSON(w, CategoriesList, nil)
}

func (h *DBHandler) PostCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory := model.NewCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&customCategory); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := validateCustomCategory(&customCategory); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	if err := checkForUniqueCustomCategory(h, &customCategory, userID); err != sql.ErrNoRows {
		if err == nil {
			responseByJSON(w, nil, NewHTTPError(http.StatusConflict, nil))
			return
		}
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	result, err := h.DBRepo.PostCustomCategory(&customCategory, userID)
	if err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory.ID = int(lastInsertId)
	responseByJSON(w, &customCategory, nil)
}

func (h *DBHandler) PutCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory := model.NewCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&customCategory); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := validateCustomCategory(&customCategory); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	if err := checkForUniqueCustomCategory(h, &customCategory, userID); err != sql.ErrNoRows {
		if err == nil {
			responseByJSON(w, nil, NewHTTPError(http.StatusConflict, nil))
			return
		}
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := h.DBRepo.PutCustomCategory(&customCategory, userID); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	responseByJSON(w, &customCategory, nil)
}

func (h *DBHandler) DeleteCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory := model.NewCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&customCategory); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := h.DBRepo.DeleteCustomCategory(&customCategory, userID); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	responseByJSON(w, &DeleteCustomCategoryMsg{"カスタムカテゴリーを削除しました。"}, nil)
}
