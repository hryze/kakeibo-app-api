package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

type HTTPError struct {
	Status    int     `json:"status"`
	ErrorList []error `json:"errors"`
}

type ValidationErrorMsg struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticationErrorMsg struct {
	Message string `json:"message"`
}

type BadRequestErrorMsg struct {
	Message string `json:"message"`
}

type InternalServerErrorMsg struct {
	Message string `json:"message"`
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	userHandler := UserHandler{userRepo: userRepo}
	return &userHandler
}

func NewHTTPError(status int, err interface{}) error {
	switch status {
	case http.StatusBadRequest:
		switch err := err.(type) {
		case *ValidationErrorMsg:
			return &HTTPError{
				Status:    status,
				ErrorList: []error{err},
			}
		default:
			return &HTTPError{
				Status:    status,
				ErrorList: []error{&BadRequestErrorMsg{"ログアウト済みです"}},
			}
		}
	case http.StatusConflict:
		return &HTTPError{
			Status:    status,
			ErrorList: []error{err.(*ValidationErrorMsg)},
		}
	case http.StatusUnauthorized:
		return &HTTPError{
			Status:    status,
			ErrorList: []error{&AuthenticationErrorMsg{"認証に失敗しました"}},
		}
	default:
		return &HTTPError{
			Status:    status,
			ErrorList: []error{&InternalServerErrorMsg{"500 Internal Server Error"}},
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

func (e *AuthenticationErrorMsg) Error() string {
	return e.Message
}

func (e *BadRequestErrorMsg) Error() string {
	return e.Message
}

func (e *InternalServerErrorMsg) Error() string {
	return e.Message
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
			validationErrorMsg.Email = "メールアドレスを正しく入力してください"
		case "Password":
			validationErrorMsg.Password = "パスワードを正しく入力してください"
		}
	}

	return &validationErrorMsg
}

func checkForUniqueUser(h *UserHandler, signUpUser *model.SignUpUser) error {
	var validationErrorMsg ValidationErrorMsg

	errID := h.userRepo.FindID(signUpUser)
	if errID != nil && errID != sql.ErrNoRows {
		return errID
	}

	errEmail := h.userRepo.FindEmail(signUpUser)
	if errEmail != nil && errEmail != sql.ErrNoRows {
		return errEmail
	}

	if errors.Is(errID, sql.ErrNoRows) && errors.Is(errEmail, sql.ErrNoRows) {
		return nil
	}

	if errID == nil && errEmail != nil {
		validationErrorMsg.ID = "このIDは既に利用されています"
		return &validationErrorMsg
	}

	if errEmail == nil && errID != nil {
		validationErrorMsg.Email = "このメールアドレスは既に利用されています"
		return &validationErrorMsg
	}

	validationErrorMsg.ID = "このIDは既に利用されています"
	validationErrorMsg.Email = "このメールアドレスは既に利用されています"
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

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpUser model.SignUpUser
	if err := json.NewDecoder(r.Body).Decode(&signUpUser); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := UserValidate(&signUpUser); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	if err := checkForUniqueUser(h, &signUpUser); err != nil {
		validationErrorMsg, ok := err.(*ValidationErrorMsg)
		if !ok {
			responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
		responseByJSON(w, nil, NewHTTPError(http.StatusConflict, validationErrorMsg))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(signUpUser.Password), 10)
	if err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	signUpUser.Password = string(hash)
	if err := h.userRepo.CreateUser(&signUpUser); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	signUpUser.Password = ""

	responseByJSON(w, &signUpUser, nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser model.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := UserValidate(&loginUser); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	password := loginUser.Password
	dbUser, err := h.userRepo.FindUser(&loginUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responseByJSON(w, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		} else if err != nil {
			responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}
	hashedPassword := dbUser.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusUnauthorized, nil))
		return
	}
	loginUser.Password = ""

	sessionID := uuid.New().String()
	expiration := 86400 * 30
	if err := h.userRepo.SetSessionID(sessionID, expiration); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(time.Duration(expiration) * time.Second),
		HttpOnly: true,
	})

	responseByJSON(w, &loginUser, nil)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		responseByJSON(w, nil, NewHTTPError(http.StatusBadRequest, nil))
		return
	}
	sessionID := cookie.Value
	if err := h.userRepo.DeleteSessionID(sessionID); err != nil {
		responseByJSON(w, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})
}
