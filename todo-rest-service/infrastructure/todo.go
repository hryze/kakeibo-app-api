package infrastructure

import (
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type TodoRepository struct {
	*MySQLHandler
}

func (r *TodoRepository) GetDailyImplementationTodoList(date time.Time, userID string) ([]model.Todo, error) {
	query := `
        SELECT
            id,
            posted_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag
        FROM
            todo_list
        WHERE
            user_id = ?
        AND
            implementation_date = ?
        ORDER BY
            implementation_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var implementationTodoList []model.Todo
	for rows.Next() {
		var implementationTodo model.Todo
		if err := rows.StructScan(&implementationTodo); err != nil {
			return nil, err
		}
		implementationTodoList = append(implementationTodoList, implementationTodo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return implementationTodoList, nil
}

func (r *TodoRepository) GetDailyDueTodoList(date time.Time, userID string) ([]model.Todo, error) {
	query := `
        SELECT
            id,
            posted_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag
        FROM
            todo_list
        WHERE
            user_id = ?
        AND
            due_date = ?
        ORDER BY
            due_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dueTodoList []model.Todo
	for rows.Next() {
		var dueTodo model.Todo
		if err := rows.StructScan(&dueTodo); err != nil {
			return nil, err
		}
		dueTodoList = append(dueTodoList, dueTodo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dueTodoList, nil
}

func (r *TodoRepository) GetMonthlyImplementationTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error) {
	query := `
        SELECT
            id,
            posted_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag
        FROM
            todo_list
        WHERE
            user_id = ?
        AND
            implementation_date >= ?
        AND
            implementation_date <= ?
        ORDER BY
            implementation_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var implementationTodoList []model.Todo
	for rows.Next() {
		var implementationTodo model.Todo
		if err := rows.StructScan(&implementationTodo); err != nil {
			return nil, err
		}
		implementationTodoList = append(implementationTodoList, implementationTodo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return implementationTodoList, nil
}

func (r *TodoRepository) GetMonthlyDueTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error) {
	query := `
        SELECT
            id,
            posted_date,
            implementation_date,
            due_date,
            todo_content,
            complete_flag
        FROM
            todo_list
        WHERE
            user_id = ?
        AND
            due_date >= ?
        AND
            due_date <= ?
        ORDER BY
            due_date`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dueTodoList []model.Todo
	for rows.Next() {
		var dueTodo model.Todo
		if err := rows.StructScan(&dueTodo); err != nil {
			return nil, err
		}
		dueTodoList = append(dueTodoList, dueTodo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dueTodoList, nil
}
