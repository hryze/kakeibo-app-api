package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/interfaces/presenter"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
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
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	out, err := h.userUsecase.SignUp(&in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in input.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	out, err := h.userUsecase.Login(&in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.Env.Cookie.Name,
		Value:    out.Cookie.SessionID,
		Expires:  time.Now().Add(config.Env.Cookie.Expiration),
		Domain:   config.Env.Cookie.Domain,
		Secure:   config.Env.Cookie.Secure,
		HttpOnly: true,
	})

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *userHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.Env.Cookie.Name)
	if xerrors.Is(err, http.ErrNoCookie) {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("ログアウト済みです")))
		return
	}

	cookieInfo := &input.CookieInfo{SessionID: cookie.Value}

	if err := h.userUsecase.Logout(cookieInfo); err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	presenter.JSON(w, http.StatusOK, presenter.NewSuccessString("ログアウトしました"))
}

func (h *userHandler) FetchLoginUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := getUserIDForContext(r)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	out, err := h.userUsecase.FetchLoginUser(authenticatedUser)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}
