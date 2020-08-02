package infrastructure

import (
	"database/sql"
	"strings"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type GroupRepository struct {
	*MySQLHandler
}

func (r *GroupRepository) GetGroupList(userID string) ([]model.Group, error) {
	query := `
        SELECT
            group_users.group_id id,
            group_names.group_name group_name
        FROM
            group_users
        LEFT JOIN
            group_names
        ON
            group_users.group_id = group_names.id
        WHERE
            group_users.user_id = ?`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupList []model.Group
	for rows.Next() {
		var group model.Group
		if err := rows.StructScan(&group); err != nil {
			return nil, err
		}
		groupList = append(groupList, group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupList, nil
}

func (r *GroupRepository) GetGroupUsersList(groupList []model.Group) ([]model.GroupUser, error) {
	sliceQuery := make([]string, len(groupList))
	for i := 0; i < len(groupList); i++ {
		sliceQuery[i] = `
            SELECT
                group_users.group_id group_id,
                group_users.user_id user_id,
                users.name user_name
            FROM
                group_users
            LEFT JOIN
                users
            ON
                group_users.user_id = users.user_id
            WHERE
                group_users.group_id = ?`
	}

	query := strings.Join(sliceQuery, " UNION ")

	groupIDList := make([]interface{}, len(groupList))
	for i := 0; i < len(groupList); i++ {
		groupIDList[i] = groupList[i].GroupID
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupIDList...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupUsersList []model.GroupUser
	for rows.Next() {
		var groupUser model.GroupUser
		if err := rows.StructScan(&groupUser); err != nil {
			return nil, err
		}
		groupUsersList = append(groupUsersList, groupUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupUsersList, nil
}

func (r *GroupRepository) GetGroupUnapprovedUsersList(groupList []model.Group) ([]model.GroupUnapprovedUser, error) {
	sliceQuery := make([]string, len(groupList))
	for i := 0; i < len(groupList); i++ {
		sliceQuery[i] = `
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
                group_unapproved_users.group_id = ?`
	}

	query := strings.Join(sliceQuery, " UNION ")

	groupIDList := make([]interface{}, len(groupList))
	for i := 0; i < len(groupList); i++ {
		groupIDList[i] = groupList[i].GroupID
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupIDList...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupUnapprovedUsersList []model.GroupUnapprovedUser
	for rows.Next() {
		var groupUnapprovedUser model.GroupUnapprovedUser
		if err := rows.StructScan(&groupUnapprovedUser); err != nil {
			return nil, err
		}
		groupUnapprovedUsersList = append(groupUnapprovedUsersList, groupUnapprovedUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupUnapprovedUsersList, nil
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

func (r *GroupRepository) PostGroupUnapprovedUser(groupUnapprovedUser *model.GroupUnapprovedUser, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_unapproved_users
            (group_id, user_id)
        VALUES
            (?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupID, groupUnapprovedUser.UserID)

	return result, err
}

func (r *GroupRepository) GetGroupUnapprovedUser(groupUnapprovedUsersID int) (*model.GroupUnapprovedUser, error) {
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
            group_unapproved_users.id = ?`

	var groupUnapprovedUser model.GroupUnapprovedUser
	if err := r.MySQLHandler.conn.QueryRowx(query, groupUnapprovedUsersID).StructScan(&groupUnapprovedUser); err != nil {
		return nil, err
	}

	return &groupUnapprovedUser, nil
}

func (r *GroupRepository) FindGroupUser(groupID int, userID string) error {
	query := `SELECT
                  id 
              FROM
                  group_users
              WHERE 
                  group_id = ? 
              AND 
                  user_id = ?`

	var groupUserID int
	err := r.MySQLHandler.conn.QueryRowx(query, groupID, userID).Scan(&groupUserID)

	return err
}

func (r *GroupRepository) FindGroupUnapprovedUser(groupID int, userID string) error {
	query := `SELECT
                  id 
              FROM
                  group_unapproved_users
              WHERE 
                  group_id = ? 
              AND 
                  user_id = ?`

	var groupUnapprovedUserID int
	err := r.MySQLHandler.conn.QueryRowx(query, groupID, userID).Scan(&groupUnapprovedUserID)

	return err
}
