package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type GroupTasksRepository struct {
	*MySQLHandler
}

func (r *GroupTasksRepository) GetGroupTasksUser(groupTasksUser model.GroupTasksUser, groupID int) (*model.GroupTasksUser, error) {
	query := `
        SELECT
            id,
            user_id,
            group_id
        FROM
            group_tasks_users
        WHERE
            user_id = ?
        AND
            group_id = ?`

	if err := r.MySQLHandler.conn.QueryRowx(query, groupTasksUser.UserID, groupID).StructScan(&groupTasksUser); err != nil {
		return nil, err
	}

	return &groupTasksUser, nil
}

func (r *GroupTasksRepository) PostGroupTasksUser(groupTasksUser model.GroupTasksUser, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_tasks_users
            (user_id, group_id)
        VALUES
            (?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupTasksUser.UserID, groupID)

	return result, err
}
