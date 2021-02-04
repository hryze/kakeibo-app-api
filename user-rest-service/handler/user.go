package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	uerrors "github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/errors"

	herrors "github.com/paypay3/kakeibo-app-api/user-rest-service/handler/errors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"golang.org/x/xerrors"

	"github.com/garyburd/redigo/redis"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type Users interface {
	ShowUser() (string, error)
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
	validate := validator.New()
	err := validate.Struct(users)
	if err == nil {
		return nil
	}

	var userValidationErrorMsg UserValidationErrorMsg
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

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *userHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var in input.SignUpUser
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		herrors.ErrorResponseByJSON(w, uerrors.NewInternalServerError(uerrors.NewErrorString("ユーザー登録に失敗しました")))
		return
	}

	out, err := h.userUsecase.SignUp(&in)
	if err != nil {
		herrors.ErrorResponseByJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(out); err != nil {
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
		if xerrors.Is(err, sql.ErrNoRows) {
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

	dbUser.Password = ""

	sessionID := uuid.New().String()
	expiration := 86400 * 30
	if err := h.UserRepo.SetSessionID(sessionID, dbUser.ID, expiration); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	domain := os.Getenv("COOKIE_DOMAIN")
	secure := true

	if domain != "shakepiper.com" {
		secure = false
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Domain:   domain,
		Expires:  time.Now().Add(time.Duration(expiration) * time.Second),
		Secure:   secure,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&dbUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if xerrors.Is(err, http.ErrNoCookie) {
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
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"ログアウトしました"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	user, err := h.UserRepo.GetUser(userID)
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"ユーザーが存在しません。"}))
			return
		} else if err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
