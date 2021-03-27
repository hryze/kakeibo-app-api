package infrastructure

import (
	"database/sql"
	"strings"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/rdb"
)

type GroupRepository struct {
	*rdb.MySQLHandler
}

func NewGroupRepository(mysqlHandler *rdb.MySQLHandler) *GroupRepository {
	return &GroupRepository{mysqlHandler}
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
	if err := r.MySQLHandler.Conn.QueryRowx(query, groupID).StructScan(&group); err != nil {
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

	_, err := r.MySQLHandler.Conn.Exec(query, group.GroupName, groupID)

	return err
}

func (r *GroupRepository) PostUnapprovedUser(unapprovedUser *model.UnapprovedUser, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_unapproved_users
            (group_id, user_id)
        VALUES
            (?,?)`

	result, err := r.MySQLHandler.Conn.Exec(query, groupID, unapprovedUser.UserID)

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
	if err := r.MySQLHandler.Conn.QueryRowx(query, groupUnapprovedUsersID).StructScan(&unapprovedUser); err != nil {
		return nil, err
	}

	return &unapprovedUser, nil
}

func (r *GroupRepository) FindApprovedUser(groupID int, userID string) error {
	query := `
        SELECT
            id 
        FROM
            group_users
        WHERE 
            group_id = ? 
        AND 
            user_id = ?`

	var groupUsersID int
	err := r.MySQLHandler.Conn.QueryRowx(query, groupID, userID).Scan(&groupUsersID)

	return err
}

func (r *GroupRepository) FindUnapprovedUser(groupID int, userID string) error {
	query := `
        SELECT
            id 
        FROM
            group_unapproved_users
        WHERE 
            group_id = ? 
        AND 
            user_id = ?`

	var groupUnapprovedUsersID int
	err := r.MySQLHandler.Conn.QueryRowx(query, groupID, userID).Scan(&groupUnapprovedUsersID)

	return err
}

func (r *GroupRepository) PostGroupApprovedUserAndDeleteGroupUnapprovedUser(groupID int, userID string, colorCode string) (sql.Result, error) {
	insertApprovedUserQuery := `
        INSERT INTO group_users
            (group_id, user_id, color_code)
        VALUES
            (?,?,?)`

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
		return nil, err
	}

	transactions := func(tx *sql.Tx) (sql.Result, error) {
		result, err := tx.Exec(insertApprovedUserQuery, groupID, userID, colorCode)
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
            users.name user_name,
            group_users.color_code color_code
        FROM
            group_users
        LEFT JOIN
            users
        ON
            group_users.user_id = users.user_id
        WHERE
            group_users.id = ?`

	var approvedUser model.ApprovedUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, approvedUsersID).StructScan(&approvedUser); err != nil {
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

	_, err := r.MySQLHandler.Conn.Exec(query, groupID, userID)

	return err
}

func (r *GroupRepository) DeleteGroupApprovedUser(groupID int, userID string) error {
	query := `
        DELETE 
        FROM
            group_users
        WHERE
            group_id = ?
        AND
            user_id = ?`

	_, err := r.MySQLHandler.Conn.Exec(query, groupID, userID)

	return err
}

func (r *GroupRepository) FindApprovedUsersList(groupID int, groupUsersList []string) (model.GroupTasksUsersListReceiver, error) {
	sliceQuery := make([]string, len(groupUsersList))
	for i := 0; i < len(groupUsersList); i++ {
		sliceQuery[i] = `
            SELECT
                user_id
            FROM
                group_users
            WHERE
                group_id = ?
            AND
                user_id = ?`
	}

	query := strings.Join(sliceQuery, " UNION ")

	var queryArgs []interface{}
	for _, userID := range groupUsersList {
		queryArgs = append(queryArgs, groupID, userID)
	}

	var dbGroupUsersList model.GroupTasksUsersListReceiver
	rows, err := r.MySQLHandler.Conn.Queryx(query, queryArgs...)
	if err != nil {
		return dbGroupUsersList, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return dbGroupUsersList, err
		}

		dbGroupUsersList.GroupUsersList = append(dbGroupUsersList.GroupUsersList, userID)
	}

	if err := rows.Err(); err != nil {
		return dbGroupUsersList, err
	}

	return dbGroupUsersList, nil
}

func (r *GroupRepository) GetGroupUsersList(groupID int) ([]string, error) {
	query := `
        SELECT
            user_id
        FROM
            group_users
        WHERE
            group_id = ?`

	rows, err := r.MySQLHandler.Conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupUserIDList := make([]string, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		groupUserIDList = append(groupUserIDList, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupUserIDList, nil
}
