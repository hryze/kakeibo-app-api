package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

type HTTPError struct {
	Status              int                     `json:"status"`
	ValidationError     *ValidationErrorMsg     `json:"errors,omitempty"`
	InternalServerError *InternalServerErrorMsg `json:"error,omitempty"`
}

type ValidationErrorMsg struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InternalServerErrorMsg struct {
	Message string `json:"message"`
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	userHandler := UserHandler{userRepo: userRepo}
	return &userHandler
}

func NewHTTPError(status int, err interface{}) error {
	switch err := err.(type) {
	case *ValidationErrorMsg:
		return &HTTPError{
			Status:          status,
			ValidationError: err,
		}
	default:
		return &HTTPError{
			Status:              status,
			InternalServerError: &InternalServerErrorMsg{"500 Internal Server Error"},
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
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func (e *InternalServerErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func UserValidate(user *model.User) error {
	var validationErrorMsg ValidationErrorMsg
	validate := validator.New()
	err := validate.Struct(user)
	if err == nil {
		return nil
	}
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		switch fieldName {
		case "ID":
			validationErrorMsg.ID = "ユーザーIDが正しくありません"
		case "Name":
			validationErrorMsg.Name = "ユーザーネームが正しくありません"
		case "Email":
			validationErrorMsg.Email = "ユーザーメールが正しくありません"
		case "Password":
			validationErrorMsg.Password = "パスワードが正しくありません"
		}
	}

	return &validationErrorMsg
}

func checkForUniqueID(h *UserHandler, user *model.User) (*ValidationErrorMsg, error) {
	var validationErrorMsg ValidationErrorMsg
	dbID, err := h.userRepo.FindID(user)
	if len(dbID) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	validationErrorMsg.ID = "このユーザーIDは登録できません"

	return &validationErrorMsg, nil
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
		return nil, NewHTTPError(http.StatusInternalServerError, err)
	}
	if validationErrorMsg := UserValidate(&user); validationErrorMsg != nil {
		return nil, NewHTTPError(http.StatusBadRequest, validationErrorMsg)
	}
	validationErrorMsg, err := checkForUniqueID(h, &user)
	if err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, err)
	}
	if validationErrorMsg != nil {
		return nil, NewHTTPError(http.StatusConflict, validationErrorMsg)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, err)
	}
	user.Password = string(hash)
	if err := h.userRepo.CreateUser(&user); err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, err)
	}
	user.Password = ""

	return &user, nil
}
