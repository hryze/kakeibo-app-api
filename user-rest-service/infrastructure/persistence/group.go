package persistence

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/groupdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/rdb"
)

type groupRepository struct {
	*rdb.MySQLHandler
}

func NewGroupRepository(mysqlHandler *rdb.MySQLHandler) *groupRepository {
	return &groupRepository{mysqlHandler}
}

func (r *groupRepository) StoreGroupAndApprovedUser(groupName groupdomain.GroupName, userID userdomain.UserID) (*groupdomain.Group, error) {
	storeGroupQuery := `
        INSERT INTO group_names
            (group_name)
        VALUES
            (?)`

	storeApprovedUserQuery := `
        INSERT INTO group_users
            (group_id, user_id, color_code)
        VALUES
            (?, ?, "#FF0000")`

	tx, err := r.MySQLHandler.Conn.Begin()
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	transactions := func(tx *sql.Tx) (int, error) {
		result, err := tx.Exec(storeGroupQuery, groupName.Value())
		if err != nil {
			return 0, err
		}

		groupID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		if _, err := tx.Exec(storeApprovedUserQuery, groupID, userID.Value()); err != nil {
			return 0, err
		}

		return int(groupID), nil
	}

	groupID, err := transactions(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	if err := tx.Commit(); err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	groupIDVo, err := groupdomain.NewGroupID(groupID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	group := groupdomain.NewGroup(groupIDVo, groupName)

	return group, nil
}

func (r *groupRepository) DeleteGroupAndApprovedUser(groupID groupdomain.GroupID) error {
	query := `
	   DELETE
	       group_users, group_names
	   FROM
	       group_users
	   LEFT JOIN
	       group_names
	   ON
	       group_users.group_id = group_names.id
	   WHERE
	       group_users.group_id = ?`

	if _, err := r.MySQLHandler.Conn.Exec(query, groupID.Value()); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}
