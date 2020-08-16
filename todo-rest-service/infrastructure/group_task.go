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

func (r *GroupTasksRepository) GetGroupTask(groupTasksID int) (*model.GroupTask, error) {
	query := `
        SELECT
            id,
            base_date,
            cycle_type,
            cycle,
            task_name,
            group_id,
            group_tasks_users_id
        FROM
            group_tasks
        WHERE
            id = ?`

	var groupTask model.GroupTask
	if err := r.MySQLHandler.conn.QueryRowx(query, groupTasksID).StructScan(&groupTask); err != nil {
		return nil, err
	}

	return &groupTask, nil
}

func (r *GroupTasksRepository) PostGroupTask(groupTask model.GroupTask, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_tasks
            (task_name, group_id)
        VALUES
            (?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupTask.TaskName, groupID)

	return result, err
}

func (r *GroupTasksRepository) PutGroupTask(groupTask *model.GroupTask, groupTodoID int) error {
	query := `
        UPDATE
            group_tasks
        SET 
            base_date = ?,
            cycle_type = ?,
            cycle = ?,
            task_name = ?,
            group_tasks_users_id = ?
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupTask.BaseDate, groupTask.CycleType, groupTask.Cycle, groupTask.TaskName, groupTask.GroupTasksUserID, groupTodoID)

	return err
}
