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

func (r *groupRepository) StoreGroupAndApprovedUser(group *groupdomain.Group, userID userdomain.UserID) (*groupdomain.Group, error) {
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
		result, err := tx.Exec(storeGroupQuery, group.GroupName().Value())
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

	group = groupdomain.NewGroup(groupIDVo, group.GroupName())

	return group, nil
}

func (r *groupRepository) DeleteGroupAndApprovedUser(group *groupdomain.Group) error {
	deleteApprovedUserQuery := `
        DELETE
        FROM
            group_users
        WHERE
            group_id = ?`

	deleteGroupQuery := `
        DELETE
        FROM
            group_names
        WHERE
            id = ?`

	tx, err := r.MySQLHandler.Conn.Begin()
	if err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	groupID, err := group.ID()
	if err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(deleteApprovedUserQuery, groupID.Value()); err != nil {
			return err
		}

		if _, err := tx.Exec(deleteGroupQuery, groupID.Value()); err != nil {
			return err
		}

		return nil
	}

	if err := transactions(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	if err := tx.Commit(); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}
