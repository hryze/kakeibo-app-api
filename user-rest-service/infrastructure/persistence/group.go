package persistence

import (
	"database/sql"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/groupdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/datasource"
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
		result, err := tx.Exec(deleteApprovedUserQuery, groupID.Value())
		if err != nil {
			return err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return err
		} else if rowsAffected != 1 {
			return xerrors.Errorf("affected rows must be a single row: %d", rowsAffected)
		}

		result, err = tx.Exec(deleteGroupQuery, groupID.Value())
		if err != nil {
			return err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return err
		} else if rowsAffected != 1 {
			return xerrors.Errorf("affected rows must be a single row: %d", rowsAffected)
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

func (r *groupRepository) UpdateGroupName(group *groupdomain.Group) error {
	query := `
        UPDATE
            group_names
        SET 
            group_name = ?
        WHERE
            id = ?`

	tx, err := r.MySQLHandler.Conn.Begin()
	if err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	transactions := func(tx *sql.Tx) error {
		groupID, err := group.ID()
		if err != nil {
			return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		result, err := tx.Exec(query, group.GroupName().Value(), groupID.Value())
		if err != nil {
			return err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return err
		} else if rowsAffected != 1 {
			return xerrors.Errorf("affected rows must be a single row: %d", rowsAffected)
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

func (r *groupRepository) StoreUnapprovedUser(unapprovedUser *groupdomain.UnapprovedUser) error {
	query := `
        INSERT INTO group_unapproved_users
            (group_id, user_id)
        VALUES
            (?, ?)`

	if _, err := r.MySQLHandler.Conn.Exec(query, unapprovedUser.GroupID().Value(), unapprovedUser.UserID().Value()); err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return nil
}

func (r *groupRepository) DeleteApprovedUser(approvedUser *groupdomain.ApprovedUser) error {
	query := `
        DELETE
        FROM
            group_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	tx, err := r.MySQLHandler.Conn.Begin()
	if err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	transactions := func(tx *sql.Tx) error {
		result, err := tx.Exec(query, approvedUser.GroupID().Value(), approvedUser.UserID().Value())
		if err != nil {
			return err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return err
		} else if rowsAffected != 1 {
			return xerrors.Errorf("affected rows must be a single row: %d", rowsAffected)
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

func (r *groupRepository) StoreApprovedUser(approvedUser *groupdomain.ApprovedUser) error {
	insertApprovedUserQuery := `
        INSERT INTO group_users
            (group_id, user_id, color_code)
        VALUES
            (?, ?, ?)`

	deleteUnapprovedUserQuery := `
        DELETE
        FROM
            group_unapproved_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	tx, err := r.MySQLHandler.Conn.Begin()
	if err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(insertApprovedUserQuery, approvedUser.GroupID().Value(), approvedUser.UserID().Value(), approvedUser.ColorCode().Value()); err != nil {
			return err
		}

		result, err := tx.Exec(deleteUnapprovedUserQuery, approvedUser.GroupID().Value(), approvedUser.UserID().Value())
		if err != nil {
			return err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return err
		} else if rowsAffected != 1 {
			return xerrors.Errorf("affected rows must be a single row: %d", rowsAffected)
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

func (r *groupRepository) DeleteUnapprovedUser(unapprovedUser *groupdomain.UnapprovedUser) error {
	query := `
        DELETE 
        FROM
            group_unapproved_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	tx, err := r.MySQLHandler.Conn.Begin()
	if err != nil {
		return apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	transactions := func(tx *sql.Tx) error {
		result, err := tx.Exec(query, unapprovedUser.GroupID().Value(), unapprovedUser.UserID().Value())
		if err != nil {
			return err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return err
		} else if rowsAffected != 1 {
			return xerrors.Errorf("affected rows must be a single row: %d", rowsAffected)
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

func (r *groupRepository) FindGroupByID(groupID *groupdomain.GroupID) (*groupdomain.Group, error) {
	query := `
        SELECT
            id,
            group_name
        FROM
            group_names
        WHERE
            id = ?`

	var groupDto datasource.Group
	if err := r.MySQLHandler.Conn.QueryRowx(query, groupID.Value()).StructScan(&groupDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("指定されたグループは存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	groupIDVo, err := groupdomain.NewGroupID(groupDto.GroupID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	groupNameVo, err := groupdomain.NewGroupName(groupDto.GroupName)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	group := groupdomain.NewGroup(groupIDVo, groupNameVo)

	return group, nil
}

func (r *groupRepository) FindApprovedUser(groupID groupdomain.GroupID, userID userdomain.UserID) (*groupdomain.ApprovedUser, error) {
	query := `
        SELECT
            group_id,
            user_id,
            color_code
        FROM
            group_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	var approvedUserDto datasource.ApprovedUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, groupID.Value(), userID.Value()).StructScan(&approvedUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	groupIDVo, err := groupdomain.NewGroupID(approvedUserDto.GroupID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	userIDVo, err := userdomain.NewUserID(approvedUserDto.UserID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	colorCodeVo, err := groupdomain.NewColorCode(approvedUserDto.ColorCode)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	approvedUser := groupdomain.NewApprovedUser(groupIDVo, userIDVo, colorCodeVo)

	return approvedUser, nil
}

func (r *groupRepository) FindUnapprovedUser(groupID groupdomain.GroupID, userID userdomain.UserID) (*groupdomain.UnapprovedUser, error) {
	query := `
        SELECT
            group_id,
            user_id
        FROM
            group_unapproved_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	var unapprovedUserDto datasource.UnapprovedUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, groupID.Value(), userID.Value()).StructScan(&unapprovedUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.NewNotFoundError(apierrors.NewErrorString("ユーザーが存在しません"))
		}

		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	groupIDVo, err := groupdomain.NewGroupID(unapprovedUserDto.GroupID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	userIDVo, err := userdomain.NewUserID(unapprovedUserDto.UserID)
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	unapprovedUser := groupdomain.NewUnapprovedUser(groupIDVo, userIDVo)

	return unapprovedUser, nil
}

func (r *groupRepository) FetchApprovedUserIDList(groupID groupdomain.GroupID) ([]userdomain.UserID, error) {
	query := `
        SELECT
            user_id
        FROM
            group_users
        WHERE
            group_id = ?`

	rows, err := r.MySQLHandler.Conn.Queryx(query, groupID.Value())
	if err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}
	defer rows.Close()

	approvedUserIDList := make([]userdomain.UserID, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		userIDVo, err := userdomain.NewUserID(userID)
		if err != nil {
			return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
		}

		approvedUserIDList = append(approvedUserIDList, userIDVo)
	}

	if err := rows.Err(); err != nil {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return approvedUserIDList, nil
}
