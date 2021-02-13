package persistence

import (
	"database/sql"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/vo"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/datasource"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/db"
)

type userRepository struct {
	*db.RedisHandler
	*db.MySQLHandler
}

func NewUserRepository(redisHandler *db.RedisHandler, mysqlHandler *db.MySQLHandler) *userRepository {
	return &userRepository{
		redisHandler,
		mysqlHandler,
	}
}

func (r *userRepository) FindSignUpUserByUserID(userID userdomain.UserID) (*userdomain.SignUpUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email
        FROM
            users
        WHERE
            user_id = ?`

	var signUpUserDto datasource.SignUpUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, userID.Value()).StructScan(&signUpUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	var userValidationError userdomain.ValidationError

	userIDVo, err := userdomain.NewUserID(signUpUserDto.UserID)
	if err != nil {
		userValidationError.UserID = "ユーザーIDが正しくありません"
	}

	nameVo, err := userdomain.NewName(signUpUserDto.Name)
	if err != nil {
		userValidationError.Name = "名前が正しくありません"
	}

	emailVo, err := vo.NewEmail(signUpUserDto.Email)
	if err != nil {
		userValidationError.Email = "メールアドレスが正しくありません"
	}

	if !userValidationError.IsEmpty() {
		return nil, apierrors.NewBadRequestError(&userValidationError)
	}

	signUpUser := userdomain.NewSignUpUserFromDataSource(userIDVo, nameVo, emailVo)

	return signUpUser, nil
}

func (r *userRepository) FindSignUpUserByEmail(email vo.Email) (*userdomain.SignUpUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email
        FROM
            users
        WHERE
            email = ?`

	var signUpUserDto datasource.SignUpUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, email.Value()).StructScan(&signUpUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	var userValidationError userdomain.ValidationError

	userIDVo, err := userdomain.NewUserID(signUpUserDto.UserID)
	if err != nil {
		userValidationError.UserID = "ユーザーIDが正しくありません"
	}

	nameVo, err := userdomain.NewName(signUpUserDto.Name)
	if err != nil {
		userValidationError.Name = "名前が正しくありません"
	}

	emailVo, err := vo.NewEmail(signUpUserDto.Email)
	if err != nil {
		userValidationError.Email = "メールアドレスが正しくありません"
	}

	if !userValidationError.IsEmpty() {
		return nil, apierrors.NewBadRequestError(&userValidationError)
	}

	signUpUser := userdomain.NewSignUpUserFromDataSource(userIDVo, nameVo, emailVo)

	return signUpUser, nil
}

func (r *userRepository) CreateSignUpUser(signUpUser *userdomain.SignUpUser) error {
	query := `
        INSERT INTO users
            (user_id, name, email, password)
        VALUES
            (?,?,?,?)`

	if _, err := r.MySQLHandler.Conn.Exec(
		query,
		signUpUser.UserID().Value(),
		signUpUser.Name().Value(),
		signUpUser.Email().Value(),
		signUpUser.Password().Value(),
	); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}

func (r *userRepository) DeleteSignUpUser(signUpUser *userdomain.SignUpUser) error {
	query := `
        DELETE
        FROM
            users
        WHERE
            user_id = ?`

	if _, err := r.MySQLHandler.Conn.Exec(query, signUpUser.UserID().Value()); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}

func (r *userRepository) FindLoginUserByEmail(email vo.Email) (*userdomain.LoginUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email,
            password
        FROM 
            users
        WHERE
            email = ?`

	var loginUserDto datasource.LoginUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, email.Value()).StructScan(&loginUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	var userValidationError userdomain.ValidationError

	userIDVo, err := userdomain.NewUserID(loginUserDto.UserID)
	if err != nil {
		userValidationError.UserID = "ユーザーIDが正しくありません"
	}

	nameVo, err := userdomain.NewName(loginUserDto.Name)
	if err != nil {
		userValidationError.Name = "名前が正しくありません"
	}

	emailVo, err := vo.NewEmail(loginUserDto.Email)
	if err != nil {
		userValidationError.Email = "メールアドレスが正しくありません"
	}

	passwordVo, err := vo.NewHashPassword(loginUserDto.Password)
	if err != nil {
		userValidationError.Password = "パスワードが正しくありません"
	}

	if !userValidationError.IsEmpty() {
		return nil, apierrors.NewBadRequestError(&userValidationError)
	}

	loginUser := userdomain.NewLoginUserFromDataSource(userIDVo, nameVo, emailVo, passwordVo)

	return loginUser, nil
}

func (r *userRepository) GetUser(userID string) (*userdomain.LoginUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email
        FROM 
            users
        WHERE
            user_id = ?`

	var user userdomain.LoginUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, userID).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) AddSessionID(sessionID string, userID userdomain.UserID, expiration int) error {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", sessionID, userID.Value(), "EX", expiration); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}

func (r *userRepository) DeleteSessionID(sessionID string) error {
	conn := r.RedisHandler.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", sessionID); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}
