package handler

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/testutil"

	"github.com/google/uuid"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type MockUserRepository struct{}

func (t MockUserRepository) FindUserID(userID string) error {
	return sql.ErrNoRows
}

func (t MockUserRepository) FindEmail(email string) error {
	return sql.ErrNoRows
}

func (t MockUserRepository) CreateUser(user *model.SignUpUser) error {
	return nil
}

func (t MockUserRepository) DeleteUser(signUpUser *model.SignUpUser) error {
	return nil
}

func (t MockUserRepository) FindUser(loginUser *model.LoginUser) (*model.LoginUser, error) {
	return &model.LoginUser{
		ID:       "testID",
		Name:     "testName",
		Email:    "test@icloud.com",
		Password: "$2a$10$teJL.9I0QfBESpaBIwlbl.VkivuHEOKhy674CW6J.4k3AnfEpcYLy",
	}, nil
}

func (t MockUserRepository) SetSessionID(sessionID string, loginUserID string, expiration int) error {
	return nil
}

func (t MockUserRepository) DeleteSessionID(sessionID string) error {
	return nil
}

func TestDBHandler_SignUp(t *testing.T) {
	if err := os.Setenv("ENVIRONMENT", "development"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	postInitStandardBudgetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: postInitStandardBudgetsHandler},
	}
	ts.Start()
	defer ts.Close()

	h := DBHandler{UserRepo: MockUserRepository{}}

	r := httptest.NewRequest("POST", "/TestDBHandler_SignUp", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	h.SignUp(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.SignUpUser{}, &model.SignUpUser{})
}

func TestDBHandler_Login(t *testing.T) {
	h := DBHandler{UserRepo: MockUserRepository{}}

	r := httptest.NewRequest("POST", "/TestDBHandler_Login", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	h.Login(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.LoginUser{}, &model.LoginUser{})
	testutil.AssertSetResponseCookie(t, res)
}

func TestDBHandler_Logout(t *testing.T) {
	h := DBHandler{UserRepo: MockUserRepository{}}

	r := httptest.NewRequest("DELETE", "/TestDBHandler_Logout", nil)
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
	testutil.AssertResponseBody(t, res, &LogoutMsg{}, &LogoutMsg{})
	testutil.AssertDeleteResponseCookie(t, res)
}
