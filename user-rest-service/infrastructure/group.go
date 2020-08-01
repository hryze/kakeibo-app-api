package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type GroupRepository struct {
	*MySQLHandler
}

func (r *GroupRepository) GetGroup(groupID int) (*model.Group, error) {
	query := `
        SELECT
            id,
            group_name
        FROM
            group_names
        WHERE
            id = ?`

	var group model.Group
	if err := r.MySQLHandler.conn.QueryRowx(query, groupID).StructScan(&group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) PostGroupAndGroupUser(group *model.Group, userID string) (sql.Result, error) {
	groupQuery := `
        INSERT INTO group_names
            (group_name)
        VALUES
            (?)`

	groupUserQuery := `
        INSERT INTO group_users
            (group_id, user_id)
        VALUES
            (?,?)`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return nil, err
	}

	transactions := func(tx *sql.Tx) (sql.Result, error) {
		result, err := tx.Exec(groupQuery, group.GroupName)
		if err != nil {
			return nil, err
		}

		groupLastInsertId, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		if _, err := tx.Exec(groupUserQuery, int(groupLastInsertId), userID); err != nil {
			return nil, err
		}

		return result, nil
	}

	result, err := transactions(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *GroupRepository) DeleteGroupAndGroupUser(groupID int, userID string) error {
	groupQuery := `
        DELETE
        FROM
            group_names
        WHERE
            id = ?`

	groupUserQuery := `
        DELETE
        FROM
            group_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {

		if _, err := tx.Exec(groupUserQuery, groupID, userID); err != nil {
			return err
		}

		if _, err := tx.Exec(groupQuery, groupID); err != nil {
			return err
		}
		return nil
	}

	if err := transactions(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *GroupRepository) PutGroup(group *model.Group, groupID int) error {
	query := `
        UPDATE
            group_names
        SET 
            group_name = ?
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, group.GroupName, groupID)
	return err
}
