package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/repository"
)

type DBHandler struct {
	DBRepo repository.DBRepository
}

type HTTPError struct {
	Status       int   `json:"status"`
	ErrorMessage error `json:"error"`
}

type AuthenticationErrorMsg struct {
	Message string `json:"message"`
}

type InternalServerErrorMsg struct {
	Message string `json:"message"`
}

func NewDBHandler(DBRepo repository.DBRepository) *DBHandler {
	DBHandler := DBHandler{DBRepo: DBRepo}
	return &DBHandler
}

func NewHTTPError(status int, err interface{}) error {
	switch status {
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

func (e *AuthenticationErrorMsg) Error() string {
	return e.Message
}

func (e *InternalServerErrorMsg) Error() string {
	return e.Message
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

func (h *DBHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		responseByJSON(w, nil, NewHTTPError(http.StatusUnauthorized, nil))
		return
	}
	sessionID := cookie.Value
	userID, err := h.DBRepo.GetUserID(sessionID)
	if err != nil {
		if err == redis.ErrNil {
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
