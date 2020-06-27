package handler

import (
	"encoding/json"
	"fmt"
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
	AuthenticationError *AuthenticationErrorMsg `json:"auth_error,omitempty"`
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

type AuthenticationErrorMsg struct {
	Message string `json:"message"`
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	userHandler := UserHandler{userRepo: userRepo}
	fmt.Println()
	return &userHandler
}

func NewHTTPError(status int, err interface{}) error {
	switch status {
	case http.StatusBadRequest, http.StatusConflict:
		return &HTTPError{
			Status:          status,
			ValidationError: err.(*ValidationErrorMsg),
		}
	case http.StatusUnauthorized:
		return &HTTPError{
			Status:              status,
			AuthenticationError: &AuthenticationErrorMsg{"認証に失敗しました"},
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

func UserValidate(user interface{}) error {
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
			validationErrorMsg.ID = "IDを正しく入力してください"
		case "Name":
			validationErrorMsg.Name = "名前を正しく入力してください"
		case "Email":
			validationErrorMsg.Email = "Eメールを正しく入力してください"
		case "Password":
			validationErrorMsg.Password = "パスワードを正しく入力してください"
		}
	}

	return &validationErrorMsg
}

func checkForUniqueID(h *UserHandler, signUpUser *model.SignUpUser) error {
	var validationErrorMsg ValidationErrorMsg
	dbID, err := h.userRepo.FindID(signUpUser)
	if len(dbID) == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	validationErrorMsg.ID = "このIDは登録できません"

	return &validationErrorMsg
}

func responseByJSON(w http.ResponseWriter, user interface{}, err error) {
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
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func ResponseByJSONMiddleware(fn interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch fn := fn.(type) {
		case func(http.ResponseWriter, *http.Request) (*model.SignUpUser, error):
			signUpUser, err := fn(w, r)
			responseByJSON(w, signUpUser, err)
		case func(http.ResponseWriter, *http.Request) (*model.LoginUser, error):
			loginUser, err := fn(w, r)
			responseByJSON(w, loginUser, err)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) (*model.SignUpUser, error) {
	var signUpUser model.SignUpUser
	if err := json.NewDecoder(r.Body).Decode(&signUpUser); err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, nil)
	}
	if err := UserValidate(&signUpUser); err != nil {
		return nil, NewHTTPError(http.StatusBadRequest, err)
	}
	if err := checkForUniqueID(h, &signUpUser); err != nil {
		validationErrorMsg, ok := err.(*ValidationErrorMsg)
		if !ok {
			return nil, NewHTTPError(http.StatusInternalServerError, nil)
		}
		return nil, NewHTTPError(http.StatusConflict, validationErrorMsg)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(signUpUser.Password), 10)
	if err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, nil)
	}
	signUpUser.Password = string(hash)
	if err := h.userRepo.CreateUser(&signUpUser); err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, nil)
	}
	signUpUser.Password = ""

	return &signUpUser, nil
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) (*model.LoginUser, error) {
	var loginUser model.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, nil)
	}
	if err := UserValidate(&loginUser); err != nil {
		return nil, NewHTTPError(http.StatusBadRequest, err)
	}
	password := loginUser.Password
	dbLoginUser, err := h.userRepo.FindUser(&loginUser)
	if err != nil {
		return nil, NewHTTPError(http.StatusInternalServerError, nil)
	}
	if dbLoginUser == nil {
		return nil, NewHTTPError(http.StatusUnauthorized, nil)
	}
	hashedPassword := dbLoginUser.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, NewHTTPError(http.StatusUnauthorized, nil)
	}
	dbLoginUser.Password = ""

	return dbLoginUser, nil
}
