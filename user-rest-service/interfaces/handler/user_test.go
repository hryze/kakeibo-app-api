package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/context"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/interfaces/presenter"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/testutil"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
)

type mockUserUsecase struct{}

func (u *mockUserUsecase) SignUp(in *input.SignUpUser) (*output.SignUpUser, error) {
	return &output.SignUpUser{
		UserID: "testID",
		Name:   "testName",
		Email:  "test@icloud.com",
	}, nil
}

func (u *mockUserUsecase) Login(in *input.LoginUser) (*output.LoginUser, error) {
	return &output.LoginUser{
		UserID: "testID",
		Name:   "testName",
		Email:  "test@icloud.com",
		Cookie: output.CookieInfo{
			SessionID: uuid.New().String(),
		},
	}, nil
}

func (u *mockUserUsecase) Logout(in *input.CookieInfo) error {
	return nil
}

func (u *mockUserUsecase) FetchLoginUser(in *input.AuthenticatedUser) (*output.LoginUser, error) {
	return &output.LoginUser{
		UserID: "testID",
		Name:   "testName",
		Email:  "test@icloud.com",
	}, nil
}

func Test_userHandler_SignUp(t *testing.T) {
	h := NewUserHandler(&mockUserUsecase{})

	r := httptest.NewRequest("POST", "/signup", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	h.SignUp(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &output.SignUpUser{}, &output.SignUpUser{})
}

func Test_userHandler_Login(t *testing.T) {
	h := NewUserHandler(&mockUserUsecase{})

	r := httptest.NewRequest("POST", "/login", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	h.Login(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &output.LoginUser{}, &output.LoginUser{})
	testutil.AssertSetResponseCookie(t, res)
}

func Test_userHandler_Logout(t *testing.T) {
	h := NewUserHandler(&mockUserUsecase{})

	r := httptest.NewRequest("DELETE", "/logout", nil)
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  config.Env.Cookie.Name,
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.Logout(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, presenter.NewSuccessString(""), presenter.NewSuccessString(""))
	testutil.AssertDeleteResponseCookie(t, res)
}

func Test_userHandler_FetchLoginUser(t *testing.T) {
	h := NewUserHandler(&mockUserUsecase{})

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()

	context.Set(r, config.Env.RequestCtx.UserID, "testID")

	h.FetchLoginUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &output.LoginUser{}, &output.LoginUser{})
}
