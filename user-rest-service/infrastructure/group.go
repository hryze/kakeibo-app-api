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

func (r *GroupRepository) PostGroup(group *model.Group) (sql.Result, error) {
	query := `
        INSERT INTO group_names
            (group_name)
        VALUES
            (?)`

	result, err := r.MySQLHandler.conn.Exec(query, group.GroupName)
	return result, err
}

func (r *GroupRepository) DeleteGroup(groupID int) error {
	query := `
        DELETE
        FROM
            group_names
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupID)
	return err
}

func (r *GroupRepository) PostGroupUser(groupID int, userID string) (sql.Result, error) {
	query := `
        INSERT INTO group_users
            (group_id, user_id)
        VALUES
            (?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupID, userID)
	return result, err
}

func (r *GroupRepository) DeleteGroupUser(groupID int, userID string) error {
	query := `
        DELETE
        FROM
            group_users
        WHERE
            id = ?
        AND
            userID = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupID, userID)
	return err
}
