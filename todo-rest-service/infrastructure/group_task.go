package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type GroupTasksRepository struct {
	*MySQLHandler
}

func NewGroupTasksRepository(mysqlHandler *MySQLHandler) *GroupTasksRepository {
	return &GroupTasksRepository{mysqlHandler}
}

func (r *GroupTasksRepository) GetGroupTasksUsersList(groupID int) ([]model.GroupTasksUser, error) {
	query := `
        SELECT
            id,
            user_id,
            group_id
        FROM
            group_tasks_users
        WHERE
            group_id = ?
        ORDER BY
            id`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupTasksUsersList := make([]model.GroupTasksUser, 0)
	for rows.Next() {
		groupTasksUser := model.GroupTasksUser{TasksList: make([]model.GroupTask, 0)}
		if err := rows.StructScan(&groupTasksUser); err != nil {
			return nil, err
		}
		groupTasksUsersList = append(groupTasksUsersList, groupTasksUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupTasksUsersList, nil
}

func (r *GroupTasksRepository) GetGroupTasksListAssignedToUser(groupID int) ([]model.GroupTask, error) {
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
            group_id = ?
        AND
            group_tasks_users_id IS NOT NULL
        ORDER BY
            id`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupTasksListAssignedToUser []model.GroupTask
	for rows.Next() {
		var groupTaskAssignedToUser model.GroupTask
		if err := rows.StructScan(&groupTaskAssignedToUser); err != nil {
			return nil, err
		}
		groupTasksListAssignedToUser = append(groupTasksListAssignedToUser, groupTaskAssignedToUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupTasksListAssignedToUser, nil
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

func (r *GroupTasksRepository) GetGroupTasksList(groupID int) ([]model.GroupTask, error) {
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
            group_id = ?
        ORDER BY
            id`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupTasksList []model.GroupTask
	for rows.Next() {
		var groupTask model.GroupTask
		if err := rows.StructScan(&groupTask); err != nil {
			return nil, err
		}
		groupTasksList = append(groupTasksList, groupTask)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupTasksList, nil
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

func (r *GroupTasksRepository) PutGroupTask(groupTask *model.GroupTask, groupTasksID int) error {
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

	_, err := r.MySQLHandler.conn.Exec(query, groupTask.BaseDate, groupTask.CycleType, groupTask.Cycle, groupTask.TaskName, groupTask.GroupTasksUserID, groupTasksID)

	return err
}

func (r *GroupTasksRepository) DeleteGroupTask(groupTasksID int) error {
	query := `
        DELETE
        FROM
            group_tasks
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupTasksID)

	return err
}
