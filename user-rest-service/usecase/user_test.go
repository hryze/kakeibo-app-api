package usecase

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/groupdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/vo"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
)

type mockUserRepository struct{}

func (r *mockUserRepository) FindSignUpUserByUserID(userID userdomain.UserID) (*userdomain.SignUpUser, error) {
	return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
}

func (r *mockUserRepository) FindSignUpUserByEmail(email vo.Email) (*userdomain.SignUpUser, error) {
	return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
}

func (r *mockUserRepository) CreateSignUpUser(user *userdomain.SignUpUser) error {
	return nil
}

func (r *mockUserRepository) DeleteSignUpUser(signUpUser *userdomain.SignUpUser) error {
	return nil
}

func (r *mockUserRepository) FindLoginUserByEmail(email vo.Email) (*userdomain.LoginUser, error) {
	loginUser := userdomain.NewLoginUserWithHashPassword("testUserID", "testName", "test@icloud.com", "$2a$10$MfTmnqbuDog.W/Kaug3vlef0ZX5OoxEbjSc9hyAB.YGNKQvfQGDd6")

	return loginUser, nil
}

type mockUserQueryService struct{}

func (q *mockUserQueryService) FindLoginUserByUserID(userID userdomain.UserID) (*userdomain.LoginUserWithoutPassword, error) {
	loginUser := userdomain.NewLoginUserWithoutPassword("testUserID", "testName", "test@icloud.com")

	return loginUser, nil
}

type mockSessionStore struct{}

func (s *mockSessionStore) StoreUserBySessionID(sessionID string, loginUserID userdomain.UserID) error {
	return nil
}

func (s *mockSessionStore) DeleteUserBySessionID(sessionID string) error {
	return nil
}

func (s *mockSessionStore) FetchUserByUserID(sessionID string) (userdomain.UserID, error) {
	userID, _ := userdomain.NewUserID("testID")

	return userID, nil
}

type mockAccountApi struct{}

func (a *mockAccountApi) PostInitStandardBudgets(userID userdomain.UserID) error {
	return nil
}

func (a *mockAccountApi) PostInitGroupStandardBudgets(groupID groupdomain.GroupID) error {
	return nil
}

func Test_userUsecase_SignUp(t *testing.T) {
	u := NewUserUsecase(&mockUserRepository{}, &mockUserQueryService{}, &mockSessionStore{}, &mockAccountApi{})

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

func Test_userUsecase_Login(t *testing.T) {
	u := NewUserUsecase(&mockUserRepository{}, &mockUserQueryService{}, &mockSessionStore{}, &mockAccountApi{})

	in := input.LoginUser{
		Email:    "test@icloud.com",
		Password: "testPassword",
	}

	gotOut, err := u.Login(&in)
	if err != nil {
		t.Errorf("unexpected error by userUsecase.Login '%#v'", err)
	}

	wantOut := &output.LoginUser{
		UserID: "testUserID",
		Name:   "testName",
		Email:  "test@icloud.com",
	}

	ignoreFieldsOption := cmpopts.IgnoreFields(output.LoginUser{}, "Cookie")

	if diff := cmp.Diff(&wantOut, &gotOut, ignoreFieldsOption); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func Test_userUsecase_Logout(t *testing.T) {
	u := NewUserUsecase(&mockUserRepository{}, &mockUserQueryService{}, &mockSessionStore{}, &mockAccountApi{})

	in := input.CookieInfo{
		SessionID: uuid.New().String(),
	}

	if err := u.Logout(&in); err != nil {
		t.Errorf("unexpected error by userUsecase.Logout '%#v'", err)
	}
}

func Test_userUsecase_FetchLoginUser(t *testing.T) {
	u := NewUserUsecase(&mockUserRepository{}, &mockUserQueryService{}, &mockSessionStore{}, &mockAccountApi{})

	in := input.AuthenticatedUser{
		UserID: "testUserID",
	}

	gotOut, err := u.FetchLoginUser(&in)
	if err != nil {
		t.Errorf("unexpected error by userUsecase.FetchUserInfo '%#v'", err)
	}

	wantOut := &output.LoginUser{
		UserID: "testUserID",
		Name:   "testName",
		Email:  "test@icloud.com",
	}

	if diff := cmp.Diff(&wantOut, &gotOut); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
