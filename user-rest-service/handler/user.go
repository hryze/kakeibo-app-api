package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/errors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"golang.org/x/xerrors"

	"github.com/garyburd/redigo/redis"
)

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
		errors.ErrorResponseByJSON(w, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error")))
		return
	}

	out, err := h.userUsecase.SignUp(&in)
	if err != nil {
		errors.ErrorResponseByJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in input.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		errors.ErrorResponseByJSON(w, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error")))
		return
	}

	out, err := h.userUsecase.Login(&in)
	if err != nil {
		errors.ErrorResponseByJSON(w, err)
		return
	}

	domain := os.Getenv("COOKIE_DOMAIN")
	secure := true

	if domain != "shakepiper.com" {
		secure = false
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    out.SessionID,
		Expires:  out.Expires,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(out); err != nil {
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

	out := &output.LoginUser{
		UserID: user.UserID(),
		Name:   user.Name(),
		Email:  user.Email(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
