package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/repository"
)

type DBHandler struct {
	HealthRepo            repository.HealthRepository
	AuthRepo              repository.AuthRepository
	TodoRepo              repository.TodoRepository
	ShoppingListRepo      repository.ShoppingListRepository
	GroupTodoRepo         repository.GroupTodoRepository
	GroupShoppingListRepo repository.GroupShoppingListRepository
	GroupTasksRepo        repository.GroupTasksRepository
	TimeManage            TimeManager
}

type TimeManager interface {
	Now() time.Time
}

type RealTime struct{}

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

type ConflictErrorMsg struct {
	Message string `json:"message"`
}

type InternalServerErrorMsg struct {
	Message string `json:"message"`
}

func NewRealTime() *RealTime {
	return &RealTime{}
}

func (r *RealTime) Now() time.Time {
	return time.Now()
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
	cookie, err := r.Cookie(config.Env.Cookie.Name)
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
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/users/%s/verify",
		config.Env.UserApi.Host, config.Env.UserApi.Port, groupID, userID,
	)

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
