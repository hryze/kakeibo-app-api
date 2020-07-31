package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
)

type DBHandler struct {
	DBRepo repository.DBRepository
}

type HTTPError struct {
	Status       int   `json:"status"`
	ErrorMessage error `json:"error"`
}

type BadRequestErrorMsg struct {
	Message string `json:"message"`
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

func NewHTTPError(status int, err error) error {
	switch status {
	case http.StatusBadRequest:
		switch err := err.(type) {
		case *ValidationErrorMsg:
			return &HTTPError{
				Status:       status,
				ErrorMessage: err,
			}
		default:
			return &HTTPError{
				Status:       status,
				ErrorMessage: &BadRequestErrorMsg{"ログアウト済みです"},
			}
		}
	case http.StatusConflict:
		return &HTTPError{
			Status:       status,
			ErrorMessage: err.(*ValidationErrorMsg),
		}
	case http.StatusUnauthorized:
		return &HTTPError{
			Status:       status,
			ErrorMessage: &AuthenticationErrorMsg{"認証に失敗しました"},
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

func (e *BadRequestErrorMsg) Error() string {
	return e.Message
}

func (e *AuthenticationErrorMsg) Error() string {
	return e.Message
}

func (e *InternalServerErrorMsg) Error() string {
	return e.Message
}

func errorResponseByJSON(w http.ResponseWriter, err error) {
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
}
