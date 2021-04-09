package query

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/rdb"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/output"
)

type groupQueryServiceImpl struct {
	*rdb.MySQLHandler
}

func NewGroupQueryServiceImpl(mysqlHandler *rdb.MySQLHandler) *groupQueryServiceImpl {
	return &groupQueryServiceImpl{mysqlHandler}
}

func (r *groupQueryServiceImpl) FetchGroupList(userID string) (*output.GroupList, error) {
	approvedGroupList, err := r.fetchApprovedGroupList(userID)
	if err != nil {
		return nil, err
	}

	unapprovedGroupList, err := r.fetchUnApprovedGroupList(userID)
	if err != nil {
		return nil, err
	}

	if len(approvedGroupList) == 0 && len(unapprovedGroupList) == 0 {
		return &output.GroupList{
			ApprovedGroupList:   approvedGroupList,
			UnapprovedGroupList: unapprovedGroupList,
		}, nil
	}

	groupIDList := generateGroupIDList(approvedGroupList, unapprovedGroupList)

	approvedUsersList, err := r.fetchApprovedUsersList(groupIDList)
	if err != nil {
		return nil, err
	}

	unapprovedUsersList, err := r.fetchUnapprovedUsersList(groupIDList)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(approvedGroupList); i++ {
		for _, approvedUser := range approvedUsersList {
			if approvedGroupList[i].GroupID == approvedUser.GroupID {
				approvedGroupList[i].ApprovedUsersList = append(approvedGroupList[i].ApprovedUsersList, approvedUser)
			}
		}

		for _, unapprovedUser := range unapprovedUsersList {
			if approvedGroupList[i].GroupID == unapprovedUser.GroupID {
				approvedGroupList[i].UnapprovedUsersList = append(approvedGroupList[i].UnapprovedUsersList, unapprovedUser)
			}
		}
	}

	for i := 0; i < len(unapprovedGroupList); i++ {
		for _, approvedUser := range approvedUsersList {
			if unapprovedGroupList[i].GroupID == approvedUser.GroupID {
				unapprovedGroupList[i].ApprovedUsersList = append(unapprovedGroupList[i].ApprovedUsersList, approvedUser)
			}
		}

		for _, unapprovedUser := range unapprovedUsersList {
			if unapprovedGroupList[i].GroupID == unapprovedUser.GroupID {
				unapprovedGroupList[i].UnapprovedUsersList = append(unapprovedGroupList[i].UnapprovedUsersList, unapprovedUser)
			}
		}
	}

	return &output.GroupList{
		ApprovedGroupList:   approvedGroupList,
		UnapprovedGroupList: unapprovedGroupList,
	}, nil
}

func (r *groupQueryServiceImpl) FetchUnapprovedUser(groupID int, userID string) (*output.UnapprovedUser, error) {
	query := `
        SELECT
            group_unapproved_users.group_id group_id,
            group_unapproved_users.user_id user_id,
            users.name user_name
        FROM
            group_unapproved_users
        INNER JOIN
            users
        ON
            group_unapproved_users.user_id = users.user_id
        WHERE
            group_unapproved_users.group_id = ?
        AND
            group_unapproved_users.user_id = ?`

	var unapprovedUser output.UnapprovedUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, groupID, userID).StructScan(&unapprovedUser); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return &unapprovedUser, nil
}

func (r *groupQueryServiceImpl) FetchApprovedUser(groupID int, userID string) (*output.ApprovedUser, error) {
	query := `
        SELECT
            group_users.group_id group_id,
            group_users.user_id user_id,
            users.name user_name,
            group_users.color_code color_code
        FROM
            group_users
        INNER JOIN
            users
        ON
            group_users.user_id = users.user_id
        WHERE
            group_users.group_id = ?
        AND
            group_users.user_id = ?`

	var approvedUser output.ApprovedUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, groupID, userID).StructScan(&approvedUser); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return &approvedUser, nil
}

func generateGroupIDList(approvedGroupList []output.ApprovedGroup, unapprovedGroupList []output.UnapprovedGroup) []interface{} {
	groupIDList := make([]interface{}, 0, len(approvedGroupList)+len(unapprovedGroupList))

	for _, approvedGroup := range approvedGroupList {
		groupIDList = append(groupIDList, approvedGroup.GroupID)
	}

	for _, unapprovedGroup := range unapprovedGroupList {
		groupIDList = append(groupIDList, unapprovedGroup.GroupID)
	}

	return groupIDList
}

