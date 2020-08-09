package infrastructure

import (
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type GroupTodoRepository struct {
	*MySQLHandler
}

func (r *GroupTodoRepository) GetDailyImplementationGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	query := `
        SELECT
            id,
            posted_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            group_id = ?
        AND
            implementation_date = ?
        ORDER BY
            implementation_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var implementationGroupTodoList []model.GroupTodo
	for rows.Next() {
		var implementationGroupTodo model.GroupTodo
		if err := rows.StructScan(&implementationGroupTodo); err != nil {
			return nil, err
		}
		implementationGroupTodoList = append(implementationGroupTodoList, implementationGroupTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return implementationGroupTodoList, nil
}

func (r *TodoRepository) GetDailyDueGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error) {
	query := `
        SELECT
            id,
            posted_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag,
            user_id
        FROM
            group_todo_list
        WHERE
            group_id = ?
        AND
            due_date = ?
        ORDER BY
            due_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dueGroupTodoList []model.GroupTodo
	for rows.Next() {
		var groupDueTodo model.GroupTodo
		if err := rows.StructScan(&groupDueTodo); err != nil {
			return nil, err
		}
		dueGroupTodoList = append(dueGroupTodoList, groupDueTodo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dueGroupTodoList, nil
}
