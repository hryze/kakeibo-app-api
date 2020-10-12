package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/repository"
)

type DBHandler struct {
	HealthRepo            repository.HealthRepository
	AuthRepo              repository.AuthRepository
	TransactionsRepo      repository.TransactionsRepository
	CategoriesRepo        repository.CategoriesRepository
	BudgetsRepo           repository.BudgetsRepository
	GroupTransactionsRepo repository.GroupTransactionsRepository
	GroupCategoriesRepo   repository.GroupCategoriesRepository
	GroupBudgetsRepo      repository.GroupBudgetsRepository
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

func verifyGroupAffiliation(groupID int, userID string) error {
	userHost := os.Getenv("USER_HOST")
	requestURL := fmt.Sprintf("http://%s:8080/groups/%d/users/%s", userHost, groupID, userID)

	request, err := http.NewRequest(
		"GET",
		requestURL,
		nil,
	)
	if err != nil {
		return err
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
		return err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode == http.StatusBadRequest {
		return &BadRequestErrorMsg{"指定されたグループに所属していません。"}
	}

	if response.StatusCode == http.StatusInternalServerError {
		return &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return nil
}
