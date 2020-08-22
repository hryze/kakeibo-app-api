package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type TestUserRepository struct{}

func (t TestUserRepository) FindUserID(userID string) error {
	return sql.ErrNoRows
}

func (t TestUserRepository) FindEmail(email string) error {
	return sql.ErrNoRows
}

func (t TestUserRepository) CreateUser(user *model.SignUpUser) error {
	return nil
}

func (t TestUserRepository) DeleteUser(signUpUser *model.SignUpUser) error {
	return nil
}

func (t TestUserRepository) FindUser(loginUser *model.LoginUser) (*model.LoginUser, error) {
	return &model.LoginUser{
		ID:       "testID",
		Name:     "testName",
		Email:    "test@icloud.com",
		Password: "$2a$10$teJL.9I0QfBESpaBIwlbl.VkivuHEOKhy674CW6J.4k3AnfEpcYLy",
	}, nil
}

func (t TestUserRepository) SetSessionID(sessionID string, loginUserID string, expiration int) error {
	return nil
}

func (t TestUserRepository) DeleteSessionID(sessionID string) error {
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

	h := DBHandler{UserRepo: TestUserRepository{}}

	requestJson := model.SignUpUser{
		ID:       "testID",
		Name:     "testName",
		Email:    "test@icloud.com",
		Password: "testPassword",
	}

	b, err := json.Marshal(&requestJson)
	if err != nil {
		t.Fatalf("error: %#v", err)
	}

	r := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.SignUp(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("want = 201, got = %d", res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Errorf("want = application/json; charset=UTF-8, got = %s", res.Header.Get("Content-Type"))
	}

	var responseJson model.SignUpUser
	if err := json.NewDecoder(res.Body).Decode(&responseJson); err != nil {
		t.Errorf("error: %#v, res: %#v", err, w)
	}

	if responseJson.ID != "testID" {
		t.Errorf("want = testID, got = %s", responseJson.ID)
	}

	if responseJson.Name != "testName" {
		t.Errorf("want = testName, got = %s", responseJson.Name)
	}

	if responseJson.Email != "test@icloud.com" {
		t.Errorf("want = test@icloud.com, got = %s", responseJson.Email)
	}
}

func TestDBHandler_Login(t *testing.T) {
	h := DBHandler{UserRepo: TestUserRepository{}}

	requestJson := model.LoginUser{
		Email:    "test@icloud.com",
		Password: "testPassword",
	}

	b, err := json.Marshal(&requestJson)
	if err != nil {
		t.Fatalf("error: %#v", err)
	}

	r := httptest.NewRequest("POST", "/login", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.Login(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("want = 201, got = %d", res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Errorf("want = application/json; charset=UTF-8, got = %s", res.Header.Get("Content-Type"))
	}

	cookies := res.Cookies()[0]

	if cookies.Name != "session_id" {
		t.Errorf("want = session_id, got = %s", cookies.Name)
	}

	if len(cookies.Value) < 0 {
		t.Errorf("want = , got = %d", len(cookies.Value))
	}

	if !time.Now().Before(cookies.Expires) {
		t.Errorf("want = true, got = %t", time.Now().Before(cookies.Expires))
	}

	if !cookies.HttpOnly {
		t.Errorf("want = true, got = %t", cookies.HttpOnly)
	}

	var responseJson model.LoginUser
	if err := json.NewDecoder(res.Body).Decode(&responseJson); err != nil {
		t.Errorf("error: %#v, res: %#v", err, w)
	}

	if responseJson.ID != "testID" {
		t.Errorf("want = testID, got = %s", responseJson.ID)
	}

	if responseJson.Name != "testName" {
		t.Errorf("want = testName, got = %s", responseJson.Name)
	}

	if responseJson.Email != "test@icloud.com" {
		t.Errorf("want = test@icloud.com, got = %s", responseJson.Email)
	}
}
