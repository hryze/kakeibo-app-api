package handler

import (
	"database/sql"
	"net"
	"net/http"
	"net/http/httptest"
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
	postInitStandardBudgetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		t.Fatalf("error: %#v", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: postInitStandardBudgetsHandler},
	}
	ts.Start()
	defer ts.Close()

	h := DBHandler{UserRepo: MockUserRepository{}}

	r := httptest.NewRequest("POST", "/signup", strings.NewReader(testutil.GetJsonFromTestData(t, "./testdata/user/signup/request.json")))
	w := httptest.NewRecorder()

	h.SignUp(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, "./testdata/user/signup/response.json.golden")
}

func TestDBHandler_Login(t *testing.T) {
	h := DBHandler{UserRepo: MockUserRepository{}}

	r := httptest.NewRequest("POST", "/login", strings.NewReader(testutil.GetJsonFromTestData(t, "./testdata/user/login/request.json")))
	w := httptest.NewRecorder()

	h.Login(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, "./testdata/user/login/response.json.golden")
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
	testutil.AssertResponseBody(t, res, "./testdata/user/logout/response.json.golden")
	testutil.AssertDeleteResponseCookie(t, res)
}
