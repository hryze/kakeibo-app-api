package infrastructure

import (
	"database/sql"
	"strings"

	"github.com/hryze/kakeibo-app-api/todo-rest-service/domain/model"
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

func (r *GroupTasksRepository) PutGroupTasksListAssignedToUser(groupTasksList []model.GroupTask, updateTaskIndexList []int) error {
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

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, index := range updateTaskIndexList {
			if _, err := r.MySQLHandler.conn.Exec(query, groupTasksList[index].BaseDate, groupTasksList[index].CycleType, groupTasksList[index].Cycle, groupTasksList[index].TaskName, groupTasksList[index].GroupTasksUserID, groupTasksList[index].ID); err != nil {
				return err
			}
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

func (r *GroupTasksRepository) PostGroupTasksUsersList(groupTasksUsersList model.GroupTasksUsersListReceiver, groupID int) error {
	query := `
        INSERT INTO group_tasks_users
            (user_id, group_id)
        VALUES`

	var queryArgs []interface{}
	for _, userID := range groupTasksUsersList.GroupUsersList {
		query += "(?, ?),"
		queryArgs = append(queryArgs, userID, groupID)
	}

	query = strings.TrimSuffix(query, ",")

	_, err := r.MySQLHandler.conn.Exec(query, queryArgs...)

	return err
}

func (r *GroupTasksRepository) GetGroupTasksIDListAssignedToUser(groupTasksUsersIdList []int, groupID int) ([]int, error) {
	sliceQuery := make([]string, len(groupTasksUsersIdList))
	for i := 0; i < len(groupTasksUsersIdList); i++ {
		sliceQuery[i] = `
            SELECT
                id
            FROM
                group_tasks
            WHERE
                group_id = ?
            AND
                group_tasks_users_id = ?`
	}

	query := strings.Join(sliceQuery, " UNION ")

	var queryArgs []interface{}
	for _, groupTasksUsersId := range groupTasksUsersIdList {
		queryArgs = append(queryArgs, groupID, groupTasksUsersId)
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupTasksIDList []int
	for rows.Next() {
		var groupTasksID int
		if err := rows.Scan(&groupTasksID); err != nil {
			return nil, err
		}

		groupTasksIDList = append(groupTasksIDList, groupTasksID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupTasksIDList, nil
}

func (r *GroupTasksRepository) DeleteGroupTasksUsersList(groupTasksUsersListReceiver model.GroupTasksUsersListReceiver, groupTasksIDList []int, groupID int) error {
	deleteQuery := `
        DELETE
        FROM
            group_tasks_users
        WHERE
            user_id = ?
        AND
            group_id = ?`

	updateQuery := `
        UPDATE
            group_tasks
        SET
            base_date = ?,
            cycle_type = ?,
            cycle = ?
        WHERE
            id = ?`

	var deleteQueryArgs []interface{}
	for _, userID := range groupTasksUsersListReceiver.GroupUsersList {
		deleteQueryArgs = append(deleteQueryArgs, userID, groupID)
	}

	var updateQueryArgs []interface{}
	for _, taskID := range groupTasksIDList {
		updateQueryArgs = append(updateQueryArgs, nil, nil, nil, taskID)
	}

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for i := 0; i < len(deleteQueryArgs); i = i + 2 {
			queryArgs := deleteQueryArgs[i : i+2]

			if _, err := tx.Exec(deleteQuery, queryArgs...); err != nil {
				return err
			}
		}

		for i := 0; i < len(updateQueryArgs); i = i + 4 {
			queryArgs := updateQueryArgs[i : i+4]

			if _, err := tx.Exec(updateQuery, queryArgs...); err != nil {
				return err
			}
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

	groupTasksList := make([]model.GroupTask, 0)
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

func (r *GroupTasksRepository) PutGroupTask(groupTask *model.GroupTask, groupTasksID int) (sql.Result, error) {
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

	result, err := r.MySQLHandler.conn.Exec(query, groupTask.BaseDate, groupTask.CycleType, groupTask.Cycle, groupTask.TaskName, groupTask.GroupTasksUserID, groupTasksID)

	return result, err
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
