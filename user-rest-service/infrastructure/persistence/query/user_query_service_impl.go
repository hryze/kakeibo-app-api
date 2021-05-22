package query

import (
	"database/sql"

	"golang.org/x/xerrors"

	"github.com/hryze/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/hryze/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/hryze/kakeibo-app-api/user-rest-service/infrastructure/persistence/datasource"
	"github.com/hryze/kakeibo-app-api/user-rest-service/infrastructure/persistence/rdb"
)

type userQueryServiceImpl struct {
	*rdb.MySQLHandler
}

func NewUserQueryService(mysqlHandler *rdb.MySQLHandler) *userQueryServiceImpl {
	return &userQueryServiceImpl{mysqlHandler}
}

func (q *userQueryServiceImpl) FindLoginUserByUserID(userID userdomain.UserID) (*userdomain.LoginUserWithoutPassword, error) {
	query := `
        SELECT
            user_id,
            name,
            email
        FROM 
            users
        WHERE
            user_id = ?`

	var loginUserDto datasource.LoginUser
	if err := q.MySQLHandler.Conn.QueryRowx(query, userID.Value()).StructScan(&loginUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	loginUser := userdomain.NewLoginUserWithoutPassword(loginUserDto.UserID, loginUserDto.Email, loginUserDto.Name)

	return loginUser, nil
}
