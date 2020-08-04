package infrastructure

import (
	"database/sql"
	"strings"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type GroupRepository struct {
	*MySQLHandler
}

func (r *GroupRepository) GetApprovedGroupList(userID string) ([]model.ApprovedGroup, error) {
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

	rows, err := r.MySQLHandler.conn.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvedGroupList []model.ApprovedGroup
	for rows.Next() {
		var approvedGroup model.ApprovedGroup
		if err := rows.StructScan(&approvedGroup); err != nil {
			return nil, err
		}
		approvedGroupList = append(approvedGroupList, approvedGroup)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return approvedGroupList, nil
}

func (r *GroupRepository) GetUnApprovedGroupList(userID string) ([]model.UnapprovedGroup, error) {
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

	rows, err := r.MySQLHandler.conn.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unapprovedGroupList []model.UnapprovedGroup
	for rows.Next() {
		var unapprovedGroup model.UnapprovedGroup
		if err := rows.StructScan(&unapprovedGroup); err != nil {
			return nil, err
		}
		unapprovedGroupList = append(unapprovedGroupList, unapprovedGroup)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return unapprovedGroupList, nil
}

func (r *GroupRepository) GetApprovedUsersList(groupIDList []interface{}) ([]model.ApprovedUser, error) {
	sliceQuery := make([]string, len(groupIDList))
	for i := 0; i < len(groupIDList); i++ {
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

	rows, err := r.MySQLHandler.conn.Queryx(query, groupIDList...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvedUsersList []model.ApprovedUser
	for rows.Next() {
		var approvedUser model.ApprovedUser
		if err := rows.StructScan(&approvedUser); err != nil {
			return nil, err
		}
		approvedUsersList = append(approvedUsersList, approvedUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return approvedUsersList, nil
}

func (r *GroupRepository) GetUnapprovedUsersList(groupIDList []interface{}) ([]model.UnapprovedUser, error) {
	sliceQuery := make([]string, len(groupIDList))
	for i := 0; i < len(groupIDList); i++ {
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

	rows, err := r.MySQLHandler.conn.Queryx(query, groupIDList...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unapprovedUsersList []model.UnapprovedUser
	for rows.Next() {
		var unapprovedUser model.UnapprovedUser
		if err := rows.StructScan(&unapprovedUser); err != nil {
			return nil, err
		}
		unapprovedUsersList = append(unapprovedUsersList, unapprovedUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return unapprovedUsersList, nil
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

func (r *GroupRepository) PostGroupAndApprovedUser(group *model.Group, userID string) (sql.Result, error) {
	groupQuery := `
        INSERT INTO group_names
            (group_name)
        VALUES
            (?)`

	approvedUserQuery := `
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

		if _, err := tx.Exec(approvedUserQuery, int(groupLastInsertId), userID); err != nil {
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

func (r *GroupRepository) DeleteGroupAndApprovedUser(groupID int, userID string) error {
	groupQuery := `
        DELETE
        FROM
            group_names
        WHERE
            id = ?`

	approvedUserQuery := `
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

		if _, err := tx.Exec(approvedUserQuery, groupID, userID); err != nil {
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

func (r *GroupRepository) PostUnapprovedUser(unapprovedUser *model.UnapprovedUser, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_unapproved_users
            (group_id, user_id)
        VALUES
            (?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupID, unapprovedUser.UserID)

	return result, err
}

func (r *GroupRepository) GetUnapprovedUser(groupUnapprovedUsersID int) (*model.UnapprovedUser, error) {
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

	var unapprovedUser model.UnapprovedUser
	if err := r.MySQLHandler.conn.QueryRowx(query, groupUnapprovedUsersID).StructScan(&unapprovedUser); err != nil {
		return nil, err
	}

	return &unapprovedUser, nil
}

func (r *GroupRepository) FindApprovedUser(groupID int, userID string) error {
	query := `SELECT
                  id 
              FROM
                  group_users
              WHERE 
                  group_id = ? 
              AND 
                  user_id = ?`

	var groupUsersID int
	err := r.MySQLHandler.conn.QueryRowx(query, groupID, userID).Scan(&groupUsersID)

	return err
}

func (r *GroupRepository) FindUnapprovedUser(groupID int, userID string) error {
	query := `SELECT
                  id 
              FROM
                  group_unapproved_users
              WHERE 
                  group_id = ? 
              AND 
                  user_id = ?`

	var groupUnapprovedUsersID int
	err := r.MySQLHandler.conn.QueryRowx(query, groupID, userID).Scan(&groupUnapprovedUsersID)

	return err
}

func (r *GroupRepository) PostGroupApprovedUserAndDeleteGroupUnapprovedUser(groupID int, userID string) (sql.Result, error) {
	insertApprovedUserQuery := `
        INSERT INTO group_users
            (group_id, user_id)
        VALUES
            (?,?)`

	deleteUnapprovedUserQuery := `
        DELETE
        FROM
            group_unapproved_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return nil, err
	}

	transactions := func(tx *sql.Tx) (sql.Result, error) {
		result, err := tx.Exec(insertApprovedUserQuery, groupID, userID)
		if err != nil {
			return nil, err
		}

		if _, err := tx.Exec(deleteUnapprovedUserQuery, groupID, userID); err != nil {
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

func (r *GroupRepository) GetApprovedUser(approvedUsersID int) (*model.ApprovedUser, error) {
	query := `
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
            group_users.id = ?`

	var approvedUser model.ApprovedUser
	if err := r.MySQLHandler.conn.QueryRowx(query, approvedUsersID).StructScan(&approvedUser); err != nil {
		return nil, err
	}

	return &approvedUser, nil
}

func (r *GroupRepository) DeleteGroupUnapprovedUser(groupID int, userID string) error {
	query := `
        DELETE 
        FROM
            group_unapproved_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupID, userID)
	return err
}
