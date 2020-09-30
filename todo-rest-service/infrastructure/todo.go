package infrastructure

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type TodoRepository struct {
	*MySQLHandler
}

func NewTodoRepository(mysqlHandler *MySQLHandler) *TodoRepository {
	return &TodoRepository{mysqlHandler}
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
            implementation_date, updated_date DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	implementationTodoList := make([]model.Todo, 0)
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
            due_date, updated_date DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dueTodoList := make([]model.Todo, 0)
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
            implementation_date, updated_date DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	implementationTodoList := make([]model.Todo, 0)
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
            due_date, updated_date DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dueTodoList := make([]model.Todo, 0)
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

func (r *TodoRepository) GetTodo(todoId int) (*model.Todo, error) {
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
            id = ?`

	var todo model.Todo
	if err := r.MySQLHandler.conn.QueryRowx(query, todoId).StructScan(&todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) PostTodo(todo *model.Todo, userID string) (sql.Result, error) {
	query := `
        INSERT INTO todo_list
            (implementation_date, due_date, todo_content, user_id)
        VALUES
            (?,?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, todo.ImplementationDate, todo.DueDate, todo.TodoContent, userID)

	return result, err
}

func (r *TodoRepository) PutTodo(todo *model.Todo, todoID int) error {
	query := `
        UPDATE
            todo_list
        SET 
            implementation_date = ?,
            due_date = ?,
            todo_content = ?,
            complete_flag = ?
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, todo.ImplementationDate, todo.DueDate, todo.TodoContent, todo.CompleteFlag, todoID)

	return err
}

func (r *TodoRepository) DeleteTodo(todoID int) error {
	query := `
        DELETE
        FROM
            todo_list
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, todoID)
	return err
}

func (r *TodoRepository) SearchTodoList(todoSqlQuery string) ([]model.Todo, error) {
	rows, err := r.MySQLHandler.conn.Queryx(todoSqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searchTodoList []model.Todo
	for rows.Next() {
		var searchTodo model.Todo
		if err := rows.StructScan(&searchTodo); err != nil {
			return nil, err
		}
		searchTodoList = append(searchTodoList, searchTodo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return searchTodoList, nil
}
