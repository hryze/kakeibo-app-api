package usecase

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	merrors "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model/errors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
	uerrors "github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/errors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/interfaces"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

type UserUsecase interface {
	SignUp(inSignUpUser *input.SignUpUser) (*output.SignUpUser, error)
}

type userUsecase struct {
	userRepository repository.UserRepository
	accountApi     interfaces.AccountApi
}

func NewUserUsecase(userRepository repository.UserRepository, accountApi interfaces.AccountApi) *userUsecase {
	return &userUsecase{
		userRepository: userRepository,
		accountApi:     accountApi,
	}
}

func (u *userUsecase) SignUp(in *input.SignUpUser) (*output.SignUpUser, error) {
	signUpUser, err := model.NewSignUpUser(in.UserID, in.Name, in.Email, in.Password)
	if err != nil {
		return nil, uerrors.NewBadRequestError(err)
	}

	if err := checkForUniqueUser(u, signUpUser); err != nil {
		var userConflictError *merrors.UserConflictError
		if xerrors.As(err, &userConflictError) {
			return nil, uerrors.NewConflictError(userConflictError)
		}

		return nil, uerrors.NewInternalServerError(uerrors.NewErrorString("ユーザー登録に失敗しました"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(signUpUser.Password()), 10)
	if err != nil {
		return nil, uerrors.NewInternalServerError(uerrors.NewErrorString("ユーザー登録に失敗しました"))
	}

	signUpUser.SetPassword(string(hash))

	if err := u.userRepository.CreateSignUpUser(signUpUser); err != nil {
		return nil, uerrors.NewInternalServerError(uerrors.NewErrorString("ユーザー登録に失敗しました"))
	}

	if err := u.accountApi.PostInitStandardBudgets(signUpUser.UserID()); err != nil {
		if err := u.userRepository.DeleteSignUpUser(signUpUser); err != nil {
			return nil, uerrors.NewInternalServerError(uerrors.NewErrorString("ユーザー登録に失敗しました"))
		}

		return nil, uerrors.NewInternalServerError(uerrors.NewErrorString("ユーザー登録に失敗しました"))
	}

	return &output.SignUpUser{
		UserID: signUpUser.UserID(),
		Name:   signUpUser.Name(),
		Email:  signUpUser.Email(),
	}, nil
}

func checkForUniqueUser(u *userUsecase, signUpUser *model.SignUpUser) error {
	var userNotFoundError *merrors.UserNotFoundError

	_, errUserID := u.userRepository.FindSignUpUserByUserID(signUpUser.UserID())
	existsUserByUserID := !xerrors.As(errUserID, &userNotFoundError)
	if errUserID != nil && existsUserByUserID {
		return errUserID
	}

	_, errEmail := u.userRepository.FindSignUpUserByEmail(signUpUser.Email())
	existsUserByEmail := !xerrors.As(errEmail, &userNotFoundError)
	if errEmail != nil && existsUserByEmail {
		return errEmail
	}

	if !existsUserByUserID && !existsUserByEmail {
		return nil
	}

	var userConflictError merrors.UserConflictError

	if existsUserByUserID {
		userConflictError.UserID = "このユーザーIDは既に利用されています"
	}

	if existsUserByEmail {
		userConflictError.Email = "このメールアドレスは既に利用されています"
	}

	return &userConflictError
}
