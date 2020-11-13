package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
)

type DBHandler struct {
	HealthRepo repository.HealthRepository
	AuthRepo   repository.AuthRepository
	UserRepo   repository.UserRepository
	GroupRepo  repository.GroupRepository
}

type DeleteContentMsg struct {
	Message string `json:"message"`
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

type NotFoundErrorMsg struct {
	Message string `json:"message"`
}

type ConflictErrorMsg struct {
	Message string `json:"message"`
}

type InternalServerErrorMsg struct {
	Message string `json:"message"`
}

func NewHTTPError(status int, err error) error {
	switch status {
	case http.StatusInternalServerError:
		return &HTTPError{
			Status:       status,
			ErrorMessage: &InternalServerErrorMsg{"500 Internal Server Error"},
		}
	default:
		return &HTTPError{
			Status:       status,
			ErrorMessage: err,
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

func (e *NotFoundErrorMsg) Error() string {
	return e.Message
}

func (e *ConflictErrorMsg) Error() string {
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

func verifySessionID(h *DBHandler, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}

	sessionID := cookie.Value
	userID, err := h.AuthRepo.GetUserID(sessionID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
