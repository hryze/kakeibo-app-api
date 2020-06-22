package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

type HTTPError struct {
	Status  int       `json:"-"`
	Message *ErrorMsg `json:"message"`
}

type ErrorMsg struct {
	ID       string `json:"error_id"`
	Name     string `json:"error_name"`
	Email    string `json:"error_email"`
	Password string `json:"error_password"`
}

func NewUserHandler(userRepo repository.UserRepository) UserHandler {
	userHandler := UserHandler{userRepo: userRepo}
	return userHandler
}

func NewHTTPError(status int, message *ErrorMsg) error {
	return &HTTPError{
		Status:  status,
		Message: message,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintln("HTTPError")
}

func UserValidate(user *model.User) *ErrorMsg {
	var errorMsg ErrorMsg
	validate := validator.New()
	err := validate.Struct(user)
	if err == nil {
		return nil
	}
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		switch fieldName {
		case "ID":
			errorMsg.ID = "ユーザーIDが正しくありません"
		case "Name":
			errorMsg.Name = "ユーザーネームが正しくありません"
		case "Email":
			errorMsg.Email = "ユーザーメールが正しくありません"
		case "Password":
			errorMsg.Password = "パスワードが正しくありません"
		}
	}

	return &errorMsg
}

func checkForUniqueID(h *UserHandler, user *model.User) (*ErrorMsg, error) {
	var errorMsg ErrorMsg
	find, err := h.userRepo.FindID(user)
	if find == true {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	errorMsg.ID = "このユーザーIDは登録できません"

	return &errorMsg, nil
}

func responseByJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

func ErrorCheckHandler(fn func(http.ResponseWriter, *http.Request) (*model.User, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := fn(w, r)
		if err == nil {
			if err := responseByJSON(w, http.StatusOK, user); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}
		httpError, ok := err.(*HTTPError)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := responseByJSON(w, httpError.Status, httpError); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) (*model.User, error) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	if errorMsg := UserValidate(&user); errorMsg != nil {
		return nil, NewHTTPError(http.StatusBadRequest, errorMsg)
	}
	errorMsg, err := checkForUniqueID(h, &user)
	if err != nil {
		return nil, err
	}
	if errorMsg != nil {
		return nil, NewHTTPError(http.StatusConflict, errorMsg)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)
	if err := h.userRepo.CreateUser(&user); err != nil {
		return nil, err
	}
	user.Password = ""

	return &user, nil
}