func (r *groupQueryServiceImpl) fetchApprovedGroupList(userID string) ([]output.ApprovedGroup, error) {
	query := `
        SELECT
            group_users.group_id group_id,
            group_names.group_name group_name
        FROM
            group_users
        LEFT JOIN
            group_names
        ON
            group_users.group_id = group_names.id
        WHERE
            group_users.user_id = ?`

	rows, err := r.MySQLHandler.Conn.Queryx(query, userID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}
	defer rows.Close()

	approvedGroupList := make([]output.ApprovedGroup, 0)
	for rows.Next() {
		approvedGroup := output.ApprovedGroup{
			ApprovedUsersList:   make([]output.ApprovedUser, 0),
			UnapprovedUsersList: make([]output.UnapprovedUser, 0),
		}

		if err := rows.StructScan(&approvedGroup); err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		approvedGroupList = append(approvedGroupList, approvedGroup)
	}

	if err := rows.Err(); err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return approvedGroupList, nil
}

func (r *groupQueryServiceImpl) fetchUnApprovedGroupList(userID string) ([]output.UnapprovedGroup, error) {
	query := `
        SELECT
            group_unapproved_users.group_id group_id,
            group_names.group_name group_name
        FROM
            group_unapproved_users
        LEFT JOIN
            group_names
        ON
            group_unapproved_users.group_id = group_names.id
        WHERE
            group_unapproved_users.user_id = ?`

	rows, err := r.MySQLHandler.Conn.Queryx(query, userID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}
	defer rows.Close()

	unapprovedGroupList := make([]output.UnapprovedGroup, 0)
	for rows.Next() {
		unapprovedGroup := output.UnapprovedGroup{
			ApprovedUsersList:   make([]output.ApprovedUser, 0),
			UnapprovedUsersList: make([]output.UnapprovedUser, 0),
		}

		if err := rows.StructScan(&unapprovedGroup); err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		unapprovedGroupList = append(unapprovedGroupList, unapprovedGroup)
	}

	if err := rows.Err(); err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return unapprovedGroupList, nil
}

func (r *groupQueryServiceImpl) fetchApprovedUsersList(groupIDList []interface{}) ([]output.ApprovedUser, error) {
	query := `
        SELECT
            group_users.group_id group_id,
            group_users.user_id user_id,
            users.name user_name,
            group_users.color_code color_code
        FROM
            group_users
        LEFT JOIN
            users
        ON
            group_users.user_id = users.user_id
        WHERE
            group_users.group_id IN (?)`

	query, queryArgs, err := sqlx.In(query, groupIDList)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error1"))
	}

	rows, err := r.MySQLHandler.Conn.Queryx(query, queryArgs...)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}
	defer rows.Close()

	var approvedUsersList []output.ApprovedUser
	for rows.Next() {
		var approvedUser output.ApprovedUser
		if err := rows.StructScan(&approvedUser); err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		approvedUsersList = append(approvedUsersList, approvedUser)
	}

	if err := rows.Err(); err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return approvedUsersList, nil
}

func (r *groupQueryServiceImpl) fetchUnapprovedUsersList(groupIDList []interface{}) ([]output.UnapprovedUser, error) {
	query := `
        SELECT
            group_unapproved_users.group_id group_id,
            group_unapproved_users.user_id user_id,
            users.name user_name
        FROM
            group_unapproved_users
        LEFT JOIN
            users
        ON
            group_unapproved_users.user_id = users.user_id
        WHERE
            group_unapproved_users.group_id IN (?)`

	query, queryArgs, err := sqlx.In(query, groupIDList)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error1"))
	}

	rows, err := r.MySQLHandler.Conn.Queryx(query, queryArgs...)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}
	defer rows.Close()

	var unapprovedUsersList []output.UnapprovedUser
	for rows.Next() {
		var unapprovedUser output.UnapprovedUser
		if err := rows.StructScan(&unapprovedUser); err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		unapprovedUsersList = append(unapprovedUsersList, unapprovedUser)
	}

	if err := rows.Err(); err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return unapprovedUsersList, nil
}
