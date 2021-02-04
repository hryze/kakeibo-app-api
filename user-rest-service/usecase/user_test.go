package usecase

import (
	"testing"

	merrors "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model/errors"

	"github.com/google/go-cmp/cmp"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type mockUserRepository struct{}

func (t *mockUserRepository) FindSignUpUserByUserID(userID string) (*model.SignUpUser, error) {
	return nil, &merrors.UserNotFoundError{
		Message: "ユーザーが見つかりませんでした",
	}
}

func (t *mockUserRepository) FindSignUpUserByEmail(email string) (*model.SignUpUser, error) {
	return nil, &merrors.UserNotFoundError{
		Message: "ユーザーが見つかりませんでした",
	}
}

func (t *mockUserRepository) CreateSignUpUser(user *model.SignUpUser) error {
	return nil
}

func (t *mockUserRepository) DeleteSignUpUser(signUpUser *model.SignUpUser) error {
	return nil
}

func (t *mockUserRepository) FindUser(loginUser *model.LoginUser) (*model.LoginUser, error) {
	return &model.LoginUser{
		ID:       "testID",
		Name:     "testName",
		Email:    "test@icloud.com",
		Password: "$2a$10$teJL.9I0QfBESpaBIwlbl.VkivuHEOKhy674CW6J.4k3AnfEpcYLy",
	}, nil
}

func (t *mockUserRepository) GetUser(userID string) (*model.LoginUser, error) {
	return &model.LoginUser{
		ID:    "testID",
		Name:  "testName",
		Email: "test@icloud.com",
	}, nil
}

func (t *mockUserRepository) SetSessionID(sessionID string, loginUserID string, expiration int) error {
	return nil
}

func (t *mockUserRepository) DeleteSessionID(sessionID string) error {
	return nil
}

type mockAccountApi struct{}

func (t *mockAccountApi) PostInitStandardBudgets(userID string) error {
	return nil
}

func Test_userUsecase_SignUp(t *testing.T) {
	u := NewUserUsecase(&mockUserRepository{}, &mockAccountApi{})

	in := input.SignUpUser{
		UserID:   "testUserID",
		Name:     "testName",
		Email:    "test@icloud.com",
		Password: "testPassword",
	}

	gotOut, err := u.SignUp(&in)
	if err != nil {
		t.Errorf("unexpected error by userUsecase.SignUp '%#v'", err)
	}

	wantOut := &output.SignUpUser{
		UserID: "testUserID",
		Name:   "testName",
		Email:  "test@icloud.com",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
