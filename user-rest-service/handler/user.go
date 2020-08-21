package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type Users interface {
	ShowUser() (string, error)
}

type LogoutMsg struct {
	Message string `json:"message"`
}

type UserValidationErrorMsg struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserConflictErrorMsg struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}

func (e *UserValidationErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func (e *UserConflictErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func validateUser(users Users) error {
	var userValidationErrorMsg UserValidationErrorMsg
	validate := validator.New()
	err := validate.Struct(users)
	if err == nil {
		return nil
	}
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		switch fieldName {
		case "ID":
			userValidationErrorMsg.ID = "IDを正しく入力してください"
		case "Name":
			userValidationErrorMsg.Name = "名前を正しく入力してください"
		case "Email":
			userValidationErrorMsg.Email = "メールアドレスを正しく入力してください"
		case "Password":
			userValidationErrorMsg.Password = "パスワードを正しく入力してください"
		}
	}

	return &userValidationErrorMsg
}

func checkForUniqueUser(h *DBHandler, signUpUser *model.SignUpUser) error {
	var userConflictErrorMsg UserConflictErrorMsg

	errID := h.UserRepo.FindUserID(signUpUser.ID)
	if errID != nil && errID != sql.ErrNoRows {
		return errID
	}

	errEmail := h.UserRepo.FindEmail(signUpUser.Email)
	if errEmail != nil && errEmail != sql.ErrNoRows {
		return errEmail
	}

	if errors.Is(errID, sql.ErrNoRows) && errors.Is(errEmail, sql.ErrNoRows) {
		return nil
	}

	if errID == nil && errEmail != nil {
		userConflictErrorMsg.ID = "このIDは既に利用されています"
		return &userConflictErrorMsg
	}

	if errEmail == nil && errID != nil {
		userConflictErrorMsg.Email = "このメールアドレスは既に利用されています"
		return &userConflictErrorMsg
	}

	userConflictErrorMsg.ID = "このIDは既に利用されています"
	userConflictErrorMsg.Email = "このメールアドレスは既に利用されています"
	return &userConflictErrorMsg
}

func postInitStandardBudgets(userID string) error {
	request, err := http.NewRequest(
		"POST",
		"http://localhost:8081/standard-budgets",
		bytes.NewBuffer([]byte(`{"user_id":"`+userID+`"}`)),
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
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode == http.StatusCreated {
		return nil
	}

	return errors.New("couldn't create a standard budget")
}

func (h *DBHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpUser model.SignUpUser
	if err := json.NewDecoder(r.Body).Decode(&signUpUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := validateUser(&signUpUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	if err := checkForUniqueUser(h, &signUpUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusConflict, err))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(signUpUser.Password), 10)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	signUpUser.Password = string(hash)
	if err := h.UserRepo.CreateUser(&signUpUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := postInitStandardBudgets(signUpUser.ID); err != nil {
		if err := h.UserRepo.DeleteUser(&signUpUser); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	signUpUser.Password = ""
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&signUpUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser model.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := validateUser(&loginUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	password := loginUser.Password
	dbUser, err := h.UserRepo.FindUser(&loginUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"認証に失敗しました"}))
			return
		} else if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}
	hashedPassword := dbUser.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"認証に失敗しました"}))
		return
	}
	loginUser.Password = ""

	sessionID := uuid.New().String()
	expiration := 86400 * 30
	if err := h.UserRepo.SetSessionID(sessionID, loginUser.ID, expiration); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(time.Duration(expiration) * time.Second),
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&loginUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"ログアウト済みです"}))
		return
	}
	sessionID := cookie.Value
	if err := h.UserRepo.DeleteSessionID(sessionID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&LogoutMsg{"ログアウトしました"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
