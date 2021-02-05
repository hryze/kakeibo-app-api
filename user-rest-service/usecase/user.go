package usecase

import (
	"time"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/errors"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/interfaces"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

type UserUsecase interface {
	SignUp(in *input.SignUpUser) (*output.SignUpUser, error)
	Login(in *input.LoginUser) (*output.LoginUser, error)
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
		return nil, errors.NewBadRequestError(err)
	}

	if err := checkForUniqueUser(u, signUpUser); err != nil {
		var userConflictError *errors.UserConflictError
		if xerrors.As(err, &userConflictError) {
			return nil, errors.NewConflictError(userConflictError)
		}

		return nil, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(signUpUser.Password()), 10)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error"))
	}

	signUpUser.SetPassword(string(hash))

	if err := u.userRepository.CreateSignUpUser(signUpUser); err != nil {
		return nil, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error"))
	}

	if err := u.accountApi.PostInitStandardBudgets(signUpUser.UserID()); err != nil {
		if err := u.userRepository.DeleteSignUpUser(signUpUser); err != nil {
			return nil, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error"))
		}

		return nil, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error"))
	}

	return &output.SignUpUser{
		UserID: signUpUser.UserID(),
		Name:   signUpUser.Name(),
		Email:  signUpUser.Email(),
	}, nil
}

func (u *userUsecase) Login(in *input.LoginUser) (*output.LoginUser, error) {
	loginUser, err := model.NewLoginUser(in.Email, in.Password)
	if err != nil {
		return nil, errors.NewBadRequestError(err)
	}

	dbLoginUser, err := u.userRepository.FindLoginUserByEmail(loginUser.Email())
	if err != nil {
		if xerrors.Is(err, errors.ErrUserNotFound) {
			return nil, errors.NewAuthenticationError(errors.NewErrorString("認証に失敗しました"))
		}

		return nil, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error"))
	}

	hashedPassword := dbLoginUser.Password()
	password := loginUser.Password()

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, errors.NewAuthenticationError(errors.NewErrorString("認証に失敗しました"))
	}

	sessionID := uuid.New().String()
	expiration := 86400 * 30

	if err := u.userRepository.AddSessionID(sessionID, dbLoginUser.UserID(), expiration); err != nil {
		return nil, errors.NewInternalServerError(errors.NewErrorString("Internal Server Error"))
	}

	return &output.LoginUser{
		UserID:    dbLoginUser.UserID(),
		Name:      dbLoginUser.Name(),
		Email:     dbLoginUser.Email(),
		SessionID: sessionID,
		Expires:   time.Now().Add(time.Duration(expiration) * time.Second),
	}, nil
}

func checkForUniqueUser(u *userUsecase, signUpUser *model.SignUpUser) error {
	_, errUserID := u.userRepository.FindSignUpUserByUserID(signUpUser.UserID())
	if errUserID != nil && !xerrors.Is(errUserID, errors.ErrUserNotFound) {
		return errUserID
	}

	_, errEmail := u.userRepository.FindSignUpUserByEmail(signUpUser.Email())
	if errEmail != nil && !xerrors.Is(errEmail, errors.ErrUserNotFound) {
		return errEmail
	}

	existsUserByUserID := !xerrors.Is(errUserID, errors.ErrUserNotFound)
	existsUserByEmail := !xerrors.Is(errEmail, errors.ErrUserNotFound)

	if !existsUserByUserID && !existsUserByEmail {
		return nil
	}

	var userConflictError errors.UserConflictError

	if existsUserByUserID {
		userConflictError.UserID = "このユーザーIDは既に利用されています"
	}

	if existsUserByEmail {
		userConflictError.Email = "このメールアドレスは既に利用されています"
	}

	return &userConflictError
}
