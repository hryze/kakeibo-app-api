package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

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
	return &model.LoginUser{}, nil
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

	if w.Code != http.StatusCreated {
		t.Errorf("want = 201, got = %d", w.Code)
	}

	var responseJson model.SignUpUser
	if err := json.NewDecoder(w.Body).Decode(&responseJson); err != nil {
		t.Errorf("error: %#v, res: %#v", err, w)
	}

	if responseJson.ID != requestJson.ID {
		t.Errorf("want = %s, got = %s", requestJson.ID, responseJson.ID)
	}

	if responseJson.Name != requestJson.Name {
		t.Errorf("want = %s, got = %s", requestJson.Name, responseJson.Name)
	}

	if responseJson.Email != requestJson.Email {
		t.Errorf("want = %s, got = %s", requestJson.Email, responseJson.Email)
	}
}
