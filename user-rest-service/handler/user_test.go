package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	merrors "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model/errors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/testutil"

	"github.com/google/uuid"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type MockUserRepository struct{}

func (t MockUserRepository) FindSignUpUserByUserID(userID string) (*model.SignUpUser, error) {
	return nil, &merrors.UserNotFoundError{
		Message: "ユーザーが見つかりませんでした",
	}
}

func (t MockUserRepository) FindSignUpUserByEmail(email string) (*model.SignUpUser, error) {
	return nil, &merrors.UserNotFoundError{
		Message: "ユーザーが見つかりませんでした",
	}
}

func (t MockUserRepository) CreateSignUpUser(user *model.SignUpUser) error {
	return nil
}

func (t MockUserRepository) DeleteSignUpUser(signUpUser *model.SignUpUser) error {
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

func (t MockUserRepository) GetUser(userID string) (*model.LoginUser, error) {
	return &model.LoginUser{
		ID:    "testID",
		Name:  "testName",
		Email: "test@icloud.com",
	}, nil
}

func (t MockUserRepository) SetSessionID(sessionID string, loginUserID string, expiration int) error {
	return nil
}

func (t MockUserRepository) DeleteSessionID(sessionID string) error {
	return nil
}

type mockUserUsecase struct{}

func (u *mockUserUsecase) SignUp(inSignUpUser *input.SignUpUser) (*output.SignUpUser, error) {
	return &output.SignUpUser{
		UserID: "testID",
		Name:   "testName",
		Email:  "test@icloud.com",
	}, nil
}

func Test_userHandler_SignUp(t *testing.T) {
	//if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
	//	t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	//}
	//
	//postInitStandardBudgetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusCreated)
	//})
	//
	//listener, err := net.Listen("tcp", "127.0.0.1:8081")
	//if err != nil {
	//	t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	//}
	//
	//ts := httptest.Server{
	//	Listener: listener,
	//	Config:   &http.Server{Handler: postInitStandardBudgetsHandler},
	//}
	//ts.Start()
	//defer ts.Close()

	h := NewUserHandler(&mockUserUsecase{})

	r := httptest.NewRequest("POST", "/signup", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	h.SignUp(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &output.SignUpUser{}, &output.SignUpUser{})
}

func TestDBHandler_Login(t *testing.T) {
	h := DBHandler{UserRepo: MockUserRepository{}}

	r := httptest.NewRequest("POST", "/login", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
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
	testutil.AssertResponseBody(t, res, &model.LoginUser{}, &model.LoginUser{})
}
