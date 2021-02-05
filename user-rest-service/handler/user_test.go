package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	merrors "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model/errors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/testutil"

	"github.com/google/uuid"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type MockUserRepository struct{}

func (t MockUserRepository) FindSignUpUserByUserID(userID string) (*model.SignUpUser, error) {
	return nil, merrors.UserNotFoundError
}

func (t MockUserRepository) FindSignUpUserByEmail(email string) (*model.SignUpUser, error) {
	return nil, merrors.UserNotFoundError
}

func (t MockUserRepository) CreateSignUpUser(user *model.SignUpUser) error {
	return nil
}

func (t MockUserRepository) DeleteSignUpUser(signUpUser *model.SignUpUser) error {
	return nil
}

func (t MockUserRepository) FindLoginUserByEmail(email string) (*model.LoginUser, error) {
	loginUser := model.NewLoginUserFromDataSource("testID", "testName", "test@icloud.com", "$2a$10$teJL.9I0QfBESpaBIwlbl.VkivuHEOKhy674CW6J.4k3AnfEpcYLy")

	return loginUser, nil
}

func (t MockUserRepository) GetUser(userID string) (*model.LoginUser, error) {
	loginUser := model.NewLoginUserFromDataSource("testID", "testName", "test@icloud.com", "$2a$10$teJL.9I0QfBESpaBIwlbl.VkivuHEOKhy674CW6J.4k3AnfEpcYLy")

	return loginUser, nil
}

func (t MockUserRepository) AddSessionID(sessionID string, userID string, expiration int) error {
	return nil
}

func (t MockUserRepository) DeleteSessionID(sessionID string) error {
	return nil
}

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
		UserID:    "testID",
		Name:      "testName",
		Email:     "test@icloud.com",
		SessionID: uuid.New().String(),
		Expires:   time.Now().Add(time.Duration(86400*30) * time.Second),
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

func TestDBHandler_Logout(t *testing.T) {
	h := DBHandler{UserRepo: MockUserRepository{}}

	r := httptest.NewRequest("DELETE", "/logout", nil)
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.Logout(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
	testutil.AssertDeleteResponseCookie(t, res)
}

func TestDBHandler_GetUser(t *testing.T) {
	h := DBHandler{
		AuthRepo: MockAuthRepository{},
		UserRepo: MockUserRepository{},
	}

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &output.LoginUser{}, &output.LoginUser{})
}
